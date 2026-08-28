import { buildShoppingList } from './shopping-list';
import type { Category, ItemStatus, PantryItem } from '../../shared/models/pantry.model';

function item(name: string, status: ItemStatus, category: Category): PantryItem {
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
    item('Tomatoes', 'IN_CART', 'FRUIT_VEG'),
    item('Milk', 'PENDING', 'DAIRY_EGGS'),
    item('Rice', 'DISCARDED', 'DRY_CANNED'),
    item('Olive oil', 'IN_CART', 'DRY_CANNED'),
  ];

  it('lists only products in the cart, grouped by category', () => {
    expect(buildShoppingList(items)).toBe(
      ['Fruit & veg', '- Tomatoes', '', 'Dry & canned', '- Olive oil'].join('\n'),
    );
  });

  it('returns an empty string when the cart is empty', () => {
    expect(buildShoppingList([item('Rice', 'DISCARDED', 'DRY_CANNED')])).toBe('');
  });
});
