# My-Context Progressive User Guide
## From Solo Work to Enterprise Coordination

> A journey from simple context tracking to sophisticated team collaboration
>
> **Designed for**: Backend devs, frontend devs, QA engineers, tech leads, product managers, scrum masters, and enterprise teams

---

## How to Use This Guide

This guide is structured in **6 progressive stages**, each building on the previous one. You don't need to read all of them:

- **Start here** if you're new to my-context
- **Jump to your stage** based on your current workflow needs
- **Each stage is ~5-10 minutes** to read and try

| Stage | You Are | Time | What You'll Do |
|-------|---------|------|---|
| **1** | Solo contributor | 5 min | Track one work task |
| **2** | Juggling projects | 8 min | Switch between multiple projects |
| **3** | Organized developer | 10 min | Group and manage many contexts |
| **4** | Team collaborator | 10 min | Share decisions with teammates |
| **5** | Automation enthusiast | 12 min | Automate your workflow |
| **6** | Enterprise team | 15 min | Coordinate across multiple teams |

---

# Stage 1: The Solo Contributor
## Your First Context (5 minutes)

You're sitting down to work on one task. You want to track what you do, take notes about decisions, and remember which files you touched.

### The Problem
- Where did you spend your time today?
- What decisions did you make?
- Which files did you modify?
- How long did it actually take?

### The Solution: Start Your First Context

**Step 1: Start a context when you begin work**

```bash
my-context start "Implement user login feature"
```

You'll see:
```
Context 'Implement user login feature' started
```

That's it. My-context is now tracking your work.

**Step 2: Take notes as you work**

When you make an important decision or discovery, jot it down:

```bash
my-context note "Using bcrypt for password hashing - 10 rounds"
my-context note "Added email validation with regex"
my-context note "Blocked users get 15-min lockout"
```

**Step 3: Associate files you modified**

As you edit files, tell my-context:

```bash
my-context file src/auth/login.go
my-context file src/auth/password.go
my-context file tests/auth_test.go
```

**Step 4: Check your progress anytime**

```bash
my-context show
```

Output:
```
Active Context: Implement user login feature
Status: Active (started 47 minutes ago)
Duration: 47 minutes

Notes:
  2025-10-22 09:15:32 - Using bcrypt for password hashing - 10 rounds
  2025-10-22 09:28:15 - Added email validation with regex
  2025-10-22 09:45:22 - Blocked users get 15-min lockout

Files:
  src/auth/login.go
  src/auth/password.go
  tests/auth_test.go
```

**Step 5: Stop when you're done**

```bash
my-context stop
```

### Why This Matters

- **Future you**: In 3 months, you'll remember why you used bcrypt instead of argon2
- **Your manager**: You have a clear record of work completed
- **Code reviews**: Context shows which files and decisions go together
- **Personal growth**: Reviewing decisions helps you improve

### Real-World Example: Frontend Developer

```bash
# Monday morning
my-context start "fe: Fix mobile responsive nav"
my-context note "Breakpoint: 768px for tablet"
my-context note "Using flexbox for menu layout"
my-context file src/components/Navigation.tsx
my-context file src/styles/responsive.css

# Check progress
my-context show

# End of day
my-context stop
```

---

# Stage 2: Multi-Project Developer
## Switching Between Projects (8 minutes)

You work on multiple projects. Sometimes you switch between them during the day. You need to keep your work organized by project.

### The Problem
- Working on 3 different codebases
- Getting context mixed up between projects
- Hard to find work from a specific project later
- Context names get confusing (which "API work" was which?)

### The Solution: Use Project Grouping

My-context has a simple solution: **label your contexts with a `--project` flag**.

**Step 1: Start contexts with project names**

```bash
# Backend work
my-context start "Implement JWT auth" --project backend-api
my-context note "Using RS256 signing with kid header"
my-context file internal/auth/jwt.go
my-context stop

# Frontend work (starts new context, auto-stops the previous one)
my-context start "Build login form" --project frontend-web
my-context note "Using React Hook Form for validation"
my-context file src/pages/LoginPage.tsx
my-context stop

# Mobile work
my-context start "Fix payment flow" --project mobile-app
my-context note "Debug: PaymentIntent not returning clientSecret"
my-context file src/screens/CheckoutScreen.tsx
my-context stop
```

**Step 2: Find work by project**

Now you can filter by project:

```bash
# See all backend work
my-context list --project backend-api

# See all frontend work
my-context list --project frontend-web
```

Output:
```
Contexts for project: frontend-web

1. Build login form (stopped, 1 hour 23 minutes)
   - React Hook Form for validation

2. Fix navigation responsive (stopped, 2 hours)
   - Mobile-first approach
```

**Step 3: Search across projects**

```bash
# Find anything mentioning "payment" across ALL projects
my-context list --search "payment"

# Find anything from frontend
my-context list --project frontend-web

# Combine: Find "bug" contexts in backend
my-context list --project backend-api --search "bug"
```

### Tips for Project Naming

| Project Type | Naming Pattern | Example |
|---|---|---|
| **Microservice** | service name | `--project payment-service` |
| **Feature branch** | feature name | `--project oauth-integration` |
| **Client work** | client code | `--project acme-corp` |
| **Team initiative** | initiative | `--project q4-performance` |

### Real-World Example: Consultant Working Multiple Clients

```bash
# Client A: Debugging database issue
my-context start "Query optimization investigation" --project acme-corp
my-context note "Found N+1 query in user reports"
my-context note "Solution: Add eager loading for relationships"
my-context file internal/queries/users.go
my-context stop

# Client B: Code review
my-context start "Review authentication PR" --project techstartup-inc
my-context note "LGTM - good error handling"
my-context file docs/review-comments.md
my-context stop

# Later: see all work for client A
my-context list --project acme-corp
```

---

# Stage 3: Organized Developer
## Managing Many Contexts (10 minutes)

You're tracking 20+ contexts. Your list is getting messy. You need to clean up old completed work and keep only active/relevant contexts visible.

### The Problem
- `my-context list` shows 50+ old contexts
- Hard to find current work
- Can't tell which contexts are done vs. in-progress
- Taking notes becomes a mess

### The Solution: Archive and Organize

**Concept 1: Stopped vs. Active**

Contexts have two states:
- **Active**: Currently being worked on (only 1 at a time)
- **Stopped**: Work is done or paused

```bash
# See ONLY your active context
my-context list --active-only

# Default list shows 10 most recent (stopped + active)
my-context list

# See all contexts
my-context list --all
```

**Concept 2: Archive Old Work**

Archive completed contexts to keep things clean:

```bash
# Archive a specific context
my-context archive "Old feature from last sprint"

# After archiving, it won't show in normal list
my-context list  # Won't see archived contexts

# But you can still find it
my-context list --archived  # See archived contexts
my-context list --search "old feature"  # Search finds it too
```

**Concept 3: Sprint-Based Organization**

Organize contexts by sprint for easy management:

```bash
# Sprint 5 work
my-context start "Sprint 5: Implement search" --project payments
my-context start "Sprint 5: Test payment flow" --project payments

# At end of sprint
my-context archive "Sprint 5: Implement search"
my-context archive "Sprint 5: Test payment flow"

# Next sprint
my-context start "Sprint 6: Optimize queries" --project payments
```

### Workflow: End-of-Sprint Cleanup

```bash
# Sunday evening: review what you did
my-context list --all

# Archive all Sprint 5 work
for context in $(my-context list --search "Sprint 5" --json | jq -r '.data.contexts[].name'); do
  my-context archive "$context"
done

# Verify Sprint 5 is archived
my-context list --archived --search "Sprint 5"

# Start fresh for Sprint 6
my-context list  # Now shows only recent/active
```

### Real-World Example: Backend Developer at Mid-Size Company

```bash
# During Sprint 5
my-context start "Sprint 5: Payment processor integration" --project payments
my-context note "Chose Stripe over PayPal - better API"
my-context file internal/payments/stripe.go
my-context stop

my-context start "Sprint 5: Add refund handler" --project payments
my-context note "Idempotency key for safety"
my-context file internal/payments/refunds.go
my-context stop

my-context start "Sprint 5: Write API tests" --project payments
my-context note "Test coverage: 87% for payment module"
my-context file tests/payments_test.go
my-context stop

# Friday before Sprint 6: cleanup
my-context list --search "Sprint 5"  # See Sprint 5 work

# Archive all Sprint 5 contexts
my-context archive "Sprint 5: Payment processor integration"
my-context archive "Sprint 5: Add refund handler"
my-context archive "Sprint 5: Write API tests"

# Monday: start fresh with Sprint 6
my-context list  # Clean slate!
```

---

# Stage 4: Team Collaborator
## Sharing Work and Decisions (10 minutes)

Your team needs to understand your work: what you built, why you chose certain approaches, what problems you solved. You want to share this without endless meetings or Slack messages.

### The Problem
- Team members don't know about your architectural decisions
- Decisions get made in multiple places (Slack, emails, meetings, code comments)
- New team members don't understand WHY code is written a certain way
- Hard to trace decisions back to their context
- Documentation is outdated or incomplete

### The Solution: Export Contexts

**Step 1: Export your work to markdown**

```bash
# Export a single context
my-context export "Implement caching layer" --to docs/sprint-5-caching.md
```

This creates a markdown file with:
- What you worked on
- All your decision notes
- Files you modified
- How long it took
- Complete audit trail

**Step 2: Share with team**

The exported markdown is perfect for:
- Pull request descriptions
- Architecture decision records (ADRs)
- Sprint retrospectives
- Onboarding new developers
- Code review context

**Step 3: Export multiple contexts at once**

```bash
# Export all Sprint 5 work to a folder
my-context export --all --to sprint-5-summary/ --search "Sprint 5"

# Creates folder with:
#   sprint-5-summary/
#   ├── context-1.md
#   ├── context-2.md
#   └── context-3.md
```

### Exported File Example

When you export a context, you get:

```markdown
# Context: Implement caching layer

**Status**: Stopped
**Duration**: 2 hours 45 minutes
**Started**: 2025-10-22 09:00:00
**Ended**: 2025-10-22 11:45:00

## Notes

- 2025-10-22 09:15:32: Using Redis for cache - 1hr TTL
- 2025-10-22 09:45:22: Pattern: Check cache → compute → store → return
- 2025-10-22 10:30:15: Added cache invalidation on user update
- 2025-10-22 11:15:00: Wrote unit tests - 45 min
- 2025-10-22 11:40:00: Ready for code review

## Files Modified

- internal/cache/redis.go
- internal/cache/middleware.go
- tests/cache_test.go
```

### Real-World Example: Tech Lead Sharing Architecture Decision

```bash
# During implementation
my-context start "Design database schema for analytics"
my-context note "Evaluated: MongoDB vs PostgreSQL"
my-context note "Decision: PostgreSQL with JSONB for flexibility"
my-context note "Reason: Complex queries needed, structured data, ACID guarantees"
my-context note "Created schema with: events, event_metadata, aggregations"
my-context file internal/schema/analytics.sql
my-context stop

# Later: export for team review
my-context export "Design database schema for analytics" \
  --to docs/adr/0005-analytics-database.md

# Share in PR
# "Context exported to docs/adr/0005-analytics-database.md"
# Team reads markdown file, understands the reasoning
```

### Exporting for Different Audiences

```bash
# For product manager: export decision notes only
my-context export "Sprint 5 work" --to product-summary.md

# For new developer: export everything
my-context export "Caching implementation" --to onboarding/caching.md

# For retrospective: export all Sprint 5
my-context export --search "Sprint 5" --to sprint-5/

# For compliance: export as JSON
my-context export "Security implementation" --json --to audit/security.json
```

---

# Stage 5: Automation Enthusiast
## Smart Workflows (12 minutes)

You want my-context to work FOR you, not just TRACK you. Use signals and watches to:
- Get notified when you finish something
- Auto-cleanup at end of day
- Remind yourself about forgotten contexts
- Trigger deployment pipelines

### Concept 1: Signals (Simple Events)

A **signal** is just a timestamped marker. Create signals to trigger events:

```bash
# Create a signal
my-context signal create "ready-for-deployment"

# List all signals
my-context signal list

# Clear a signal when done
my-context signal clear "ready-for-deployment"
```

### Concept 2: Watches (Monitor and React)

A **watch** monitors your active context and runs a command when it changes:

```bash
# Watch for notes containing "done" or "complete"
my-context watch --new-notes --pattern="done|complete" \
  --exec="notify-send 'Task complete!'"

# Watch for any note and print it
my-context watch --new-notes

# Watch with a custom polling interval (default 5 sec)
my-context watch --interval 10s
```

### Pattern 1: Completion Notifications

Get notified when you note a milestone:

```bash
# In one terminal: watch for completion
my-context watch --new-notes --pattern="complete|done|finished" \
  --exec="notify-send 'Work milestone reached!'"

# In another terminal: work and note completion
my-context start "Code review for payment PR"
# ... review code ...
my-context note "Code review complete - 2 issues found"
# → Watch detects "complete" → notification appears!
```

### Pattern 2: End-of-Day Automation

Export and archive your work automatically:

```bash
# Create a shell script: ~/.local/bin/eod.sh
#!/bin/bash
export MY_CONTEXT_HOME=~/.my-context

# Export today's work
TODAY=$(date +%Y-%m-%d)
my-context list --json | \
  jq -r '.data.contexts[].name' | while read ctx; do
    my-context export "$ctx" --to ~/work-logs/$TODAY/$ctx.md
  done

# Stop any active context
my-context stop

echo "✅ Work exported to ~/work-logs/$TODAY/"
echo "✅ Active context stopped"

# Create signal for completion
my-context signal create "eod-$TODAY"
```

Run it daily:

```bash
# Add to ~/.bashrc or manual end-of-day ritual
eod.sh
```

### Pattern 3: Work-in-Progress Warnings

Warn yourself if a context gets too large:

```bash
# Watch: alert if more than 10 notes in active context
my-context watch --exec='
  COUNT=$(my-context show --json | jq .data.context.notes | wc -l)
  if [ "$COUNT" -gt 10 ]; then
    notify-send "WIP too large!" "Context has $COUNT notes - consider stopping"
  fi
'
```

### Pattern 4: Project Change Notifications

Get notified when you switch projects:

```bash
# Create a signal when changing project
alias cx="my-context"

start-project() {
  PROJECT="$1"
  CONTEXT="$2"

  # Stop previous and signal the change
  my-context signal create "project-changed-$(date +%s)"
  my-context start "$CONTEXT" --project "$PROJECT"

  notify-send "Switched to $PROJECT"
}

# Usage:
start-project "frontend-web" "Fix login button"
```

### Real-World Example: Backend Developer's Daily Workflow

```bash
# Morning: start your sprint
my-context start "Daily standup notes" --project payments
my-context note "Assigned: Payment processor refactor"
my-context stop

my-context start "Implement Stripe API client" --project payments
# [work for 2 hours]
my-context note "API client complete - ready for review"

# Watch alerts you to the "complete" note → lunch break
# You're notified: "Work milestone reached!"

# Afternoon: different task
my-context start "Review payment PRs" --project payments

# End of day: automation kicks in
eod.sh  # Exports all work, stops context, creates signal
# Creates: ~/work-logs/2025-10-22/
# Contains all your context exports

# Tomorrow: see yesterday's work
cat ~/work-logs/2025-10-22/
ls ~/work-logs/2025-10-22/  # See all contexts from yesterday
```

---

# Stage 6: Enterprise Team Coordination
## Multi-Team Synchronization (15 minutes)

You're coordinating work across multiple teams. You need to:
- Signal when deployments are ready
- Wait for other teams to complete work
- Coordinate across timezones
- Track shared dependencies

### Concept 1: Cross-Team Signals

Signals can coordinate between teams:

```bash
# Backend team: signals that API is ready
my-context note "API deployed to staging"
my-context signal create "backend-api-ready-for-qa"

# QA team: watches for the signal, knows when to start testing
my-context signal wait "backend-api-ready-for-qa" --timeout 1h
# Blocks until backend team creates the signal
```

### Concept 2: Multi-Context Watches

Monitor multiple contexts across the team:

```bash
# Scrum master: watch all team members' contexts
# (requires shared context home or scanning script)

my-context watch --exec='
  ACTIVE=$(my-context list --active-only --json | jq .data.contexts | wc -l)
  if [ "$ACTIVE" -eq 0 ]; then
    notify-send "No active contexts - team idle?"
  fi
'
```

### Pattern 1: Deployment Coordination

```bash
# Backend team: implements and signals ready
my-context start "Release payment API v2"
my-context note "All tests pass - ready for deploy"
my-context signal create "payment-v2-backend-ready"
my-context stop

# QA team: waits for backend signal
my-context start "QA: Payment API v2"
my-context note "Waiting for backend-ready signal..."

# Watch for the signal (in separate terminal)
my-context signal wait "payment-v2-backend-ready" --timeout 2h

# When backend creates the signal, QA proceeds
my-context note "Backend ready! Starting regression tests"
```

### Pattern 2: Dependency Tracking

```bash
# Team A: working on feature
my-context start "Implement payment webhooks"
my-context note "Ready for Team B to integrate"
my-context signal create "team-a-webhooks-ready"

# Team B: waits for Team A
my-context start "Integration: Handle payment webhooks"
my-context note "Blocked: waiting for team-a-webhooks-ready"

# Watch for signal with timeout
if my-context signal wait "team-a-webhooks-ready" --timeout 8h; then
  my-context note "Team A ready! Starting integration"
else
  my-context note "Timeout waiting for Team A"
fi
```

### Pattern 3: Sprint Completion Tracking

```bash
# Each team: marks when their sprint work is done
my-context start "Sprint 5 Complete: Payment Module"
my-context note "All stories closed, tests passing"
my-context signal create "sprint-5-team-payments-complete"
my-context stop

my-context start "Sprint 5 Complete: Frontend"
my-context note "All stories closed, E2E tests passing"
my-context signal create "sprint-5-team-frontend-complete"
my-context stop

# Release manager: waits for all signals
TEAMS=("payments" "frontend" "mobile")
for team in "${TEAMS[@]}"; do
  echo "Waiting for $team..."
  my-context signal wait "sprint-5-team-$team-complete" --timeout 24h
done

echo "All teams complete! Ready for release."
```

### Real-World Example: Multi-Team Microservices Release

```bash
# --- TEAM A: Payment Service ---
my-context start "Sprint 5: Payment service"
my-context note "Implemented new payment processor"
my-context note "Database migrations: OK"
my-context note "Load tests: PASS"
my-context note "Code review: approved"
my-context stop

# Signal: ready for QA
my-context signal create "service-payments-ready-for-qa"

# --- TEAM B: API Gateway ---
my-context start "Integration: Payment service"

# Wait for payment service to be ready
echo "Waiting for payment service..."
my-context signal wait "service-payments-ready-for-qa" --timeout 4h

# When signal is created, proceed
my-context note "Payment service ready, integrating endpoints"
my-context note "Integration tests: PASS"
my-context stop

# Signal: ready for QA
my-context signal create "service-gateway-ready-for-qa"

# --- QA TEAM ---
my-context start "QA: Full system testing"

# Wait for all services
for service in "payments" "gateway" "auth"; do
  my-context signal wait "service-$service-ready-for-qa" --timeout 8h
done

my-context note "All services ready! Starting E2E tests"
my-context note "E2E tests: PASS"
my-context signal create "qa-approved-for-production"

# --- DEVOPS ---
my-context watch --signal "qa-approved-for-production" \
  --exec="./scripts/deploy-production.sh"

# When QA signals ready, auto-deploy!
```

### Enterprise Dashboard (Advanced)

```bash
# Script: show team status
#!/bin/bash

echo "=== SPRINT 5 STATUS ==="
echo ""

for team in payments gateway auth frontend qa; do
  STATUS=$(my-context list --project "sprint-5-$team" | head -1)
  READY=$(my-context signal list | grep "sprint-5-$team-ready" | wc -l)

  if [ "$READY" -gt 0 ]; then
    echo "✅ $team: READY"
  else
    echo "⏳ $team: IN PROGRESS"
    echo "   $STATUS"
  fi
done

echo ""
echo "Signals:"
my-context signal list
```

---

# Reference: Quick Command Guide

## By Role

### Backend Developer
```bash
# Start work
my-context start "Feature name" --project myservice

# Take notes about decisions
my-context note "Using database X because..."

# Track files
my-context file internal/package/file.go

# End work
my-context stop

# Export for PR description
my-context export "Feature name" --to pr-context.md
```

### Frontend Developer
```bash
# Track UI work
my-context start "Component: Login Form" --project frontend-web

# Note design decisions
my-context note "Using Tailwind for responsive design"
my-context note "Mobile-first breakpoint at 768px"

# Track modified files
my-context file src/components/LoginForm.tsx
my-context file src/styles/form.css

# View your work
my-context show

# Export for design review
my-context export "Component: Login Form" --to design-review.md
```

### QA Engineer
```bash
# Track test sessions
my-context start "Test: Payment flow" --project payments

# Note test findings
my-context note "Testing with 100 concurrent users"
my-context note "Found: Race condition in order processing"
my-context note "Bug report: #5847"

# Track test files
my-context file tests/e2e/payment_test.go

# Export test report
my-context export "Test: Payment flow" --to test-report.md
```

### Tech Lead
```bash
# Track architectural work
my-context start "Architecture: Caching layer" --project infrastructure

# Document decisions
my-context note "Evaluated: Redis vs Memcached"
my-context note "Decision: Redis - persistence needed"
my-context note "Replication: 3 nodes with failover"

# Track design documents
my-context file docs/caching-architecture.md

# Export for ADR (Architecture Decision Record)
my-context export "Architecture: Caching layer" --to adr-0005.md

# Share with team
# Push adr-0005.md to repository
```

### Scrum Master
```bash
# Track sprint planning
my-context start "Sprint 6 planning"
my-context note "Stories: 34 points"
my-context note "Team velocity: 32 points/sprint"
my-context stop

# Track mid-sprint check-ins
my-context list --search "Sprint 6"

# Export sprint summary
my-context export --all --to sprint-reports/sprint-6/ --search "Sprint 6"

# Review completion for retrospective
my-context list --archived --search "Sprint 6"
```

### Product Manager
```bash
# Track feature definition
my-context start "Feature: User preferences"
my-context note "User stories: 5 requirements"
my-context note "Acceptance criteria defined"
my-context note "Design review: approved"

# Get developer context on implementation
my-context list --search "preferences"  # See all related work

# Export for stakeholder update
my-context export "Feature: User preferences" --to stakeholder-update.md
```

---

# Common Workflows by Scenario

## Daily Workflow
```bash
# Morning
my-context start "Daily standup prep"
my-context note "Today's priorities: Payment API integration"
my-context stop

# During day
my-context start "Payment API integration" --project backend
# [work]
my-context stop

my-context start "Code review"
my-context note "Reviewed PR #123"
my-context stop

# End of day
my-context list  # Review what you did
my-context stop  # Ensure nothing is active
```

## Weekly Retrospective
```bash
# Collect all week's work
my-context export --all --to weekly-summary/ --search "week-oct-22"

# Review
cat weekly-summary/*.md

# Archive completed
my-context archive "Old context"
```

## Sprint Close
```bash
# See all sprint work
my-context list --search "Sprint 5" --all

# Export for retrospective
my-context export --all --to retrospective/sprint-5/ --search "Sprint 5"

# Archive all sprint contexts
my-context list --search "Sprint 5" --json | \
  jq -r '.data.contexts[].name' | \
  xargs -I {} my-context archive "{}"

# Verify archive
my-context list --archived --search "Sprint 5"
```

## Onboarding New Developer
```bash
# Show new developer company's work context structure
my-context list --all

# Show their team's work
my-context list --project their-team

# Export recent work for learning
my-context export --all --to onboarding/ --search "recent"

# Share folder with new developer
# Let them read past decisions and implementation details
```

---

# Tips & Tricks

## Alias Shortcuts
Add to `~/.bashrc` or `~/.zshrc`:

```bash
alias cx="my-context"
alias cxs="my-context start"
alias cxn="my-context note"
alias cxf="my-context file"
alias cxl="my-context list"
alias cxw="my-context show"
alias cxe="my-context export"
```

Then: `cxs "My task"` instead of `my-context start "My task"`

## Shell Prompt with Context
Add to `~/.bashrc`:

```bash
# Show active context in shell prompt
if command -v my-context &> /dev/null; then
  export PS1='[\$(my-context show --json 2>/dev/null | jq -r ".data.context.name // \"idle\"")]$ '
fi
```

## Git Hook Integration
Create `.git/hooks/post-commit`:

```bash
#!/bin/bash
if command -v my-context &> /dev/null; then
  COMMIT_MSG=$(git log -1 --pretty=%B)
  my-context note "Committed: $COMMIT_MSG"
fi
```

Make executable: `chmod +x .git/hooks/post-commit`

## Quick Note Format
Consistent note format for searchability:

```bash
my-context note "DECISION: Using Postgres over MongoDB"
my-context note "BUG: Found race condition in cache invalidation"
my-context note "QUESTION: Is this performance acceptable?"
my-context note "BLOCKER: Waiting for API docs from Team A"
```

Then search easily: `my-context list --search "DECISION"`

---

# Next Steps

- **Stage 1 done?** → Try Stage 2 (multi-project)
- **Want automation?** → Jump to Stage 5 (signals/watches)
- **Managing a team?** → Go to Stage 6 (coordination)
- **Need help?** → See TROUBLESHOOTING.md for your platform

## Official Tutorials

- **[TRIGGERS-TUTORIAL.md](TRIGGERS-TUTORIAL.md)** - Deep dive into signals and watches
- **[README.md](../README.md)** - Complete command reference
- **[TROUBLESHOOTING.md](TROUBLESHOOTING.md)** - Platform-specific help

---

# Summary

| Stage | Key Idea | Key Command |
|---|---|---|
| 1 | Track single task | `my-context start "task"` |
| 2 | Organize by project | `--project flag` |
| 3 | Clean up old work | `my-context archive` |
| 4 | Share decisions | `my-context export` |
| 5 | Automate workflows | `my-context watch` |
| 6 | Coordinate teams | `my-context signal` |

Start at Stage 1 and progress at your own pace. Each stage builds on the previous one, but you can skip ahead if you only need certain features.

Happy tracking! 🚀
