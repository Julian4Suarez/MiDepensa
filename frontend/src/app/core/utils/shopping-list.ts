import { CATEGORIES, CATEGORY_META } from '../../shared/models/pantry.meta';
import type { PantryItem } from '../../shared/models/pantry.model';

/**
 * Builds the plain-text shopping list, grouped by category.
 *
 * Only products explicitly placed in the cart are included. Returns an empty
 * string when the cart is empty.
 */
export function buildShoppingList(items: PantryItem[]): string {
  const inCart = items.filter((item) => item.status === 'IN_CART');

  return CATEGORIES.map((category) => {
    const lines = inCart
      .filter((item) => item.category === category)
      .flatMap((item) => {
        if (item.product.variants.length === 0) {
          return [`- ${item.product.name}`];
        }
        const selected = new Set(item.selectedVariantIds);
        if (selected.size === 0) {
          return [`- ${item.product.name}`];
        }
        return item.product.variants
          .filter((variant) => selected.has(variant.id))
          .map((variant) => `- ${variant.name}`);
      });

    return lines.length ? `${CATEGORY_META[category].label}\n${lines.join('\n')}` : '';
  })
    .filter(Boolean)
    .join('\n\n');
}
