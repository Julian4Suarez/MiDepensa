import { TestBed } from '@angular/core/testing';
import { of, throwError } from 'rxjs';

import { PantryApiService } from '../../core/services/pantry-api.service';
import type { Category, PantryItem, ProductType } from '../../shared/models/pantry.model';
import { PantryStore } from './pantry.store';

function item(
  name: string,
  type: ProductType,
  category: Category,
  status: PantryItem['status'] = 'OK',
): PantryItem {
  return {
    product: { id: name, code: name, name, image: `${name}.svg` },
    status,
    type,
    category,
    updatedAt: '2026-01-01T00:00:00Z',
  };
}

describe('PantryStore', () => {
  const items = [
    item('Tomatoes', 'ESSENTIAL', 'FRUIT_VEG'),
    item('Rice', 'ESSENTIAL', 'DRY_CANNED'),
    item('Wine', 'SECONDARY', 'DRINKS', 'ARCHIVED'),
    item('Dish soap', 'SECONDARY', 'HOME_CARE'),
  ];

  let api: jest.Mocked<Pick<PantryApiService, 'getPantry' | 'updateItem'>>;
  let store: PantryStore;

  beforeEach(() => {
    api = {
      getPantry: jest.fn(),
      updateItem: jest.fn(),
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

  it('shows only active products for All and only archived products for Archived', async () => {
    await loadPantry();

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

    store.category.set('DRY_CANNED');

    expect(store.visibleItems().map((i) => i.product.name)).toEqual(['Rice']);
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

  it('cycles the status and persists it', async () => {
    await loadPantry();
    api.updateItem.mockReturnValue(of({ ...items[0], status: 'LOW' }));

    await store.cycleStatus(store.items()[0]);

    expect(api.updateItem).toHaveBeenCalledWith('familia', 'Tomatoes', { status: 'LOW' });
    expect(store.items()[0].status).toBe('LOW');
  });

  it('does not cycle the stock status of an archived product', async () => {
    await loadPantry();

    await store.cycleStatus(store.items()[2]);

    expect(api.updateItem).not.toHaveBeenCalled();
    expect(store.items()[2].status).toBe('ARCHIVED');
  });

  it('rolls the optimistic update back when the request fails', async () => {
    await loadPantry();
    api.updateItem.mockReturnValue(throwError(() => new Error('offline')));

    await store.cycleStatus(store.items()[0]);

    expect(store.items()[0].status).toBe('OK');
    expect(store.error()).not.toBeNull();
  });
});
