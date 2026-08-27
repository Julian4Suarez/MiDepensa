#!/usr/bin/env bash
# Downloads the product icons from OpenMoji (CC BY-SA 4.0) into
# src/assets/products/, named after the product codes used by the API seed.
#
# The icons are committed to the repository, so this only needs to be re-run
# when a product is added to backend/internal/infrastructure/migrations.
#
# Usage: ./scripts/fetch-product-icons.sh [--missing-only]
set -euo pipefail

missing_only=0
if [ "${1:-}" = "--missing-only" ]; then
    missing_only=1
elif [ "$#" -gt 0 ]; then
    echo "Usage: $0 [--missing-only]" >&2
    exit 2
fi

BASE_URL="https://raw.githubusercontent.com/hfg-gmuend/openmoji/master/color/svg"
TARGET_DIR="$(cd "$(dirname "$0")/.." && pwd)/src/assets/products"

# product_code emoji_codepoint  — keep in sync with the seed migration.
ICONS="
tomato 1F345
onion 1F9C5
garlic 1F9C4
potato 1F954
carrot 1F955
leafy_greens 1F96C
bell_pepper 1FAD1
cucumber 1F952
broccoli 1F966
mushroom 1F344
avocado 1F951
banana 1F34C
apple 1F34E
orange 1F34A
lemon 1F34B
strawberry 1F353
milk 1F95B
cheese 1F9C0
butter 1F9C8
eggs 1F95A
chicken 1F357
red_meat 1F969
bacon 1F953
fish 1F41F
bread 1F35E
pasta 1F35D
rice 1F35A
flour 1F33E
salt 1F9C2
olive_oil 1FAD2
canned_food 1F96B
beans 1FAD8
sweetcorn 1F33D
nuts 1F95C
honey 1F36F
cereal 1F963
cookies 1F36A
chocolate 1F36B
popcorn 1F37F
spices 1F336
water 1F4A7
coffee 2615
tea 1F375
juice 1F9C3
soft_drinks 1F964
beer 1F37A
wine 1F377
hand_soap 1F9FC
shampoo 1F9F4
toothpaste 1FAA5
dish_soap 1FAE7
sponges 1F9FD
floor_cleaner 1F9F9
toilet_paper 1F9FB
bin_bags 1F5D1
razors 1FA92
batteries 1F50B
lettuce 1F96C
spinach 1F96C
courgette 1F952
aubergine 1F346
cauliflower 1F966
cabbage 1F96C
green_beans 1FAD8
asparagus 1F96C
peas 1FAD8
sweet_potato 1F360
pumpkin 1F383
celery 1F96C
leek 1F9C5
ginger 1FADA
parsley 1F33F
coriander 1F33F
pear 1F350
grapes 1F347
kiwi 1F95D
pineapple 1F34D
watermelon 1F349
peaches 1F351
blueberries 1FAD0
frozen_vegetables 1F9CA
turkey 1F983
minced_meat 1F969
pork_chops 1F356
sausages 1F32D
cooked_ham 1F356
serrano_ham 1F356
sliced_turkey 1F983
salmon 1F41F
hake 1F41F
tuna_steaks 1F41F
prawns 1F990
frozen_fish 1F9CA
tofu 1F372
veggie_burgers 1F354
natural_yogurt 1F963
greek_yogurt 1F963
fruit_yogurt 1F963
plant_milk 1F95B
cream 1F95B
cottage_cheese 1F9C0
mozzarella 1F9C0
grated_cheese 1F9C0
cream_cheese 1F9C0
margarine 1F9C8
egg_whites 1F95A
dairy_desserts 1F36E
oats 1F33E
breadcrumbs 1F35E
tortilla_wraps 1F32F
crackers 1F9C7
quinoa 1F35A
couscous 1F35A
noodles 1F35C
lentils 1FAD8
chickpeas 1FAD8
canned_peas 1F96B
canned_tuna 1F96B
canned_sardines 1F96B
tomato_sauce 1F345
crushed_tomato 1F345
tomato_paste 1F345
mayonnaise 1FAD9
ketchup 1F345
mustard 1F32D
vinegar 1FAD2
sunflower_oil 1FAD2
sugar 1F36C
brown_sugar 1F36C
baking_powder 1F33E
yeast 1F35E
cocoa_powder 1F36B
jam 1F36F
peanut_butter 1F95C
dried_fruit 1F347
cereal_bars 1F36B
crisps 1F35F
frozen_pizza 1F355
sparkling_water 1F4A7
decaf_coffee 2615
herbal_tea 1F375
iced_tea 1F9CB
energy_drinks 1F964
tonic_water 1F964
isotonic_drink 1F964
kombucha 1F9C3
cider 1F37A
sparkling_wine 1F942
spirits 1F943
laundry_detergent 1F9F4
fabric_softener 1F9F4
stain_remover 1F9F4
bleach 1F9F4
disinfectant 1F9F4
multipurpose_cleaner 1F9F4
glass_cleaner 1F9F4
degreaser 1F9F4
toilet_cleaner 1F9F4
dishwasher_tablets 1F9FC
dishwasher_salt 1F9C2
dishcloths 1F9FD
rubber_gloves 1F9E4
aluminium_foil 1F9FB
baking_paper 1F4DC
cling_film 1F9FB
paper_towels 1F9FB
napkins 1F9FB
tissues 1F9FB
deodorant 1F9F4
shower_gel 1F9FC
conditioner 1F9F4
dental_floss 1F9F5
mouthwash 1F9F4
toothbrush_heads 1FAA5
cotton_pads 26AA
feminine_hygiene 1FA78
sunscreen 1F9F4
insect_repellent 1F99F
first_aid_kit 1FA79
light_bulbs 1F4A1
matches 1F525
candles 1F56F
air_freshener 1F33C
"

mkdir -p "$TARGET_DIR"

failed=0
while read -r code codepoint; do
    [ -z "$code" ] && continue
    if [ "$missing_only" -eq 1 ] && [ -f "$TARGET_DIR/$code.svg" ]; then
        continue
    fi
    if curl -fsSL "$BASE_URL/$codepoint.svg" -o "$TARGET_DIR/$code.svg"; then
        printf '  %-15s <- %s\n' "$code" "$codepoint"
    else
        echo "  FAILED: $code ($codepoint)" >&2
        rm -f "$TARGET_DIR/$code.svg"
        failed=$((failed + 1))
    fi
done <<< "$ICONS"

if [ "$failed" -gt 0 ]; then
    echo "$failed icon(s) could not be downloaded." >&2
    exit 1
fi

echo "Icons written to $TARGET_DIR"
