import { ChangeDetectionStrategy, Component, type OnInit, computed, inject, input } from '@angular/core';
import { RouterLink } from '@angular/router';
import {
  ActionSheetController,
  AlertController,
  IonButton,
  IonButtons,
  IonContent,
  IonHeader,
  IonIcon,
  IonSpinner,
  IonTitle,
  IonToolbar,
  ModalController,
} from '@ionic/angular/standalone';
import { addIcons } from 'ionicons';
import {
  appsOutline,
  archiveOutline,
  bookmarkOutline,
  cartOutline,
  eggOutline,
  fileTrayStackedOutline,
  filterOutline,
  fishOutline,
  homeOutline,
  nutritionOutline,
  refreshOutline,
  sparklesOutline,
  starOutline,
  swapVerticalOutline,
  wineOutline,
} from 'ionicons/icons';

import {
  FilterBarComponent,
  type FilterOption,
} from '../../../../shared/components/filter-bar/filter-bar.component';
import {
  ProductCardComponent,
  type StatusChange,
} from '../../../../shared/components/product-card/product-card.component';
import {
  ALL_ICON,
  CATEGORY_META,
  SORT_META,
  SORT_MODES,
  SHOPPING_STATUSES,
  STATUS_META,
  TYPES,
  TYPE_META,
} from '../../../../shared/models/pantry.meta';
import {
  ALL,
  ARCHIVED,
  type CategoryFilter,
  type PantryItem,
  type TypeFilter,
} from '../../../../shared/models/pantry.model';
import { ItemSettingsModalComponent } from '../../components/item-settings-modal/item-settings-modal.component';
import { ShoppingListModalComponent } from '../../components/shopping-list-modal/shopping-list-modal.component';
import { VariantSelectorModalComponent } from '../../components/variant-selector-modal/variant-selector-modal.component';
import { PantryStore } from '../../pantry.store';

/** The pantry screen: toolbar controls, two filter bars and the product grid. */
@Component({
  selector: 'app-pantry',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  // Provided here rather than in root so the state is discarded on navigation.
  providers: [PantryStore],
  imports: [
    FilterBarComponent,
    ProductCardComponent,
    RouterLink,
    IonButton,
    IonButtons,
    IonContent,
    IonHeader,
    IonIcon,
    IonSpinner,
    IonTitle,
    IonToolbar,
  ],
  templateUrl: './pantry.page.html',
  styleUrl: './pantry.page.scss',
})
export class PantryPage implements OnInit {
  protected readonly store = inject(PantryStore);
  private readonly modals = inject(ModalController);
  private readonly actionSheets = inject(ActionSheetController);
  private readonly alerts = inject(AlertController);

  protected readonly sortMeta = SORT_META;

  /** Bound from the `:slug` route parameter by withComponentInputBinding. */
  readonly slug = input.required<string>();

  private readonly allOption: FilterOption = { value: ALL, label: 'All', icon: ALL_ICON };

  protected readonly typeOptions: FilterOption[] = [
    this.allOption,
    ...TYPES.map((type) => ({ value: type, ...TYPE_META[type] })),
    { value: ARCHIVED, label: 'Archived', icon: 'archive-outline' },
  ];

  /** Rebuilt whenever the type filter changes the set of available categories. */
  protected readonly categoryOptions = computed<FilterOption[]>(() => [
    this.allOption,
    ...this.store
      .availableCategories()
      .map((category) => ({ value: category, ...CATEGORY_META[category] })),
  ]);

  constructor() {
    addIcons({
      starOutline,
      bookmarkOutline,
      archiveOutline,
      nutritionOutline,
      refreshOutline,
      fishOutline,
      eggOutline,
      fileTrayStackedOutline,
      filterOutline,
      wineOutline,
      sparklesOutline,
      appsOutline,
      cartOutline,
      homeOutline,
      swapVerticalOutline,
    });
  }

  ngOnInit(): void {
    void this.store.load(this.slug());
  }

  // The filter bar is value-agnostic, so the page narrows the emitted string.
  protected selectType(value: string): void {
    this.store.selectType(value as TypeFilter);
  }

  protected selectCategory(value: string): void {
    this.store.category.set(value as CategoryFilter);
  }

  protected async openSort(): Promise<void> {
    const sheet = await this.actionSheets.create({
      header: 'Sort by',
      buttons: [
        ...SORT_MODES.map((mode) => ({
          text: SORT_META[mode].label,
          role: this.store.sort() === mode ? 'selected' : undefined,
          handler: () => this.store.sort.set(mode),
        })),
        { text: 'Cancel', role: 'cancel' },
      ],
    });
    await sheet.present();
  }

  protected async openStatusFilter(): Promise<void> {
    const sheet = await this.actionSheets.create({
      header: 'Filter by status',
      buttons: [
        {
          text: 'All statuses',
          role: this.store.status() === ALL ? 'selected' : undefined,
          handler: () => this.store.status.set(ALL),
        },
        ...SHOPPING_STATUSES.map((status) => ({
          text: STATUS_META[status].label,
          role: this.store.status() === status ? 'selected' : undefined,
          handler: () => this.store.status.set(status),
        })),
        { text: 'Cancel', role: 'cancel' },
      ],
    });
    await sheet.present();
  }

  protected async confirmReset(): Promise<void> {
    const alert = await this.alerts.create({
      header: 'Reset products?',
      message: 'Every active product will return to pending. Archived products will not change.',
      buttons: [
        { text: 'Cancel', role: 'cancel' },
        {
          text: 'Reset',
          role: 'confirm',
          handler: () => void this.store.resetActiveItems(),
        },
      ],
    });
    await alert.present();
  }

  protected async openSettings(item: PantryItem): Promise<void> {
    const modal = await this.modals.create({
      component: ItemSettingsModalComponent,
      componentProps: { item },
      breakpoints: [0, 0.7],
      initialBreakpoint: 0.7,
    });
    await modal.present();

    const { data } = await modal.onWillDismiss();
    if (data) {
      await this.store.patch(item, data);
    }
  }

  protected async handleStatusChange(change: StatusChange): Promise<void> {
    await this.store.setStatus(change.item, change.status);
  }

  protected async openVariants(item: PantryItem): Promise<void> {
    const modal = await this.modals.create({
      component: VariantSelectorModalComponent,
      componentProps: { item },
      breakpoints: [0, 0.7],
      initialBreakpoint: 0.7,
    });
    await modal.present();
    const { data } = await modal.onWillDismiss();
    if (data) {
      await this.store.patch(item, data);
    }
  }

  protected async openShoppingList(): Promise<void> {
    const modal = await this.modals.create({
      component: ShoppingListModalComponent,
      componentProps: { items: this.store.cartItems() },
      breakpoints: [0, 0.9],
      initialBreakpoint: 0.9,
    });
    await modal.present();
  }
}
