import { ChangeDetectionStrategy, Component, type OnInit, computed, inject, input } from '@angular/core';
import { RouterLink } from '@angular/router';
import {
  ActionSheetController,
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
  fishOutline,
  homeOutline,
  nutritionOutline,
  sparklesOutline,
  starOutline,
  swapVerticalOutline,
  wineOutline,
} from 'ionicons/icons';

import {
  FilterBarComponent,
  type FilterOption,
} from '../../../../shared/components/filter-bar/filter-bar.component';
import { ProductCardComponent } from '../../../../shared/components/product-card/product-card.component';
import {
  ALL_ICON,
  CATEGORY_META,
  SORT_META,
  SORT_MODES,
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
import { PantryStore } from '../../pantry.store';

/** The pantry screen: two filter bars and the product grid. */
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
      fishOutline,
      eggOutline,
      fileTrayStackedOutline,
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

  protected async openShoppingList(): Promise<void> {
    const modal = await this.modals.create({
      component: ShoppingListModalComponent,
      componentProps: { items: this.store.visibleItems() },
      breakpoints: [0, 0.9],
      initialBreakpoint: 0.9,
    });
    await modal.present();
  }
}
