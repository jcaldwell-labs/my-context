#!/usr/bin/env bash
# BATS Test Helper for my-context integration tests
# Provides setup/teardown and assertion helpers for isolated testing

# Project root (where go.mod lives)
BATS_TEST_PROJECT_ROOT="$(cd "$(dirname "${BATS_TEST_FILENAME}")/../.." && pwd)"

# Binary path (built before tests)
MY_CONTEXT_BIN="${BATS_TEST_PROJECT_ROOT}/bin/my-context-test"

# Build the binary once for all tests
build_test_binary() {
    if [[ ! -f "$MY_CONTEXT_BIN" ]]; then
        echo "# Building test binary..." >&3
        go build -o "$MY_CONTEXT_BIN" "${BATS_TEST_PROJECT_ROOT}/cmd/my-context/" || {
            echo "# Failed to build test binary" >&3
            return 1
        }
    fi
}

# Setup isolated file-based test environment
setup_file_mode() {
    export MY_CONTEXT_TEST_DIR="$(mktemp -d)"
    export MY_CONTEXT_HOME="$MY_CONTEXT_TEST_DIR"
    unset DATABASE_URL
}

# Setup isolated database mode test environment
# Requires DATABASE_URL to be set externally (CI or local postgres)
setup_db_mode() {
    if [[ -z "$DATABASE_URL_TEST" ]]; then
        skip "DATABASE_URL_TEST not set - skipping database mode tests"
    fi
    export MY_CONTEXT_HOME="db"
    export DATABASE_URL="$DATABASE_URL_TEST"

    # Create unique schema for test isolation
    export MY_CONTEXT_TEST_SCHEMA="test_$(date +%s)_$$"
}

# Cleanup test environment
cleanup_test_env() {
    if [[ -n "$MY_CONTEXT_TEST_DIR" && -d "$MY_CONTEXT_TEST_DIR" ]]; then
        rm -rf "$MY_CONTEXT_TEST_DIR"
    fi
}

# Run my-context command and capture output
# Usage: run_my_context start "Test context"
run_my_context() {
    run "$MY_CONTEXT_BIN" "$@"
}

# Run my-context with JSON output
run_my_context_json() {
    run "$MY_CONTEXT_BIN" --json "$@"
}

# Assert output contains substring
assert_output_contains() {
    local expected="$1"
    if [[ "$output" != *"$expected"* ]]; then
        echo "Expected output to contain: $expected"
        echo "Actual output: $output"
        return 1
    fi
}

# Assert output does NOT contain substring
assert_output_not_contains() {
    local unexpected="$1"
    if [[ "$output" == *"$unexpected"* ]]; then
        echo "Expected output NOT to contain: $unexpected"
        echo "Actual output: $output"
        return 1
    fi
}

# Assert JSON field equals value
# Usage: assert_json_field ".data.context" "test-context"
assert_json_field() {
    local field="$1"
    local expected="$2"
    local actual
    actual=$(echo "$output" | jq -r "$field")
    if [[ "$actual" != "$expected" ]]; then
        echo "Expected $field = '$expected', got '$actual'"
        echo "Full output: $output"
        return 1
    fi
}

# Assert command succeeded (exit code 0)
assert_success() {
    if [[ "$status" -ne 0 ]]; then
        echo "Expected success (exit 0), got exit $status"
        echo "Output: $output"
        return 1
    fi
}

# Assert command failed (exit code != 0)
assert_failure() {
    if [[ "$status" -eq 0 ]]; then
        echo "Expected failure (exit != 0), got exit 0"
        echo "Output: $output"
        return 1
    fi
}

# Count contexts in list output
count_contexts() {
    run_my_context_json list --all
    echo "$output" | jq -r '.data.contexts | length'
}

# Wait for condition with timeout
wait_for() {
    local timeout="$1"
    local condition="$2"
    local interval="${3:-1}"
    local elapsed=0

    while [[ $elapsed -lt $timeout ]]; do
        if eval "$condition"; then
            return 0
        fi
        sleep "$interval"
        elapsed=$((elapsed + interval))
    done
    return 1
}
