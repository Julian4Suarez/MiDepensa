import { ChangeDetectionStrategy, Component, input, output } from '@angular/core';
import { IonButton, IonIcon } from '@ionic/angular/standalone';
import { addIcons } from 'ionicons';
import { ellipsisHorizontal } from 'ionicons/icons';

import { StatusPillComponent } from '../status-pill/status-pill.component';
import type { PantryItem } from '../../models/pantry.model';

/**
 * One product in the pantry grid: image, name, status pill and a settings
 * button. Tapping the card cycles the status; the button opens the settings.
 */
@Component({
  selector: 'app-product-card',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [IonButton, IonIcon, StatusPillComponent],
  templateUrl: './product-card.component.html',
  styleUrl: './product-card.component.scss',
})
export class ProductCardComponent {
  readonly item = input.required<PantryItem>();

  /** Emitted when the user taps the card to change the stock status. */
  readonly statusToggled = output<PantryItem>();

  /** Emitted when the user opens the item settings. */
  readonly settingsOpened = output<PantryItem>();

  constructor() {
    // Standalone Ionic requires every icon to be registered explicitly, which
    // keeps unused icons out of the bundle.
    addIcons({ ellipsisHorizontal });
  }

  protected imageUrl(): string {
    return `assets/products/${this.item().product.image}`;
  }
}
