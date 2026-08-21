import type { PantryItem, SortMode, StockStatus } from '../../shared/models/pantry.model';

/** Most urgent first, so a shopping trip starts at the top of the grid. */
const STATUS_ORDER: Record<StockStatus, number> = { OUT: 0, LOW: 1, OK: 2 };

/**
 * Returns a new array ordered by `mode`.
 *
 * `DEFAULT` keeps the catalog order the API returned. Array sort is stable, so
 * sorting by stock level preserves the catalog order within each status.
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
