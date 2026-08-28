-- Categories now describe what a product is and where it is normally found,
-- rather than how it is packaged or preserved.
ALTER TABLE products DROP CONSTRAINT IF EXISTS products_default_category_check;
ALTER TABLE pantry_items DROP CONSTRAINT IF EXISTS pantry_items_category_check;

UPDATE products SET default_category = CASE default_category
    WHEN 'FRUIT_VEG' THEN 'VEGETABLES'
    WHEN 'DRY_CANNED' THEN 'COOKING_CONDIMENTS'
    WHEN 'HOME_CARE' THEN 'HOUSEHOLD_CLEANING'
    ELSE default_category
END;

UPDATE products SET default_category = 'FRUIT' WHERE code IN (
    'banana', 'apple', 'orange', 'lemon', 'strawberry', 'pear', 'grapes',
    'kiwi', 'pineapple', 'watermelon', 'peaches', 'blueberries'
);

UPDATE products SET default_category = 'VEGETABLES' WHERE code IN (
    'tomato', 'onion', 'garlic', 'potato', 'carrot', 'leafy_greens',
    'bell_pepper', 'cucumber', 'broccoli', 'mushroom', 'avocado', 'lettuce',
    'spinach', 'courgette', 'aubergine', 'cauliflower', 'cabbage',
    'green_beans', 'asparagus', 'peas', 'sweet_potato', 'pumpkin', 'celery',
    'leek', 'ginger', 'parsley', 'coriander', 'frozen_vegetables',
    'sweetcorn', 'canned_peas'
);

UPDATE products SET default_category = 'MEAT_FISH' WHERE code IN (
    'chicken', 'red_meat', 'bacon', 'fish', 'turkey', 'minced_meat',
    'pork_chops', 'sausages', 'cooked_ham', 'sliced_turkey', 'prawns',
    'tofu', 'veggie_burgers', 'canned_tuna', 'canned_sardines'
);

UPDATE products SET default_category = 'DAIRY_EGGS' WHERE code IN (
    'milk', 'cheese', 'butter', 'eggs', 'natural_yogurt', 'cream',
    'margarine', 'egg_whites', 'dairy_desserts'
);

UPDATE products SET default_category = 'BAKERY_BREAKFAST_SNACKS' WHERE code IN (
    'bread', 'nuts', 'honey', 'cereal', 'cookies', 'chocolate', 'popcorn',
    'oats', 'tortilla_wraps', 'crackers', 'jam', 'peanut_butter',
    'dried_fruit', 'cereal_bars', 'crisps', 'frozen_pizza'
);

UPDATE products SET default_category = 'PASTA_RICE_LEGUMES' WHERE code IN (
    'pasta', 'rice', 'beans', 'quinoa', 'couscous', 'noodles'
);

UPDATE products SET default_category = 'COOKING_CONDIMENTS' WHERE code IN (
    'flour', 'salt', 'olive_oil', 'canned_food', 'spices', 'breadcrumbs',
    'tomato_sauce', 'mayonnaise', 'ketchup', 'mustard', 'vinegar', 'sugar',
    'baking_powder', 'yeast', 'cocoa_powder'
);

UPDATE products SET default_category = 'PERSONAL_CARE' WHERE code IN (
    'hand_soap', 'shampoo', 'toothpaste', 'razors', 'deodorant', 'shower_gel',
    'conditioner', 'dental_floss', 'mouthwash', 'toothbrush_heads',
    'cotton_pads', 'feminine_hygiene', 'sunscreen', 'insect_repellent',
    'first_aid_kit'
);

UPDATE pantry_items AS item
SET category = product.default_category
FROM products AS product
WHERE product.id = item.product_id;

ALTER TABLE products ADD CONSTRAINT products_default_category_check CHECK (
    default_category IN (
        'FRUIT', 'VEGETABLES', 'MEAT_FISH', 'DAIRY_EGGS',
        'BAKERY_BREAKFAST_SNACKS', 'PASTA_RICE_LEGUMES', 'COOKING_CONDIMENTS',
        'DRINKS', 'HOUSEHOLD_CLEANING', 'PERSONAL_CARE'
    )
);
ALTER TABLE pantry_items ADD CONSTRAINT pantry_items_category_check CHECK (
    category IN (
        'FRUIT', 'VEGETABLES', 'MEAT_FISH', 'DAIRY_EGGS',
        'BAKERY_BREAKFAST_SNACKS', 'PASTA_RICE_LEGUMES', 'COOKING_CONDIMENTS',
        'DRINKS', 'HOUSEHOLD_CLEANING', 'PERSONAL_CARE'
    )
);
