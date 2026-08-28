import { TestBed } from '@angular/core/testing';
import { of, throwError } from 'rxjs';

import { PantryApiService } from '../../core/services/pantry-api.service';
import type { Category, PantryItem, ProductType } from '../../shared/models/pantry.model';
import { PantryStore } from './pantry.store';

function item(
  name: string,
  type: ProductType,
  category: Category,
  status: PantryItem['status'] = 'PENDING',
): PantryItem {
  return {
    product: { id: name, code: name, name, image: `${name}.svg`, variants: [] },
    status,
    type,
    category,
    selectedVariantIds: [],
    updatedAt: '2026-01-01T00:00:00Z',
  };
}

describe('PantryStore', () => {
  const items = [
    item('Tomatoes', 'ESSENTIAL', 'FRUIT_VEG'),
    item('Rice', 'ESSENTIAL', 'DRY_CANNED', 'DISCARDED'),
    item('Wine', 'SECONDARY', 'DRINKS', 'ARCHIVED'),
    item('Dish soap', 'SECONDARY', 'HOME_CARE', 'IN_CART'),
  ];

  let api: jest.Mocked<
    Pick<PantryApiService, 'getPantry' | 'updateItem' | 'resetActiveItems'>
  >;
  let store: PantryStore;

  beforeEach(() => {
    api = {
      getPantry: jest.fn(),
      updateItem: jest.fn(),
      resetActiveItems: jest.fn(),
    };

    TestBed.configureTestingModule({
      providers: [PantryStore, { provide: PantryApiService, useValue: api }],
    });
    store = TestBed.inject(PantryStore);
  });

  async function loadPantry(): Promise<void> {
    api.getPantry.mockReturnValue(
      of({
        id: 'id',
        slug: 'familia',
        name: 'Familia',
        createdAt: '',
        updatedAt: '',
        items,
      }),
    );
    await store.load('familia');
  }

  it('uses pending as the default shopping-status filter', async () => {
    await loadPantry();

    expect(store.status()).toBe('PENDING');
    expect(store.visibleItems().map((i) => i.product.name)).toEqual(['Tomatoes']);
  });

  it('shows only active products for All and only archived products for Archived', async () => {
    await loadPantry();
    store.status.set('ALL');

    expect(store.visibleItems().map((i) => i.product.name)).toEqual(['Tomatoes', 'Rice']);

    store.selectType('ARCHIVED');
    expect(store.visibleItems().map((i) => i.product.name)).toEqual(['Wine']);

    store.selectType('ALL');
    expect(store.visibleItems().map((i) => i.product.name)).toEqual([
      'Tomatoes',
      'Rice',
      'Dish soap',
    ]);
  });

  it('filters by the selected category inside the type', async () => {
    await loadPantry();
    store.status.set('ALL');

    store.category.set('DRY_CANNED');

    expect(store.visibleItems().map((i) => i.product.name)).toEqual(['Rice']);
  });

  it('filters active products by shopping status', async () => {
    await loadPantry();
    store.selectType('ALL');

    store.status.set('IN_CART');

    expect(store.visibleItems().map((i) => i.product.name)).toEqual(['Dish soap']);
  });

  it('collects every cart product independently of screen filters', async () => {
    await loadPantry();

    expect(store.cartItems().map((i) => i.product.name)).toEqual(['Dish soap']);
  });

  it('resets every active product to pending and preserves archived products', async () => {
    await loadPantry();
    api.resetActiveItems.mockReturnValue(of(undefined));

    await store.resetActiveItems();

    expect(api.resetActiveItems).toHaveBeenCalledWith('familia');
    expect(store.items().map((i) => i.status)).toEqual([
      'PENDING',
      'PENDING',
      'ARCHIVED',
      'PENDING',
    ]);
  });

  it('only offers chips for categories present under the current type', async () => {
    await loadPantry();

    expect(store.availableCategories()).toEqual(['FRUIT_VEG', 'DRY_CANNED']);
  });

  it('keeps the selected category when the type changes', async () => {
    await loadPantry();
    store.category.set('DRY_CANNED');

    store.selectType('ARCHIVED');

    expect(store.category()).toBe('DRY_CANNED');
    expect(store.availableCategories()).toEqual(['DRY_CANNED', 'DRINKS']);
    expect(store.visibleItems()).toEqual([]);
  });

  it('also keeps Home and Care selected and visible when the type changes', async () => {
    await loadPantry();
    store.selectType('SECONDARY');
    store.category.set('HOME_CARE');

    store.selectType('ARCHIVED');

    expect(store.category()).toBe('HOME_CARE');
    expect(store.availableCategories()).toEqual(['DRINKS', 'HOME_CARE']);
    expect(store.visibleItems()).toEqual([]);
  });

  it('changes the shopping status and persists it', async () => {
    await loadPantry();
    api.updateItem.mockReturnValue(of({ ...items[0], status: 'IN_CART' }));

    await store.setStatus(store.items()[0], 'IN_CART');

    expect(api.updateItem).toHaveBeenCalledWith('familia', 'Tomatoes', {
      status: 'IN_CART',
    });
    expect(store.items()[0].status).toBe('IN_CART');
  });

  it('rejects transitions that skip the pending state', async () => {
    await loadPantry();

    await store.setStatus(store.items()[1], 'IN_CART');

    expect(api.updateItem).not.toHaveBeenCalled();
    expect(store.items()[1].status).toBe('DISCARDED');
  });

  it('rolls the optimistic update back when the request fails', async () => {
    await loadPantry();
    api.updateItem.mockReturnValue(throwError(() => new Error('offline')));

    await store.setStatus(store.items()[0], 'DISCARDED');

    expect(store.items()[0].status).toBe('PENDING');
    expect(store.error()).not.toBeNull();
  });
});
