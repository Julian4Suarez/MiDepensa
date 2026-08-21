ALTER TABLE products DROP CONSTRAINT IF EXISTS products_default_type_check;
ALTER TABLE pantry_items DROP CONSTRAINT IF EXISTS pantry_items_product_type_check;

UPDATE products SET default_type = CASE default_type
    WHEN 'ESSENTIAL' THEN 'PRIMARY'
    WHEN 'OTHERS' THEN 'OTHER'
    ELSE default_type
END;

UPDATE pantry_items SET product_type = CASE product_type
    WHEN 'ESSENTIAL' THEN 'PRIMARY'
    WHEN 'OTHERS' THEN 'OTHER'
    ELSE product_type
END;

ALTER TABLE products RENAME COLUMN default_type TO default_view;
ALTER TABLE pantry_items RENAME COLUMN product_type TO pantry_view;

ALTER TABLE products ADD CONSTRAINT products_default_view_check
    CHECK (default_view IN ('PRIMARY', 'SECONDARY', 'OTHER'));

ALTER TABLE pantry_items ADD CONSTRAINT pantry_items_pantry_view_check
    CHECK (pantry_view IN ('PRIMARY', 'SECONDARY', 'OTHER'));
