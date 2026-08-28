import type { ItemStatus, PantryItem, SortMode } from '../../shared/models/pantry.model';

/** Pending decisions first, followed by cart, discarded and archived items. */
const STATUS_ORDER: Record<ItemStatus, number> = {
  PENDING: 0,
  IN_CART: 1,
  DISCARDED: 2,
  ARCHIVED: 3,
};

/**
 * Returns a new array ordered by `mode`.
 *
 * `DEFAULT` keeps the catalog order the API returned. Array sort is stable, so
 * sorting by status preserves the catalog order within each state.
 */
export function sortItems(items: PantryItem[], mode: SortMode): PantryItem[] {
  switch (mode) {
    case 'NAME':
      return [...items].sort((a, b) => a.product.name.localeCompare(b.product.name));
    case 'STATUS':
      return [...items].sort((a, b) => STATUS_ORDER[a.status] - STATUS_ORDER[b.status]);
    default:
      return items;
  }
}
