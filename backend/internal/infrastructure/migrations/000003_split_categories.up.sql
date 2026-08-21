-- Replaces the four original categories with six supermarket-aisle ones.
--
-- The old FRESH split into three aisles, and PANTRY was renamed because it
-- collided with the name of the app itself. Item categories are re-derived from
-- the catalog: a per-item override under the old taxonomy has no equivalent
-- under the new one. Stock status and view are untouched.

ALTER TABLE products DROP CONSTRAINT IF EXISTS products_default_category_check;
ALTER TABLE pantry_items DROP CONSTRAINT IF EXISTS pantry_items_category_check;

UPDATE products SET default_category = 'FRUIT_VEG' WHERE code IN (
    'tomato', 'onion', 'garlic', 'potato', 'carrot', 'leafy_greens', 'bell_pepper',
    'cucumber', 'broccoli', 'mushroom', 'avocado', 'banana', 'apple', 'orange',
    'lemon', 'strawberry'
);

UPDATE products SET default_category = 'MEAT_FISH' WHERE code IN (
    'chicken', 'red_meat', 'bacon', 'fish'
);

UPDATE products SET default_category = 'DAIRY_EGGS' WHERE code IN (
    'milk', 'cheese', 'butter', 'eggs'
);

UPDATE products SET default_category = 'DRY_CANNED' WHERE code IN (
    'bread', 'pasta', 'rice', 'flour', 'salt', 'olive_oil', 'canned_food', 'beans',
    'sweetcorn', 'nuts', 'honey', 'cereal', 'cookies', 'chocolate', 'popcorn', 'spices'
);

-- DRINKS and HOME_CARE keep their codes and their members.

UPDATE pantry_items i
SET category = p.default_category
FROM products p
WHERE p.id = i.product_id;

ALTER TABLE products ADD CONSTRAINT products_default_category_check
    CHECK (default_category IN ('FRUIT_VEG', 'MEAT_FISH', 'DAIRY_EGGS', 'DRY_CANNED', 'DRINKS', 'HOME_CARE'));

ALTER TABLE pantry_items ADD CONSTRAINT pantry_items_category_check
    CHECK (category IN ('FRUIT_VEG', 'MEAT_FISH', 'DAIRY_EGGS', 'DRY_CANNED', 'DRINKS', 'HOME_CARE'));
