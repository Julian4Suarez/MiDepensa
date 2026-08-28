-- Replaces the catch-all OTHERS type with an explicit archived state. Archived
-- catalog products remain secondary so restoring one puts it in a useful view.

ALTER TABLE products ADD COLUMN default_status TEXT NOT NULL DEFAULT 'OK';

ALTER TABLE products DROP CONSTRAINT IF EXISTS products_default_type_check;
ALTER TABLE pantry_items DROP CONSTRAINT IF EXISTS pantry_items_product_type_check;
ALTER TABLE pantry_items DROP CONSTRAINT IF EXISTS pantry_items_status_check;

UPDATE products
SET default_type = 'SECONDARY', default_status = 'ARCHIVED'
WHERE default_type = 'OTHERS';

UPDATE pantry_items
SET product_type = 'SECONDARY', status = 'ARCHIVED'
WHERE product_type = 'OTHERS';

ALTER TABLE products ADD CONSTRAINT products_default_type_check
    CHECK (default_type IN ('ESSENTIAL', 'SECONDARY'));
ALTER TABLE products ADD CONSTRAINT products_default_status_check
    CHECK (default_status IN ('OUT', 'LOW', 'OK', 'ARCHIVED'));
ALTER TABLE pantry_items ADD CONSTRAINT pantry_items_product_type_check
    CHECK (product_type IN ('ESSENTIAL', 'SECONDARY'));
ALTER TABLE pantry_items ADD CONSTRAINT pantry_items_status_check
    CHECK (status IN ('OUT', 'LOW', 'OK', 'ARCHIVED'));
