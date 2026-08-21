import { CATEGORIES, CATEGORY_META } from '../../shared/models/pantry.meta';
import type { PantryItem } from '../../shared/models/pantry.model';

/**
 * Builds the plain-text shopping list, grouped by category.
 *
 * Products that are out of stock are always included; running-low products are
 * added only when `includeLow` is set and are suffixed with "(low)".
 * Returns an empty string when nothing needs buying.
 */
export function buildShoppingList(items: PantryItem[], includeLow: boolean): string {
  const needed = items.filter(
    (item) => item.status === 'OUT' || (includeLow && item.status === 'LOW'),
  );

  return CATEGORIES.map((category) => {
    const lines = needed
      .filter((item) => item.category === category)
      .map((item) => `- ${item.product.name}${item.status === 'LOW' ? ' (low)' : ''}`);

    return lines.length ? `${CATEGORY_META[category].label}\n${lines.join('\n')}` : '';
  })
    .filter(Boolean)
    .join('\n\n');
}
