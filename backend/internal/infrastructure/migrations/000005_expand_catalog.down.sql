-- Removing catalog products cascades to their per-pantry items.
DELETE FROM products
WHERE code IN (
    'lettuce', 'spinach', 'courgette', 'aubergine', 'cauliflower', 'cabbage',
    'green_beans', 'asparagus', 'peas', 'sweet_potato', 'pumpkin', 'celery',
    'leek', 'ginger', 'parsley', 'coriander', 'pear', 'grapes',
    'kiwi', 'pineapple', 'watermelon', 'peaches', 'blueberries', 'frozen_vegetables',
    'turkey', 'minced_meat', 'pork_chops', 'sausages', 'cooked_ham', 'serrano_ham',
    'sliced_turkey', 'salmon', 'hake', 'tuna_steaks', 'prawns', 'frozen_fish',
    'tofu', 'veggie_burgers', 'natural_yogurt', 'greek_yogurt', 'fruit_yogurt', 'plant_milk',
    'cream', 'cottage_cheese', 'mozzarella', 'grated_cheese', 'cream_cheese', 'margarine',
    'egg_whites', 'dairy_desserts', 'oats', 'breadcrumbs', 'tortilla_wraps', 'crackers',
    'quinoa', 'couscous', 'noodles', 'lentils', 'chickpeas', 'canned_peas',
    'canned_tuna', 'canned_sardines', 'tomato_sauce', 'crushed_tomato', 'tomato_paste', 'mayonnaise',
    'ketchup', 'mustard', 'vinegar', 'sunflower_oil', 'sugar', 'brown_sugar',
    'baking_powder', 'yeast', 'cocoa_powder', 'jam', 'peanut_butter', 'dried_fruit',
    'cereal_bars', 'crisps', 'frozen_pizza', 'sparkling_water', 'decaf_coffee', 'herbal_tea',
    'iced_tea', 'energy_drinks', 'tonic_water', 'isotonic_drink', 'kombucha', 'cider',
    'sparkling_wine', 'spirits', 'laundry_detergent', 'fabric_softener', 'stain_remover', 'bleach',
    'disinfectant', 'multipurpose_cleaner', 'glass_cleaner', 'degreaser', 'toilet_cleaner', 'dishwasher_tablets',
    'dishwasher_salt', 'dishcloths', 'rubber_gloves', 'aluminium_foil', 'baking_paper', 'cling_film',
    'paper_towels', 'napkins', 'tissues', 'deodorant', 'shower_gel', 'conditioner',
    'dental_floss', 'mouthwash', 'toothbrush_heads', 'cotton_pads', 'feminine_hygiene', 'sunscreen',
    'insect_repellent', 'first_aid_kit', 'light_bulbs', 'matches', 'candles', 'air_freshener'
);

