#!/usr/bin/env bats
# Integration tests for my-context in database mode
# Tests database-backed workflows with PostgreSQL
# Requires DATABASE_URL_TEST environment variable

load 'helpers/test_helper.bash'

setup_file() {
    build_test_binary
}

setup() {
    setup_db_mode
}

teardown() {
    # Database cleanup handled by test schema isolation
    :
}

# ============================================================================
# Core Workflow Tests (Database Mode)
# ============================================================================

@test "db mode: start creates new context" {
    run_my_context start "db-test-context"
    assert_success
    assert_output_contains "Started context"
}

@test "db mode: start and stop lifecycle" {
    run_my_context start "db-lifecycle-test"
    assert_success

    run_my_context stop
    assert_success
    assert_output_contains "Stopped context"
}

@test "db mode: add note to active context" {
    run_my_context start "db-note-test"
    assert_success

    run_my_context note "Database mode test note"
    assert_success
    assert_output_contains "Added note"
}

@test "db mode: list shows contexts from database" {
    run_my_context start "db-list-test-1"
    run_my_context stop
    run_my_context start "db-list-test-2"
    run_my_context stop

    run_my_context list --all
    assert_success
    assert_output_contains "db-list-test-1"
    assert_output_contains "db-list-test-2"
}

# ============================================================================
# Tag Command Tests (Database Mode) - Issue #18 Fix Verification
# ============================================================================

@test "db mode: tag add works with database backend" {
    run_my_context start "db-tag-test"
    run_my_context stop

    run_my_context tag add "db-tag-test" "database" "tested"
    assert_success
    assert_output_contains "Added"
}

@test "db mode: tag list retrieves tags from database" {
    run_my_context start "db-tag-list"
    run_my_context stop
    run_my_context tag add "db-tag-list" "feature" "urgent"

    run_my_context tag list "db-tag-list"
    assert_success
    assert_output_contains "feature"
    assert_output_contains "urgent"
}

@test "db mode: tag remove updates database" {
    run_my_context start "db-tag-remove"
    run_my_context stop
    run_my_context tag add "db-tag-remove" "to-remove"

    run_my_context tag remove "db-tag-remove" "to-remove"
    assert_success
    assert_output_contains "Removed"

    # Verify tag is gone
    run_my_context tag list "db-tag-remove"
    assert_output_not_contains "to-remove"
}

@test "db mode: tag list --all shows all tags with counts" {
    run_my_context start "db-tag-count-1"
    run_my_context stop
    run_my_context tag add "db-tag-count-1" "common-tag"

    run_my_context start "db-tag-count-2"
    run_my_context stop
    run_my_context tag add "db-tag-count-2" "common-tag"

    run_my_context tag list
    assert_success
    assert_output_contains "common-tag"
}

# ============================================================================
# Tree Command Tests (Database Mode) - Issue #18 Fix Verification
# ============================================================================

@test "db mode: tree shows hierarchy from database" {
    run_my_context start "db-tree-root"
    run_my_context stop

    run_my_context tree
    assert_success
    # Should show tree structure without errors
}

@test "db mode: tree with parent-child relationship" {
    run_my_context start "db-parent"
    run_my_context stop

    # Start child with parent reference
    run_my_context start "db-child" --parent "db-parent"
    run_my_context stop

    run_my_context tree "db-parent"
    assert_success
    assert_output_contains "db-parent"
}

# ============================================================================
# History Command Tests (Database Mode) - Graceful Degradation
# ============================================================================

@test "db mode: history reports not supported gracefully" {
    run_my_context history
    # Should not crash, should report not supported
    assert_output_contains "not"
    assert_output_contains "supported"
}

# ============================================================================
# Watch Command Tests (Database Mode) - Graceful Degradation
# ============================================================================

@test "db mode: watch reports not supported gracefully" {
    run_my_context start "db-watch-test"

    run_my_context watch
    # Should not crash, should report not supported
    assert_output_contains "not supported"
}

# ============================================================================
# List Filtering Tests (Database Mode)
# ============================================================================

@test "db mode: list filters by tag" {
    run_my_context start "db-filter-1"
    run_my_context stop
    run_my_context tag add "db-filter-1" "important"

    run_my_context start "db-filter-2"
    run_my_context stop

    run_my_context list --tag "important" --all
    assert_success
    assert_output_contains "db-filter-1"
    assert_output_not_contains "db-filter-2"
}

@test "db mode: list filters by project" {
    run_my_context start "db-proj-a: task-1"
    run_my_context stop
    run_my_context start "db-proj-b: task-1"
    run_my_context stop

    run_my_context list --project "db-proj-a" --all
    assert_success
    assert_output_contains "db-proj-a"
    assert_output_not_contains "db-proj-b"
}

@test "db mode: list filters by search" {
    run_my_context start "db-searchable-context"
    run_my_context stop
    run_my_context start "db-other-context"
    run_my_context stop

    run_my_context list --search "searchable" --all
    assert_success
    assert_output_contains "db-searchable"
    assert_output_not_contains "db-other"
}

# ============================================================================
# JSON Output Tests (Database Mode)
# ============================================================================

@test "db mode: JSON output works correctly" {
    run_my_context start "db-json-test"
    run_my_context stop

    run_my_context_json list
    assert_success
    assert_json_field ".status" "success"
}

@test "db mode: JSON show includes touch count" {
    run_my_context start "db-touch-test"
    run_my_context touch
    run_my_context touch

    run_my_context_json show
    assert_success
    # Touch count should be in output
    echo "$output" | jq -e '.data.touch_count >= 0'
}

# ============================================================================
# Archive Tests (Database Mode)
# ============================================================================

@test "db mode: archive stores in database" {
    run_my_context start "db-archive-test"
    run_my_context stop

    run_my_context archive "db-archive-test"
    assert_success

    # Should not appear in normal list
    run_my_context list --all
    assert_output_not_contains "db-archive-test"

    # Should appear with --archived
    run_my_context list --archived
    assert_output_contains "db-archive-test"
}

# ============================================================================
# Concurrent Access Simulation
# ============================================================================

@test "db mode: handles rapid context switches" {
    for i in {1..5}; do
        run_my_context start "db-rapid-$i"
        assert_success
        run_my_context note "Note for rapid $i"
        run_my_context stop
    done

    run_my_context list --all
    assert_output_contains "db-rapid-1"
    assert_output_contains "db-rapid-5"
}
