import { sortItems } from './sort-items';
import type { PantryItem, StockStatus } from '../../shared/models/pantry.model';

function item(name: string, status: StockStatus): PantryItem {
  return {
    product: { id: name, code: name, name, image: `${name}.svg` },
    status,
    view: 'PRIMARY',
    category: 'FRUIT_VEG',
    updatedAt: '2026-01-01T00:00:00Z',
  };
}

describe('sortItems', () => {
  const items = [item('Tomatoes', 'OK'), item('Apples', 'OUT'), item('Milk', 'LOW')];

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

  it('sorts by stock level, most urgent first', () => {
    expect(sortItems(items, 'STATUS').map((i) => i.status)).toEqual(['OUT', 'LOW', 'OK']);
  });

  it('keeps the catalog order within the same stock level', () => {
    const sameStatus = [item('Zucchini', 'OUT'), item('Apples', 'OUT')];

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
