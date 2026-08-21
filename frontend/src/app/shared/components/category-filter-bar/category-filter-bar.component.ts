import { ChangeDetectionStrategy, Component, input, output } from '@angular/core';
import { IonIcon } from '@ionic/angular/standalone';
import { addIcons } from 'ionicons';
import {
  appsOutline,
  eggOutline,
  fileTrayStackedOutline,
  fishOutline,
  nutritionOutline,
  sparklesOutline,
  wineOutline,
} from 'ionicons/icons';

import { CATEGORY_META } from '../../models/pantry.meta';
import type { Category, CategoryFilter } from '../../models/pantry.model';

/**
 * Horizontal, single-select row of category chips. Faster to hit than a select
 * on a phone, and it shows every option at once.
 */
@Component({
  selector: 'app-category-filter-bar',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [IonIcon],
  templateUrl: './category-filter-bar.component.html',
  styleUrl: './category-filter-bar.component.scss',
})
export class CategoryFilterBarComponent {
  readonly categories = input.required<Category[]>();
  readonly selected = input.required<CategoryFilter>();
  readonly selectedChange = output<CategoryFilter>();

  protected readonly meta = CATEGORY_META;

  constructor() {
    addIcons({
      nutritionOutline,
      fishOutline,
      eggOutline,
      fileTrayStackedOutline,
      wineOutline,
      sparklesOutline,
      appsOutline,
    });
  }
}
