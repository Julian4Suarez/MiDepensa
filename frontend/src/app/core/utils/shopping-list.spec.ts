import { buildShoppingList } from './shopping-list';
import type { Category, PantryItem, StockStatus } from '../../shared/models/pantry.model';

function item(name: string, status: StockStatus, category: Category): PantryItem {
  return {
    product: { id: name, code: name, name, image: `${name}.svg` },
    status,
    type: 'ESSENTIAL',
    category,
    updatedAt: '2026-01-01T00:00:00Z',
  };
}

describe('buildShoppingList', () => {
  const items = [
    item('Tomatoes', 'OUT', 'FRUIT_VEG'),
    item('Milk', 'LOW', 'DAIRY_EGGS'),
    item('Rice', 'OK', 'DRY_CANNED'),
    item('Olive oil', 'OUT', 'DRY_CANNED'),
    item('Dish soap', 'LOW', 'HOME_CARE'),
  ];

  it('lists only out-of-stock products, grouped by category', () => {
    expect(buildShoppingList(items, false)).toBe(
      ['Fruit & veg', '- Tomatoes', '', 'Dry & canned', '- Olive oil'].join('\n'),
    );
  });

  it('adds running-low products marked as low when requested', () => {
    expect(buildShoppingList(items, true)).toBe(
      [
        'Fruit & veg',
        '- Tomatoes',
        '',
        'Dairy & eggs',
        '- Milk (low)',
        '',
        'Dry & canned',
        '- Olive oil',
        '',
        'Home & care',
        '- Dish soap (low)',
      ].join('\n'),
    );
  });

  it('returns an empty string when nothing needs buying', () => {
    expect(buildShoppingList([item('Rice', 'OK', 'DRY_CANNED')], true)).toBe('');
  });
});
