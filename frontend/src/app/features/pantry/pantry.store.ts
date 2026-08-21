import { Injectable, computed, inject, signal } from '@angular/core';
import { firstValueFrom } from 'rxjs';

import { PantryApiService } from '../../core/services/pantry-api.service';
import { sortItems } from '../../core/utils/sort-items';
import { CATEGORIES, NEXT_STATUS } from '../../shared/models/pantry.meta';
import {
  ALL,
  type Category,
  type CategoryFilter,
  type ItemPatch,
  type Pantry,
  type PantryItem,
  type SortMode,
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
  readonly sort = signal<SortMode>('DEFAULT');
  readonly loading = signal(false);
  readonly error = signal<string | null>(null);

  // ── Derived state ─────────────────────────────────────────

  /** Items matching the type filter, before narrowing down by category. */
  readonly itemsOfType = computed(() => {
    const type = this.type();
    return type === ALL ? this.items() : this.items().filter((item) => item.type === type);
  });

  /** Items shown in the grid: type, then category, then sort order. */
  readonly visibleItems = computed(() => {
    const category = this.category();
    const filtered =
      category === ALL
        ? this.itemsOfType()
        : this.itemsOfType().filter((item) => item.category === category);

    return sortItems(filtered, this.sort());
  });

  /** Only the categories actually present under the current type get a chip. */
  readonly availableCategories = computed<Category[]>(() =>
    CATEGORIES.filter((category) => this.itemsOfType().some((item) => item.category === category)),
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

  /** Narrowing by type can leave the category filter pointing at nothing. */
  selectType(type: TypeFilter): void {
    this.type.set(type);
    this.category.set(ALL);
  }

  /** Advances an item to the next stock status. */
  cycleStatus(item: PantryItem): Promise<void> {
    return this.patch(item, { status: NEXT_STATUS[item.status] });
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
