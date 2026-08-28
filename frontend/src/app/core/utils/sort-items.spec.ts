import { sortItems } from './sort-items';
import type { ItemStatus, PantryItem } from '../../shared/models/pantry.model';

function item(name: string, status: ItemStatus): PantryItem {
  return {
    product: { id: name, code: name, name, image: `${name}.svg` },
    status,
    type: 'ESSENTIAL',
    category: 'FRUIT_VEG',
    updatedAt: '2026-01-01T00:00:00Z',
  };
}

describe('sortItems', () => {
  const items = [
    item('Tomatoes', 'DISCARDED'),
    item('Apples', 'IN_CART'),
    item('Milk', 'PENDING'),
  ];

  it('keeps the catalog order by default', () => {
    expect(sortItems(items, 'DEFAULT')).toBe(items);
  });

  it('sorts alphabetically by product name', () => {
    expect(sortItems(items, 'NAME').map((i) => i.product.name)).toEqual([
      'Apples',
      'Milk',
      'Tomatoes',
    ]);
  });

  it('sorts by shopping status, pending decisions first', () => {
    expect(sortItems(items, 'STATUS').map((i) => i.status)).toEqual([
      'PENDING',
      'IN_CART',
      'DISCARDED',
    ]);
  });

  it('keeps the catalog order within the same shopping status', () => {
    const sameStatus = [item('Zucchini', 'IN_CART'), item('Apples', 'IN_CART')];

    expect(sortItems(sameStatus, 'STATUS').map((i) => i.product.name)).toEqual([
      'Zucchini',
      'Apples',
    ]);
  });

  it('does not mutate the input', () => {
    const original = [...items];

    sortItems(items, 'NAME');

    expect(items).toEqual(original);
  });
});
