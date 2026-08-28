import type {
  Category,
  ItemStatus,
  ProductType,
  ShoppingStatus,
  SortMode,
} from './pantry.model';

/** Icon used by the "All" chip of every filter bar. */
export const ALL_ICON = 'apps-outline';

/**
 * Display order of the ten categories. They follow the order in which you walk
 * a supermarket, so the generated shopping list matches the route you take.
 */
export const CATEGORIES: readonly Category[] = [
  'FRUIT',
  'VEGETABLES',
  'MEAT_FISH',
  'DAIRY_EGGS',
  'BAKERY_BREAKFAST_SNACKS',
  'PASTA_RICE_LEGUMES',
  'COOKING_CONDIMENTS',
  'DRINKS',
  'HOUSEHOLD_CLEANING',
  'PERSONAL_CARE',
];

export const CATEGORY_META: Record<Category, { label: string; icon: string }> = {
  FRUIT: { label: 'Fruit', icon: 'nutrition-outline' },
  VEGETABLES: { label: 'Vegetables', icon: 'leaf-outline' },
  MEAT_FISH: { label: 'Meat & fish', icon: 'fish-outline' },
  DAIRY_EGGS: { label: 'Dairy & eggs', icon: 'egg-outline' },
  BAKERY_BREAKFAST_SNACKS: { label: 'Bakery, breakfast & snacks', icon: 'cafe-outline' },
  PASTA_RICE_LEGUMES: { label: 'Pasta, rice & legumes', icon: 'restaurant-outline' },
  COOKING_CONDIMENTS: { label: 'Cooking essentials & condiments', icon: 'flame-outline' },
  DRINKS: { label: 'Drinks', icon: 'wine-outline' },
  HOUSEHOLD_CLEANING: { label: 'Household & cleaning', icon: 'construct-outline' },
  PERSONAL_CARE: { label: 'Personal care', icon: 'person-outline' },
};

/** Display order and labelling of the two active product types. */
export const TYPES: readonly ProductType[] = ['ESSENTIAL', 'SECONDARY'];

export const TYPE_META: Record<ProductType, { label: string; icon: string }> = {
  ESSENTIAL: { label: 'Essential', icon: 'star-outline' },
  SECONDARY: { label: 'Secondary', icon: 'bookmark-outline' },
};

export const SHOPPING_STATUSES: readonly ShoppingStatus[] = ['DISCARDED', 'PENDING', 'IN_CART'];

export const STATUS_META: Record<ItemStatus, { label: string; color: string }> = {
  DISCARDED: { label: 'Discarded', color: 'medium' },
  PENDING: { label: 'Pending', color: 'warning' },
  IN_CART: { label: 'In cart', color: 'success' },
  ARCHIVED: { label: 'Archived', color: 'medium' },
};

/** Sort options offered by the toolbar, in the order they are listed. */
export const SORT_MODES: readonly SortMode[] = ['DEFAULT', 'NAME', 'STATUS'];

export const SORT_META: Record<SortMode, { label: string }> = {
  DEFAULT: { label: 'Catalog order' },
  NAME: { label: 'Name (A–Z)' },
  STATUS: { label: 'Shopping status' },
};
