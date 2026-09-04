import { ChangeDetectionStrategy, Component, input, output } from '@angular/core';
import { IonButton, IonIcon } from '@ionic/angular/standalone';
import { addIcons } from 'ionicons';
import {
  cartOutline,
  ellipsisHorizontal,
  layersOutline,
  timeOutline,
  trashOutline,
} from 'ionicons/icons';

import { StatusPillComponent } from '../status-pill/status-pill.component';
import type { PantryItem, ShoppingStatus } from '../../models/pantry.model';
import { productImageUrl } from '../../utils/product-image-url';

export interface StatusChange {
  item: PantryItem;
  status: ShoppingStatus;
}

/**
 * One product in the pantry grid: image, name, status pill, settings button and
 * the valid actions for its current shopping state.
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

  /** Emitted when the user chooses one of the allowed shopping transitions. */
  readonly statusChanged = output<StatusChange>();

  /** Emitted when the user opens the item settings. */
  readonly settingsOpened = output<PantryItem>();

  /** Emitted from the explicit variants control on products that have choices. */
  readonly variantsOpened = output<PantryItem>();

  constructor() {
    // Standalone Ionic requires every icon to be registered explicitly, which
    // keeps unused icons out of the bundle.
    addIcons({ cartOutline, ellipsisHorizontal, layersOutline, timeOutline, trashOutline });
  }

  protected imageUrl(): string {
    return productImageUrl(this.item().product.image);
  }

  protected changeStatus(status: ShoppingStatus): void {
    this.statusChanged.emit({ item: this.item(), status });
  }
}
