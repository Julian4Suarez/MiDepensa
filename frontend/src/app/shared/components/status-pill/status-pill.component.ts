import { ChangeDetectionStrategy, Component, input } from '@angular/core';

import { STATUS_META } from '../../models/pantry.meta';
import type { ItemStatus } from '../../models/pantry.model';

/** Coloured pill showing the shopping status of a product. */
@Component({
  selector: 'app-status-pill',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `<span class="pill" [attr.data-status]="status()">{{ label() }}</span>`,
  styles: `
    .pill {
      display: inline-block;
      padding: 4px 12px;
      border-radius: 999px;
      font-size: 0.75rem;
      font-weight: 600;
      line-height: 1.2;
      white-space: nowrap;
    }

    .pill[data-status='DISCARDED'] {
      background: var(--ion-color-light-shade);
      color: var(--ion-color-medium-shade);
    }

    .pill[data-status='PENDING'] {
      background: var(--app-status-low-soft);
      color: var(--ion-color-warning-shade);
    }

    .pill[data-status='IN_CART'] {
      background: var(--app-status-ok-soft);
      color: var(--ion-color-success-shade);
    }

    .pill[data-status='ARCHIVED'] {
      background: var(--ion-color-light-shade);
      color: var(--ion-color-medium-shade);
    }
  `,
})
export class StatusPillComponent {
  /** Status to render. `input()` is the signal-based replacement for `@Input()`. */
  readonly status = input.required<ItemStatus>();

  protected label(): string {
    return STATUS_META[this.status()].label;
  }
}
