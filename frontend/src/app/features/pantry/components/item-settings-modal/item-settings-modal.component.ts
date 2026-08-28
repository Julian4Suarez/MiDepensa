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

import { CATEGORIES, CATEGORY_META, TYPES, TYPE_META } from '../../../../shared/models/pantry.meta';
import {
  ARCHIVED,
  type Category,
  type ItemPatch,
  type PantryItem,
  type ProductType,
} from '../../../../shared/models/pantry.model';

/**
 * Bottom sheet to change how often a product is bought, or which aisle it
 * belongs to.
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

  protected readonly types = TYPES;
  protected readonly typeMeta = TYPE_META;
  protected readonly categories = CATEGORIES;
  protected readonly categoryMeta = CATEGORY_META;

  protected readonly selectedType = signal<ProductType | null>(null);
  protected readonly selectedCategory = signal<Category | null>(null);

  constructor() {
    addIcons({ closeOutline });
  }

  protected currentType(): ProductType {
    return this.selectedType() ?? this.item.type;
  }

  protected currentCategory(): Category {
    return this.selectedCategory() ?? this.item.category;
  }

  protected close(): void {
    void this.modals.dismiss();
  }

  protected save(): void {
    const patch: ItemPatch = {};
    if (this.currentType() !== this.item.type) {
      patch.type = this.currentType();
    }
    if (this.currentCategory() !== this.item.category) {
      patch.category = this.currentCategory();
    }
    void this.modals.dismiss(Object.keys(patch).length ? patch : undefined);
  }

  protected toggleArchived(): void {
    void this.modals.dismiss({
      status: this.item.status === ARCHIVED ? 'OK' : ARCHIVED,
    } satisfies ItemPatch);
  }
}
