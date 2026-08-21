import { ChangeDetectionStrategy, Component, input } from '@angular/core';

import { STATUS_META } from '../../models/pantry.meta';
import type { StockStatus } from '../../models/pantry.model';

/** Coloured pill showing the stock status of a product. */
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

    .pill[data-status='OUT'] {
      background: var(--app-status-out-soft);
      color: var(--ion-color-danger-shade);
    }

    .pill[data-status='LOW'] {
      background: var(--app-status-low-soft);
      color: var(--ion-color-warning-shade);
    }

    .pill[data-status='OK'] {
      background: var(--app-status-ok-soft);
      color: var(--ion-color-success-shade);
    }
  `,
})
export class StatusPillComponent {
  /** Status to render. `input()` is the signal-based replacement for `@Input()`. */
  readonly status = input.required<StockStatus>();

  protected label(): string {
    return STATUS_META[this.status()].label;
  }
}
