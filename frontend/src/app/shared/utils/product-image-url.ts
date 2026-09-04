/**
 * Cache revision for catalog images whose stable file names are stored in the
 * database. Bump it whenever existing SVG artwork changes.
 */
const PRODUCT_IMAGE_REVISION = '20260905';

export function productImageUrl(image: string): string {
  return `assets/products/${image}?v=${PRODUCT_IMAGE_REVISION}`;
}
