import { Injectable, computed, inject, signal } from '@angular/core';
import { firstValueFrom } from 'rxjs';

import { PantryApiService } from '../../core/services/pantry-api.service';
import { sortItems } from '../../core/utils/sort-items';
import { CATEGORIES } from '../../shared/models/pantry.meta';
import {
  ALL,
  ARCHIVED,
  type Category,
  type CategoryFilter,
  type ItemPatch,
  type Pantry,
  type PantryItem,
  type ShoppingStatus,
  type SortMode,
  type StatusFilter,
  type TypeFilter,
} from '../../shared/models/pantry.model';

/**
 * State container for a single pantry page.
 *
 * Provided by the pantry route rather than in root, so navigating away disposes
 * the state. Everything the template needs is a signal or a computed signal:
 * changing a filter re-renders the grid without any subscription.
 */
@Injectable()
export class PantryStore {
  private readonly api = inject(PantryApiService);

  // ── Writable state ────────────────────────────────────────
  readonly pantry = signal<Pantry | null>(null);
  readonly items = signal<PantryItem[]>([]);
  // Opens on the products bought every time, which is the common case.
  readonly type = signal<TypeFilter>('ESSENTIAL');
  readonly category = signal<CategoryFilter>(ALL);
  // Pending is the working view: products that still need a decision.
  readonly status = signal<StatusFilter>('PENDING');
  readonly sort = signal<SortMode>('DEFAULT');
  readonly loading = signal(false);
  readonly resetting = signal(false);
  readonly error = signal<string | null>(null);

  // ── Derived state ─────────────────────────────────────────

  /** Items matching the type filter, before narrowing down by category. */
  readonly itemsOfType = computed(() => {
    const type = this.type();
    if (type === ARCHIVED) {
      return this.items().filter((item) => item.status === ARCHIVED);
    }
    if (type === ALL) {
      return this.items().filter((item) => item.status !== ARCHIVED);
    }
    return this.items().filter((item) => item.type === type && item.status !== ARCHIVED);
  });

  /** Items shown in the grid: type, then category, then sort order. */
  readonly visibleItems = computed(() => {
    const category = this.category();
    const categoryItems =
      category === ALL
        ? this.itemsOfType()
        : this.itemsOfType().filter((item) => item.category === category);
    const status = this.status();
    const filtered =
      this.type() === ARCHIVED || status === ALL
        ? categoryItems
        : categoryItems.filter((item) => item.status === status);

    return sortItems(filtered, this.sort());
  });

  /** Every product on the shopping list, regardless of the screen filters. */
  readonly cartItems = computed(() => this.items().filter((item) => item.status === 'IN_CART'));

  /**
   * Only categories present under the current type get a chip, except for the
   * active category, which stays visible while changing type.
   */
  readonly availableCategories = computed<Category[]>(() =>
    CATEGORIES.filter(
      (category) =>
        this.category() === category ||
        this.itemsOfType().some((item) => item.category === category),
    ),
  );

  // ── Commands ──────────────────────────────────────────────

  /** Loads the pantry addressed by the URL slug. */
  async load(slug: string): Promise<void> {
    this.loading.set(true);
    this.error.set(null);
    try {
      const detail = await firstValueFrom(this.api.getPantry(slug));
      const { items, ...pantry } = detail;
      this.pantry.set(pantry);
      this.items.set(items);
    } catch {
      this.error.set('This pantry could not be loaded.');
    } finally {
      this.loading.set(false);
    }
  }

  /** Changes the type without discarding the independently selected category. */
  selectType(type: TypeFilter): void {
    this.type.set(type);
  }

  /** Moves an item along one of the explicitly allowed shopping transitions. */
  setStatus(item: PantryItem, status: ShoppingStatus): Promise<void> {
    const allowed =
      (item.status === 'PENDING' && (status === 'DISCARDED' || status === 'IN_CART')) ||
      ((item.status === 'DISCARDED' || item.status === 'IN_CART') && status === 'PENDING');
    if (!allowed) {
      return Promise.resolve();
    }
    return this.patch(item, { status });
  }

  /** Resets every active product to pending, preserving archived products. */
  async resetActiveItems(): Promise<void> {
    const slug = this.pantry()?.slug;
    if (!slug || this.resetting()) {
      return;
    }

    const previous = this.items();
    this.resetting.set(true);
    this.items.update((items) =>
      items.map((item) => (item.status === ARCHIVED ? item : { ...item, status: 'PENDING' })),
    );
    try {
      await firstValueFrom(this.api.resetActiveItems(slug));
    } catch {
      this.items.set(previous);
      this.error.set('The products could not be reset.');
    } finally {
      this.resetting.set(false);
    }
  }

  /** Applies a settings change coming from the item modal. */
  patch(item: PantryItem, patch: ItemPatch): Promise<void> {
    return this.applyOptimistically(item, patch);
  }

  /**
   * Updates the signal first so the UI reacts instantly, then persists. The
   * previous value is restored if the request fails.
   */
  private async applyOptimistically(item: PantryItem, patch: ItemPatch): Promise<void> {
    const slug = this.pantry()?.slug;
    if (!slug) {
      return;
    }

    const previous = this.items();
    this.items.update((items) =>
      items.map((current) =>
        current.product.id === item.product.id ? { ...current, ...patch } : current,
      ),
    );

    try {
      const updated = await firstValueFrom(this.api.updateItem(slug, item.product.id, patch));
      this.items.update((items) =>
        items.map((current) => (current.product.id === updated.product.id ? updated : current)),
      );
    } catch {
      this.items.set(previous);
      this.error.set('The change could not be saved.');
    }
  }
}
