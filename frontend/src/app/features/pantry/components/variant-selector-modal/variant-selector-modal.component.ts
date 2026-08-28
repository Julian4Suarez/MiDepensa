import { ChangeDetectionStrategy, Component, type OnInit, inject, signal } from '@angular/core';
import {
  IonButton,
  IonButtons,
  IonCheckbox,
  IonContent,
  IonHeader,
  IonIcon,
  IonItem,
  IonList,
  IonTitle,
  IonToolbar,
  ModalController,
} from '@ionic/angular/standalone';
import { addIcons } from 'ionicons';
import { closeOutline } from 'ionicons/icons';

import type { ItemPatch, PantryItem } from '../../../../shared/models/pantry.model';

/** Selects one or more concrete variants before adding a general product. */
@Component({
  selector: 'app-variant-selector-modal',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [
    IonButton,
    IonButtons,
    IonCheckbox,
    IonContent,
    IonHeader,
    IonIcon,
    IonItem,
    IonList,
    IonTitle,
    IonToolbar,
  ],
  templateUrl: './variant-selector-modal.component.html',
  styleUrl: './variant-selector-modal.component.scss',
})
export class VariantSelectorModalComponent implements OnInit {
  private readonly modals = inject(ModalController);

  item!: PantryItem;
  protected readonly selected = signal<Set<string>>(new Set());

  constructor() {
    addIcons({ closeOutline });
  }

  ngOnInit(): void {
    this.selected.set(new Set(this.item.selectedVariantIds));
  }

  protected toggle(id: string, checked: boolean): void {
    const selected = new Set(this.selected());
    if (checked) {
      selected.add(id);
    } else {
      selected.delete(id);
    }
    this.selected.set(selected);
  }

  protected close(): void {
    void this.modals.dismiss();
  }

  protected confirm(): void {
    const patch: ItemPatch = {
      selectedVariantIds: this.item.product.variants
        .filter((variant) => this.selected().has(variant.id))
        .map((variant) => variant.id),
    };
    void this.modals.dismiss(patch);
  }
}
