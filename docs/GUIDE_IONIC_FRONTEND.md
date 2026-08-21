# Ionic frontend guide

Everything here is explained with code from this repository. Paths are relative
to `frontend/`.

---

## 1. Ionic and Angular: who does what

They are two separate things stacked on top of each other.

| | Angular | Ionic |
| --- | --- | --- |
| Provides | Components, dependency injection, routing, change detection, HTTP | A library of UI web components (`ion-*`) plus overlays and gestures |
| You use it for | Application structure and logic | Everything the user sees and touches |

Ionic components are **web components**, not Angular components.
`@ionic/angular` is a thin wrapper that makes them behave like Angular ones
(bindings, forms integration, router integration). This matters in practice: an
`ion-input` emits a DOM `ionChange` event, which is why you see
`(ionChange)="..."` rather than `(change)`.

You could learn Ionic with React or Vue instead. Angular is used here because it
is what the sibling `ductifact` project uses.

---

## 2. How the app starts

Three files, in order.

**`src/main.ts`** — the entry point. Modern Angular bootstraps a *component*, not
a module.

```ts
bootstrapApplication(AppComponent, appConfig).catch((error) => console.error(error));
```

**`src/app/app.config.ts`** — every application-wide provider.

```ts
export const appConfig: ApplicationConfig = {
  providers: [
    provideZoneChangeDetection({ eventCoalescing: true }),
    provideRouter(routes, withComponentInputBinding()),
    provideIonicAngular({ mode: 'md' }),
    provideHttpClient(),
    provideAppInitializer(() => inject(ConfigService).load()),
  ],
};
```

- `provideIonicAngular({ mode: 'md' })` forces Material styling everywhere.
  Without it Ionic switches to iOS styling on Apple devices — great for a native
  app, confusing for a web app that should look the same to everyone.
- `provideAppInitializer` returns a promise; Angular waits for it before
  rendering anything. That is how the API URL is known before the first request.

**`src/app/app.component.ts`** — the shell. `ion-app` establishes the layout
context (safe areas, overlay containers) that every other Ionic component
assumes exists.

```ts
template: `
  <ion-app>
    <ion-router-outlet />
  </ion-app>
`,
```

`ion-router-outlet` replaces Angular's `router-outlet`. It keeps pages in the DOM
and adds native-feeling transitions.

---

## 3. Standalone components — why there are no NgModules

Every component declares its own dependencies:

```ts
@Component({
  selector: 'app-product-card',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [IonButton, IonIcon, StatusPillComponent],
  templateUrl: './product-card.component.html',
  styleUrl: './product-card.component.scss',
})
export class ProductCardComponent { ... }
```

Three things worth noticing:

- `imports` lists exactly what the template uses. Forget one and the build fails
  with a clear message. Import something unused and the bundler drops it.
- `OnPush` tells Angular to re-render only when an input changes or a signal it
  reads changes. Combined with signals it is essentially free correctness.
- `styleUrl` is component-scoped: those styles cannot leak out.

---

## 4. Routing

**`src/app/app.routes.ts`**

```ts
export const routes: Routes = [
  {
    path: '',
    loadComponent: () => import('./features/home/pages/home/home.page').then((m) => m.HomePage),
  },
  {
    path: 'pantries/:slug',
    loadComponent: () =>
      import('./features/pantry/pages/pantry/pantry.page').then((m) => m.PantryPage),
  },
  { path: '**', redirectTo: '' },
];
```

`loadComponent` makes each page a **lazy chunk**. The production build confirms
it:

```
Lazy chunk files   | Names        | Raw size
chunk-T4K4XSD3.js  | pantry-page  | 18.09 kB
chunk-FIZJBOXL.js  | home-page    |  4.33 kB
```

Someone landing on the home page never downloads the pantry screen.

### Route parameters as inputs

`withComponentInputBinding()` in `app.config.ts` binds route params straight onto
matching inputs:

```ts
/** Bound from the `:slug` route parameter by withComponentInputBinding. */
readonly slug = input.required<string>();

ngOnInit(): void {
  void this.store.load(this.slug());
}
```

No `ActivatedRoute`, no subscription, no `takeUntilDestroyed`.

---

## 5. Signals — the state model

A signal is a value you read by calling it. Read one inside a template or a
`computed`, and Angular records the dependency and re-runs only what needs
re-running.

**`src/app/features/pantry/pantry.store.ts`** is the whole state layer:

```ts
// Writable state — the only things that actually change.
readonly items = signal<PantryItem[]>([]);
readonly view = signal<PantryView>('PRIMARY');
readonly category = signal<CategoryFilter>('ALL');
readonly sort = signal<SortMode>('DEFAULT');

// Derived state — recomputed automatically, cached until a dependency changes.
readonly itemsInView = computed(() => this.items().filter((item) => item.view === this.view()));

readonly visibleItems = computed(() => {
  const category = this.category();
  const filtered =
    category === 'ALL'
      ? this.itemsInView()
      : this.itemsInView().filter((item) => item.category === category);

  return sortItems(filtered, this.sort());
});

readonly availableCategories = computed<Category[]>(() =>
  CATEGORIES.filter((category) => this.itemsInView().some((item) => item.category === category)),
);
```

Tapping a category chip calls `store.category.set('FRUIT_VEG')`; picking a sort
mode calls `store.sort.set('STATUS')`. That is the entire filtering and ordering
feature: `visibleItems` recomputes, the grid re-renders, and
`availableCategories` recomputes only if `itemsInView` actually changed.

Compare with RxJS: no `BehaviorSubject`, no `combineLatest`, no `async` pipe, no
unsubscribing.

### Signals in components

```ts
protected readonly name = signal('');
protected readonly slug = computed(() => slugify(this.name()));
protected readonly canSubmit = computed(() => this.slug().length > 0 && !this.submitting());
```

The URL preview and the disabled state of the button both fall out of one
`computed`. The template just reads them:

```html
@if (slug()) {
  <ion-note class="preview">/pantries/{{ slug() }}</ion-note>
}
<ion-button type="submit" expand="block" [disabled]="!canSubmit()">…</ion-button>
```

### Optimistic updates

Tapping a status pill must feel instant even on a slow connection. The store
updates the signal first and rolls back if the request fails:

```ts
const previous = this.items();
this.items.update((items) =>
  items.map((current) =>
    current.product.id === item.product.id ? { ...current, ...patch } : current,
  ),
);

try {
  const updated = await firstValueFrom(this.api.updateItem(slug, item.product.id, patch));
  this.items.update((items) =>
    items.map((current) => (current.product.id === updated.product.id ? updated : current)),
  );
} catch {
  this.items.set(previous);
  this.error.set('The change could not be saved.');
}
```

`update()` must return a **new array**; mutating the existing one would not
notify anything.

### Why the store is provided by the route

```ts
@Component({
  providers: [PantryStore],   // not providedIn: 'root'
})
export class PantryPage { … }
```

The store lives and dies with the page. Navigating away disposes it, so opening
a different pantry can never show stale items.

---

## 6. Template control flow

Angular 17+ has built-in blocks — no `*ngIf`, no `*ngFor`, nothing to import.

```html
@if (store.loading()) {
  <div class="state"><ion-spinner name="crescent" /></div>
} @else if (store.error()) {
  <div class="state"><p class="text-muted">{{ store.error() }}</p></div>
} @else if (store.visibleItems().length === 0) {
  <div class="state"><p class="text-muted">No products in this view yet.</p></div>
} @else {
  <div class="grid">
    @for (item of store.visibleItems(); track item.product.id) {
      <app-product-card
        [item]="item"
        (statusToggled)="store.cycleStatus($event)"
        (settingsOpened)="openSettings($event)"
      />
    }
  </div>
}
```

`track` is mandatory. It tells Angular which DOM node belongs to which item, so
changing one product's status re-renders one card instead of 57.

---

## 7. Inputs and outputs, signal-style

```ts
readonly item = input.required<PantryItem>();          // was @Input({ required: true })
readonly statusToggled = output<PantryItem>();         // was @Output() … = new EventEmitter()
```

`input()` produces a read-only signal, so `item()` inside a `computed` just works.
`output()` is a plain emitter with better typing and automatic cleanup.

---

## 8. The Ionic components used here

| Component | Where | What it gives you |
| --- | --- | --- |
| `ion-app`, `ion-router-outlet` | `app.component.ts` | Layout context and page transitions |
| `ion-content` | every page | The scrollable area, with safe-area padding |
| `ion-header`, `ion-toolbar`, `ion-title` | pantry page, modals | Top bar that stays put while content scrolls |
| `ion-split-pane` | pantry page | Sidebar on desktop, drawer on mobile |
| `ion-menu`, `ion-menu-button` | pantry page | The drawer itself and its hamburger |
| `ion-list`, `ion-item`, `ion-label`, `ion-badge` | side menu | Menu entries with pending-item counters |
| `ion-input`, `ion-note` | home page | Text field with floating label, helper text |
| `ion-button`, `ion-icon`, `ion-spinner` | everywhere | Buttons, icons, loading state |
| `ion-segment`, `ion-segment-button` | settings modal | Three mutually exclusive options |
| `ion-checkbox` | shopping list modal | The "include low stock" toggle |
| `ModalController`, `ToastController` | pantry page, modals | Bottom sheets and transient messages |
| `ActionSheetController` | pantry page | The "sort by" menu, built from an array of buttons |

Two Ionic conventions that trip people up:

**Slots.** `slot` places an element in a named region of its parent.

```html
<ion-buttons slot="start"><ion-menu-button /></ion-buttons>
<ion-title>{{ viewMeta[store.view()].label }}</ion-title>
<ion-buttons slot="end">…</ion-buttons>
```

**Shadow-DOM styling.** You cannot reach inside an Ionic component with a normal
CSS selector. Use the custom properties it exposes:

```scss
ion-menu ion-item {
  --border-radius: var(--app-radius-sm);   /* an Ionic CSS variable */
  margin: 2px 10px;                        /* a normal property, applies to the host */
}
```

---

## 9. Icons

Standalone Ionic does not ship all 1,300 icons. You register the ones you use:

```ts
import { addIcons } from 'ionicons';
import { leafOutline, wineOutline, sparklesOutline } from 'ionicons/icons';

constructor() {
  addIcons({ leafOutline, wineOutline, sparklesOutline });
}
```

Then reference them in kebab-case:

```html
<ion-icon name="leaf-outline" aria-hidden="true" />
```

Anything you do not import is not in the bundle.

---

## 10. Modals — and the one real gotcha

Ionic modals are created imperatively:

```ts
const modal = await this.modals.create({
  component: ItemSettingsModalComponent,
  componentProps: { item },
  breakpoints: [0, 0.7],
  initialBreakpoint: 0.7,     // opens as a bottom sheet at 70% height
});
await modal.present();

const { data } = await modal.onWillDismiss();
if (data) {
  await this.store.patch(item, data);
}
```

The modal closes by returning data:

```ts
void this.modals.dismiss(Object.keys(patch).length ? patch : undefined);
```

> **The gotcha.** `componentProps` assigns values **directly onto the component
> instance**. If the target is declared with `input()`, the assignment overwrites
> the InputSignal function itself and you get `TypeError: this.items is not a
> function` at runtime. Modal inputs must be plain properties:
>
> ```ts
> /**
>  * Items of the current view, assigned through `componentProps`.
>  *
>  * A plain property, not `input()`: ModalController assigns componentProps
>  * straight onto the instance, which would overwrite an InputSignal.
>  */
> items: PantryItem[] = [];
> ```

Toasts follow the same pattern:

```ts
const toast = await this.toasts.create({ message, duration: 2000, position: 'bottom' });
await toast.present();
```

---

## 11. Theming

Ionic is styled entirely through CSS custom properties. Override them once in
`src/theme/variables.scss` and every component follows.

Each Ionic colour needs six variables — base, `-rgb`, `-contrast`,
`-contrast-rgb`, `-shade` and `-tint` — because Ionic derives hover, active and
text colours from them:

```scss
--ion-color-primary: #7c9082;              /* muted sage */
--ion-color-primary-rgb: 124, 144, 130;    /* used for rgba() opacity */
--ion-color-primary-contrast: #ffffff;     /* text drawn on top */
--ion-color-primary-contrast-rgb: 255, 255, 255;
--ion-color-primary-shade: #6d7f72;        /* pressed state */
--ion-color-primary-tint: #899b8f;         /* hover state */
```

The palette maps the domain onto Ionic's semantic colours, which is why
`color="danger"` on a badge is automatically the same red as an out-of-stock
pill:

| Ionic role | Hex | Meaning here |
| --- | --- | --- |
| `primary` | `#7c9082` | Sage — buttons, active chips |
| `success` | `#7fa37f` | Enough stock |
| `warning` | `#d9a84e` | Running low |
| `danger` | `#c4756b` | Out of stock |

App-specific tokens live alongside them:

```scss
--app-radius: 16px;
--app-shadow: 0 1px 2px rgba(51, 55, 47, 0.04), 0 8px 24px rgba(51, 55, 47, 0.05);
--app-status-out-soft: #f6e3e0;      /* behind the status pill */
--app-status-out-surface: #f9f1f0;   /* the whole card, lighter so the pill still reads */
```

The status pill then needs no logic at all — just an attribute selector:

```scss
.pill[data-status='OUT'] {
  background: var(--app-status-out-soft);
  color: var(--ion-color-danger-shade);
}
```

---

## 12. Responsive layout

Two mechanisms, zero media queries.

**`ion-split-pane`** shows the menu permanently from the `lg` breakpoint up and
collapses it into a swipeable drawer below it:

```html
<ion-split-pane contentId="pantry-content" when="lg">
  <ion-menu contentId="pantry-content" type="overlay"> … </ion-menu>

  <div class="ion-page" id="pantry-content"> … </div>
</ion-split-pane>
```

Both `contentId` values must match the `id` of an element carrying the
`ion-page` class. That is how Ionic knows what to push aside.

**CSS Grid with `auto-fill`** sizes the product grid:

```scss
.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 12px;
  max-width: 1100px;
  margin: 0 auto;
}
```

Two columns on a phone, six on a laptop, and nothing to maintain in between.

---

## 13. Runtime configuration

A built Angular app is static files. Baking the API URL in at build time would
mean one image per environment. Instead:

1. `entrypoint.sh` writes `config.json` into the nginx root at container start:

   ```sh
   : "${BACKEND_URL:?[entrypoint] BACKEND_URL is required but not set}"

   cat > /usr/share/nginx/html/config.json <<EOF
   { "BACKEND_URL": "${BACKEND_URL}" }
   EOF

   exec nginx -g "daemon off;"
   ```

2. `ConfigService.load()` fetches and validates it before the app bootstraps:

   ```ts
   const loaded = await firstValueFrom(this.http.get<AppConfig>('config.json'));

   if (!loaded?.BACKEND_URL || !/^https?:\/\/.+/.test(loaded.BACKEND_URL)) {
     throw new Error('config.json: BACKEND_URL must be an absolute http(s) URL');
   }
   ```

3. `nginx.conf` marks it `no-store`, so a redeploy is picked up immediately:

   ```nginx
   location = /config.json {
       add_header Cache-Control "no-store" always;
   }
   ```

`BACKEND_URL` is resolved **by the browser**, so it must be the public URL —
`http://localhost:8080/v1`, not `http://backend:8080/v1`.

---

## 14. SPA routing in nginx

`/pantries/familia-suarez` is not a file. Without this, a refresh returns 404:

```nginx
location / {
    try_files $uri $uri/ /index.html;
}
```

---

## 15. Tests

Three suites, all fast, none requiring a browser.

```bash
make -C frontend test
```

Pure functions are trivial to test:

```ts
it.each([
  ['Familia Suárez', 'familia-suarez'],
  ['La Peña', 'la-pena'],
])('turns %s into %s', (input, expected) => {
  expect(slugify(input)).toBe(expected);
});
```

The store is tested through `TestBed` with a fake API — no HTTP, no components:

```ts
TestBed.configureTestingModule({
  providers: [PantryStore, { provide: PantryApiService, useValue: api }],
});
store = TestBed.inject(PantryStore);

it('rolls the optimistic update back when the request fails', async () => {
  await loadPantry();
  api.updateItem.mockReturnValue(throwError(() => new Error('offline')));

  await store.cycleStatus(store.items()[0]);

  expect(store.items()[0].status).toBe('OK');
  expect(store.error()).not.toBeNull();
});
```

Note `setup-jest.ts`: since `jest-preset-angular` v15 the entry point is
`setup-env/zone`, not the old `setup-jest`.

```ts
import { setupZoneTestEnv } from 'jest-preset-angular/setup-env/zone';

setupZoneTestEnv();
```

---

## 16. Turning this into a mobile app

`capacitor.config.ts` is already committed:

```ts
const config: CapacitorConfig = {
  appId: 'work.midepensa.app',
  appName: 'MiDepensa',
  webDir: 'www',
};
```

When you want an Android build:

```bash
cd frontend
npm run build
npx cap add android     # generates android/ (git-ignored)
npx cap sync
npx cap open android    # opens Android Studio
```

Two things will need attention:

1. `BACKEND_URL` — a phone cannot reach `localhost`. Point it at the machine's
   LAN address or a real domain.
2. Capacitor serves the app from `capacitor://localhost`, so that origin has to
   be in the backend's `CORS_ORIGINS`.

---

## 17. Where to go next

Small, useful exercises on this codebase:

- **Add a product.** One row in a new migration under
  `backend/internal/infrastructure/migrations/`, one line in
  `frontend/scripts/fetch-product-icons.sh`.
- **Add a category.** Extend the Go `Category` enum and the `CHECK` constraints
  in a new migration, then `CATEGORIES` and `CATEGORY_META`. The chip bar and
  the shopping list pick it up with no other change.
- **Add a "mark everything as enough" button.** A new method on `PantryStore`
  and a bulk endpoint on the backend.
- **Add dark mode.** Import `@ionic/angular/css/palettes/dark.system.css` and
  define the same variables under `@media (prefers-color-scheme: dark)`.
- **Add i18n.** Install `@ngx-translate/core`, move the inline English strings
  into `assets/i18n/en.json` and `es.json`. The structure already supports it.

### Reference

- [Ionic Angular components](https://ionicframework.com/docs/components)
- [Ionic theming](https://ionicframework.com/docs/theming/basics)
- [Angular signals](https://angular.dev/guide/signals)
- [Angular standalone components](https://angular.dev/guide/components/importing)
- [Capacitor](https://capacitorjs.com/docs)
