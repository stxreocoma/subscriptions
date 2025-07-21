#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

export TERM=xterm-256color

# Colors
RED=$(tput setaf 1)
GREEN=$(tput setaf 2)
YELLOW=$(tput setaf 3)
BLUE=$(tput setaf 4)
BOLD=$(tput bold)
NC=$(tput sgr0)

# Error codes
ERR_MIGRATE=2
ERR_VALIDATION=3

# Color output function
log() {
    local level="$1"
    local message="$2"
    local color icon

    case "$level" in
        ERROR)   color="$RED";    icon="❌";;
        WARN)    color="$YELLOW"; icon="⚠️ ";;
        INFO)    color="$BLUE";   icon="ℹ️ ";;
        SUCCESS) color="$GREEN";  icon="✅";;
        *)       color="";        icon="";;
    esac

    local timestamp
    timestamp=$(date '+%Y-%m-%d %H:%M:%S')

    if [ -t 1 ]; then
        echo -e "${BOLD}${color}${icon} [$timestamp] $level:${NC} $message"
    else
        echo "[$timestamp] $level: $message"
    fi
}

# Load .env if exists, without overriding existing vars
if [ -f ".env" ]; then
    log INFO "Loading environment from .env"
    while IFS='=' read -r key value; do
        [[ "$key" =~ ^#.*$ || -z "$key" ]] && continue
        if [ -z "${!key:-}" ]; then
            export "$key=$value" || {
                log ERROR "Failed to export $key"
                exit $ERR_ENV
            }
        fi
    done < .env
fi

# DB settings
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-postgres}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-1}
DB_SSLMODE=${DB_SSLMODE:-disable}
MIGRATIONS_DIR="internal/storage/database/migrations"

# Generate DB URL
DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"
MIGRATE_CMD="migrate -source file://${MIGRATIONS_DIR} -database ${DB_URL}"

COMMAND=${1:-}

case "$COMMAND" in
    up)
        log INFO "Applying migrations (up)..."
        $MIGRATE_CMD up || {
            log ERROR "Error applying migrations"
            exit $ERR_MIGRATE
        }
        log SUCCESS "Migrations applied successfully"
        ;;

    down)
        log INFO "Rolling back last migration (down)..."
        $MIGRATE_CMD down 1 || {
            log ERROR "Error rolling back migration"
            exit $ERR_MIGRATE
        }
        log SUCCESS "Migration rolled back successfully"
        ;;

    force)
        VERSION_ARG=${2:-}
        if [[ -z "$VERSION_ARG" || ! "$VERSION_ARG" =~ ^[0-9]+$ ]]; then
            log ERROR "You must provide a numeric version: $0 force <version>"
            exit $ERR_VALIDATION
        fi
        log INFO "Forcing migration version to $VERSION_ARG..."
        $MIGRATE_CMD force "$VERSION_ARG" || {
            log ERROR "Error forcing migration version"
            exit $ERR_MIGRATE
        }
        log SUCCESS "Migration version forced to $VERSION_ARG"
        ;;

    version)
        log INFO "Current migration version:"
        $MIGRATE_CMD version || {
            log ERROR "Error retrieving migration version"
            exit $ERR_MIGRATE
        }
        ;;

    clean)
        log WARN "Dropping all migrations (this is destructive)..."
        $MIGRATE_CMD drop || {
            log ERROR "Error dropping migrations"
            exit $ERR_MIGRATE
        }
        log SUCCESS "All migrations have been dropped"
        ;;

    help|"")
        echo -e "${BOLD}Usage:${NC} $0 [command] [args]"
        echo ""
        echo "Commands:"
        echo "  up                  - Apply all pending migrations"
        echo "  down                - Roll back the last migration"
        echo "  force <version>     - Force set migration version (e.g. 3)"
        echo "  version             - Show current migration version"
        echo "  clean               - Drop all migrations (⚠ destructive!)"
        echo "  help                - Show this help message"
        echo ""
        echo "Example:"
        echo "  $0 up"
        echo "  $0 force 5"
        ;;

    *)
        log ERROR "Unknown command: $COMMAND"
        echo "Run '$0 help' to see available commands"
        exit $ERR_VALIDATION
        ;;
esac