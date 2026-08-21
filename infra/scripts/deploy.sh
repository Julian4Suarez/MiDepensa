#!/usr/bin/env bash
# Deploys or stops a MiDepensa environment.
#
# Usage: ./scripts/deploy.sh <local|staging|production> [deploy|stop]
#
# local      builds the images from this working copy
# staging /  pull the images pinned in environments/<env>.manifest.env
# production
set -euo pipefail

INFRA_DIR="$(cd "$(dirname "$0")/.." && pwd)"
ENVIRONMENT="${1:-}"
ACTION="${2:-deploy}"

case "$ENVIRONMENT" in
    local | staging | production) ;;
    *)
        echo "Usage: $0 <local|staging|production> [deploy|stop]" >&2
        exit 1
        ;;
esac

ENV_FILE="$INFRA_DIR/.env.$ENVIRONMENT"
COMPOSE=(docker compose --env-file "$ENV_FILE" -f "$INFRA_DIR/docker-compose.yml")

# Local is meant to be zero-setup: seed the env file from the template.
if [[ "$ENVIRONMENT" == "local" && ! -f "$ENV_FILE" ]]; then
    cp "$INFRA_DIR/.env.example" "$ENV_FILE"
    echo "Created $ENV_FILE from .env.example"
fi

if [[ "$ENVIRONMENT" != "local" && ! -f "$ENV_FILE" ]]; then
    echo "Missing $ENV_FILE — run ./scripts/build_env.sh $ENVIRONMENT first." >&2
    exit 1
fi

if [[ "$ACTION" == "stop" ]]; then
    "${COMPOSE[@]}" down
    echo "Environment '$ENVIRONMENT' stopped."
    exit 0
fi

"$INFRA_DIR/scripts/validate_env.sh" "$ENVIRONMENT"

if [[ "$ENVIRONMENT" == "local" ]]; then
    "${COMPOSE[@]}" up -d --build --wait
else
    # Remote environments only ever run images built by CI.
    "${COMPOSE[@]}" pull
    "${COMPOSE[@]}" up -d --remove-orphans --wait
fi

docker image prune -f >/dev/null

echo ""
"$INFRA_DIR/scripts/smoke.sh" "$ENVIRONMENT"
