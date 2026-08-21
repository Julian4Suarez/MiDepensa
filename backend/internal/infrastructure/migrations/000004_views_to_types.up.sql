-- Turns the three "views" into a product type, filtered on the same screen as
-- the categories instead of through a side menu.
--
-- PRIMARY -> ESSENTIAL and OTHER -> OTHERS are a pure rename: the old OTHER
-- bucket already meant "rarely used", which is what OTHERS means.

ALTER TABLE products RENAME COLUMN default_view TO default_type;
ALTER TABLE pantry_items RENAME COLUMN pantry_view TO product_type;

-- Renaming a column keeps its CHECK constraint under the original name.
ALTER TABLE products DROP CONSTRAINT IF EXISTS products_default_view_check;
ALTER TABLE pantry_items DROP CONSTRAINT IF EXISTS pantry_items_pantry_view_check;

UPDATE products SET default_type = CASE default_type
    WHEN 'PRIMARY' THEN 'ESSENTIAL'
    WHEN 'OTHER' THEN 'OTHERS'
    ELSE default_type
END;

UPDATE pantry_items SET product_type = CASE product_type
    WHEN 'PRIMARY' THEN 'ESSENTIAL'
    WHEN 'OTHER' THEN 'OTHERS'
    ELSE product_type
END;

ALTER TABLE products ADD CONSTRAINT products_default_type_check
    CHECK (default_type IN ('ESSENTIAL', 'SECONDARY', 'OTHERS'));

ALTER TABLE pantry_items ADD CONSTRAINT pantry_items_product_type_check
    CHECK (product_type IN ('ESSENTIAL', 'SECONDARY', 'OTHERS'));
