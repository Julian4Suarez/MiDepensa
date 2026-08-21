import { Injectable, computed, inject, signal } from '@angular/core';
import { firstValueFrom } from 'rxjs';

import { PantryApiService } from '../../core/services/pantry-api.service';
import { CATEGORIES, NEXT_STATUS } from '../../shared/models/pantry.meta';
import type {
  Category,
  CategoryFilter,
  ItemPatch,
  Pantry,
  PantryItem,
  PantryView,
} from '../../shared/models/pantry.model';

/**
 * State container for a single pantry page.
 *
 * Provided by the pantry route rather than in root, so navigating away disposes
 * the state. Everything the template needs is a signal or a computed signal:
 * changing `view` or `category` re-renders the grid without any subscription.
 */
@Injectable()
export class PantryStore {
  private readonly api = inject(PantryApiService);

  // ── Writable state ────────────────────────────────────────
  readonly pantry = signal<Pantry | null>(null);
  readonly items = signal<PantryItem[]>([]);
  readonly view = signal<PantryView>('PRIMARY');
  readonly category = signal<CategoryFilter>('ALL');
  readonly loading = signal(false);
  readonly error = signal<string | null>(null);

  // ── Derived state ─────────────────────────────────────────

  /** Items filed under the selected view, ignoring the category filter. */
  readonly itemsInView = computed(() => this.items().filter((item) => item.view === this.view()));

  /** Items shown in the grid: selected view, then selected category. */
  readonly visibleItems = computed(() => {
    const category = this.category();
    return category === 'ALL'
      ? this.itemsInView()
      : this.itemsInView().filter((item) => item.category === category);
  });

  /** Only the categories actually present in the current view get a chip. */
  readonly availableCategories = computed<Category[]>(() =>
    CATEGORIES.filter((category) => this.itemsInView().some((item) => item.category === category)),
  );

  /** How many products each view needs restocking, shown in the side menu. */
  readonly missingByView = computed<Record<PantryView, number>>(() => {
    const counts: Record<PantryView, number> = { PRIMARY: 0, SECONDARY: 0, OTHER: 0 };
    for (const item of this.items()) {
      if (item.status !== 'OK') {
        counts[item.view] += 1;
      }
    }
    return counts;
  });

  /** Number of products in the current view that need buying. */
  readonly missingInView = computed(
    () => this.itemsInView().filter((item) => item.status !== 'OK').length,
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
