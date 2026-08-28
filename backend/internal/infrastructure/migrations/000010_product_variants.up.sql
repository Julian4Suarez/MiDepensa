CREATE TABLE product_variants (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID    NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    code       TEXT    NOT NULL,
    name       TEXT    NOT NULL,
    image      TEXT    NOT NULL,
    sort_order INTEGER NOT NULL,
    UNIQUE (product_id, code),
    UNIQUE (id, product_id)
);

CREATE TABLE pantry_item_variants (
    pantry_id  UUID NOT NULL,
    product_id UUID NOT NULL,
    variant_id UUID NOT NULL,
    PRIMARY KEY (pantry_id, product_id, variant_id),
    FOREIGN KEY (pantry_id, product_id)
        REFERENCES pantry_items (pantry_id, product_id) ON DELETE CASCADE,
    FOREIGN KEY (variant_id, product_id)
        REFERENCES product_variants (id, product_id) ON DELETE CASCADE
);

CREATE TEMP TABLE variant_mapping (
    parent_code TEXT,
    child_code TEXT,
    variant_code TEXT,
    variant_name TEXT,
    variant_order INTEGER
) ON COMMIT DROP;

INSERT INTO variant_mapping VALUES
    ('natural_yogurt', 'natural_yogurt', 'natural_yogurt', 'Natural yogurt', 1),
    ('natural_yogurt', 'greek_yogurt', 'greek_yogurt', 'Greek yogurt', 2),
    ('natural_yogurt', 'fruit_yogurt', 'fruit_yogurt', 'Fruit yogurt', 3),
    ('cheese', 'cheese', 'regular_cheese', 'Regular cheese', 1),
    ('cheese', 'cottage_cheese', 'cottage_cheese', 'Fresh cheese', 2),
    ('cheese', 'mozzarella', 'mozzarella', 'Mozzarella', 3),
    ('cheese', 'grated_cheese', 'grated_cheese', 'Grated cheese', 4),
    ('cheese', 'cream_cheese', 'cream_cheese', 'Cream cheese', 5),
    ('milk', 'milk', 'whole_milk', 'Whole milk', 1),
    ('milk', 'plant_milk', 'plant_milk', 'Plant-based milk', 5),
    ('fish', 'fish', 'white_fish', 'White fish', 1),
    ('fish', 'salmon', 'salmon', 'Salmon', 2),
    ('fish', 'hake', 'hake', 'Hake', 3),
    ('fish', 'tuna_steaks', 'tuna_steaks', 'Tuna steaks', 4),
    ('fish', 'frozen_fish', 'frozen_fish', 'Frozen fish', 5),
    ('water', 'water', 'still_water', 'Still water', 1),
    ('water', 'sparkling_water', 'sparkling_water', 'Sparkling water', 2),
    ('coffee', 'coffee', 'regular_coffee', 'Regular coffee', 1),
    ('coffee', 'decaf_coffee', 'decaf_coffee', 'Decaf coffee', 2),
    ('tea', 'tea', 'black_tea', 'Black tea', 1),
    ('tea', 'herbal_tea', 'herbal_tea', 'Herbal tea', 2),
    ('sugar', 'sugar', 'white_sugar', 'White sugar', 1),
    ('sugar', 'brown_sugar', 'brown_sugar', 'Brown sugar', 2),
    ('olive_oil', 'olive_oil', 'olive_oil', 'Olive oil', 1),
    ('olive_oil', 'sunflower_oil', 'sunflower_oil', 'Sunflower oil', 2),
    ('tomato_sauce', 'tomato_sauce', 'tomato_sauce', 'Tomato sauce', 1),
    ('tomato_sauce', 'crushed_tomato', 'crushed_tomato', 'Crushed tomatoes', 2),
    ('tomato_sauce', 'tomato_paste', 'tomato_paste', 'Tomato paste', 3),
    ('cooked_ham', 'cooked_ham', 'cooked_ham', 'Cooked ham', 1),
    ('cooked_ham', 'serrano_ham', 'serrano_ham', 'Serrano ham', 2),
    ('beans', 'beans', 'beans', 'Beans', 1),
    ('beans', 'lentils', 'lentils', 'Lentils', 2),
    ('beans', 'chickpeas', 'chickpeas', 'Chickpeas', 3);

INSERT INTO product_variants (product_id, code, name, image, sort_order)
SELECT parent.id, mapping.variant_code, mapping.variant_name, child.image, mapping.variant_order
FROM variant_mapping AS mapping
JOIN products AS parent ON parent.code = mapping.parent_code
JOIN products AS child ON child.code = mapping.child_code;

-- Extra obvious choices that did not exist as individual catalog products.
INSERT INTO product_variants (product_id, code, name, image, sort_order)
SELECT product.id, variant.code, variant.name, product.image, variant.sort_order
FROM products AS product
JOIN (VALUES
    ('bread', 'white_bread', 'White bread', 1),
    ('bread', 'wholemeal_bread', 'Wholemeal bread', 2),
    ('bread', 'sliced_bread', 'Sliced bread', 3),
    ('bread', 'baguette', 'Baguette', 4),
    ('milk', 'semi_skimmed_milk', 'Semi-skimmed milk', 2),
    ('milk', 'skimmed_milk', 'Skimmed milk', 3),
    ('milk', 'lactose_free_milk', 'Lactose-free milk', 4)
) AS variant(parent_code, code, name, sort_order) ON variant.parent_code = product.code;

-- Preserve child products that were already in the cart.
INSERT INTO pantry_item_variants (pantry_id, product_id, variant_id)
SELECT item.pantry_id, parent.id, variant.id
FROM variant_mapping AS mapping
JOIN products AS parent ON parent.code = mapping.parent_code
JOIN products AS child ON child.code = mapping.child_code
JOIN pantry_items AS item ON item.product_id = child.id AND item.status = 'IN_CART'
JOIN product_variants AS variant
  ON variant.product_id = parent.id AND variant.code = mapping.variant_code;

-- Bread had no concrete subtype before; preserve an in-cart bread as sliced bread.
INSERT INTO pantry_item_variants (pantry_id, product_id, variant_id)
SELECT item.pantry_id, product.id, variant.id
FROM products AS product
JOIN pantry_items AS item ON item.product_id = product.id AND item.status = 'IN_CART'
JOIN product_variants AS variant
  ON variant.product_id = product.id AND variant.code = 'sliced_bread'
WHERE product.code = 'bread';

-- The general product represents the most relevant state among its old rows.
UPDATE pantry_items AS parent_item
SET status = grouped.status, updated_at = now()
FROM (
    SELECT item.pantry_id, parent.id AS product_id,
           CASE
             WHEN bool_or(item.status = 'IN_CART') THEN 'IN_CART'
             WHEN bool_or(item.status = 'PENDING') THEN 'PENDING'
             WHEN bool_or(item.status = 'DISCARDED') THEN 'DISCARDED'
             ELSE 'ARCHIVED'
           END AS status
    FROM variant_mapping AS mapping
    JOIN products AS parent ON parent.code = mapping.parent_code
    JOIN products AS child ON child.code = mapping.child_code
    JOIN pantry_items AS item ON item.product_id = child.id
    GROUP BY item.pantry_id, parent.id
) AS grouped
WHERE parent_item.pantry_id = grouped.pantry_id
  AND parent_item.product_id = grouped.product_id;

DELETE FROM pantry_items AS item
USING products AS child, variant_mapping AS mapping
WHERE item.product_id = child.id
  AND child.code = mapping.child_code
  AND mapping.child_code <> mapping.parent_code;

DELETE FROM products AS child
USING variant_mapping AS mapping
WHERE child.code = mapping.child_code
  AND mapping.child_code <> mapping.parent_code;

UPDATE products SET name = CASE code
    WHEN 'natural_yogurt' THEN 'Yogurt'
    WHEN 'olive_oil' THEN 'Cooking oil'
    WHEN 'tomato_sauce' THEN 'Cooking tomato'
    WHEN 'cooked_ham' THEN 'Ham'
    WHEN 'beans' THEN 'Legumes'
    ELSE name
END
WHERE code IN ('natural_yogurt', 'olive_oil', 'tomato_sauce', 'cooked_ham', 'beans');

CREATE INDEX idx_product_variants_product_id ON product_variants (product_id);
CREATE INDEX idx_pantry_item_variants_item ON pantry_item_variants (pantry_id, product_id);
