import { buildShoppingList } from './shopping-list';
import type { Category, ItemStatus, PantryItem } from '../../shared/models/pantry.model';

function item(name: string, status: ItemStatus, category: Category): PantryItem {
  return {
    product: { id: name, code: name, name, image: `${name}.svg`, variants: [] },
    status,
    type: 'ESSENTIAL',
    category,
    selectedVariantIds: [],
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

  it('lists selected variants instead of the general product', () => {
    const yogurt = item('Yogurt', 'IN_CART', 'DAIRY_EGGS');
    yogurt.product.variants = [
      { id: 'natural', code: 'natural', name: 'Natural yogurt', image: 'yogurt.svg' },
      { id: 'greek', code: 'greek', name: 'Greek yogurt', image: 'yogurt.svg' },
      { id: 'fruit', code: 'fruit', name: 'Fruit yogurt', image: 'yogurt.svg' },
    ];
    yogurt.selectedVariantIds = ['greek', 'fruit'];

    expect(buildShoppingList([yogurt])).toBe(
      ['Dairy & eggs', '- Greek yogurt', '- Fruit yogurt'].join('\n'),
    );
  });

  it('lists the general product when it has variants but none are selected', () => {
    const yogurt = item('Yogurt', 'IN_CART', 'DAIRY_EGGS');
    yogurt.product.variants = [
      { id: 'natural', code: 'natural', name: 'Natural yogurt', image: 'yogurt.svg' },
    ];

    expect(buildShoppingList([yogurt])).toBe(['Dairy & eggs', '- Yogurt'].join('\n'));
  });
});
