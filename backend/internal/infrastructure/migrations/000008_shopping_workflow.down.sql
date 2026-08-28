ALTER TABLE products DROP CONSTRAINT IF EXISTS products_default_status_check;
ALTER TABLE pantry_items DROP CONSTRAINT IF EXISTS pantry_items_status_check;

UPDATE products SET default_status = CASE default_status
    WHEN 'IN_CART' THEN 'OUT'
    WHEN 'PENDING' THEN 'LOW'
    WHEN 'DISCARDED' THEN 'OK'
    ELSE default_status
END;

UPDATE pantry_items SET status = CASE status
    WHEN 'IN_CART' THEN 'OUT'
    WHEN 'PENDING' THEN 'LOW'
    WHEN 'DISCARDED' THEN 'OK'
    ELSE status
END;

ALTER TABLE products ALTER COLUMN default_status SET DEFAULT 'OK';
ALTER TABLE products ADD CONSTRAINT products_default_status_check
    CHECK (default_status IN ('OUT', 'LOW', 'OK', 'ARCHIVED'));
ALTER TABLE pantry_items ADD CONSTRAINT pantry_items_status_check
    CHECK (status IN ('OUT', 'LOW', 'OK', 'ARCHIVED'));
