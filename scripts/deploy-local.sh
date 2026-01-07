#!/usr/bin/env bash
# deploy-local.sh - Build and install my-context for all local users
# Triggered by GitHub Actions self-hosted runner on push to main
#
# Usage: ./scripts/deploy-local.sh [--dry-run]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
BINARY_NAME="my-context"

# Users to install for (must have ~/.local/bin in their PATH)
USERS=(godev be-dev-agent cdev)

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

DRY_RUN=false
if [[ "$1" == "--dry-run" ]]; then
    DRY_RUN=true
    echo -e "${YELLOW}DRY RUN MODE - no changes will be made${NC}"
fi

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Change to project directory
cd "$PROJECT_DIR"

# Get current version from git
GIT_SHA=$(git rev-parse --short HEAD)
GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

log_info "Building my-context from $GIT_BRANCH ($GIT_SHA)"

# Build the binary
if [[ "$DRY_RUN" == false ]]; then
    log_info "Compiling binary..."
    CGO_ENABLED=0 go build -ldflags "-s -w" -o "$BINARY_NAME" ./cmd/my-context/

    if [[ ! -f "$BINARY_NAME" ]]; then
        log_error "Build failed - binary not created"
        exit 1
    fi

    # Verify it runs
    VERSION=$("./$BINARY_NAME" --version 2>/dev/null || echo "unknown")
    log_info "Built version: $VERSION"
else
    log_info "[DRY RUN] Would compile binary"
fi

# Install for each user
INSTALL_COUNT=0
FAIL_COUNT=0

for USER in "${USERS[@]}"; do
    USER_HOME=$(getent passwd "$USER" | cut -d: -f6 2>/dev/null || echo "")

    if [[ -z "$USER_HOME" ]]; then
        log_warn "User $USER not found, skipping"
        continue
    fi

    INSTALL_DIR="$USER_HOME/.local/bin"
    TARGET="$INSTALL_DIR/$BINARY_NAME"

    # Check if install directory exists
    if [[ ! -d "$INSTALL_DIR" ]]; then
        log_warn "Install directory $INSTALL_DIR does not exist for $USER, skipping"
        continue
    fi

    if [[ "$DRY_RUN" == false ]]; then
        # Backup existing binary if present
        if [[ -f "$TARGET" ]]; then
            OLD_VERSION=$("$TARGET" --version 2>/dev/null || echo "unknown")
            log_info "Upgrading $USER: $OLD_VERSION -> $VERSION"
            cp "$TARGET" "$TARGET.backup" 2>/dev/null || true
        else
            log_info "Installing for $USER (new installation)"
        fi

        # Copy and set permissions
        if cp "$BINARY_NAME" "$TARGET" 2>/dev/null; then
            chmod +x "$TARGET"

            # Verify installation
            if "$TARGET" --version >/dev/null 2>&1; then
                log_info "✓ Installed for $USER"
                rm -f "$TARGET.backup" 2>/dev/null || true
                INSTALL_COUNT=$((INSTALL_COUNT + 1))
            else
                log_error "✗ Verification failed for $USER, restoring backup"
                if [[ -f "$TARGET.backup" ]]; then
                    mv "$TARGET.backup" "$TARGET"
                fi
                FAIL_COUNT=$((FAIL_COUNT + 1))
            fi
        else
            log_error "✗ Failed to copy binary for $USER (permission denied?)"
            FAIL_COUNT=$((FAIL_COUNT + 1))
        fi
    else
        log_info "[DRY RUN] Would install to $TARGET"
        INSTALL_COUNT=$((INSTALL_COUNT + 1))
    fi
done

# Cleanup build artifact
if [[ "$DRY_RUN" == false ]] && [[ -f "$BINARY_NAME" ]]; then
    rm "$BINARY_NAME"
fi

# Summary
echo ""
log_info "Deployment complete: $INSTALL_COUNT installed, $FAIL_COUNT failed"

if [[ $FAIL_COUNT -gt 0 ]]; then
    exit 1
fi

exit 0
