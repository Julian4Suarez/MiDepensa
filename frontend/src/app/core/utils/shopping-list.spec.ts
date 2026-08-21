import { buildShoppingList } from './shopping-list';
import type { Category, PantryItem, StockStatus } from '../../shared/models/pantry.model';

function item(name: string, status: StockStatus, category: Category): PantryItem {
  return {
    product: { id: name, code: name, name, image: `${name}.svg` },
    status,
    view: 'PRIMARY',
    category,
    updatedAt: '2026-01-01T00:00:00Z',
  };
}

describe('buildShoppingList', () => {
  const items = [
    item('Tomatoes', 'OUT', 'FRESH'),
    item('Milk', 'LOW', 'FRESH'),
    item('Rice', 'OK', 'PANTRY'),
    item('Olive oil', 'OUT', 'PANTRY'),
    item('Dish soap', 'LOW', 'HOME_CARE'),
  ];

  it('lists only out-of-stock products, grouped by category', () => {
    expect(buildShoppingList(items, false)).toBe(
      ['Fresh', '- Tomatoes', '', 'Pantry', '- Olive oil'].join('\n'),
    );
  });

  it('adds running-low products marked as low when requested', () => {
    expect(buildShoppingList(items, true)).toBe(
      [
        'Fresh',
        '- Tomatoes',
        '- Milk (low)',
        '',
        'Pantry',
        '- Olive oil',
        '',
        'Home & care',
        '- Dish soap (low)',
      ].join('\n'),
    );
  });

  it('returns an empty string when nothing needs buying', () => {
    expect(buildShoppingList([item('Rice', 'OK', 'PANTRY')], true)).toBe('');
  });
});
