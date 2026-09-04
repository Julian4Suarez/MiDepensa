-- This rollback restores removed choices as independent products where
-- possible. Per-pantry non-cart child states cannot be reconstructed exactly.

CREATE TEMP TABLE rollback_variant_mapping (
    parent_code TEXT,
    child_code TEXT,
    child_name TEXT,
    variant_code TEXT,
    sort_order INTEGER
) ON COMMIT DROP;

INSERT INTO rollback_variant_mapping VALUES
    ('cooked_ham', 'sliced_turkey', 'Sliced turkey', 'sliced_turkey', 88),
    ('fish', 'canned_tuna', 'Canned tuna', 'tuna', 118),
    ('fish', 'canned_sardines', 'Canned sardines', 'sardines', 119),
    ('fish', 'prawns', 'Prawns', 'prawns', 92),
    ('leafy_greens', 'lettuce', 'Lettuce', 'lettuce', 58),
    ('leafy_greens', 'spinach', 'Spinach', 'spinach', 59),
    ('onion', 'leek', 'Leek', 'leek', 70),
    ('potato', 'sweet_potato', 'Sweet potato', 'sweet_potato', 67),
    ('cereal', 'oats', 'Oats', 'oats', 108),
    ('aluminium_foil', 'cling_film', 'Cling film', 'cling_film', 165),
    ('aluminium_foil', 'baking_paper', 'Baking paper', 'baking_paper', 164),
    ('laundry_detergent', 'fabric_softener', 'Fabric softener', 'fabric_softener', 151),
    ('laundry_detergent', 'stain_remover', 'Stain remover', 'stain_remover', 152),
    ('laundry_detergent', 'bleach', 'Bleach', 'bleach', 153),
    ('multipurpose_cleaner', 'disinfectant', 'Disinfectant', 'disinfectant', 154),
    ('multipurpose_cleaner', 'toilet_cleaner', 'Toilet cleaner', 'toilet_cleaner', 158),
    ('multipurpose_cleaner', 'floor_cleaner', 'Floor cleaner', 'floor_cleaner', 53),
    ('multipurpose_cleaner', 'glass_cleaner', 'Glass cleaner', 'glass_cleaner', 156),
    ('multipurpose_cleaner', 'degreaser', 'Degreaser', 'degreaser', 157),
    ('dish_soap', 'dishwasher_tablets', 'Dishwasher tablets', 'dishwasher_tablets', 159),
    ('dish_soap', 'dishwasher_salt', 'Dishwasher salt', 'dishwasher_salt', 160),
    ('sponges', 'dishcloths', 'Dishcloths', 'dishcloths', 161),
    ('sponges', 'rubber_gloves', 'Rubber gloves', 'rubber_gloves', 162),
    ('hand_soap', 'shower_gel', 'Shower gel', 'shower_gel', 170),
    ('hand_soap', 'shampoo', 'Shampoo', 'shampoo', 49),
    ('hand_soap', 'conditioner', 'Conditioner', 'conditioner', 171);

INSERT INTO products (id, code, name, image, default_category, default_type, default_status, sort_order)
SELECT gen_random_uuid(), mapping.child_code, mapping.child_name, variant.image,
       parent.default_category, parent.default_type, parent.default_status, mapping.sort_order
FROM rollback_variant_mapping AS mapping
JOIN products AS parent ON parent.code = mapping.parent_code
JOIN product_variants AS variant
  ON variant.product_id = parent.id AND variant.code = mapping.variant_code
ON CONFLICT (code) DO NOTHING;

INSERT INTO pantry_items (pantry_id, product_id, status, product_type, category)
SELECT pantry.id, child.id,
       CASE WHEN selected.variant_id IS NULL THEN child.default_status ELSE 'IN_CART' END,
       child.default_type, child.default_category
FROM pantries AS pantry
JOIN products AS child ON child.code IN (SELECT child_code FROM rollback_variant_mapping)
JOIN rollback_variant_mapping AS mapping ON mapping.child_code = child.code
JOIN products AS parent ON parent.code = mapping.parent_code
JOIN product_variants AS variant
  ON variant.product_id = parent.id AND variant.code = mapping.variant_code
LEFT JOIN pantry_item_variants AS selected
  ON selected.pantry_id = pantry.id AND selected.variant_id = variant.id
ON CONFLICT (pantry_id, product_id) DO NOTHING;

-- Restore products removed outright by the consolidation.
INSERT INTO products (id, code, name, image, default_category, default_type, default_status, sort_order) VALUES
    (gen_random_uuid(), 'canned_food', 'Canned food', 'canned_food.svg', 'COOKING_CONDIMENTS', 'SECONDARY', 'PENDING', 31),
    (gen_random_uuid(), 'canned_peas', 'Canned peas', 'canned_peas.svg', 'VEGETABLES', 'SECONDARY', 'ARCHIVED', 117),
    (gen_random_uuid(), 'frozen_vegetables', 'Frozen vegetables', 'frozen_vegetables.svg', 'VEGETABLES', 'SECONDARY', 'PENDING', 81),
    (gen_random_uuid(), 'margarine', 'Margarine', 'margarine.svg', 'DAIRY_EGGS', 'SECONDARY', 'ARCHIVED', 105),
    (gen_random_uuid(), 'egg_whites', 'Egg whites', 'egg_whites.svg', 'DAIRY_EGGS', 'SECONDARY', 'ARCHIVED', 106)
ON CONFLICT (code) DO NOTHING;

INSERT INTO pantry_items (pantry_id, product_id, status, product_type, category)
SELECT pantry.id, product.id, product.default_status, product.default_type, product.default_category
FROM pantries AS pantry
JOIN products AS product ON product.code IN (
    'canned_food', 'canned_peas', 'frozen_vegetables', 'margarine', 'egg_whites'
)
ON CONFLICT (pantry_id, product_id) DO NOTHING;

DELETE FROM product_variants
WHERE code IN (
    'whole_chicken', 'chicken_breast', 'chicken_thighs', 'chicken_wings',
    'beef_steaks', 'stewing_beef', 'pork_chops', 'pork_ribs', 'turkey_breast',
    'turkey_pieces', 'whole_rabbit', 'rabbit_pieces', 'minced_beef', 'minced_pork',
    'minced_chicken_turkey', 'minced_mixed', 'pork_sausages',
    'chicken_turkey_sausages', 'frankfurters', 'sardines', 'prawns',
    'mixed_greens', 'lettuce', 'spinach', 'yellow_onion', 'red_onion',
    'sweet_onion', 'leek', 'white_potato', 'red_potato', 'sweet_potato',
    'green_pepper', 'red_pepper', 'yellow_pepper', 'black_beans', 'kidney_beans',
    'salted_butter', 'unsalted_butter', 'ghee', 'corn_flakes', 'oats', 'muesli',
    'granola', 'aluminium_foil', 'cling_film', 'baking_paper', 'laundry_detergent',
    'fabric_softener', 'stain_remover', 'bleach', 'multipurpose_cleaner',
    'disinfectant', 'toilet_cleaner', 'floor_cleaner', 'glass_cleaner', 'degreaser',
    'dish_soap', 'dishwasher_tablets', 'dishwasher_salt', 'sponges', 'dishcloths',
    'rubber_gloves', 'hand_soap', 'shower_gel', 'shampoo', 'conditioner'
);

-- Restore fish variants removed or renamed by the consolidation.
INSERT INTO product_variants (product_id, code, name, image, sort_order)
SELECT product.id, 'frozen_fish', 'Frozen fish', 'frozen_fish.svg', 5
FROM products AS product WHERE product.code = 'fish'
ON CONFLICT (product_id, code) DO NOTHING;

UPDATE products SET name = CASE code
    WHEN 'red_meat' THEN 'Red meat'
    WHEN 'pork_chops' THEN 'Pork chops'
    WHEN 'cooked_ham' THEN 'Ham'
    WHEN 'aluminium_foil' THEN 'Aluminium foil'
    WHEN 'laundry_detergent' THEN 'Laundry detergent'
    WHEN 'multipurpose_cleaner' THEN 'Multipurpose cleaner'
    WHEN 'dish_soap' THEN 'Dish soap'
    WHEN 'sponges' THEN 'Sponges'
    WHEN 'hand_soap' THEN 'Hand soap'
    WHEN 'coriander' THEN 'Coriander'
    ELSE name
END;

UPDATE product_variants
SET code = 'tuna_steaks', name = 'Tuna steaks', sort_order = 4
WHERE code = 'tuna'
  AND product_id = (SELECT id FROM products WHERE code = 'fish');

UPDATE product_variants
SET name = 'Beans', sort_order = 1
WHERE code = 'beans'
  AND product_id = (SELECT id FROM products WHERE code = 'beans');

UPDATE products SET default_status = 'ARCHIVED', default_type = 'SECONDARY'
WHERE code IN ('pork_chops', 'sausages', 'bacon', 'avocado', 'coriander', 'honey', 'mustard');

UPDATE products SET default_status = 'PENDING'
WHERE code IN ('breadcrumbs', 'tomato_sauce', 'toothbrush_heads', 'cotton_pads', 'feminine_hygiene');

DELETE FROM products WHERE code IN ('rabbit', 'mango', 'mandarin', 'cherries', 'raspberries');
