ALTER TABLE products DROP CONSTRAINT IF EXISTS products_default_type_check;
ALTER TABLE products DROP CONSTRAINT IF EXISTS products_default_status_check;
ALTER TABLE pantry_items DROP CONSTRAINT IF EXISTS pantry_items_product_type_check;
ALTER TABLE pantry_items DROP CONSTRAINT IF EXISTS pantry_items_status_check;

UPDATE products SET default_type = 'OTHERS' WHERE default_status = 'ARCHIVED';
UPDATE pantry_items SET product_type = 'OTHERS', status = 'OK' WHERE status = 'ARCHIVED';

ALTER TABLE products DROP COLUMN default_status;

ALTER TABLE products ADD CONSTRAINT products_default_type_check
    CHECK (default_type IN ('ESSENTIAL', 'SECONDARY', 'OTHERS'));
ALTER TABLE pantry_items ADD CONSTRAINT pantry_items_product_type_check
    CHECK (product_type IN ('ESSENTIAL', 'SECONDARY', 'OTHERS'));
ALTER TABLE pantry_items ADD CONSTRAINT pantry_items_status_check
    CHECK (status IN ('OUT', 'LOW', 'OK'));
