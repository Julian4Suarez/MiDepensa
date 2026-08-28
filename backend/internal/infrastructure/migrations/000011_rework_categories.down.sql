ALTER TABLE products DROP CONSTRAINT IF EXISTS products_default_category_check;
ALTER TABLE pantry_items DROP CONSTRAINT IF EXISTS pantry_items_category_check;

UPDATE products SET default_category = CASE default_category
    WHEN 'FRUIT' THEN 'FRUIT_VEG'
    WHEN 'VEGETABLES' THEN 'FRUIT_VEG'
    WHEN 'BAKERY_BREAKFAST_SNACKS' THEN 'DRY_CANNED'
    WHEN 'PASTA_RICE_LEGUMES' THEN 'DRY_CANNED'
    WHEN 'COOKING_CONDIMENTS' THEN 'DRY_CANNED'
    WHEN 'HOUSEHOLD_CLEANING' THEN 'HOME_CARE'
    WHEN 'PERSONAL_CARE' THEN 'HOME_CARE'
    ELSE default_category
END;

UPDATE pantry_items AS item
SET category = product.default_category
FROM products AS product
WHERE product.id = item.product_id;

ALTER TABLE products ADD CONSTRAINT products_default_category_check
    CHECK (default_category IN ('FRUIT_VEG', 'MEAT_FISH', 'DAIRY_EGGS', 'DRY_CANNED', 'DRINKS', 'HOME_CARE'));
ALTER TABLE pantry_items ADD CONSTRAINT pantry_items_category_check
    CHECK (category IN ('FRUIT_VEG', 'MEAT_FISH', 'DAIRY_EGGS', 'DRY_CANNED', 'DRINKS', 'HOME_CARE'));
