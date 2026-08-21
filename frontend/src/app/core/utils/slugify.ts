/** Maximum slug length accepted by the API. */
export const MAX_SLUG_LENGTH = 60;

/**
 * Turns a free-text pantry name into the URL-safe slug the API will derive.
 *
 * Kept in sync with the Go implementation so the home page can preview the
 * final URL while the user types.
 */
export function slugify(raw: string): string {
  return raw
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '')
    .slice(0, MAX_SLUG_LENGTH)
    .replace(/-+$/g, '');
}
