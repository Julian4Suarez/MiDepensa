-- Applies refinements requested after migration 000012 had already shipped.
-- This migration is intentionally safe both after the original 000012 and on
-- a clean installation where some rows may already have been consolidated.

CREATE TEMP TABLE paper_variant_mapping (
    child_code TEXT,
    variant_code TEXT,
    variant_name TEXT,
    variant_order INTEGER
) ON COMMIT DROP;

INSERT INTO paper_variant_mapping VALUES
    ('paper_towels', 'kitchen_roll', 'Kitchen roll', 1),
    ('napkins', 'napkins', 'Napkins', 2),
    ('tissues', 'tissues', 'Tissues', 3),
    ('toilet_paper', 'toilet_paper', 'Toilet paper', 4);

INSERT INTO product_variants (product_id, code, name, image, sort_order)
SELECT parent.id, mapping.variant_code, mapping.variant_name, child.image, mapping.variant_order
FROM products AS parent
JOIN paper_variant_mapping AS mapping ON true
JOIN products AS child ON child.code = mapping.child_code
WHERE parent.code = 'paper_towels'
ON CONFLICT (product_id, code) DO UPDATE
SET name = EXCLUDED.name, sort_order = EXCLUDED.sort_order;

-- Preserve concrete paper products that were already in the cart.
INSERT INTO pantry_item_variants (pantry_id, product_id, variant_id)
SELECT item.pantry_id, parent.id, variant.id
FROM paper_variant_mapping AS mapping
JOIN products AS parent ON parent.code = 'paper_towels'
JOIN products AS child ON child.code = mapping.child_code
JOIN pantry_items AS item ON item.product_id = child.id AND item.status = 'IN_CART'
JOIN product_variants AS variant
  ON variant.product_id = parent.id AND variant.code = mapping.variant_code
WHERE mapping.child_code <> 'paper_towels'
ON CONFLICT DO NOTHING;

UPDATE pantry_items AS parent_item
SET status = grouped.status, updated_at = now()
FROM (
    SELECT item.pantry_id, parent.id AS product_id,
           CASE
             WHEN bool_or(item.status = 'IN_CART') THEN 'IN_CART'
             WHEN bool_or(item.status = 'PENDING') THEN 'PENDING'
             WHEN bool_or(item.status = 'DISCARDED') THEN 'DISCARDED'
             ELSE 'ARCHIVED'
           END AS status
    FROM paper_variant_mapping AS mapping
    JOIN products AS parent ON parent.code = 'paper_towels'
    JOIN products AS child ON child.code = mapping.child_code
    JOIN pantry_items AS item ON item.product_id = child.id
    GROUP BY item.pantry_id, parent.id
) AS grouped
WHERE parent_item.pantry_id = grouped.pantry_id
  AND parent_item.product_id = grouped.product_id;

DELETE FROM products WHERE code IN ('napkins', 'tissues', 'toilet_paper');
UPDATE products SET name = 'Paper products' WHERE code = 'paper_towels';

-- New choices that never existed as independent products use their parent's
-- illustration and therefore do not need data migration.
INSERT INTO product_variants (product_id, code, name, image, sort_order)
SELECT parent.id, variant.code, variant.name, parent.image, variant.sort_order
FROM products AS parent
JOIN (VALUES
    ('tomato', 'cherry_tomatoes', 'Cherry tomatoes', 1),
    ('tomato', 'roma_tomatoes', 'Roma tomatoes', 2),
    ('tomato', 'beefsteak_tomatoes', 'Beefsteak tomatoes', 3),
    ('tomato', 'vine_tomatoes', 'Vine tomatoes', 4),
    ('spices', 'black_pepper', 'Black pepper', 1),
    ('spices', 'paprika', 'Paprika', 2),
    ('spices', 'oregano', 'Oregano', 3),
    ('spices', 'cumin', 'Cumin', 4),
    ('spices', 'cinnamon', 'Cinnamon', 5),
    ('spices', 'curry_powder', 'Curry powder', 6)
) AS variant(parent_code, code, name, sort_order) ON variant.parent_code = parent.code
ON CONFLICT (product_id, code) DO UPDATE
SET name = EXCLUDED.name, sort_order = EXCLUDED.sort_order;

-- Active refinements. Previously archived items are reactivated as pending,
-- while an existing active decision is preserved.
UPDATE products
SET default_type = CASE WHEN code IN ('orange', 'lemon') THEN 'ESSENTIAL' ELSE 'SECONDARY' END,
    default_status = 'PENDING'
WHERE code IN (
    'orange', 'lemon', 'grapes', 'blueberries', 'peaches', 'watermelon',
    'cucumber', 'spices'
);

UPDATE pantry_items AS item
SET product_type = CASE WHEN product.code IN ('orange', 'lemon') THEN 'ESSENTIAL' ELSE 'SECONDARY' END,
    status = CASE WHEN item.status = 'ARCHIVED' THEN 'PENDING' ELSE item.status END,
    updated_at = now()
FROM products AS product
WHERE item.product_id = product.id
  AND product.code IN (
      'orange', 'lemon', 'grapes', 'blueberries', 'peaches', 'watermelon',
      'cucumber', 'spices'
  );

-- Products intentionally hidden from the normal workflow.
UPDATE products SET default_status = 'ARCHIVED'
WHERE code IN (
    'cherries', 'raspberries', 'green_beans', 'bacon', 'turkey', 'tofu', 'mouthwash'
);

UPDATE pantry_items AS item
SET status = 'ARCHIVED', updated_at = now()
FROM products AS product
WHERE item.product_id = product.id
  AND product.code IN (
      'cherries', 'raspberries', 'green_beans', 'bacon', 'turkey', 'tofu', 'mouthwash'
  );
