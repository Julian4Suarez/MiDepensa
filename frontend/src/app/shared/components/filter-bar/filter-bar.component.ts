import { ChangeDetectionStrategy, Component, input, output } from '@angular/core';
import { IonIcon } from '@ionic/angular/standalone';

/** One selectable chip. `value` is opaque to the bar. */
export interface FilterOption {
  value: string;
  label: string;
  icon: string;
}

/**
 * Horizontal, single-select row of chips. Faster to hit than a select on a
 * phone, and it shows every option at once.
 *
 * Icons are registered by the page that uses the bar, so unused ones stay out
 * of the bundle.
 */
@Component({
  selector: 'app-filter-bar',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [IonIcon],
  templateUrl: './filter-bar.component.html',
  styleUrl: './filter-bar.component.scss',
})
export class FilterBarComponent {
  readonly options = input.required<FilterOption[]>();
  readonly selected = input.required<string>();
  readonly label = input.required<string>();
  readonly selectedChange = output<string>();

  /** Makes a regular desktop mouse wheel useful on the horizontal chip row. */
  protected scrollWithWheel(event: WheelEvent, bar: HTMLElement): void {
    if (bar.scrollWidth <= bar.clientWidth || Math.abs(event.deltaX) >= Math.abs(event.deltaY)) {
      return;
    }

    bar.scrollLeft += event.deltaY;
    event.preventDefault();
  }
}
