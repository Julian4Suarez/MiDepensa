import { ChangeDetectionStrategy, Component, inject, signal } from '@angular/core';
import {
  IonButton,
  IonButtons,
  IonContent,
  IonHeader,
  IonIcon,
  IonLabel,
  IonSegment,
  IonSegmentButton,
  IonTitle,
  IonToolbar,
  ModalController,
} from '@ionic/angular/standalone';
import { addIcons } from 'ionicons';
import { closeOutline } from 'ionicons/icons';

import { CATEGORIES, CATEGORY_META, VIEWS, VIEW_META } from '../../../../shared/models/pantry.meta';
import type {
  Category,
  ItemPatch,
  PantryItem,
  PantryView,
} from '../../../../shared/models/pantry.model';

/**
 * Bottom sheet to move a product to another view or category.
 *
 * Dismisses with an {@link ItemPatch} when something changed, or with no data
 * when the user cancels.
 */
@Component({
  selector: 'app-item-settings-modal',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [
    IonButton,
    IonButtons,
    IonContent,
    IonHeader,
    IonIcon,
    IonLabel,
    IonSegment,
    IonSegmentButton,
    IonTitle,
    IonToolbar,
  ],
  templateUrl: './item-settings-modal.component.html',
  styleUrl: './item-settings-modal.component.scss',
})
export class ItemSettingsModalComponent {
  private readonly modals = inject(ModalController);

  /**
   * Assigned through `componentProps`. A plain property, not `input()`:
   * ModalController writes componentProps straight onto the instance.
   */
  item!: PantryItem;

  protected readonly views = VIEWS;
  protected readonly viewMeta = VIEW_META;
  protected readonly categories = CATEGORIES;
  protected readonly categoryMeta = CATEGORY_META;

  protected readonly selectedView = signal<PantryView | null>(null);
  protected readonly selectedCategory = signal<Category | null>(null);

  constructor() {
    addIcons({ closeOutline });
  }

  protected currentView(): PantryView {
    return this.selectedView() ?? this.item.view;
  }

  protected currentCategory(): Category {
    return this.selectedCategory() ?? this.item.category;
  }

  protected close(): void {
    void this.modals.dismiss();
  }

  protected save(): void {
    const patch: ItemPatch = {};
    if (this.currentView() !== this.item.view) {
      patch.view = this.currentView();
    }
    if (this.currentCategory() !== this.item.category) {
      patch.category = this.currentCategory();
    }
    void this.modals.dismiss(Object.keys(patch).length ? patch : undefined);
  }
}
