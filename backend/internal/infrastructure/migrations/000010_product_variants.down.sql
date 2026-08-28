-- Restore variant rows as independent products where their code is available.
INSERT INTO products (id, code, name, image, default_category, default_type, default_status, sort_order)
SELECT gen_random_uuid(), variant.code, variant.name, variant.image,
       parent.default_category, parent.default_type, parent.default_status,
       1000 + row_number() OVER (ORDER BY parent.sort_order, variant.sort_order)
FROM product_variants AS variant
JOIN products AS parent ON parent.id = variant.product_id
WHERE NOT EXISTS (SELECT 1 FROM products existing WHERE existing.code = variant.code)
  AND variant.code NOT IN (
    'regular_cheese', 'whole_milk', 'white_fish', 'still_water', 'regular_coffee',
    'black_tea', 'white_sugar', 'white_bread', 'wholemeal_bread', 'sliced_bread',
    'baguette', 'semi_skimmed_milk', 'skimmed_milk', 'lactose_free_milk'
  );

INSERT INTO pantry_items (pantry_id, product_id, status, product_type, category)
SELECT pantry.id, product.id,
       CASE WHEN selected.variant_id IS NULL THEN 'PENDING' ELSE 'IN_CART' END,
       product.default_type, product.default_category
FROM pantries AS pantry
CROSS JOIN products AS product
JOIN product_variants AS variant ON variant.code = product.code
LEFT JOIN pantry_item_variants AS selected
  ON selected.pantry_id = pantry.id AND selected.variant_id = variant.id
ON CONFLICT (pantry_id, product_id) DO NOTHING;

UPDATE products SET name = CASE code
    WHEN 'natural_yogurt' THEN 'Natural yogurt'
    WHEN 'olive_oil' THEN 'Olive oil'
    WHEN 'tomato_sauce' THEN 'Tomato sauce'
    WHEN 'cooked_ham' THEN 'Cooked ham'
    WHEN 'beans' THEN 'Beans'
    ELSE name
END;

DROP TABLE pantry_item_variants;
DROP TABLE product_variants;
