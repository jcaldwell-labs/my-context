#!/usr/bin/env python3
"""
Generate Tutorial Context Homes

Creates all 15+ context homes with realistic data for the my-context visual tutorials.
Each context home represents a different role or scenario showing how my-context is used.
"""

import subprocess
import os
import sys
from pathlib import Path
from datetime import datetime, timedelta
import time

# Base directory for context homes
CONTEXT_HOMES_DIR = Path(__file__).parent / "context-homes"
MY_CONTEXT_BIN = "/home/be-dev-agent/.local/bin/my-context"

def run_my_context(context_home: Path, command: list):
    """Execute my-context command with specified context home"""
    env = os.environ.copy()
    env["MY_CONTEXT_HOME"] = str(context_home)
    result = subprocess.run(
        [MY_CONTEXT_BIN] + command,
        env=env,
        capture_output=True,
        text=True
    )
    if result.returncode != 0:
        print(f"❌ Error: {result.stderr}")
    return result

def create_context(context_home: Path, name: str, project: str = None, labels: str = None, created_by: str = None):
    """Start a new context"""
    cmd = ["start", name]
    if project:
        cmd.extend(["--project", project])
    if labels:
        cmd.extend(["--labels", labels])
    if created_by:
        cmd.extend(["--created-by", created_by])
    print(f"  Creating context: {name}")
    return run_my_context(context_home, cmd)

def add_note(context_home: Path, note: str):
    """Add a note to active context"""
    return run_my_context(context_home, ["note", note])

def add_file(context_home: Path, filepath: str):
    """Add file association to active context"""
    return run_my_context(context_home, ["file", filepath])

def stop_context(context_home: Path):
    """Stop the active context"""
    return run_my_context(context_home, ["stop"])

def create_signal(context_home: Path, signal_name: str):
    """Create a signal"""
    return run_my_context(context_home, ["signal", "create", signal_name])

# ============================================================================
# Tutorial 1: Backend Developer Solo (Alice)
# ============================================================================

def tutorial_01_backend_solo():
    """Tutorial 1: Backend Developer implementing payment retry logic"""
    print("\n📘 Tutorial 1: Backend Developer Solo (Alice)")
    context_home = CONTEXT_HOMES_DIR / "tutorial-01-backend-solo"

    # Start payment retry feature
    create_context(context_home, "payment-retry-logic", project="payment-service", created_by="alice")
    add_note(context_home, "DECISION: Exponential backoff strategy - 1s, 2s, 4s, 8s")
    add_note(context_home, "Using exponential backoff to handle transient failures")
    add_note(context_home, "Max retries: 3 attempts before marking payment as failed")
    add_note(context_home, "Using context.Context for cancellation support")
    add_note(context_home, "Database transaction rolled back on timeout")
    add_file(context_home, "internal/payments/retry.go")
    add_file(context_home, "internal/payments/backoff.go")
    add_file(context_home, "tests/payments/retry_test.go")
    stop_context(context_home)

    print("  ✅ Tutorial 1 complete")

# ============================================================================
# Tutorial 2: Frontend Developer Solo (Bob)
# ============================================================================

def tutorial_02_frontend_solo():
    """Tutorial 2: Frontend Developer implementing responsive checkout UI"""
    print("\n📘 Tutorial 2: Frontend Developer Solo (Bob)")
    context_home = CONTEXT_HOMES_DIR / "tutorial-02-frontend-solo"

    create_context(context_home, "checkout-ui-responsive", project="web-app", created_by="bob")
    add_note(context_home, "DECISION: Using CSS Grid for 3-column layout")
    add_note(context_home, "Mobile-first approach - 1 column on mobile, 3 on desktop")
    add_note(context_home, "A11Y: Added ARIA labels to all payment form fields")
    add_note(context_home, "A11Y: Keyboard navigation fully supported (Tab order tested)")
    add_note(context_home, "Using semantic HTML for better screen reader support")
    add_note(context_home, "Color contrast ratio: 4.5:1 (WCAG AA compliant)")
    add_file(context_home, "src/components/Checkout.tsx")
    add_file(context_home, "src/components/PaymentForm.tsx")
    add_file(context_home, "src/styles/checkout.css")
    add_file(context_home, "tests/components/Checkout.test.tsx")
    stop_context(context_home)

    print("  ✅ Tutorial 2 complete")

# ============================================================================
# Tutorial 3: QA Engineer Solo (Carol)
# ============================================================================

def tutorial_03_qa_solo():
    """Tutorial 3: QA Engineer testing payment flow"""
    print("\n📘 Tutorial 3: QA Engineer Solo (Carol)")
    context_home = CONTEXT_HOMES_DIR / "tutorial-03-qa-solo"

    create_context(context_home, "payment-flow-testing", project="qa-suite", created_by="carol")
    add_note(context_home, "Test scope: Payment flow across Chrome, Firefox, Safari, Edge")
    add_note(context_home, "✅ Test passed: Chrome 120 - payment successful")
    add_note(context_home, "✅ Test passed: Firefox 121 - payment successful")
    add_note(context_home, "✅ Test passed: Edge 120 - payment successful")
    add_note(context_home, "❌ BUG: Safari 16.2 - card validation fails on submit")
    add_note(context_home, "BUG DETAILS: JS error in form validation - RegEx incompatibility")
    add_note(context_home, "WORKAROUND: Need to update regex pattern for Safari compatibility")
    add_file(context_home, "tests/e2e/payment-flow.spec.js")
    add_file(context_home, "tests/fixtures/test-cards.json")
    add_file(context_home, "bug-reports/safari-validation-bug.md")
    stop_context(context_home)

    print("  ✅ Tutorial 3 complete")

# ============================================================================
# Tutorial 4: Multi-Project Consultant (Alice)
# ============================================================================

def tutorial_04_multi_project():
    """Tutorial 4: Consultant juggling 3 client projects"""
    print("\n📘 Tutorial 4: Multi-Project Consultant (Alice)")
    context_home = CONTEXT_HOMES_DIR / "tutorial-04-multi-project"

    # Client 1: ACME Corp
    create_context(context_home, "api-optimization", project="client-acme", created_by="alice")
    add_note(context_home, "Client context: E-commerce platform, high traffic")
    add_note(context_home, "Performance issue: API latency >500ms on product search")
    add_note(context_home, "SOLUTION: Added Redis caching layer for product catalog")
    add_file(context_home, "acme/internal/cache/redis.go")
    stop_context(context_home)

    # Client 2: TechCorp
    create_context(context_home, "database-migration", project="client-techcorp", created_by="alice")
    add_note(context_home, "Client context: B2B SaaS, migrating MySQL → PostgreSQL")
    add_note(context_home, "Migration strategy: Blue-green deployment with dual-write period")
    add_note(context_home, "Data validation: Compare row counts after migration")
    add_file(context_home, "techcorp/migrations/001_initial_schema.sql")
    stop_context(context_home)

    # Client 3: Startup
    create_context(context_home, "security-audit", project="client-startup", created_by="alice")
    add_note(context_home, "Client context: Fintech startup preparing for SOC2 audit")
    add_note(context_home, "FINDING: API keys stored in plaintext in config files")
    add_note(context_home, "RECOMMENDATION: Move to AWS Secrets Manager")
    add_note(context_home, "FINDING: No rate limiting on public API endpoints")
    add_note(context_home, "RECOMMENDATION: Implement token bucket algorithm")
    add_file(context_home, "startup/security/audit-report.md")
    stop_context(context_home)

    # More contexts showing switching
    create_context(context_home, "acme-payment-integration", project="client-acme", created_by="alice")
    add_note(context_home, "Integrating Stripe payment gateway")
    add_note(context_home, "Using Stripe API v2023-10-16")
    stop_context(context_home)

    create_context(context_home, "techcorp-ci-pipeline", project="client-techcorp", created_by="alice")
    add_note(context_home, "Setting up GitHub Actions for automated testing")
    add_note(context_home, "Pipeline stages: lint → test → build → deploy")
    stop_context(context_home)

    create_context(context_home, "startup-monitoring", project="client-startup", created_by="alice")
    add_note(context_home, "Setting up Datadog for application monitoring")
    add_note(context_home, "Alerts configured for: API latency, error rate, throughput")
    stop_context(context_home)

    print("  ✅ Tutorial 4 complete")

# ============================================================================
# Tutorial 5: Scrum Master Sprint Management (Dave)
# ============================================================================

def tutorial_05_scrum_master():
    """Tutorial 5: Scrum Master managing sprint"""
    print("\n📘 Tutorial 5: Scrum Master (Dave)")
    context_home = CONTEXT_HOMES_DIR / "tutorial-05-scrum-master"

    # Sprint 5 planning
    create_context(context_home, "sprint-5-planning", project="team-alpha", created_by="dave")
    add_note(context_home, "Sprint 5 goals: Payment integration + responsive UI")
    add_note(context_home, "Team capacity: 5 developers × 8 days = 40 dev-days")
    add_note(context_home, "Team velocity: 42 story points (3-sprint average: 38)")
    add_note(context_home, "Sprint commitment: 40 story points (conservative)")
    add_note(context_home, "BLOCKER: API dependency on Platform team (payment gateway)")
    add_note(context_home, "Mitigation: Daily sync with Platform team lead")
    add_file(context_home, "sprint-5/planning/sprint-backlog.md")
    add_file(context_home, "sprint-5/planning/capacity-plan.xlsx")
    stop_context(context_home)

    # Daily standup tracking
    create_context(context_home, "sprint-5-day-3-standup", project="team-alpha", created_by="dave")
    add_note(context_home, "Alice: Payment retry logic - 80% complete")
    add_note(context_home, "Bob: Checkout UI - blocked on design review")
    add_note(context_home, "Carol: E2E tests - Safari bug found")
    add_note(context_home, "ACTION ITEM: Schedule design review for Bob today")
    stop_context(context_home)

    # Old sprint to archive
    create_context(context_home, "sprint-4-retrospective", project="team-alpha", created_by="dave")
    add_note(context_home, "Sprint 4: Completed 38 / 40 story points")
    add_note(context_home, "What went well: Good collaboration, clear requirements")
    add_note(context_home, "What to improve: Earlier QA involvement, better estimation")
    stop_context(context_home)

    create_context(context_home, "sprint-4-planning", project="team-alpha", created_by="dave")
    add_note(context_home, "Sprint 4 goals: User authentication + profile management")
    stop_context(context_home)

    print("  ✅ Tutorial 5 complete")

# ============================================================================
# Tutorial 6: Team Handoff - Alice (Backend) to Bob (Frontend)
# ============================================================================

def tutorial_06_team_handoff():
    """Tutorial 6: Team collaboration with async context sharing"""
    print("\n📘 Tutorial 6: Team Handoff (Alice → Bob)")

    # Alice's side
    context_home_alice = CONTEXT_HOMES_DIR / "tutorial-06-team-alice"
    create_context(context_home_alice, "payment-api-v2", project="backend-services", created_by="alice")
    add_note(context_home_alice, "API SPEC: POST /api/v2/payments/process")
    add_note(context_home_alice, "Request body: { amount, currency, payment_method, customer_id }")
    add_note(context_home_alice, "Response: { payment_id, status, transaction_id }")
    add_note(context_home_alice, "DECISION: Using idempotency keys to prevent duplicate charges")
    add_note(context_home_alice, "Idempotency key header: X-Idempotency-Key (UUID)")
    add_note(context_home_alice, "Error handling: Return 400 for validation, 402 for payment failure")
    add_note(context_home_alice, "Rate limiting: 100 requests/minute per customer")
    add_file(context_home_alice, "internal/api/v2/payments.go")
    add_file(context_home_alice, "docs/api/payment-endpoint-spec.md")
    stop_context(context_home_alice)

    # Bob's side (references Alice's work)
    context_home_bob = CONTEXT_HOMES_DIR / "tutorial-06-team-bob"
    create_context(context_home_bob, "payment-ui-integration", project="web-app", created_by="bob")
    add_note(context_home_bob, "REF: Alice's payment API spec in docs/api/payment-endpoint-spec.md")
    add_note(context_home_bob, "Implementing UI integration with POST /api/v2/payments/process")
    add_note(context_home_bob, "Added UUID generation for X-Idempotency-Key header")
    add_note(context_home_bob, "Error handling: Display user-friendly message for 400/402 errors")
    add_note(context_home_bob, "Loading state: Show spinner during payment processing")
    add_file(context_home_bob, "src/services/payment-api.ts")
    add_file(context_home_bob, "src/components/PaymentProcessor.tsx")
    stop_context(context_home_bob)

    print("  ✅ Tutorial 6 complete")

# ============================================================================
# Tutorial 7: Signal Coordination (4 roles)
# ============================================================================

def tutorial_07_signal_coordination():
    """Tutorial 7: Real-time coordination using signals"""
    print("\n📘 Tutorial 7: Signal Coordination (Alice, Bob, Carol, Eve)")

    # Alice (Backend)
    context_home_alice = CONTEXT_HOMES_DIR / "tutorial-07-release-alice"
    create_context(context_home_alice, "payment-api-release", project="backend-services", created_by="alice")
    add_note(context_home_alice, "Release: Payment API v2.0 to staging")
    add_note(context_home_alice, "Endpoints ready: /process, /refund, /status")
    add_note(context_home_alice, "Database migrations applied successfully")
    add_note(context_home_alice, "Health check: ✅ All endpoints responding")
    create_signal(context_home_alice, "api-v2-staging-ready")
    stop_context(context_home_alice)

    # Bob (Frontend)
    context_home_bob = CONTEXT_HOMES_DIR / "tutorial-07-release-bob"
    create_context(context_home_bob, "frontend-integration", project="web-app", created_by="bob")
    add_note(context_home_bob, "Waiting for api-v2-staging-ready signal...")
    add_note(context_home_bob, "Signal received! Starting integration work")
    add_note(context_home_bob, "Integrated payment API v2.0 endpoints")
    add_note(context_home_bob, "Manual testing: Payment flow working end-to-end")
    add_note(context_home_bob, "Deployed to staging environment")
    create_signal(context_home_bob, "frontend-staging-ready")
    stop_context(context_home_bob)

    # Carol (QA)
    context_home_carol = CONTEXT_HOMES_DIR / "tutorial-07-release-carol"
    create_context(context_home_carol, "integration-testing", project="qa-suite", created_by="carol")
    add_note(context_home_carol, "Waiting for frontend-staging-ready signal...")
    add_note(context_home_carol, "Signal received! Starting E2E tests")
    add_note(context_home_carol, "✅ Test suite: payment-flow-e2e (15 tests passed)")
    add_note(context_home_carol, "✅ Test: Process payment with valid card")
    add_note(context_home_carol, "✅ Test: Handle invalid card gracefully")
    add_note(context_home_carol, "✅ Test: Refund processed payment")
    add_note(context_home_carol, "✅ Test: Check payment status")
    add_note(context_home_carol, "All tests passed - staging environment approved")
    create_signal(context_home_carol, "qa-approved-staging")
    stop_context(context_home_carol)

    # Eve (Product Owner)
    context_home_eve = CONTEXT_HOMES_DIR / "tutorial-07-release-eve"
    create_context(context_home_eve, "release-coordination", project="product", created_by="eve")
    add_note(context_home_eve, "Release: Payment v2.0 feature")
    add_note(context_home_eve, "Waiting for qa-approved-staging signal...")
    add_note(context_home_eve, "QA approval received!")
    add_note(context_home_eve, "DECISION: Release window - Friday 2pm PST")
    add_note(context_home_eve, "Stakeholder notification sent")
    add_note(context_home_eve, "Marketing: Blog post scheduled for Monday")
    add_note(context_home_eve, "Support team: Training completed on new payment flow")
    stop_context(context_home_eve)

    print("  ✅ Tutorial 7 complete")

# ============================================================================
# Tutorial 8: Agent Workflows (Alice + 3 Agents)
# ============================================================================

def tutorial_08_agent_workflows():
    """Tutorial 8: AI agents and automation as first-class context users"""
    print("\n📘 Tutorial 8: Agent Workflows (Alice + Claude + CI/CD + QA Bot)")

    # Alice (Human)
    context_home_alice = CONTEXT_HOMES_DIR / "tutorial-08-human-alice"
    create_context(context_home_alice, "oauth-integration", project="backend-services", created_by="alice", labels="feature,backend")
    add_note(context_home_alice, "Implementing OAuth 2.0 client flow")
    add_note(context_home_alice, "Providers: Google, GitHub")
    add_note(context_home_alice, "Using authorization code flow with PKCE")
    stop_context(context_home_alice)

    # Claude Code Agent
    context_home_claude = CONTEXT_HOMES_DIR / "tutorial-08-agent-claude"
    create_context(context_home_claude, "oauth-code-assistance", project="backend-services", created_by="claude-agent")
    add_note(context_home_claude, "Parent context: oauth-integration (Alice)")
    add_note(context_home_claude, "Generated OAuth client boilerplate code")
    add_note(context_home_claude, "DECISION: Using golang.org/x/oauth2 library (official Google package)")
    add_note(context_home_claude, "Implemented token storage with encryption")
    add_note(context_home_claude, "Added automatic token refresh logic")
    add_file(context_home_claude, "internal/auth/oauth_client.go")
    add_file(context_home_claude, "internal/auth/token_storage.go")
    create_signal(context_home_claude, "code-ready-for-review")
    stop_context(context_home_claude)

    # Resume Alice's context (reviewing agent work)
    create_context(context_home_alice, "oauth-integration-review", project="backend-services", created_by="alice")
    add_note(context_home_alice, "Reviewed Claude agent's generated code")
    add_note(context_home_alice, "Code quality: Excellent, follows Go best practices")
    add_note(context_home_alice, "Added error handling for network failures")
    add_note(context_home_alice, "Added logging for OAuth flow debugging")
    create_signal(context_home_alice, "feature-ready-for-ci")
    stop_context(context_home_alice)

    # CI/CD Agent
    context_home_cicd = CONTEXT_HOMES_DIR / "tutorial-08-agent-cicd"
    create_context(context_home_cicd, "build-oauth-feature", project="backend-services", created_by="cicd-agent")
    add_note(context_home_cicd, "Parent context: oauth-integration (Alice)")
    add_note(context_home_cicd, "Build #5598 triggered by commit abc123def456")
    add_note(context_home_cicd, "✅ Stage 1: Lint - No issues found")
    add_note(context_home_cicd, "✅ Stage 2: Unit tests - 127/127 passed (0.8s)")
    add_note(context_home_cicd, "✅ Stage 3: Integration tests - 43/43 passed (12.3s)")
    add_note(context_home_cicd, "✅ Stage 4: Code coverage - 94.2% (target: 90%)")
    add_note(context_home_cicd, "✅ Stage 5: Security scan - No vulnerabilities")
    add_note(context_home_cicd, "Build artifacts uploaded to S3")
    create_signal(context_home_cicd, "ci-build-passed")
    stop_context(context_home_cicd)

    # QA Automation Bot
    context_home_qa = CONTEXT_HOMES_DIR / "tutorial-08-agent-qa"
    create_context(context_home_qa, "e2e-oauth-flow", project="qa-automation", created_by="qa-bot")
    add_note(context_home_qa, "Parent context: oauth-integration (Alice)")
    add_note(context_home_qa, "Test suite: OAuth end-to-end flows")
    add_note(context_home_qa, "✅ Test: Authorization code flow (Google)")
    add_note(context_home_qa, "✅ Test: Authorization code flow (GitHub)")
    add_note(context_home_qa, "✅ Test: Token refresh on expiry")
    add_note(context_home_qa, "✅ Test: Handle invalid authorization code")
    add_note(context_home_qa, "✅ Test: Handle expired refresh token")
    add_note(context_home_qa, "Test duration: 45.2 seconds")
    add_note(context_home_qa, "All tests passed - feature ready for production")
    create_signal(context_home_qa, "qa-automated-passed")
    stop_context(context_home_qa)

    # Alice final review
    create_context(context_home_alice, "oauth-feature-complete", project="backend-services", created_by="alice")
    add_note(context_home_alice, "All agents completed successfully:")
    add_note(context_home_alice, "  ✅ Claude agent: Code generation")
    add_note(context_home_alice, "  ✅ CI/CD agent: Build and tests")
    add_note(context_home_alice, "  ✅ QA bot: E2E validation")
    add_note(context_home_alice, "Feature ready for production deployment")
    add_note(context_home_alice, "Deployment scheduled for: Tomorrow 10am PST")
    stop_context(context_home_alice)

    print("  ✅ Tutorial 8 complete")

# ============================================================================
# Main Execution
# ============================================================================

def main():
    """Generate all tutorial context homes"""
    print("=" * 70)
    print("MY-CONTEXT TUTORIAL CONTEXT GENERATION")
    print("=" * 70)
    print(f"Context homes directory: {CONTEXT_HOMES_DIR}")
    print(f"My-context binary: {MY_CONTEXT_BIN}")
    print()

    # Verify my-context is available
    if not Path(MY_CONTEXT_BIN).exists():
        print(f"❌ Error: my-context not found at {MY_CONTEXT_BIN}")
        sys.exit(1)

    # Generate all tutorials
    tutorial_01_backend_solo()
    tutorial_02_frontend_solo()
    tutorial_03_qa_solo()
    tutorial_04_multi_project()
    tutorial_05_scrum_master()
    tutorial_06_team_handoff()
    tutorial_07_signal_coordination()
    tutorial_08_agent_workflows()

    print()
    print("=" * 70)
    print("✅ ALL TUTORIAL CONTEXTS GENERATED SUCCESSFULLY")
    print("=" * 70)
    print()
    print(f"Context homes created in: {CONTEXT_HOMES_DIR}")
    print("Next step: Run export-tutorial-panels.py to generate HTML exports")

if __name__ == "__main__":
    main()
