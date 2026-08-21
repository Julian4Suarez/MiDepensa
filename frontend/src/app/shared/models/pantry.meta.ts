import type { Category, PantryView, StockStatus } from './pantry.model';

/** Display order and labelling of the four product categories. */
export const CATEGORIES: readonly Category[] = ['FRESH', 'PANTRY', 'DRINKS', 'HOME_CARE'];

export const CATEGORY_META: Record<Category, { label: string; icon: string }> = {
  FRESH: { label: 'Fresh', icon: 'leaf-outline' },
  PANTRY: { label: 'Pantry', icon: 'file-tray-stacked-outline' },
  DRINKS: { label: 'Drinks', icon: 'wine-outline' },
  HOME_CARE: { label: 'Home & care', icon: 'sparkles-outline' },
};

/** Display order and labelling of the three pantry views. */
export const VIEWS: readonly PantryView[] = ['PRIMARY', 'SECONDARY', 'OTHER'];

export const VIEW_META: Record<PantryView, { label: string; icon: string }> = {
  PRIMARY: { label: 'Primary', icon: 'star-outline' },
  SECONDARY: { label: 'Secondary', icon: 'bookmark-outline' },
  OTHER: { label: 'Other', icon: 'apps-outline' },
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
