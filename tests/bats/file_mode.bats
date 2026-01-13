#!/usr/bin/env bats
# Integration tests for my-context in file-based mode
# Tests core workflows with isolated file system storage

load 'helpers/test_helper.bash'

setup_file() {
    build_test_binary
}

setup() {
    setup_file_mode
}

teardown() {
    cleanup_test_env
}

# ============================================================================
# Core Workflow Tests
# ============================================================================

@test "file mode: start creates new context" {
    run_my_context start "test-context"
    assert_success
    assert_output_contains "Started"
    assert_output_contains "test-context"
}

@test "file mode: start and stop lifecycle" {
    run_my_context start "lifecycle-test"
    assert_success

    run_my_context stop
    assert_success
    assert_output_contains "Stopped context"
}

@test "file mode: add note to active context" {
    run_my_context start "note-test"
    assert_success

    run_my_context note "This is a test note"
    assert_success
    assert_output_contains "Note added"
}

@test "file mode: show displays context details" {
    run_my_context start "show-test"
    run_my_context note "Test note 1"
    run_my_context note "Test note 2"

    run_my_context show
    assert_success
    assert_output_contains "show-test"
    assert_output_contains "Test note 1"
}

@test "file mode: list shows all contexts" {
    run_my_context start "list-test-1"
    run_my_context stop
    run_my_context start "list-test-2"
    run_my_context stop

    run_my_context list --all
    assert_success
    assert_output_contains "list-test-1"
    assert_output_contains "list-test-2"
}

# ============================================================================
# Tag Command Tests
# ============================================================================

@test "file mode: tag add adds tags to context" {
    run_my_context start "tag-test"
    run_my_context stop

    run_my_context tag add "tag-test" "bug" "urgent"
    assert_success
    assert_output_contains "Added"
    assert_output_contains "bug"
}

@test "file mode: tag list shows context tags" {
    run_my_context start "tag-list-test"
    run_my_context stop
    run_my_context tag add "tag-list-test" "feature" "priority"

    run_my_context tag list "tag-list-test"
    assert_success
    assert_output_contains "feature"
    assert_output_contains "priority"
}

@test "file mode: tag remove removes tags" {
    run_my_context start "tag-remove-test"
    run_my_context stop
    run_my_context tag add "tag-remove-test" "temp-tag"

    run_my_context tag remove "tag-remove-test" "temp-tag"
    assert_success
    assert_output_contains "Removed"
}

# ============================================================================
# Tree Command Tests
# ============================================================================

@test "file mode: tree shows context hierarchy" {
    run_my_context start "parent-ctx" --parent ""
    run_my_context stop

    run_my_context tree
    assert_success
    assert_output_contains "parent-ctx"
}

# ============================================================================
# History Command Tests
# ============================================================================

@test "file mode: history shows transitions" {
    run_my_context start "history-test-1"
    run_my_context stop
    run_my_context start "history-test-2"
    run_my_context stop

    run_my_context history
    assert_success
    assert_output_contains "history-test"
}

# ============================================================================
# Export Command Tests
# ============================================================================

@test "file mode: export generates markdown" {
    run_my_context start "export-test"
    run_my_context note "Export note 1"
    run_my_context note "Export note 2"
    run_my_context stop

    run_my_context export "export-test"
    assert_success
    # Export writes to file, confirm success message
    assert_output_contains "Exported"
}

# ============================================================================
# JSON Output Tests
# ============================================================================

@test "file mode: JSON output for list" {
    run_my_context start "json-test"
    run_my_context stop

    run_my_context_json list
    assert_success
    assert_json_field ".command" "list"
}

@test "file mode: JSON output for show" {
    run_my_context start "json-show-test"

    run_my_context_json show
    assert_success
    assert_json_field ".command" "show"
}

# ============================================================================
# Error Handling Tests
# ============================================================================

@test "file mode: stop with no active context fails gracefully" {
    run_my_context stop
    # Should fail but not crash - case insensitive check
    [[ "$status" -ne 0 ]] || assert_output_contains "No active context"
}

@test "file mode: show nonexistent context fails" {
    run_my_context show "nonexistent-context"
    assert_failure
}

# ============================================================================
# Project Filter Tests
# ============================================================================

@test "file mode: list filters by project" {
    run_my_context start "project-a: task-1"
    run_my_context stop
    run_my_context start "project-a: task-2"
    run_my_context stop
    run_my_context start "project-b: task-1"
    run_my_context stop

    run_my_context list --project "project-a" --all
    assert_success
    assert_output_contains "project-a"
    assert_output_not_contains "project-b"
}

# ============================================================================
# Archive and Delete Tests
# ============================================================================

@test "file mode: archive context" {
    run_my_context start "archive-test"
    run_my_context stop

    run_my_context archive "archive-test"
    assert_success
    assert_output_contains "Archived"

    # Should not appear in normal list
    run_my_context list --all
    assert_output_not_contains "archive-test"

    # Should appear with --archived flag
    run_my_context list --archived
    assert_output_contains "archive-test"
}
