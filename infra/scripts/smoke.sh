#!/usr/bin/env bash
# Post-deploy health verification. Run automatically by deploy.sh.
#
# Usage: ./scripts/smoke.sh <local|staging|production>
set -uo pipefail

INFRA_DIR="$(cd "$(dirname "$0")/.." && pwd)"
ENVIRONMENT="${1:-local}"
ENV_FILE="$INFRA_DIR/.env.$ENVIRONMENT"

[[ -f "$ENV_FILE" ]] || { echo "Missing $ENV_FILE" >&2; exit 1; }

# shellcheck disable=SC1090  # path is built at runtime
set -a && source "$ENV_FILE" && set +a

BACKEND="midepensa_${ENVIRONMENT}_backend"
FRONTEND="midepensa_${ENVIRONMENT}_frontend"
POSTGRES="midepensa_${ENVIRONMENT}_postgres"

failures=0

pass() { echo "  ok    $1"; }
fail() { echo "  FAIL  $1" >&2; failures=$((failures + 1)); }

check_running() {
    if [[ "$(docker inspect -f '{{.State.Running}}' "$1" 2>/dev/null)" == "true" ]]; then
        pass "$1 is running"
    else
        fail "$1 is not running"
    fi
}

echo "Smoke test for '$ENVIRONMENT'"

check_running "$POSTGRES"
check_running "$BACKEND"
check_running "$FRONTEND"

if docker exec "$POSTGRES" pg_isready -U "$DB_USER" >/dev/null 2>&1; then
    pass "postgres accepts connections"
else
    fail "postgres is not accepting connections"
fi

if curl -fsS "http://127.0.0.1:${BACKEND_PORT}/healthz" >/dev/null; then
    pass "backend /healthz"
else
    fail "backend /healthz"
fi

if curl -fsS "http://127.0.0.1:${BACKEND_PORT}/readyz" | grep -q '"status":"ready"'; then
    pass "backend /readyz reports ready"
else
    fail "backend /readyz"
fi

if curl -fsS "http://127.0.0.1:${BACKEND_PORT}/v1/catalog" | grep -q '"products"'; then
    pass "catalog is seeded"
else
    fail "catalog is empty or unreachable"
fi

if curl -fsS "http://127.0.0.1:${FRONTEND_PORT}/" >/dev/null; then
    pass "frontend responds"
else
    fail "frontend does not respond"
fi

if curl -fsS "http://127.0.0.1:${FRONTEND_PORT}/config.json" | grep -q 'BACKEND_URL'; then
    pass "frontend runtime config is present"
else
    fail "frontend config.json is missing"
fi

echo ""
if [[ "$failures" -gt 0 ]]; then
    echo "$failures check(s) failed." >&2
    exit 1
fi
echo "All checks passed. Open http://localhost:${FRONTEND_PORT}"
