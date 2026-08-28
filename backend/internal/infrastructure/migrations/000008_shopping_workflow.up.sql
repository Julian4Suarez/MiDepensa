-- Replaces inventory levels with explicit shopping workflow states while
-- preserving the intent of existing pantry data.

ALTER TABLE products DROP CONSTRAINT IF EXISTS products_default_status_check;
ALTER TABLE pantry_items DROP CONSTRAINT IF EXISTS pantry_items_status_check;

UPDATE products SET default_status = CASE default_status
    WHEN 'OUT' THEN 'IN_CART'
    WHEN 'LOW' THEN 'PENDING'
    WHEN 'OK' THEN 'DISCARDED'
    ELSE default_status
END;

UPDATE pantry_items SET status = CASE status
    WHEN 'OUT' THEN 'IN_CART'
    WHEN 'LOW' THEN 'PENDING'
    WHEN 'OK' THEN 'DISCARDED'
    ELSE status
END;

ALTER TABLE products ALTER COLUMN default_status SET DEFAULT 'DISCARDED';
ALTER TABLE products ADD CONSTRAINT products_default_status_check
    CHECK (default_status IN ('DISCARDED', 'PENDING', 'IN_CART', 'ARCHIVED'));
ALTER TABLE pantry_items ADD CONSTRAINT pantry_items_status_check
    CHECK (status IN ('DISCARDED', 'PENDING', 'IN_CART', 'ARCHIVED'));
