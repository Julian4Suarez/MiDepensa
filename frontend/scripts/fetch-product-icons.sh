#!/usr/bin/env bash
# Downloads the product icons from OpenMoji (CC BY-SA 4.0) into
# src/assets/products/, named after the product codes used by the API seed.
#
# The icons are committed to the repository, so this only needs to be re-run
# when a product is added to backend/internal/infrastructure/migrations.
#
# Usage: ./scripts/fetch-product-icons.sh
set -euo pipefail

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
"

mkdir -p "$TARGET_DIR"

failed=0
while read -r code codepoint; do
    [ -z "$code" ] && continue
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
