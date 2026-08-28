import { ChangeDetectionStrategy, Component, computed, inject } from '@angular/core';
import {
  IonButton,
  IonButtons,
  IonContent,
  IonHeader,
  IonIcon,
  IonTitle,
  IonToolbar,
  ModalController,
  ToastController,
} from '@ionic/angular/standalone';
import { addIcons } from 'ionicons';
import { closeOutline, copyOutline } from 'ionicons/icons';

import { buildShoppingList } from '../../../../core/utils/shopping-list';
import type { PantryItem } from '../../../../shared/models/pantry.model';

/** Bottom sheet that renders the shopping list and copies it to the clipboard. */
@Component({
  selector: 'app-shopping-list-modal',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [
    IonButton,
    IonButtons,
    IonContent,
    IonHeader,
    IonIcon,
    IonTitle,
    IonToolbar,
  ],
  templateUrl: './shopping-list-modal.component.html',
  styleUrl: './shopping-list-modal.component.scss',
})
export class ShoppingListModalComponent {
  private readonly modals = inject(ModalController);
  private readonly toasts = inject(ToastController);

  /**
   * Products currently in the cart, assigned through `componentProps`.
   *
   * A plain property, not `input()`: ModalController assigns componentProps
   * straight onto the instance, which would overwrite an InputSignal.
   */
  items: PantryItem[] = [];

  protected readonly text = computed(() => buildShoppingList(this.items));

  constructor() {
    addIcons({ closeOutline, copyOutline });
  }

  protected close(): void {
    void this.modals.dismiss();
  }

  protected async copy(): Promise<void> {
    try {
      await navigator.clipboard.writeText(this.text());
      await this.notify('Shopping list copied');
    } catch {
      // Blocked outside a secure context (plain http on a phone, for example).
      await this.notify('Copying is not available here — select the text instead');
    }
  }

  private async notify(message: string): Promise<void> {
    const toast = await this.toasts.create({ message, duration: 2000, position: 'bottom' });
    await toast.present();
  }
}
