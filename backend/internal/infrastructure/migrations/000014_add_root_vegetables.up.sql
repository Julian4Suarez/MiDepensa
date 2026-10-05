-- Adds three requested standalone products to the active catalogue.
INSERT INTO products (
    id, code, name, image, default_category, default_type, default_status, sort_order
) VALUES
    (gen_random_uuid(), 'cassava', 'Yuca', 'cassava.svg', 'VEGETABLES', 'SECONDARY', 'PENDING', 189),
    (gen_random_uuid(), 'plantain', 'Plantains', 'plantain.svg', 'FRUIT', 'SECONDARY', 'PENDING', 190),
    (gen_random_uuid(), 'beetroot', 'Beetroot', 'beetroot.svg', 'VEGETABLES', 'SECONDARY', 'PENDING', 191)
ON CONFLICT (code) DO NOTHING;

-- Existing pantries receive the new products with their catalogue defaults.
INSERT INTO pantry_items (pantry_id, product_id, status, product_type, category)
SELECT pantry.id, product.id, product.default_status, product.default_type, product.default_category
FROM pantries AS pantry
CROSS JOIN products AS product
WHERE product.code IN ('cassava', 'plantain', 'beetroot')
ON CONFLICT (pantry_id, product_id) DO NOTHING;
