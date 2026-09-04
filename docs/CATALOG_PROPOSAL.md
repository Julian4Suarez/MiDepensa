# Catálogo ampliado de productos

> Documento histórico de la ampliación inicial. Las categorías actuales se
> definen en la migración `000011_rework_categories` y la consolidación en
> productos generales con subtipos en `000012_consolidate_catalog`.

Este documento recoge el catálogo original y la ampliación implementada en la
migración `000005_expand_catalog`. En ese momento el catálogo contenía 183
productos; las migraciones posteriores agrupan los duplicados como subtipos.

## Cómo leer y editar la propuesta

- **Código**: identificador técnico único, en minúsculas y con guiones bajos.
- **Tipo**: frecuencia de compra sugerida. `ESSENTIAL` = habitual, `SECONDARY` = periódica y `OTHERS` = ocasional.
- Los nombres actuales se conservan tal cual están en la aplicación. Puedes cambiar los nombres propuestos, eliminar filas o moverlas de categoría/tipo.
- Al incorporarlos al catálogo habrá que crear también el icono SVG correspondiente a cada código nuevo.

## Catálogo actual (57 productos)

### Fruta y verdura (`FRUIT_VEG`)

| Código | Producto actual | Tipo |
| --- | --- | --- |
| tomato | Tomatoes | ESSENTIAL |
| onion | Onions | ESSENTIAL |
| garlic | Garlic | SECONDARY |
| potato | Potatoes | ESSENTIAL |
| carrot | Carrots | SECONDARY |
| leafy_greens | Leafy greens | SECONDARY |
| bell_pepper | Bell peppers | SECONDARY |
| cucumber | Cucumber | OTHERS |
| broccoli | Broccoli | OTHERS |
| mushroom | Mushrooms | OTHERS |
| avocado | Avocado | OTHERS |
| banana | Bananas | ESSENTIAL |
| apple | Apples | ESSENTIAL |
| orange | Oranges | SECONDARY |
| lemon | Lemons | SECONDARY |
| strawberry | Strawberries | OTHERS |

### Carne y pescado (`MEAT_FISH`)

| Código | Producto actual | Tipo |
| --- | --- | --- |
| chicken | Chicken | ESSENTIAL |
| red_meat | Red meat | SECONDARY |
| bacon | Bacon | OTHERS |
| fish | Fish | SECONDARY |

### Lácteos y huevos (`DAIRY_EGGS`)

| Código | Producto actual | Tipo |
| --- | --- | --- |
| milk | Milk | ESSENTIAL |
| cheese | Cheese | ESSENTIAL |
| butter | Butter | SECONDARY |
| eggs | Eggs | ESSENTIAL |

### Despensa y conservas (`DRY_CANNED`)

| Código | Producto actual | Tipo |
| --- | --- | --- |
| bread | Bread | ESSENTIAL |
| pasta | Pasta | ESSENTIAL |
| rice | Rice | ESSENTIAL |
| flour | Flour | SECONDARY |
| salt | Salt | SECONDARY |
| olive_oil | Olive oil | ESSENTIAL |
| canned_food | Canned food | SECONDARY |
| beans | Beans | SECONDARY |
| sweetcorn | Sweetcorn | OTHERS |
| nuts | Nuts | OTHERS |
| honey | Honey | OTHERS |
| cereal | Cereal | SECONDARY |
| cookies | Cookies | OTHERS |
| chocolate | Chocolate | OTHERS |
| popcorn | Popcorn | OTHERS |
| spices | Spices | OTHERS |

### Bebidas (`DRINKS`)

| Código | Producto actual | Tipo |
| --- | --- | --- |
| water | Water | ESSENTIAL |
| coffee | Coffee | ESSENTIAL |
| tea | Tea | SECONDARY |
| juice | Juice | SECONDARY |
| soft_drinks | Soft drinks | OTHERS |
| beer | Beer | OTHERS |
| wine | Wine | OTHERS |

### Hogar y cuidado personal (`HOME_CARE`)

| Código | Producto actual | Tipo |
| --- | --- | --- |
| hand_soap | Hand soap | ESSENTIAL |
| shampoo | Shampoo | SECONDARY |
| toothpaste | Toothpaste | SECONDARY |
| dish_soap | Dish soap | ESSENTIAL |
| sponges | Sponges | SECONDARY |
| floor_cleaner | Floor cleaner | OTHERS |
| toilet_paper | Toilet paper | ESSENTIAL |
| bin_bags | Bin bags | SECONDARY |
| razors | Razors | OTHERS |
| batteries | Batteries | OTHERS |

## Productos añadidos

### Fruta y verdura (`FRUIT_VEG`)

| Código propuesto | Producto | Tipo sugerido |
| --- | --- | --- |
| lettuce | Lettuce | ESSENTIAL |
| spinach | Spinach | SECONDARY |
| courgette | Courgette | SECONDARY |
| aubergine | Aubergine | OTHERS |
| cauliflower | Cauliflower | OTHERS |
| cabbage | Cabbage | OTHERS |
| green_beans | Green beans | SECONDARY |
| asparagus | Asparagus | OTHERS |
| peas | Peas | SECONDARY |
| sweet_potato | Sweet potatoes | OTHERS |
| pumpkin | Pumpkin | OTHERS |
| celery | Celery | OTHERS |
| leek | Leeks | SECONDARY |
| ginger | Ginger | SECONDARY |
| parsley | Parsley | SECONDARY |
| coriander | Coriander | OTHERS |
| pear | Pears | SECONDARY |
| grapes | Grapes | OTHERS |
| kiwi | Kiwis | SECONDARY |
| pineapple | Pineapple | OTHERS |
| watermelon | Watermelon | OTHERS |
| peaches | Peaches | OTHERS |
| blueberries | Blueberries | OTHERS |
| frozen_vegetables | Frozen vegetables | SECONDARY |

### Carne y pescado (`MEAT_FISH`)

| Código propuesto | Producto | Tipo sugerido |
| --- | --- | --- |
| turkey | Turkey | SECONDARY |
| minced_meat | Minced meat | SECONDARY |
| pork_chops | Pork chops | OTHERS |
| sausages | Sausages | OTHERS |
| cooked_ham | Cooked ham | SECONDARY |
| serrano_ham | Serrano ham | OTHERS |
| sliced_turkey | Sliced turkey | SECONDARY |
| salmon | Salmon | SECONDARY |
| hake | Hake | SECONDARY |
| tuna_steaks | Tuna steaks | OTHERS |
| prawns | Prawns | OTHERS |
| frozen_fish | Frozen fish | SECONDARY |
| tofu | Tofu | SECONDARY |
| veggie_burgers | Veggie burgers | OTHERS |

### Lácteos y huevos (`DAIRY_EGGS`)

| Código propuesto | Producto | Tipo sugerido |
| --- | --- | --- |
| natural_yogurt | Natural yogurt | ESSENTIAL |
| greek_yogurt | Greek yogurt | SECONDARY |
| fruit_yogurt | Fruit yogurt | SECONDARY |
| plant_milk | Plant-based milk | SECONDARY |
| cream | Cooking cream | SECONDARY |
| cottage_cheese | Fresh cheese | SECONDARY |
| mozzarella | Mozzarella | SECONDARY |
| grated_cheese | Grated cheese | SECONDARY |
| cream_cheese | Cream cheese | OTHERS |
| margarine | Margarine | OTHERS |
| egg_whites | Egg whites | OTHERS |
| dairy_desserts | Dairy desserts | OTHERS |

### Despensa y conservas (`DRY_CANNED`)

| Código propuesto | Producto | Tipo sugerido |
| --- | --- | --- |
| oats | Oats | SECONDARY |
| breadcrumbs | Breadcrumbs | SECONDARY |
| tortilla_wraps | Tortilla wraps | SECONDARY |
| crackers | Crackers | OTHERS |
| quinoa | Quinoa | SECONDARY |
| couscous | Couscous | SECONDARY |
| noodles | Noodles | SECONDARY |
| lentils | Lentils | SECONDARY |
| chickpeas | Chickpeas | SECONDARY |
| canned_peas | Canned peas | OTHERS |
| canned_tuna | Canned tuna | SECONDARY |
| canned_sardines | Canned sardines | OTHERS |
| tomato_sauce | Tomato sauce | ESSENTIAL |
| crushed_tomato | Crushed tomatoes | SECONDARY |
| tomato_paste | Tomato paste | OTHERS |
| mayonnaise | Mayonnaise | SECONDARY |
| ketchup | Ketchup | SECONDARY |
| mustard | Mustard | OTHERS |
| vinegar | Vinegar | SECONDARY |
| sunflower_oil | Sunflower oil | SECONDARY |
| sugar | Sugar | SECONDARY |
| brown_sugar | Brown sugar | OTHERS |
| baking_powder | Baking powder | OTHERS |
| yeast | Baker | OTHERS |
| cocoa_powder | Cocoa powder | OTHERS |
| jam | Jam | SECONDARY |
| peanut_butter | Peanut butter | SECONDARY |
| dried_fruit | Dried fruit | OTHERS |
| cereal_bars | Cereal bars | OTHERS |
| crisps | Crisps | OTHERS |
| frozen_pizza | Frozen pizza | OTHERS |

### Bebidas (`DRINKS`)

| Código propuesto | Producto | Tipo sugerido |
| --- | --- | --- |
| sparkling_water | Sparkling water | SECONDARY |
| decaf_coffee | Decaf coffee | SECONDARY |
| herbal_tea | Herbal tea | SECONDARY |
| iced_tea | Iced tea | OTHERS |
| energy_drinks | Energy drinks | OTHERS |
| tonic_water | Tonic water | OTHERS |
| isotonic_drink | Isotonic drink | OTHERS |
| kombucha | Kombucha | OTHERS |
| cider | Cider | OTHERS |
| sparkling_wine | Sparkling wine | OTHERS |
| spirits | Spirits | OTHERS |

### Hogar y cuidado personal (`HOME_CARE`)

| Código propuesto | Producto | Tipo sugerido |
| --- | --- | --- |
| laundry_detergent | Laundry detergent | ESSENTIAL |
| fabric_softener | Fabric softener | SECONDARY |
| stain_remover | Stain remover | OTHERS |
| bleach | Bleach | SECONDARY |
| disinfectant | Disinfectant | SECONDARY |
| multipurpose_cleaner | Multipurpose cleaner | ESSENTIAL |
| glass_cleaner | Glass cleaner | OTHERS |
| degreaser | Degreaser | SECONDARY |
| toilet_cleaner | Toilet cleaner | SECONDARY |
| dishwasher_tablets | Dishwasher tablets | ESSENTIAL |
| dishwasher_salt | Dishwasher salt | OTHERS |
| dishcloths | Cleaning cloths | SECONDARY |
| rubber_gloves | Rubber gloves | OTHERS |
| aluminium_foil | Aluminium foil | SECONDARY |
| baking_paper | Baking paper | SECONDARY |
| cling_film | Cling film | SECONDARY |
| paper_towels | Kitchen roll | ESSENTIAL |
| napkins | Napkins | SECONDARY |
| tissues | Tissues | SECONDARY |
| deodorant | Deodorant | ESSENTIAL |
| shower_gel | Shower gel | ESSENTIAL |
| conditioner | Conditioner | SECONDARY |
| dental_floss | Dental floss | SECONDARY |
| mouthwash | Mouthwash | SECONDARY |
| toothbrush_heads | Toothbrushes or replacement heads | SECONDARY |
| cotton_pads | Cotton pads | SECONDARY |
| feminine_hygiene | Feminine hygiene | ESSENTIAL |
| sunscreen | Sunscreen | OTHERS |
| insect_repellent | Insect repellent | OTHERS |
| first_aid_kit | First aid supplies | OTHERS |
| light_bulbs | Light bulbs | OTHERS |
| matches | Matches or lighter | OTHERS |
| candles | Candles | OTHERS |
| air_freshener | Air freshener | OTHERS |

## Posibles revisiones futuras

- Los nombres mostrados por la aplicación se mantienen en inglés.
- ¿Quieres separar en más categorías algunos productos (por ejemplo, congelados, panadería, mascotas o bebé) en una fase posterior?
- Revisa especialmente el tipo sugerido: depende de cada hogar y no debe tratarse como una recomendación rígida.
