/** How much of a product is left. Mirrors the backend enum. */
export type StockStatus = 'OUT' | 'LOW' | 'OK';

/** How often a product is bought. The coarse filter of the pantry screen. */
export type ProductType = 'ESSENTIAL' | 'SECONDARY' | 'OTHERS';

/** Product group used for quick filtering, mirroring supermarket aisles. */
export type Category =
  | 'FRUIT_VEG'
  | 'MEAT_FISH'
  | 'DAIRY_EGGS'
  | 'DRY_CANNED'
  | 'DRINKS'
  | 'HOME_CARE';

/** Sentinel used by the filter bars when nothing is narrowed down. */
export const ALL = 'ALL';
export type CategoryFilter = Category | typeof ALL;
export type TypeFilter = ProductType | typeof ALL;

/** How the product grid is ordered. */
export type SortMode = 'DEFAULT' | 'NAME' | 'STATUS';

/** A catalog entry. */
export interface Product {
  id: string;
  code: string;
  name: string;
  /** File name under assets/products, e.g. "tomato.svg". */
  image: string;
}

/** The per-pantry state of a catalog product. */
export interface PantryItem {
  product: Product;
  status: StockStatus;
  type: ProductType;
  category: Category;
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
  status?: StockStatus;
  type?: ProductType;
  category?: Category;
}

/** Error envelope returned by the API. */
export interface ApiError {
  error: { code: string; message: string };
}
