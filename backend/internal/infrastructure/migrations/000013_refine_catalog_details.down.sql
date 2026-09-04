CREATE TEMP TABLE paper_variant_mapping (
    child_code TEXT,
    child_name TEXT,
    variant_code TEXT,
    sort_order INTEGER,
    default_type TEXT
) ON COMMIT DROP;

INSERT INTO paper_variant_mapping VALUES
    ('napkins', 'Napkins', 'napkins', 167, 'SECONDARY'),
    ('tissues', 'Tissues', 'tissues', 168, 'SECONDARY'),
    ('toilet_paper', 'Toilet paper', 'toilet_paper', 54, 'ESSENTIAL');

INSERT INTO products (id, code, name, image, default_category, default_type, default_status, sort_order)
SELECT gen_random_uuid(), mapping.child_code, mapping.child_name, variant.image,
       'HOUSEHOLD_CLEANING', mapping.default_type, 'PENDING', mapping.sort_order
FROM paper_variant_mapping AS mapping
JOIN products AS parent ON parent.code = 'paper_towels'
JOIN product_variants AS variant
  ON variant.product_id = parent.id AND variant.code = mapping.variant_code
ON CONFLICT (code) DO NOTHING;

INSERT INTO pantry_items (pantry_id, product_id, status, product_type, category)
SELECT pantry.id, child.id,
       CASE WHEN selected.variant_id IS NULL THEN child.default_status ELSE 'IN_CART' END,
       child.default_type, child.default_category
FROM pantries AS pantry
JOIN products AS child ON child.code IN (SELECT child_code FROM paper_variant_mapping)
JOIN paper_variant_mapping AS mapping ON mapping.child_code = child.code
JOIN products AS parent ON parent.code = 'paper_towels'
JOIN product_variants AS variant
  ON variant.product_id = parent.id AND variant.code = mapping.variant_code
LEFT JOIN pantry_item_variants AS selected
  ON selected.pantry_id = pantry.id AND selected.variant_id = variant.id
ON CONFLICT (pantry_id, product_id) DO NOTHING;

DELETE FROM product_variants
WHERE code IN (
    'kitchen_roll', 'napkins', 'tissues', 'toilet_paper',
    'cherry_tomatoes', 'roma_tomatoes', 'beefsteak_tomatoes', 'vine_tomatoes',
    'black_pepper', 'paprika', 'oregano', 'cumin', 'cinnamon', 'curry_powder'
);

UPDATE products SET name = 'Kitchen roll' WHERE code = 'paper_towels';

UPDATE products SET default_type = 'SECONDARY'
WHERE code IN ('orange', 'lemon');

UPDATE products SET default_status = 'ARCHIVED', default_type = 'SECONDARY'
WHERE code IN ('grapes', 'blueberries', 'peaches', 'watermelon', 'cucumber', 'spices');

UPDATE products SET default_status = 'PENDING'
WHERE code IN ('cherries', 'raspberries', 'green_beans', 'bacon', 'turkey', 'tofu', 'mouthwash');
