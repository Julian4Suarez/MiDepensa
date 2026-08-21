#!/usr/bin/env bash
# PostgreSQL backup and restore.
#
# Usage:
#   ./scripts/db.sh backup  <environment>
#   ./scripts/db.sh restore <environment> [backup-file]
#   ./scripts/db.sh list    <environment>
set -euo pipefail

INFRA_DIR="$(cd "$(dirname "$0")/.." && pwd)"
COMMAND="${1:-}"
ENVIRONMENT="${2:-}"
BACKUP_FILE="${3:-}"

if [[ -z "$COMMAND" || -z "$ENVIRONMENT" ]]; then
    echo "Usage: $0 <backup|restore|list> <environment> [backup-file]" >&2
    exit 1
fi

ENV_FILE="$INFRA_DIR/.env.$ENVIRONMENT"
[[ -f "$ENV_FILE" ]] || { echo "Missing $ENV_FILE" >&2; exit 1; }

# shellcheck disable=SC1090  # path is built at runtime
set -a && source "$ENV_FILE" && set +a

CONTAINER="midepensa_${ENVIRONMENT}_postgres"
BACKUP_DIR="${BACKUP_DIR:-$HOME/backups/midepensa/$ENVIRONMENT}"
RETENTION_DAYS="${RETENTION_DAYS:-7}"

mkdir -p "$BACKUP_DIR"

case "$COMMAND" in
    backup)
        target="$BACKUP_DIR/$(date +%Y%m%d_%H%M%S).dump.gz"
        echo "Creating $target"
        docker exec "$CONTAINER" pg_dump -U "$DB_USER" -d "$DB_NAME" -Fc | gzip >"$target"
        [[ -s "$target" ]] || { echo "Backup is empty" >&2; rm -f "$target"; exit 1; }
        find "$BACKUP_DIR" -name '*.dump.gz' -mtime "+$RETENTION_DAYS" -delete
        echo "Done ($(du -h "$target" | cut -f1))"
        ;;
    restore)
        if [[ -z "$BACKUP_FILE" ]]; then
            BACKUP_FILE="$(find "$BACKUP_DIR" -name '*.dump.gz' | sort | tail -1)"
        fi
        [[ -n "$BACKUP_FILE" && -f "$BACKUP_FILE" ]] || { echo "No backup found" >&2; exit 1; }
        echo "Restoring $BACKUP_FILE into $CONTAINER"
        gunzip -c "$BACKUP_FILE" |
            docker exec -i "$CONTAINER" pg_restore -U "$DB_USER" -d "$DB_NAME" --clean --if-exists
        echo "Done"
        ;;
    list)
        find "$BACKUP_DIR" -name '*.dump.gz' -printf '%f\t%s bytes\n' | sort
        ;;
    *)
        echo "Unknown command: $COMMAND" >&2
        exit 1
        ;;
esac
