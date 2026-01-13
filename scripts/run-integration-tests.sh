#!/usr/bin/env bash
# Integration test runner for my-context
# Runs BATS tests in isolated environments

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
BATS_DIR="$PROJECT_ROOT/tests/bats"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Check for BATS installation
check_bats() {
    if command -v bats &> /dev/null; then
        log_info "Using system BATS: $(bats --version)"
        return 0
    fi

    # Check for local installation
    if [[ -x "$PROJECT_ROOT/test/bats/bin/bats" ]]; then
        export PATH="$PROJECT_ROOT/test/bats/bin:$PATH"
        log_info "Using local BATS installation"
        return 0
    fi

    log_error "BATS not found. Install with:"
    echo "  brew install bats-core    # macOS"
    echo "  apt install bats          # Debian/Ubuntu"
    echo "  # Or clone as submodule:"
    echo "  git submodule add https://github.com/bats-core/bats-core.git test/bats"
    return 1
}

# Build test binary
build_binary() {
    log_info "Building test binary..."
    mkdir -p "$PROJECT_ROOT/bin"
    go build -o "$PROJECT_ROOT/bin/my-context-test" "$PROJECT_ROOT/cmd/my-context/"
    log_info "Binary built: $PROJECT_ROOT/bin/my-context-test"
}

# Run file mode tests
run_file_tests() {
    log_info "Running file mode integration tests..."
    echo ""
    bats "$BATS_DIR/file_mode.bats" --timing
}

# Run database mode tests
run_db_tests() {
    if [[ -z "${DATABASE_URL_TEST:-}" ]]; then
        log_warn "DATABASE_URL_TEST not set - skipping database mode tests"
        log_info "To run database tests, set DATABASE_URL_TEST to a PostgreSQL connection string"
        return 0
    fi

    log_info "Running database mode integration tests..."
    echo ""
    bats "$BATS_DIR/database_mode.bats" --timing
}

# Run all tests
run_all_tests() {
    local failed=0

    run_file_tests || failed=1

    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""

    run_db_tests || failed=1

    return $failed
}

# Print usage
usage() {
    cat << EOF
Usage: $0 [OPTIONS] [TEST_SUITE]

Run my-context integration tests using BATS.

Test Suites:
  file      Run file mode tests only
  db        Run database mode tests only (requires DATABASE_URL_TEST)
  all       Run all tests (default)

Options:
  -h, --help     Show this help message
  -v, --verbose  Verbose output
  --tap          Output in TAP format (for CI)

Environment Variables:
  DATABASE_URL_TEST    PostgreSQL connection string for database tests
                       Example: postgres://user:pass@localhost:5432/test_db

Examples:
  $0                          # Run all tests
  $0 file                     # Run file mode tests only
  $0 db                       # Run database mode tests only
  DATABASE_URL_TEST="..." $0  # Run with database tests enabled

EOF
}

# Main
main() {
    local test_suite="all"
    local bats_opts=()

    while [[ $# -gt 0 ]]; do
        case "$1" in
            -h|--help)
                usage
                exit 0
                ;;
            -v|--verbose)
                bats_opts+=("--verbose-run")
                shift
                ;;
            --tap)
                bats_opts+=("--tap")
                shift
                ;;
            file|db|all)
                test_suite="$1"
                shift
                ;;
            *)
                log_error "Unknown option: $1"
                usage
                exit 1
                ;;
        esac
    done

    check_bats || exit 1
    build_binary

    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "  my-context Integration Tests"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""

    case "$test_suite" in
        file)
            run_file_tests
            ;;
        db)
            run_db_tests
            ;;
        all)
            run_all_tests
            ;;
    esac
}

main "$@"
