/** Current state in the shopping workflow, or archived when hidden. */
export type ItemStatus = 'DISCARDED' | 'PENDING' | 'IN_CART' | 'ARCHIVED';
export type ShoppingStatus = Exclude<ItemStatus, 'ARCHIVED'>;

/** How often a product is bought. The coarse filter of the pantry screen. */
export type ProductType = 'ESSENTIAL' | 'SECONDARY';

/** Product group used for quick filtering, mirroring supermarket aisles. */
export type Category =
  | 'FRUIT_VEG'
  | 'MEAT_FISH'
  | 'DAIRY_EGGS'
  | 'DRY_CANNED'
  | 'DRINKS'
  | 'HOME_CARE';

/** Sentinels used by the filter bars for active and archived views. */
export const ALL = 'ALL';
export const ARCHIVED = 'ARCHIVED';
export type CategoryFilter = Category | typeof ALL;
export type TypeFilter = ProductType | typeof ALL | typeof ARCHIVED;
export type StatusFilter = ShoppingStatus | typeof ALL;

/** How the product grid is ordered. */
export type SortMode = 'DEFAULT' | 'NAME' | 'STATUS';

/** A catalog entry. */
export interface ProductVariant {
  id: string;
  code: string;
  name: string;
  image: string;
}

export interface Product {
  id: string;
  code: string;
  name: string;
  /** File name under assets/products, e.g. "tomato.svg". */
  image: string;
  variants: ProductVariant[];
}

/** The per-pantry state of a catalog product. */
export interface PantryItem {
  product: Product;
  status: ItemStatus;
  type: ProductType;
  category: Category;
  selectedVariantIds: string[];
  updatedAt: string;
}

/** A pantry without its items. */
export interface Pantry {
  id: string;
  slug: string;
  name: string;
  createdAt: string;
  updatedAt: string;
}

/** A pantry with all of its items. */
export interface PantryDetail extends Pantry {
  items: PantryItem[];
}

/** Partial update of an item; omitted fields are left untouched. */
export interface ItemPatch {
  status?: ItemStatus;
  type?: ProductType;
  category?: Category;
  selectedVariantIds?: string[];
}

/** Error envelope returned by the API. */
export interface ApiError {
  error: { code: string; message: string };
}
