import { TestBed } from '@angular/core/testing';
import { of, throwError } from 'rxjs';

import { PantryApiService } from '../../core/services/pantry-api.service';
import type { Category, PantryItem, PantryView } from '../../shared/models/pantry.model';
import { PantryStore } from './pantry.store';

function item(name: string, view: PantryView, category: Category): PantryItem {
  return {
    product: { id: name, code: name, name, image: `${name}.svg` },
    status: 'OK',
    view,
    category,
    updatedAt: '2026-01-01T00:00:00Z',
  };
}

describe('PantryStore', () => {
  const items = [
    item('Tomatoes', 'PRIMARY', 'FRUIT_VEG'),
    item('Rice', 'PRIMARY', 'DRY_CANNED'),
    item('Wine', 'OTHER', 'DRINKS'),
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

  it('filters by the selected view', async () => {
    await loadPantry();

    expect(store.visibleItems().map((i) => i.product.name)).toEqual(['Tomatoes', 'Rice']);

    store.view.set('OTHER');
    expect(store.visibleItems().map((i) => i.product.name)).toEqual(['Wine']);
  });

  it('filters by the selected category inside the view', async () => {
    await loadPantry();

    store.category.set('DRY_CANNED');

    expect(store.visibleItems().map((i) => i.product.name)).toEqual(['Rice']);
  });

  it('only offers chips for categories present in the view', async () => {
    await loadPantry();

    expect(store.availableCategories()).toEqual(['FRUIT_VEG', 'DRY_CANNED']);
  });

  it('cycles the status and persists it', async () => {
    await loadPantry();
    api.updateItem.mockReturnValue(of({ ...items[0], status: 'LOW' }));

    await store.cycleStatus(store.items()[0]);

    expect(api.updateItem).toHaveBeenCalledWith('familia', 'Tomatoes', { status: 'LOW' });
    expect(store.items()[0].status).toBe('LOW');
  });

  it('rolls the optimistic update back when the request fails', async () => {
    await loadPantry();
    api.updateItem.mockReturnValue(throwError(() => new Error('offline')));

    await store.cycleStatus(store.items()[0]);

    expect(store.items()[0].status).toBe('OK');
    expect(store.error()).not.toBeNull();
  });
});
