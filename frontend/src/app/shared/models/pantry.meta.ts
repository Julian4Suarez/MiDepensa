import type { Category, ProductType, SortMode, StockStatus } from './pantry.model';

/** Icon used by the "All" chip of every filter bar. */
export const ALL_ICON = 'apps-outline';

/**
 * Display order of the six categories. They follow the order in which you walk
 * a supermarket, so the generated shopping list matches the route you take.
 */
export const CATEGORIES: readonly Category[] = [
  'FRUIT_VEG',
  'MEAT_FISH',
  'DAIRY_EGGS',
  'DRY_CANNED',
  'DRINKS',
  'HOME_CARE',
];

export const CATEGORY_META: Record<Category, { label: string; icon: string }> = {
  FRUIT_VEG: { label: 'Fruit & veg', icon: 'nutrition-outline' },
  MEAT_FISH: { label: 'Meat & fish', icon: 'fish-outline' },
  DAIRY_EGGS: { label: 'Dairy & eggs', icon: 'egg-outline' },
  DRY_CANNED: { label: 'Dry & canned', icon: 'file-tray-stacked-outline' },
  DRINKS: { label: 'Drinks', icon: 'wine-outline' },
  HOME_CARE: { label: 'Home & care', icon: 'sparkles-outline' },
};

/** Display order and labelling of the three product types. */
export const TYPES: readonly ProductType[] = ['ESSENTIAL', 'SECONDARY', 'OTHERS'];

export const TYPE_META: Record<ProductType, { label: string; icon: string }> = {
  ESSENTIAL: { label: 'Essential', icon: 'star-outline' },
  SECONDARY: { label: 'Secondary', icon: 'bookmark-outline' },
  OTHERS: { label: 'Others', icon: 'ellipsis-horizontal-outline' },
};

export const STATUS_META: Record<StockStatus, { label: string; color: string }> = {
  OUT: { label: 'Out', color: 'danger' },
  LOW: { label: 'Low', color: 'warning' },
  OK: { label: 'Enough', color: 'success' },
};

/**
 * Tapping the status pill walks stock downwards and wraps around, which is the
 * direction a pantry is actually used: enough -> low -> out -> restocked.
 */
export const NEXT_STATUS: Record<StockStatus, StockStatus> = {
  OK: 'LOW',
  LOW: 'OUT',
  OUT: 'OK',
};

/** Sort options offered by the toolbar, in the order they are listed. */
export const SORT_MODES: readonly SortMode[] = ['DEFAULT', 'NAME', 'STATUS'];

export const SORT_META: Record<SortMode, { label: string }> = {
  DEFAULT: { label: 'Catalog order' },
  NAME: { label: 'Name (A–Z)' },
  STATUS: { label: 'Stock level' },
};
