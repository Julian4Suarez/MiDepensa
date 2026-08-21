ALTER TABLE products DROP CONSTRAINT IF EXISTS products_default_category_check;
ALTER TABLE pantry_items DROP CONSTRAINT IF EXISTS pantry_items_category_check;

UPDATE products SET default_category = 'FRESH'
WHERE default_category IN ('FRUIT_VEG', 'MEAT_FISH', 'DAIRY_EGGS')
   OR code = 'bread';

UPDATE products SET default_category = 'PANTRY' WHERE default_category = 'DRY_CANNED';

UPDATE pantry_items i
SET category = p.default_category
FROM products p
WHERE p.id = i.product_id;

ALTER TABLE products ADD CONSTRAINT products_default_category_check
    CHECK (default_category IN ('FRESH', 'PANTRY', 'DRINKS', 'HOME_CARE'));

ALTER TABLE pantry_items ADD CONSTRAINT pantry_items_category_check
    CHECK (category IN ('FRESH', 'PANTRY', 'DRINKS', 'HOME_CARE'));
