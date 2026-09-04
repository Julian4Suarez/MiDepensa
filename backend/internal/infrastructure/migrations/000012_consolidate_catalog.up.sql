-- Consolidates products that are normally chosen as variations of the same
-- shopping need. Categories remain unchanged: these are products and variants.

-- New general products and fruit requested for the catalog.
INSERT INTO products (
    id, code, name, image, default_category, default_type, default_status, sort_order
) VALUES
    (gen_random_uuid(), 'rabbit', 'Rabbit', 'rabbit.svg', 'MEAT_FISH', 'SECONDARY', 'PENDING', 184),
    (gen_random_uuid(), 'mango', 'Mangoes', 'mango.svg', 'FRUIT', 'SECONDARY', 'PENDING', 185),
    (gen_random_uuid(), 'mandarin', 'Mandarins', 'mandarin.svg', 'FRUIT', 'SECONDARY', 'PENDING', 186),
    (gen_random_uuid(), 'cherries', 'Cherries', 'cherries.svg', 'FRUIT', 'SECONDARY', 'PENDING', 187),
    (gen_random_uuid(), 'raspberries', 'Raspberries', 'raspberries.svg', 'FRUIT', 'SECONDARY', 'PENDING', 188)
ON CONFLICT (code) DO NOTHING;

INSERT INTO pantry_items (pantry_id, product_id, status, product_type, category)
SELECT pantry.id, product.id, product.default_status, product.default_type, product.default_category
FROM pantries AS pantry
CROSS JOIN products AS product
WHERE product.code IN ('rabbit', 'mango', 'mandarin', 'cherries', 'raspberries')
ON CONFLICT (pantry_id, product_id) DO NOTHING;

-- Each row maps an old independent product to a variant of a general product.
-- Parent/self rows create the obvious first choice but are not auto-selected.
CREATE TEMP TABLE catalog_variant_mapping (
    parent_code TEXT,
    child_code TEXT,
    variant_code TEXT,
    variant_name TEXT,
    variant_order INTEGER
) ON COMMIT DROP;

INSERT INTO catalog_variant_mapping VALUES
    ('chicken', 'chicken', 'whole_chicken', 'Whole chicken', 1),
    ('pork_chops', 'pork_chops', 'pork_chops', 'Pork chops', 1),
    ('turkey', 'turkey', 'turkey_breast', 'Turkey breast', 1),
    ('rabbit', 'rabbit', 'whole_rabbit', 'Whole rabbit', 1),
    ('minced_meat', 'minced_meat', 'minced_beef', 'Minced beef', 1),
    ('sausages', 'sausages', 'pork_sausages', 'Pork sausages', 1),
    ('cooked_ham', 'cooked_ham', 'cooked_ham', 'Cooked ham', 1),
    ('cooked_ham', 'sliced_turkey', 'sliced_turkey', 'Sliced turkey', 3),
    ('fish', 'fish', 'white_fish', 'White fish', 1),
    ('fish', 'canned_tuna', 'tuna', 'Tuna', 4),
    ('fish', 'canned_sardines', 'sardines', 'Sardines', 5),
    ('fish', 'prawns', 'prawns', 'Prawns', 6),
    ('leafy_greens', 'leafy_greens', 'mixed_greens', 'Mixed greens', 1),
    ('leafy_greens', 'lettuce', 'lettuce', 'Lettuce', 2),
    ('leafy_greens', 'spinach', 'spinach', 'Spinach', 3),
    ('onion', 'onion', 'yellow_onion', 'Yellow onion', 1),
    ('onion', 'leek', 'leek', 'Leek', 4),
    ('potato', 'potato', 'white_potato', 'White potato', 1),
    ('potato', 'sweet_potato', 'sweet_potato', 'Sweet potato', 3),
    ('bell_pepper', 'bell_pepper', 'green_pepper', 'Green pepper', 1),
    ('cereal', 'cereal', 'corn_flakes', 'Corn flakes', 1),
    ('cereal', 'oats', 'oats', 'Oats', 2),
    ('aluminium_foil', 'aluminium_foil', 'aluminium_foil', 'Aluminium foil', 1),
    ('aluminium_foil', 'cling_film', 'cling_film', 'Cling film', 2),
    ('aluminium_foil', 'baking_paper', 'baking_paper', 'Baking paper', 3),
    ('laundry_detergent', 'laundry_detergent', 'laundry_detergent', 'Laundry detergent', 1),
    ('laundry_detergent', 'fabric_softener', 'fabric_softener', 'Fabric softener', 2),
    ('laundry_detergent', 'stain_remover', 'stain_remover', 'Stain remover', 3),
    ('laundry_detergent', 'bleach', 'bleach', 'Bleach', 4),
    ('multipurpose_cleaner', 'multipurpose_cleaner', 'multipurpose_cleaner', 'Multipurpose cleaner', 1),
    ('multipurpose_cleaner', 'disinfectant', 'disinfectant', 'Disinfectant', 2),
    ('multipurpose_cleaner', 'toilet_cleaner', 'toilet_cleaner', 'Toilet cleaner', 3),
    ('multipurpose_cleaner', 'floor_cleaner', 'floor_cleaner', 'Floor cleaner', 4),
    ('multipurpose_cleaner', 'glass_cleaner', 'glass_cleaner', 'Glass cleaner', 5),
    ('multipurpose_cleaner', 'degreaser', 'degreaser', 'Degreaser', 6),
    ('dish_soap', 'dish_soap', 'dish_soap', 'Dish soap', 1),
    ('dish_soap', 'dishwasher_tablets', 'dishwasher_tablets', 'Dishwasher tablets', 2),
    ('dish_soap', 'dishwasher_salt', 'dishwasher_salt', 'Dishwasher salt', 3),
    ('sponges', 'sponges', 'sponges', 'Sponges', 1),
    ('sponges', 'dishcloths', 'dishcloths', 'Dishcloths', 2),
    ('sponges', 'rubber_gloves', 'rubber_gloves', 'Rubber gloves', 3),
    ('hand_soap', 'hand_soap', 'hand_soap', 'Hand soap', 1),
    ('hand_soap', 'shower_gel', 'shower_gel', 'Shower gel', 2),
    ('hand_soap', 'shampoo', 'shampoo', 'Shampoo', 3),
    ('hand_soap', 'conditioner', 'conditioner', 'Conditioner', 4);

-- Tuna already existed as a fish variant under a preparation-specific name.
UPDATE product_variants
SET code = 'tuna', name = 'Tuna', sort_order = 4
WHERE code = 'tuna_steaks'
  AND product_id = (SELECT id FROM products WHERE code = 'fish');

-- Frozen is a storage format, not a fish subtype.
DELETE FROM product_variants
WHERE code = 'frozen_fish'
  AND product_id = (SELECT id FROM products WHERE code = 'fish');

INSERT INTO product_variants (product_id, code, name, image, sort_order)
SELECT parent.id, mapping.variant_code, mapping.variant_name, child.image, mapping.variant_order
FROM catalog_variant_mapping AS mapping
JOIN products AS parent ON parent.code = mapping.parent_code
JOIN products AS child ON child.code = mapping.child_code
ON CONFLICT (product_id, code) DO UPDATE
SET name = EXCLUDED.name, sort_order = EXCLUDED.sort_order;

-- Variants without a previous independent catalog product reuse the parent's
-- illustration, keeping visually related choices together.
INSERT INTO product_variants (product_id, code, name, image, sort_order)
SELECT parent.id, variant.code, variant.name, parent.image, variant.sort_order
FROM products AS parent
JOIN (VALUES
    ('chicken', 'chicken_breast', 'Chicken breast', 2),
    ('chicken', 'chicken_thighs', 'Chicken thighs', 3),
    ('chicken', 'chicken_wings', 'Chicken wings', 4),
    ('red_meat', 'beef_steaks', 'Beef steaks', 1),
    ('red_meat', 'stewing_beef', 'Stewing beef', 2),
    ('pork_chops', 'pork_ribs', 'Pork ribs', 2),
    ('turkey', 'turkey_pieces', 'Turkey pieces', 2),
    ('rabbit', 'rabbit_pieces', 'Rabbit pieces', 2),
    ('minced_meat', 'minced_pork', 'Minced pork', 2),
    ('minced_meat', 'minced_chicken_turkey', 'Minced chicken or turkey', 3),
    ('minced_meat', 'minced_mixed', 'Mixed minced meat', 4),
    ('sausages', 'chicken_turkey_sausages', 'Chicken or turkey sausages', 2),
    ('sausages', 'frankfurters', 'Frankfurters', 3),
    ('onion', 'red_onion', 'Red onion', 2),
    ('onion', 'sweet_onion', 'Sweet onion', 3),
    ('potato', 'red_potato', 'Red potato', 2),
    ('bell_pepper', 'red_pepper', 'Red pepper', 2),
    ('bell_pepper', 'yellow_pepper', 'Yellow pepper', 3),
    ('beans', 'black_beans', 'Black beans', 4),
    ('beans', 'kidney_beans', 'Kidney beans', 5),
    ('butter', 'salted_butter', 'Salted butter', 1),
    ('butter', 'unsalted_butter', 'Unsalted butter', 2),
    ('butter', 'ghee', 'Ghee', 3),
    ('cereal', 'muesli', 'Muesli', 3),
    ('cereal', 'granola', 'Granola', 4)
) AS variant(parent_code, code, name, sort_order) ON variant.parent_code = parent.code
ON CONFLICT (product_id, code) DO UPDATE
SET name = EXCLUDED.name, sort_order = EXCLUDED.sort_order;

-- Beans is now the neutral white-bean choice within Legumes.
UPDATE product_variants
SET name = 'White beans', sort_order = 1
WHERE code = 'beans'
  AND product_id = (SELECT id FROM products WHERE code = 'beans');

-- A child already in the cart becomes an explicit selection on its parent.
-- A generic parent that was in the cart remains generic (no guessed subtype).
INSERT INTO pantry_item_variants (pantry_id, product_id, variant_id)
SELECT item.pantry_id, parent.id, variant.id
FROM catalog_variant_mapping AS mapping
JOIN products AS parent ON parent.code = mapping.parent_code
JOIN products AS child ON child.code = mapping.child_code
JOIN pantry_items AS item ON item.product_id = child.id AND item.status = 'IN_CART'
JOIN product_variants AS variant
  ON variant.product_id = parent.id AND variant.code = mapping.variant_code
WHERE mapping.child_code <> mapping.parent_code
ON CONFLICT DO NOTHING;

-- The general product inherits the most relevant state among all merged rows.
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
    FROM catalog_variant_mapping AS mapping
    JOIN products AS parent ON parent.code = mapping.parent_code
    JOIN products AS child ON child.code = mapping.child_code
    JOIN pantry_items AS item ON item.product_id = child.id
    GROUP BY item.pantry_id, parent.id
) AS grouped
WHERE parent_item.pantry_id = grouped.pantry_id
  AND parent_item.product_id = grouped.product_id;

DELETE FROM pantry_items AS item
USING products AS child, catalog_variant_mapping AS mapping
WHERE item.product_id = child.id
  AND child.code = mapping.child_code
  AND mapping.child_code <> mapping.parent_code;

DELETE FROM products AS child
USING catalog_variant_mapping AS mapping
WHERE child.code = mapping.child_code
  AND mapping.child_code <> mapping.parent_code;

-- Canned peas are the same shopping need as peas, so preserve the strongest
-- state without exposing packaging as a separate choice.
UPDATE pantry_items AS peas_item
SET status = merged.status, updated_at = now()
FROM (
    SELECT item.pantry_id,
           CASE
             WHEN bool_or(item.status = 'IN_CART') THEN 'IN_CART'
             WHEN bool_or(item.status = 'PENDING') THEN 'PENDING'
             WHEN bool_or(item.status = 'DISCARDED') THEN 'DISCARDED'
             ELSE 'ARCHIVED'
           END AS status
    FROM pantry_items AS item
    JOIN products AS product ON product.id = item.product_id
    WHERE product.code IN ('peas', 'canned_peas')
    GROUP BY item.pantry_id
) AS merged
WHERE peas_item.pantry_id = merged.pantry_id
  AND peas_item.product_id = (SELECT id FROM products WHERE code = 'peas');

-- Remove generic, packaging-based or deliberately redundant catalog entries.
DELETE FROM products WHERE code IN (
    'canned_food', 'canned_peas', 'frozen_vegetables', 'margarine', 'egg_whites'
);

UPDATE products SET name = CASE code
    WHEN 'red_meat' THEN 'Beef'
    WHEN 'pork_chops' THEN 'Pork'
    WHEN 'cooked_ham' THEN 'Cold cuts'
    WHEN 'aluminium_foil' THEN 'Food wrapping'
    WHEN 'laundry_detergent' THEN 'Laundry care'
    WHEN 'multipurpose_cleaner' THEN 'Household cleaners'
    WHEN 'dish_soap' THEN 'Dishwashing'
    WHEN 'sponges' THEN 'Cleaning tools'
    WHEN 'hand_soap' THEN 'Bath & hair care'
    WHEN 'coriander' THEN 'Cilantro'
    ELSE name
END
WHERE code IN (
    'red_meat', 'pork_chops', 'cooked_ham', 'aluminium_foil',
    'laundry_detergent', 'multipurpose_cleaner', 'dish_soap', 'sponges', 'hand_soap',
    'coriander'
);

-- Explicit catalogue defaults requested by the product design. Archived rows
-- are reactivated only for products the user asked to bring back.
UPDATE products
SET default_type = CASE WHEN code = 'avocado' THEN 'ESSENTIAL' ELSE 'SECONDARY' END,
    default_status = 'PENDING'
WHERE code IN ('pork_chops', 'sausages', 'bacon', 'avocado', 'coriander', 'honey', 'mustard');

UPDATE pantry_items AS item
SET product_type = CASE WHEN product.code = 'avocado' THEN 'ESSENTIAL' ELSE 'SECONDARY' END,
    status = CASE WHEN item.status = 'ARCHIVED' THEN 'PENDING' ELSE item.status END,
    updated_at = now()
FROM products AS product
WHERE item.product_id = product.id
  AND product.code IN ('pork_chops', 'sausages', 'bacon', 'avocado', 'coriander', 'honey', 'mustard');

UPDATE products SET default_status = 'ARCHIVED'
WHERE code IN ('breadcrumbs', 'tomato_sauce', 'toothbrush_heads', 'cotton_pads', 'feminine_hygiene');

UPDATE pantry_items AS item
SET status = 'ARCHIVED', updated_at = now()
FROM products AS product
WHERE item.product_id = product.id
  AND product.code IN ('breadcrumbs', 'tomato_sauce', 'toothbrush_heads', 'cotton_pads', 'feminine_hygiene');
