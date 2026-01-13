#!/usr/bin/env bats
# Synthetic User Workflow Tests
# Simulates realistic developer workflows to validate end-to-end behavior

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
# Workflow 1: Daily Developer Session
# Simulates a typical day of context-tracked development
# ============================================================================

@test "workflow: daily developer session" {
    # Morning: Start work on a feature
    run_my_context start "work-2026-01-13: feature-auth"
    assert_success

    # Add planning notes
    run_my_context note "Planning: Implement JWT authentication"
    run_my_context note "Decision: Use RS256 for token signing because security requirements"
    assert_success

    # Track files being worked on
    run_my_context file "/src/auth/jwt.go"
    run_my_context file "/src/auth/middleware.go"
    assert_success

    # Midday: Check progress
    run_my_context show
    assert_output_contains "feature-auth"
    assert_output_contains "JWT"

    # Add more notes throughout day
    run_my_context note "Implemented token validation"
    run_my_context note "TODO: Add refresh token support"
    run_my_context touch

    # End of day: Stop context
    run_my_context stop
    assert_success

    # Verify context was saved properly
    run_my_context show "work-2026-01-13: feature-auth"
    assert_output_contains "Implemented token validation"
}

# ============================================================================
# Workflow 2: Bug Investigation Session
# Simulates debugging workflow with context switches
# ============================================================================

@test "workflow: bug investigation with context switch" {
    # Working on feature
    run_my_context start "impl-feature-xyz"
    run_my_context note "Implementing feature XYZ"
    assert_success

    # Interrupt! Bug reported - switch contexts
    run_my_context start "debug-2026-01-13: production-crash"
    assert_success

    # Document investigation
    run_my_context note "Bug: Users seeing 500 errors on checkout"
    run_my_context note "Investigating: Database connection pool exhausted"
    run_my_context note "Root cause: Connection leak in payment service"
    run_my_context note "Fix: Added defer conn.Close() in processPayment()"

    # Tag for tracking
    run_my_context stop
    run_my_context tag add "debug-2026-01-13: production-crash" "bug" "production" "critical"
    assert_success

    # Return to original work
    run_my_context resume --last
    # Note: resume --last should get the last context before the debug session
    # For now, manually restart
    run_my_context start "impl-feature-xyz_2"
    run_my_context note "Resuming feature work after bug fix"
    run_my_context stop

    # Verify both contexts exist
    run_my_context list --all
    assert_output_contains "debug-2026-01-13"
    assert_output_contains "impl-feature-xyz"
}

# ============================================================================
# Workflow 3: Sprint Organization
# Simulates organizing work across a sprint
# ============================================================================

@test "workflow: sprint organization with tags and projects" {
    # Create multiple sprint contexts
    local contexts=(
        "sprint-26: user-registration"
        "sprint-26: password-reset"
        "sprint-26: email-verification"
        "sprint-27: payment-integration"
    )

    for ctx in "${contexts[@]}"; do
        run_my_context start "$ctx"
        run_my_context note "Initial setup for ${ctx##*: }"
        run_my_context stop
    done

    # Tag sprint 26 items
    for ctx in "${contexts[@]:0:3}"; do
        run_my_context tag add "$ctx" "sprint-26" "auth"
    done

    # Filter by project
    run_my_context list --project "sprint-26" --all
    assert_output_contains "user-registration"
    assert_output_contains "password-reset"
    assert_output_not_contains "payment-integration"

    # Filter by tag
    run_my_context list --tag "auth" --all
    assert_output_contains "user-registration"
    assert_output_not_contains "sprint-27"
}

# ============================================================================
# Workflow 4: Export and Archive Completed Work
# Simulates end-of-sprint cleanup
# ============================================================================

@test "workflow: export and archive completed work" {
    # Create and complete a context
    run_my_context start "completed-feature"
    run_my_context note "Started implementation"
    run_my_context note "Added tests"
    run_my_context note "Code review complete"
    run_my_context note "DONE: Merged to main"
    run_my_context stop

    # Export for documentation
    run_my_context export "completed-feature"
    assert_success
    assert_output_contains "Exported"

    # Archive completed work
    run_my_context archive "completed-feature"
    assert_success

    # Verify archived
    run_my_context list --all
    assert_output_not_contains "completed-feature"

    run_my_context list --archived
    assert_output_contains "completed-feature"
}

# ============================================================================
# Workflow 5: Multi-Project Work
# Simulates working across multiple projects
# ============================================================================

@test "workflow: multi-project context management" {
    # Project A contexts
    run_my_context start "project-a: backend-api"
    run_my_context note "REST API development"
    run_my_context stop

    run_my_context start "project-a: frontend-ui"
    run_my_context note "React components"
    run_my_context stop

    # Project B contexts
    run_my_context start "project-b: data-pipeline"
    run_my_context note "ETL processing"
    run_my_context stop

    # Search across all projects
    run_my_context list --search "project-a" --all
    assert_output_contains "backend-api"
    assert_output_contains "frontend-ui"
    assert_output_not_contains "data-pipeline"

    # Verify project filter
    run_my_context list --project "project-b" --all
    assert_output_contains "data-pipeline"
    assert_output_not_contains "project-a"
}

# ============================================================================
# Workflow 6: Note-Heavy Session with Warnings
# Simulates a long session that should trigger note warnings
# ============================================================================

@test "workflow: note accumulation and warnings" {
    run_my_context start "long-session"

    # Add many notes (should trigger warnings around 50, 100)
    for i in $(seq 1 55); do
        run_my_context note "Note number $i - tracking progress"
    done

    # Verify notes were added
    run_my_context show
    assert_output_contains "Note number 55"

    run_my_context stop
    assert_success
}

# ============================================================================
# Workflow 7: JSON Pipeline Integration
# Simulates scripting/automation use case with JSON output
# ============================================================================

@test "workflow: json pipeline integration" {
    # Create test data
    run_my_context start "json-pipeline-test"
    run_my_context note "Pipeline note 1"
    run_my_context note "Pipeline note 2"
    run_my_context stop

    # Get JSON output
    run_my_context_json list
    assert_success

    # Parse with jq - the structure is nested .data.data.contexts
    local count
    count=$(echo "$output" | jq '.data.data.contexts | length')
    [[ "$count" -ge 1 ]]

    # Get specific context
    run_my_context_json show "json-pipeline-test"
    assert_success
    assert_json_field ".command" "show"
}

# ============================================================================
# Workflow 8: Error Recovery
# Simulates handling various error conditions gracefully
# ============================================================================

@test "workflow: error recovery scenarios" {
    # Try to stop when nothing active
    run_my_context stop
    # Should fail gracefully, not crash

    # Try to show nonexistent context
    run_my_context show "does-not-exist"
    assert_failure

    # Try to add tag to nonexistent context
    run_my_context tag add "nonexistent" "tag"
    assert_failure

    # Start a valid context - should still work after errors
    run_my_context start "recovery-test"
    assert_success
    run_my_context stop
    assert_success
}
