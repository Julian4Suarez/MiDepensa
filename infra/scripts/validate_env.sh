#!/usr/bin/env bash
# Verifies that .env.<environment> defines every variable listed in .env.example.
#
# Usage: ./scripts/validate_env.sh <local|staging|production>
set -euo pipefail

INFRA_DIR="$(cd "$(dirname "$0")/.." && pwd)"
ENVIRONMENT="${1:-}"

if [[ -z "$ENVIRONMENT" ]]; then
    echo "Usage: $0 <local|staging|production>" >&2
    exit 1
fi

TEMPLATE="$INFRA_DIR/.env.example"
TARGET="$INFRA_DIR/.env.$ENVIRONMENT"

[[ -f "$TEMPLATE" ]] || { echo "Missing $TEMPLATE" >&2; exit 1; }
[[ -f "$TARGET" ]] || { echo "Missing $TARGET" >&2; exit 1; }

failures=0

# Every KEY= line of the template must exist with a non-empty value.
while IFS= read -r key; do
    value="$(grep -E "^${key}=" "$TARGET" | head -1 | cut -d= -f2-)"
    if ! grep -qE "^${key}=" "$TARGET"; then
        echo "  missing: $key" >&2
        failures=$((failures + 1))
    elif [[ -z "$value" ]]; then
        echo "  empty:   $key" >&2
        failures=$((failures + 1))
    fi
done < <(grep -E '^[A-Z_][A-Z0-9_]*=' "$TEMPLATE" | cut -d= -f1)

# Reject unresolved secret placeholders outside local.
if [[ "$ENVIRONMENT" != "local" ]] && grep -q 'CHANGE_ME' "$TARGET"; then
    echo "  unresolved CHANGE_ME placeholder(s) in $TARGET" >&2
    failures=$((failures + 1))
fi

if [[ "$failures" -gt 0 ]]; then
    echo "$failures check(s) failed." >&2
    exit 1
fi

echo "Environment '$ENVIRONMENT' is valid."
