import { ChangeDetectionStrategy, Component, type OnInit, inject, input } from '@angular/core';
import { RouterLink } from '@angular/router';
import {
  IonBadge,
  IonButton,
  IonButtons,
  IonContent,
  IonHeader,
  IonIcon,
  IonItem,
  IonLabel,
  IonList,
  IonListHeader,
  IonMenu,
  IonMenuButton,
  IonSpinner,
  IonSplitPane,
  IonTitle,
  IonToolbar,
  ModalController,
} from '@ionic/angular/standalone';
import { addIcons } from 'ionicons';
import {
  appsOutline,
  bookmarkOutline,
  cartOutline,
  homeOutline,
  starOutline,
} from 'ionicons/icons';

import { CategoryFilterBarComponent } from '../../../../shared/components/category-filter-bar/category-filter-bar.component';
import { ProductCardComponent } from '../../../../shared/components/product-card/product-card.component';
import { VIEWS, VIEW_META } from '../../../../shared/models/pantry.meta';
import type { PantryItem, PantryView } from '../../../../shared/models/pantry.model';
import { ItemSettingsModalComponent } from '../../components/item-settings-modal/item-settings-modal.component';
import { ShoppingListModalComponent } from '../../components/shopping-list-modal/shopping-list-modal.component';
import { PantryStore } from '../../pantry.store';

/** The pantry screen: side menu with views, category filter and product grid. */
@Component({
  selector: 'app-pantry',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  // Provided here rather than in root so the state is discarded on navigation.
  providers: [PantryStore],
  imports: [
    CategoryFilterBarComponent,
    ProductCardComponent,
    RouterLink,
    IonBadge,
    IonButton,
    IonButtons,
    IonContent,
    IonHeader,
    IonIcon,
    IonItem,
    IonLabel,
    IonList,
    IonListHeader,
    IonMenu,
    IonMenuButton,
    IonSpinner,
    IonSplitPane,
    IonTitle,
    IonToolbar,
  ],
  templateUrl: './pantry.page.html',
  styleUrl: './pantry.page.scss',
})
export class PantryPage implements OnInit {
  protected readonly store = inject(PantryStore);
  private readonly modals = inject(ModalController);

  protected readonly views = VIEWS;
  protected readonly viewMeta = VIEW_META;

  /** Bound from the `:slug` route parameter by withComponentInputBinding. */
  readonly slug = input.required<string>();

  constructor() {
    addIcons({ starOutline, bookmarkOutline, appsOutline, cartOutline, homeOutline });
  }

  ngOnInit(): void {
    void this.store.load(this.slug());
  }

  protected selectView(view: PantryView): void {
    this.store.view.set(view);
    this.store.category.set('ALL');
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
      componentProps: { items: this.store.itemsInView() },
      breakpoints: [0, 0.9],
      initialBreakpoint: 0.9,
    });
    await modal.present();
  }
}
