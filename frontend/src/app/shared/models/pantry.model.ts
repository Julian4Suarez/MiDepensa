/** How much of a product is left. Mirrors the backend enum. */
export type StockStatus = 'OUT' | 'LOW' | 'OK';

/** Section of the pantry a product is filed under. */
export type PantryView = 'PRIMARY' | 'SECONDARY' | 'OTHER';

/** Product group used for quick filtering. */
export type Category = 'FRESH' | 'PANTRY' | 'DRINKS' | 'HOME_CARE';

/** Sentinel used by the filter bar when no category is selected. */
export type CategoryFilter = Category | 'ALL';

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
  view: PantryView;
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
  view?: PantryView;
  category?: Category;
}

/** Error envelope returned by the API. */
export interface ApiError {
  error: { code: string; message: string };
}
