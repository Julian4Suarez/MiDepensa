# MiDepensa — Implementation Plan

> Status: **IMPLEMENTED**. Every phase below is done; `make up` brings up the full stack and
> `make test` is green. Kept as the record of what was decided and why.
> All source code, comments, tests, guides **and UI copy** are written in English.

---

## 0. Confirmed decisions

| Topic | Decision |
| --- | --- |
| Backend | **Go 1.26 + Gin + PostgreSQL** (mirrors `ductifact/backend`) |
| Categories | **4**: Fresh / Pantry / Drinks / Home & Care — later replaced, see §11 |
| Product images | **Local OpenMoji SVGs** (offline, no external requests) |
| Item config | Catalog ships **default category + default view**; **one category per item** |
| UI language | **English only** (no i18n layer) |
| Pantry URL | **`/pantries/:slug`** |
| Infra scope | **Local + `staging`/`production` manifests + `deploy.sh`** |
| Observability | **None** — no Prometheus, no Grafana |
| Catalog | **Fixed, ~50 seeded products**; no user-created products in v1 |

---

## 1. Product summary

A small, anonymous, multi-platform pantry-stock tracker.

- No authentication. Anyone with the URL can read and write.
- Screen 1 (Home): app title, illustration, short description, a text input to name a pantry, and a
  "Create" button. The name is slugified (`Familia Suárez` → `familia-suarez`) and the app navigates
  to the pantry URL.
- Screen 2 (Pantry): a grid/list of predefined products. Each product shows an image, a name and a
  stock status with 3 values: `OUT` (red), `LOW` (amber), `OK` (green).
- Side menu with 3 **views**: Primary, Secondary, Other. Each product belongs to exactly one view.
- Inside a view, a horizontal row of **category** icon-chips filters the products (no combobox).
- A "Generate shopping list" button produces copy-to-clipboard plain text grouped by category,
  listing all `OUT` products, with a checkbox to also include `LOW` products.

---

## 2. Recommended technology stack

### 2.1 Frontend — Ionic 8 + Angular 20 (standalone) ✅ recommended

| Concern | Choice | Rationale |
| --- | --- | --- |
| UI framework | **Ionic 8** (`@ionic/angular/standalone`) | Requested. Mobile-first components, `ion-menu`, `ion-chip`, `ion-card` cover ~100% of the UI. |
| App framework | **Angular 20**, standalone components + signals | Ionic's most mature binding. **Identical to `ductifact/frontend`**, so everything you learn transfers. |
| State | **Angular signals inside a service** | No NgRx. A single `PantryStore` service with `signal()` + `computed()` is enough and is the modern Angular way. |
| Styling | **Ionic CSS variables + SCSS** | Ionic theming alone delivers the muted palette. **No Tailwind** — one less thing to learn while you focus on Ionic. |
| i18n | **None** — English strings inline | Confirmed decision. `ngx-translate` can be added later without restructuring. |
| HTTP | `HttpClient` + runtime `config.json` via `APP_INITIALIZER` | Same pattern as ductifact: `entrypoint.sh` writes `config.json` at container start, so one image works in every environment. |
| Mobile later | **Capacitor 8** config committed from day 1, `android/`/`ios/` generated later | Zero cost now, no rework later. |
| Tests | **Jest** — minimal | ~4 tests: slugify util, `PantryStore` filtering, shopping-list text generator, one component smoke test. |
| Lint/format | ESLint + Prettier + Husky + commitlint | Same as ductifact. |

**Alternatives considered:** Ionic + React (you would not reuse ductifact knowledge), Ionic + Vue
(smallest ecosystem for Angular-shaped problems). Recommendation: stick with Angular.

### 2.2 Backend — Go 1.26 + Gin + PostgreSQL ✅ recommended

A backend **is required**: a pantry URL must be openable from any device (phone + laptop), which
rules out `localStorage`-only persistence.

| Concern | Choice | Rationale |
| --- | --- | --- |
| Language | **Go 1.26** | Same as `ductifact/backend`. You already read that code. |
| HTTP | **Gin** | Same as ductifact. |
| DB | **PostgreSQL 16-alpine** | Same as ductifact; one container. |
| Migrations | **golang-migrate**, `//go:embed`-ed SQL | Same as ductifact — self-contained binary. |
| Data access | **`pgx` + plain SQL** (no GORM) | 3 tables, ~8 queries. GORM would be over-engineering here. *(Say the word if you prefer GORM for consistency.)* |
| Architecture | **Light hexagonal** — same folder names as ductifact, fewer layers | See §5. |
| Config | `config.Load()` from env vars, panics on missing | Same as ductifact. |
| Tests | Minimal: unit tests for slug + shopping-list domain rules, 1 integration test per repository | No Schemathesis, no E2E suite. |

**Deliberately dropped from the ductifact stack** (over-engineering for this app):
Redis, MinIO, SMTP/MailPit, JWT/auth, rate limiting, Prometheus, Grafana, Schemathesis contract
testing, E2E suite.

### 2.3 Infra — Docker Compose + Makefile

Same methodology as `ductifact/infra`, scaled down:

- Pinned images by `@sha256` digest.
- `.env.example` as single source of truth, `environments/` manifests for **local, staging and
  production** (`images.manifest.env`, `<env>.manifest.env`, `<env>.config.env`).
- `scripts/deploy.sh`, `scripts/smoke.sh`, `scripts/validate_env.sh`, `scripts/db.sh`
  (shellcheck-clean, same CLI shape as ductifact: `./scripts/deploy.sh <env> [deploy|stop]`).
- Healthchecks + `depends_on: condition: service_healthy` everywhere.
- Ports bound to `127.0.0.1` in staging/production so a host reverse proxy can front them.
- Non-root users in both Dockerfiles, multi-stage builds.
- Conventional Commits + `feat/`/`fix/`/`chore/` branches + git-cliff changelog.
- **No GitHub Actions pipelines** — deploys are run manually via `deploy.sh`, as documented in
  `ductifact/infra/CONTRIBUTING.md`.

---

## 3. Repository layout (monorepo)

```
MiDepensa/
├── Makefile                     # root orchestrator: make up / down / logs / test / lint
├── README.md
├── CONTRIBUTING.md              # branch + commit conventions
├── .env.example
├── backend/
│   ├── Makefile                 # dev, build, test, lint, docker-*
│   ├── Dockerfile               # multi-stage, non-root
│   ├── docker-compose.yml       # postgres (+ app under `smoke` profile)
│   ├── go.mod
│   ├── cmd/api/main.go          # bootstrap + wiring
│   ├── internal/
│   │   ├── config/config.go
│   │   ├── domain/
│   │   │   ├── entities/        # pantry.go, product.go, pantry_item.go, category.go, view.go
│   │   │   └── repositories/    # pantry_repository.go, product_repository.go (interfaces)
│   │   ├── application/
│   │   │   ├── usecases/        # inbound port interfaces
│   │   │   └── services/        # pantry_service.go, catalog_service.go
│   │   └── infrastructure/
│   │       ├── adapters/inbound/http/
│   │       │   ├── router.go
│   │       │   ├── handlers/    # pantry_handler.go, catalog_handler.go, health_handler.go
│   │       │   ├── middleware/  # request_id.go, logger.go, recovery.go, cors.go
│   │       │   └── helpers/error_handler.go
│   │       ├── adapters/outbound/persistence/
│   │       │   ├── postgres_connection.go
│   │       │   ├── postgres_pantry_repository.go
│   │       │   └── postgres_product_repository.go
│   │       └── migrations/      # 000001_init.up.sql, 000002_seed_catalog.up.sql
│   └── test/{unit,integration}
├── frontend/
│   ├── Makefile
│   ├── Dockerfile               # node build → nginx runtime, non-root
│   ├── docker-compose.yml       # `dev` (ng serve) + `frontend` (nginx)
│   ├── nginx.conf               # SPA fallback + cache + security headers
│   ├── entrypoint.sh            # writes config.json from env
│   ├── angular.json / package.json / jest.config.js / capacitor.config.ts
│   ├── public/config.json       # dev-time runtime config
│   └── src/
│       ├── theme/variables.scss # Ionic palette (muted colors)
│       ├── styles.scss
│       ├── assets/products/     # ~50 OpenMoji SVGs + ATTRIBUTION.md
│       └── app/
│           ├── app.config.ts / app.routes.ts / app.component.*
│           ├── core/
│           │   ├── config/      # config.service.ts, config.model.ts
│           │   ├── services/    # pantry-api.service.ts, catalog.service.ts
│           │   └── utils/       # slugify.ts, shopping-list.ts
│           ├── shared/
│           │   ├── models/      # pantry.model.ts, product.model.ts, enums
│           │   └── components/  # product-card, status-badge, category-filter-bar
│           └── features/
│               ├── home/pages/home/
│               └── pantry/
│                   ├── pages/pantry/
│                   ├── components/product-config-modal/
│                   ├── components/shopping-list-modal/
│                   └── pantry.routes.ts
├── infra/
│   ├── docker-compose.yml       # full stack: postgres + backend + frontend
│   ├── environments/
│   │   ├── images.manifest.env  # postgres image pinned by digest
│   │   ├── staging.manifest.env / staging.config.env
│   │   └── production.manifest.env / production.config.env
│   ├── scripts/                 # deploy.sh, smoke.sh, validate_env.sh, db.sh
│   ├── README.md
│   └── MAINTENANCE.md
├── contracts/openapi/openapi.yaml   # hand-written spec, served at /v1/docs
└── docs/
    ├── GUIDE_SETUP.md
    ├── GUIDE_ARCHITECTURE.md
    └── GUIDE_IONIC_FRONTEND.md  # ← the learning guide you asked for
```

---

## 4. Domain model

### 4.1 Concepts

| Concept | Scope | Notes |
| --- | --- | --- |
| **Pantry** | one row per created URL | `id`, `slug` (unique), `name`, timestamps |
| **Product** | global catalog, seeded | `id`, `code`, `name`, `image`, `default_category`, `default_view` |
| **PantryItem** | per pantry × product | `status`, `view`, `category` — overridable per pantry |
| **Category** | fixed enum | see §4.3 |
| **View** | fixed enum | `PRIMARY`, `SECONDARY`, `OTHER` |
| **Status** | fixed enum | `OUT` (red), `LOW` (amber), `OK` (green) |

### 4.2 Schema

```sql
CREATE TABLE pantries (
  id          UUID PRIMARY KEY,
  slug        TEXT NOT NULL UNIQUE,
  name        TEXT NOT NULL,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE products (            -- global, seeded by migration
  id                UUID PRIMARY KEY,
  code              TEXT NOT NULL UNIQUE,   -- 'tomato'
  name              TEXT NOT NULL,          -- 'Tomatoes'
  image             TEXT NOT NULL,          -- 'tomato.svg'
  default_category  TEXT NOT NULL,
  default_view      TEXT NOT NULL,
  sort_order        INT  NOT NULL
);

CREATE TABLE pantry_items (
  pantry_id   UUID NOT NULL REFERENCES pantries(id) ON DELETE CASCADE,
  product_id  UUID NOT NULL REFERENCES products(id),
  status      TEXT NOT NULL,   -- OUT | LOW | OK
  view        TEXT NOT NULL,   -- PRIMARY | SECONDARY | OTHER
  category    TEXT NOT NULL,
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (pantry_id, product_id)
);
```

Creating a pantry inserts one `pantry_items` row per catalog product, initialised from the
product's `default_*` values with `status = 'OK'`.

### 4.3 Proposed categories (4)

You asked me to pick 3–5. My recommendation — **4 categories**, because views already encode
priority, so categories should encode *where the product lives / what it is*:

| # | Category | Ionicon | Covers |
| --- | --- | --- | --- |
| 1 | **Fresh** | `leaf-outline` | fruit, vegetables, meat, fish, dairy, eggs, bread |
| 2 | **Pantry** | `file-tray-stacked-outline` | pasta, rice, legumes, canned food, oil, flour, sugar, spices, snacks |
| 3 | **Drinks** | `wine-outline` | water, juice, coffee, tea, soda, beer, wine |
| 4 | **Home & Care** | `sparkles-outline` | detergent, soap, bleach, sponges, toilet paper, shampoo, toothpaste |

### 4.4 Seed catalog

~50 products across the 4 categories (e.g. Tomatoes, Onions, Potatoes, Milk, Eggs, Chicken, Pasta,
Rice, Olive oil, Coffee, Toilet paper, Dish soap…). Enough to be useful, small enough to review.
Each product ships with a `default_category` and a `default_view`, so a brand-new pantry is
immediately usable without configuring anything.

### 4.5 Product artwork

**OpenMoji** (CC BY-SA 4.0) SVGs, downloaded once and committed to
`frontend/src/assets/products/<code>.svg`. Attribution in `README.md` and
`frontend/src/assets/products/ATTRIBUTION.md`. No external CDN, works offline, ~2 KB per file.

---

## 5. API design

Base path `/v1`. Errors as `{ "error": { "code": "...", "message": "..." } }`.

| Method | Path | Purpose |
| --- | --- | --- |
| `GET` | `/healthz` | liveness |
| `GET` | `/readyz` | readiness (DB ping, bounded context timeout) |
| `GET` | `/v1/catalog` | products + categories + views metadata |
| `POST` | `/v1/pantries` | `{ name }` → creates pantry + items, returns `{ id, slug, name }` |
| `GET` | `/v1/pantries/{slug}` | pantry + all items joined with product data |
| `PATCH` | `/v1/pantries/{slug}/items/{productId}` | `{ status?, view?, category? }` |
| `PATCH` | `/v1/pantries/{slug}/items` | bulk status update (used by "mark all OK") — *optional* |

Slug rules (enforced on **both** sides): lowercase, NFD-normalised (á→a, ñ→n), non-alphanumerics →
`-`, collapse repeats, trim `-`, max 60 chars, must match `^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`.
On collision: return `409` and let the UI suggest `familia-suarez-2`.

Minimal-but-sane security (as you asked — not bulletproof):
parameterised SQL, body size limit, request timeout, CORS allowlist, slug regex validation,
`ON DELETE CASCADE`, no secrets in the frontend.

---

## 6. Frontend design

### 6.1 Routes

```
/                        → HomePage
/pantries/:slug          → PantryPage   (lazy-loaded route)
/**                      → redirect to /
```

### 6.2 Pantry page layout

- `ion-split-pane` — permanent sidebar on desktop, swipe/burger `ion-menu` on mobile.
- Sidebar: pantry name + the 3 views as `ion-item` entries with counters (`Primary · 12`).
- Header: view title, "Shopping list" button.
- Sub-header: horizontal scrollable `ion-chip` row = category filter (icon + label, single-select,
  plus an "All" chip).
- Content: responsive `ion-grid` of product cards. Card = image, name, and a status pill.
  **Tap the pill cycles OUT → LOW → OK → OUT** (fast, one tap, no modal).
  **Long-press / ⋯ button opens the config modal** (change view, change category).
- FAB or header button → shopping-list modal: checkbox "Include low stock", `ion-textarea` preview,
  "Copy" button using the Clipboard API.

### 6.3 State

One `PantryStore` service per route:

```ts
readonly items    = signal<PantryItem[]>([]);
readonly view     = signal<PantryView>('PRIMARY');
readonly category = signal<Category | 'ALL'>('ALL');
readonly visible  = computed(() => /* filter by view + category */);
```

Updates are **optimistic**: mutate the signal, fire `PATCH`, roll back + toast on failure.

### 6.4 Visual style — muted palette

Light, warm, low-saturation. Proposal (Ionic CSS variables in `theme/variables.scss`):

| Token | Hex | Use |
| --- | --- | --- |
| `--ion-background-color` | `#F6F5F1` | warm off-white page background |
| surface | `#FFFFFF` | cards |
| `--ion-color-primary` | `#7C9082` | muted sage — buttons, active chips |
| text | `#33372F` | soft near-black |
| `--ion-color-danger` (OUT) | `#C4756B` | muted terracotta |
| `--ion-color-warning` (LOW) | `#D9A84E` | muted amber |
| `--ion-color-success` (OK) | `#7FA37F` | muted green |

Rounded corners (12–16px), soft shadows, generous whitespace, no gradients.

---

## 7. Makefile targets

Root `Makefile` (thin wrapper, mirrors ductifact naming):

| Target | Effect |
| --- | --- |
| `make help` | list targets (default goal) |
| `make up` | build + start the full stack (postgres + backend + frontend) |
| `make down` | stop everything |
| `make logs` | tail all logs |
| `make dev` | postgres in Docker + backend via `air` + `ng serve` |
| `make test` | backend + frontend tests |
| `make lint` / `make fmt` | both projects |
| `make smoke` | run `infra/scripts/smoke.sh local` |
| `make deploy ENV=staging` | run `infra/scripts/deploy.sh staging` |
| `make db-backup` / `make db-restore` | `infra/scripts/db.sh` |
| `make clean` | remove build artifacts |

`backend/Makefile` and `frontend/Makefile` keep their own full target sets, matching ductifact.

---

## 8. Documentation deliverables

1. `README.md` — what it is, quickstart (`make up`), screenshots.
2. `docs/GUIDE_SETUP.md` — prerequisites, env vars, dev loop.
3. `docs/GUIDE_ARCHITECTURE.md` — hexagonal layers, request lifecycle, why each folder exists.
4. **`docs/GUIDE_IONIC_FRONTEND.md`** — the main learning deliverable:
   - Ionic vs Angular: who does what.
   - Standalone components & why there are no NgModules.
   - Every Ionic component used, with a snippet from *this* codebase.
   - Signals: `signal` / `computed` / `effect` explained through `PantryStore`.
   - Routing + lazy-loaded routes.
   - Ionic theming: CSS variables, dark mode, `ion-color`.
   - Responsive: `ion-split-pane`, `ion-grid`, breakpoints.
   - Runtime config via `APP_INITIALIZER` + `config.json`.
   - Capacitor: how to turn this into an Android/iOS app later.
5. `CONTRIBUTING.md` — branch names, Conventional Commits, changelog.
6. `infra/README.md` + `infra/MAINTENANCE.md` — deploy, backup, rollback.

---

## 9. Implementation phases

| Phase | Deliverable |
| --- | --- |
| 0 | Repo scaffolding, root `Makefile`, `.env.example`, `CONTRIBUTING.md`, `.gitignore` |
| 1 | Backend: config, domain, migrations + catalog seed, repositories, services |
| 2 | Backend: HTTP handlers, middleware, router, health probes, `Dockerfile`, compose |
| 3 | Backend: minimal unit + integration tests |
| 4 | Frontend: Ionic/Angular scaffold, theme, runtime config, `Dockerfile`, `nginx.conf`, `entrypoint.sh` |
| 5 | Frontend: Home page (slugify + create + navigate) |
| 6 | Frontend: Pantry page (menu, views, category chips, product grid, status cycling) |
| 7 | Frontend: config modal + shopping-list generator + clipboard |
| 8 | Frontend: minimal Jest tests |
| 9 | `infra/`: compose, environments (local/staging/production), `deploy.sh`, `smoke.sh`, `validate_env.sh`, `db.sh` |
| 10 | Docs: README + 4 guides |

---

## 10. Assumptions to confirm before starting

1. Target repository is `/home/jsuarez/workspace/ductifact/MiDepensa` (already cloned; currently
   contains only `README.md`, `LICENSE`, `prompt.md`). Everything is created there.
2. Work happens on the current branch; commits follow Conventional Commits, and nothing is pushed
   without asking you first.
3. Node 22 and Go 1.26 are available locally — otherwise everything still runs through Docker.

---

## 11. Deviations from the plan

Things that changed while building, and why:

| Change | Reason |
| --- | --- |
| Categories reworked into **6 supermarket aisles** (migration `000003`) | The original four did not survive contact with real use: `PANTRY` collided with the name of the app, and `FRESH` held 44% of the catalog. Now: Fruit & veg / Meat & fish / Dairy & eggs / Dry & canned / Drinks / Home & care. |
| Catalog has **57** products, not ~50 | Four categories needed enough entries each to be useful. |
| No `contracts/openapi/` folder | The API has six endpoints, fully described in `docs/GUIDE_ARCHITECTURE.md`. A hand-maintained spec with no code generation and no contract tests would have been dead weight. |
| No Tailwind | Ionic CSS variables covered the whole design. One less thing to learn. |
| Frontend runs through Docker | Node is not installed on this machine; `frontend/Makefile` detects that and falls back to `node:22-alpine`. |
| Compose project names pinned (`name: midepensa-backend`) | Without it, the `backend/` directory name collided with the sibling ductifact project and reused its PostgreSQL 15 volume. |
| Modal inputs are plain properties, not `input()` | `ModalController` assigns `componentProps` directly onto the instance, which overwrites an `InputSignal`. |
| Healthchecks use `wget -O /dev/null` | `wget --spider` sends HEAD, and Gin only routes the GET handlers. |
