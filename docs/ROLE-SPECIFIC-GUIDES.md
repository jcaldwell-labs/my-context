# My-Context Role-Specific Guides
## Tailored workflows for your job title

Quick guides for specific roles. Each guide shows:
- **The problem** you solve with my-context
- **Daily workflow** example
- **Must-know commands**
- **Tips specific to your role**

---

## Backend Developer / API Engineer

### Your Challenge
You maintain multiple microservices. You need to track:
- Which service you're working on
- Architectural decisions
- Performance improvements
- Bug fixes and hotfixes
- Integration points between services

### Daily Workflow

```bash
# Morning standup
my-context start "Daily standup"
my-context note "Today: Payment service refactor"
my-context stop

# Start feature work
my-context start "Add payment retry logic" --project payment-service
my-context note "Using exponential backoff: 1s, 2s, 4s, 8s"
my-context note "Max retries: 3 attempts"
my-context file internal/payments/retry.go
my-context file tests/payments/retry_test.go

# [2 hours of development]

# Note important decisions
my-context note "Using context.Context for cancellation"
my-context note "Database transaction rolled back on timeout"

# Check progress
my-context show

# Code review time - export context
my-context export "Add payment retry logic" --to pr-description.md
my-context stop

# Quick hotfix
my-context start "Hotfix: Payment timeout issue" --project payment-service
my-context note "Bug: Timeout not propagating to client"
my-context file internal/payments/client.go
my-context note "Quick fix: wrap client.Do() with timeout"
my-context stop

# End of day
my-context list --project payment-service  # See all payment work
```

### Your Key Commands

```bash
# Always tag by service
my-context start "Task description" --project service-name

# Document architectural decisions IMMEDIATELY
my-context note "DECISION: Using Postgres JSONB for flexibility"
my-context note "RATIONALE: Schema evolution without migrations"
my-context note "TRADE-OFF: Query complexity vs flexibility"

# Track exactly which code changed
my-context file internal/path/file.go tests/path/file_test.go

# Review before PR
my-context show  # See all your notes and decision rationale

# Export for code review context
my-context export "Feature name" --to feature-summary.md
```

### Tips

1. **Use decision prefixes**: Start notes with `DECISION:`, `BUG:`, `REFACTOR:`, etc.
   ```bash
   my-context note "DECISION: Using gRPC over REST for service communication"
   my-context note "BUG: Connection pool leak in idle state"
   my-context note "REFACTOR: Extract auth middleware to separate package"
   ```

2. **Track both happy path and edge cases**:
   ```bash
   my-context note "Happy path: Normal payment flow - 2ms latency"
   my-context note "Edge case: Retry after network failure - 5s total"
   ```

3. **Document trade-offs**:
   ```bash
   my-context note "Chose cache-aside over write-through"
   my-context note "Pro: Simple, no consistency burden"
   my-context note "Con: Cold cache on restart"
   ```

4. **Performance context**:
   ```bash
   my-context note "Before: Query took 500ms, scanned 100k rows"
   my-context note "After: Query takes 50ms, uses index, scans 10 rows"
   my-context note "Improvement: 10x faster"
   ```

---

## Frontend Developer / UI Engineer

### Your Challenge
You juggle components, styling, accessibility, and responsiveness. You need to track:
- UI component implementations
- Design system updates
- Responsive design decisions
- Browser compatibility fixes
- Performance optimizations

### Daily Workflow

```bash
# Morning standup
my-context start "Standup prep" --project frontend
my-context note "Focus: Login page responsive design"
my-context stop

# Start UI work
my-context start "Build responsive login form" --project frontend-web
my-context note "Using Tailwind CSS - mobile-first"
my-context note "Breakpoint: 768px for tablet layout"
my-context note "Accessibility: ARIA labels for form inputs"
my-context file src/components/LoginForm.tsx
my-context file src/styles/forms.css

# [1 hour of development]

# Test and note findings
my-context note "Tested: Chrome, Firefox, Safari"
my-context note "Mobile: iPhone 12, Pixel 5 - both responsive"
my-context note "Screen reader: VoiceOver reads all fields correctly"

# Design refinement
my-context note "Added loading state with spinner"
my-context note "Error handling: Display inline validation"
my-context note "Password field: Show/hide toggle"

# Export for design review
my-context export "Build responsive login form" --to design-review.md
my-context stop

# Accessibility audit task
my-context start "A11y audit: Accessibility check" --project frontend-web
my-context note "Lighthouse score: 92/100"
my-context note "Color contrast: All text WCAG AA compliant"
my-context note "Focus order: Tab navigation works left-to-right"
my-context file src/components/LoginForm.tsx
my-context stop
```

### Your Key Commands

```bash
# Always tag by project
my-context start "Feature description" --project frontend-web

# Document design decisions
my-context note "DESIGN: Using card-based layout"
my-context note "RATIONALE: Cards provide visual hierarchy"

# Track UI files
my-context file src/components/ComponentName.tsx
my-context file src/styles/component.css

# Test findings
my-context note "TESTED: Chrome 119, Firefox 120, Safari 17"
my-context note "ISSUE: Button text wraps on small screens - fixed"

# Export for review
my-context export "Component name" --to review.md
```

### Tips

1. **Browser/device testing notes**:
   ```bash
   my-context note "TESTED: Desktop (1920x1080), Tablet (768x1024), Mobile (375x812)"
   my-context note "Responsive: Works on 320px to 1920px widths"
   my-context note "ISSUE: Overflow on 320px screens - padding issue"
   ```

2. **Accessibility checks**:
   ```bash
   my-context note "A11y: Tested with VoiceOver (macOS)"
   my-context note "A11y: All images have alt text"
   my-context note "A11y: Color contrast ratio 4.5:1 (WCAG AA)"
   ```

3. **Design system consistency**:
   ```bash
   my-context note "Using design system: Buttons, spacing, colors"
   my-context note "Breaking change: Button API - 'variant' not 'size'"
   my-context note "Deprecated: Old ButtonGroup component"
   ```

4. **Performance observations**:
   ```bash
   my-context note "Bundle size: +15KB (icon set)"
   my-context note "Performance: First contentful paint 1.2s"
   my-context note "Optimization: Lazy load images below fold"
   ```

---

## QA Engineer / Test Automation Engineer

### Your Challenge
You ensure quality across:
- Functional testing
- Regression testing
- Performance testing
- Accessibility testing
- Load/stress testing

### Daily Workflow

```bash
# Start test session
my-context start "Regression test: Payment flow" --project payments

# Test cases and findings
my-context note "TEST: Valid payment with visa card"
my-context note "  Status: PASS - charged $99.99, receipt generated"

my-context note "TEST: Payment retry on network timeout"
my-context note "  Status: PASS - auto-retried, succeeded"

my-context note "TEST: Invalid card number"
my-context note "  Status: PASS - error displayed, no charge"

# Track bugs found
my-context note "BUG: Receipt not generating sometimes"
my-context note "  Steps: Pay via PayPal, quick refresh"
my-context note "  Expected: Receipt displays"
my-context note "  Actual: Receipt blank"
my-context note "  Severity: High - affected 5% of users"
my-context note "  Bug ID: #5847"

# Environment and tools
my-context note "Environment: Staging"
my-context note "Build: v2.5.1"
my-context note "Tools: Cypress 13.6, Manual testing"

my-context file tests/e2e/payment_flow_test.js
my-context file test-results/regression-2025-10-22.json

my-context stop

# Load testing
my-context start "Load test: 100 concurrent users" --project payments
my-context note "Test: 100 concurrent payment requests"
my-context note "Duration: 5 minutes"
my-context note "Results:"
my-context note "  Success: 95%"
my-context note "  Failures: 5 (timeout)"
my-context note "  Avg response: 250ms"
my-context note "  P95: 500ms"
my-context note "CONCERN: 5% failure rate unacceptable"
my-context note "Recommendation: Scale database connections"

my-context export "Load test: 100 concurrent users" --to load-test-report.md
my-context stop
```

### Your Key Commands

```bash
# Test session tracking
my-context start "Test: Feature name" --project service

# Systematic test notes
my-context note "TEST: Scenario description"
my-context note "  Input: ..."
my-context note "  Expected: ..."
my-context note "  Actual: ..."
my-context note "  Result: PASS/FAIL"

# Bug tracking
my-context note "BUG: Bug ID #123 - description"
my-context note "  Severity: Critical/High/Medium/Low"
my-context note "  Reproduction: Steps to reproduce"
my-context note "  Environment: Staging/v2.5.1"

# Test results export
my-context export "Regression test session" --to test-report.md
```

### Tips

1. **Systematic test documentation**:
   ```bash
   my-context note "TEST SUITE: Payment Flow"
   my-context note "  1. Valid payment: PASS"
   my-context note "  2. Invalid card: PASS"
   my-context note "  3. Timeout recovery: PASS"
   my-context note "  4. Refund processing: FAIL (see BUG #123)"
   my-context note "OVERALL: 3/4 tests passed"
   ```

2. **Environment documentation**:
   ```bash
   my-context note "Environment: Staging"
   my-context note "Build version: v2.5.1"
   my-context note "Database: Fresh seed data"
   my-context note "External services: Mocked"
   ```

3. **Performance metrics**:
   ```bash
   my-context note "Performance metrics:"
   my-context note "  Page load: 1.2s"
   my-context note "  API response: 200ms average"
   my-context note "  Database query: 50ms (index used)"
   ```

4. **Regression impact analysis**:
   ```bash
   my-context note "Regression impact:"
   my-context note "  Changed: Payment processor integration"
   my-context note "  Risk: Payment flow, refunds, reporting"
   my-context note "  Test coverage: 95% of affected code"
   ```

---

## Tech Lead / Architect

### Your Challenge
You make architectural decisions that affect:
- Multiple teams
- System scalability
- Technical debt
- Security and compliance
- Long-term maintainability

### Daily Workflow

```bash
# Architecture decision session
my-context start "Architecture: Redis caching layer" --project infrastructure

# Problem statement
my-context note "PROBLEM: API responses slow (500ms avg)"
my-context note "ROOT CAUSE: Repeated database queries"
my-context note "SCOPE: Affects 30% of API endpoints"

# Solution evaluation
my-context note "EVALUATED OPTIONS:"
my-context note "  1. Redis - In-memory, fast, scalable"
my-context note "  2. Memcached - Simpler, stateless"
my-context note "  3. Database query optimization - Limited improvement"

# Decision
my-context note "DECISION: Implement Redis caching"
my-context note "RATIONALE: Significant perf improvement, team familiar"
my-context note "TRADE-OFFS: Cache invalidation complexity, operational burden"

# Architecture details
my-context note "ARCHITECTURE:"
my-context note "  Pattern: Cache-aside (check cache → compute → store)"
my-context note "  TTL: 1 hour for most queries, 5 min for user data"
my-context note "  Replication: 3-node cluster with automatic failover"
my-context note "  Monitoring: Alert on cache miss rate > 20%"

# Implementation plan
my-context note "IMPLEMENTATION:"
my-context note "  Phase 1: Redis infrastructure (Ops) - Week 1"
my-context note "  Phase 2: Middleware integration (Backend) - Week 2"
my-context note "  Phase 3: Deployment to production (Devops) - Week 3"

my-context file docs/adr-0008-redis-caching.md
my-context file architecture/caching-architecture.png
my-context stop

# Code review leadership
my-context start "Code review: Payment processor abstraction" --project backend

my-context note "REVIEW FOCUS:"
my-context note "  - Abstraction design (interfaces)"
my-context note "  - Error handling and retries"
my-context note "  - Testing and testability"
my-context note "  - Performance (API limits)"

my-context note "REVIEW COMMENTS:"
my-context note "  1. Interface looks good - supports multiple providers"
my-context note "  2. Error handling: Consider exponential backoff for retries"
my-context note "  3. Tests: Need more edge cases (timeout, rate limiting)"
my-context note "  4. LGTM with suggested changes"

my-context export "Code review: Payment processor abstraction" --to review-summary.md
my-context stop

# Team sync prep
my-context list --search "architecture"  # Review recent decisions
my-context list --project "backend" --limit 5  # See team's work
```

### Your Key Commands

```bash
# Architecture decision records
my-context start "ADR: Architecture decision" --project infrastructure

# Problem/Solution/Decision structure
my-context note "PROBLEM: Clear problem statement"
my-context note "ROOT CAUSE: Why it's happening"
my-context note "EVALUATED: Option A, Option B, Option C"
my-context note "DECISION: We choose Option A"
my-context note "RATIONALE: Why we chose it"
my-context note "TRADE-OFFS: Pros and cons"

# Export for team
my-context export "ADR: Decision name" --to docs/adr-000X.md

# Track team health
my-context list --project team-name
my-context list --search "technical-debt"
```

### Tips

1. **Architecture Decision Records (ADRs)**:
   ```bash
   my-context start "ADR: Microservices vs Monolith"
   my-context note "PROBLEM: Scaling to 100 services"
   my-context note "EVALUATED:"
   my-context note "  - Monolith: Simpler, shared DB, tight coupling"
   my-context note "  - Microservices: Scalable, independent, complex"
   my-context note "  - Hybrid: Services for compute-heavy, monolith for core"
   my-context note "DECISION: Hybrid approach"
   my-context note "CONSEQUENCES: Operational complexity increases"
   my-context export "ADR: Microservices vs Monolith" --to adr.md
   ```

2. **Technical debt tracking**:
   ```bash
   my-context start "Tech debt: Legacy auth system" --project infrastructure
   my-context note "SCOPE: Authentication module from 2018"
   my-context note "ISSUES: No multi-factor auth, deprecated crypto"
   my-context note "IMPACT: Security risk, hard to extend"
   my-context note "EFFORT: 4 weeks to replace"
   my-context note "PRIORITY: High - security critical"
   ```

3. **Cross-team coordination**:
   ```bash
   my-context note "DEPENDENCY: Team B needs this API by end of Q4"
   my-context note "BLOCKER: Waiting for design review from Product"
   my-context note "SIGNAL: team-b-payment-api-ready"
   my-context signal create "team-a-db-migration-complete"
   ```

4. **Retrospective insights**:
   ```bash
   my-context list --all --search "Q3"
   # Review all decisions made in Q3
   # Analyze what worked, what didn't
   # Share learnings with team
   ```

---

## Scrum Master / Agile Coach

### Your Challenge
You facilitate:
- Sprint planning and execution
- Team velocity tracking
- Dependency management
- Sprint retrospectives
- Blocker identification

### Daily Workflow

```bash
# Sprint planning prep
my-context start "Sprint 6 Planning" --project scrum

my-context note "PLANNING NOTES:"
my-context note "  Team capacity: 32 story points"
my-context note "  Previous velocity: 28 points/sprint"
my-context note "  Stories available: 45 points"
my-context note "  Committed: 32 points"

my-context stop

# Standup facilitation
my-context start "Daily standup - Oct 22" --project scrum

my-context note "TEAM UPDATES:"
my-context note "  Backend: Payment API 80% complete"
my-context note "  Frontend: Login form in review"
my-context note "  QA: Regression testing in progress"

my-context note "BLOCKERS:"
my-context note "  - Frontend waiting on API docs (expected today)"
my-context note "  - QA waiting on staging deployment"

my-context note "RISKS:"
my-context note "  - Database migration complexity may cause delay"
my-context note "  - One team member out sick"

my-context note "ACTIONS:"
my-context note "  - Backend to provide API docs by 3pm"
my-context note "  - Devops to deploy to staging by 5pm"

my-context stop

# Sprint progress tracking
my-context start "Sprint 6 Mid-Sprint Check" --project scrum

my-context note "PROGRESS:"
my-context note "  Completed: 12 story points (37%)"
my-context note "  In Progress: 18 story points (56%)"
my-context note "  Blocked: 2 story points (6%)"

my-context note "TREND:"
my-context note "  On track for 30/32 points"
my-context note "  Velocity: 15/16 points/day (75% expected pace)"

my-context stop

# Retrospective prep
my-context start "Sprint 6 Retrospective" --project scrum

my-context note "WHAT WENT WELL:"
my-context note "  - Quick API integration (1 day vs 3 days estimated)"
my-context note "  - Fewer production issues than previous sprints"
my-context note "  - Team collaborated well across silos"

my-context note "WHAT DIDN'T GO WELL:"
my-context note "  - Database migration took 4 days vs 2 estimated"
my-context note "  - QA testing blocked by late deployments"
my-context note "  - No clear acceptance criteria on Story #42"

my-context note "IMPROVEMENTS:"
my-context note "  1. Define acceptance criteria in refinement"
my-context note "  2. Deploy to staging daily, not at end"
my-context note "  3. Database team should estimate more conservatively"

my-context note "ACTION ITEMS:"
my-context note "  - PO: Refine acceptance criteria earlier"
my-context note "  - Devops: Implement daily staging deployments"
my-context note "  - Backend: Get DB migration expert for estimation"

my-context export "Sprint 6 Retrospective" --to retrospective.md
my-context stop
```

### Your Key Commands

```bash
# Sprint tracking
my-context start "Sprint N: Planning/Progress/Retro"

# Structured notes
my-context note "METRIC: X points completed"
my-context note "BLOCKER: Issue description and impact"
my-context note "RISK: Potential problem and mitigation"
my-context note "ACTION: Who will do what by when"

# Export sprint summary
my-context export "Sprint 6 Retrospective" --to sprint-retro.md

# Review team's work
my-context list --project "sprint-6" --all
my-context list --search "blocker"
```

### Tips

1. **Velocity tracking**:
   ```bash
   my-context start "Velocity Trend"
   my-context note "Sprint 1: 24 points (baseline)"
   my-context note "Sprint 2: 28 points (+17%)"
   my-context note "Sprint 3: 26 points (-7%)"
   my-context note "Sprint 4: 32 points (+23%)"
   my-context note "Trend: Increasing - team stabilizing"
   ```

2. **Blocker and risk tracking**:
   ```bash
   my-context note "BLOCKER: Cannot start UI testing"
   my-context note "  Reason: API not deployed to staging"
   my-context note "  Impact: QA team idle (4 person-days)"
   my-context note "  Dependency: Devops deployment"
   my-context note "  ETA: End of day"
   ```

3. **Sprint health assessment**:
   ```bash
   my-context note "HEALTH CHECK - Day 5 of 10"
   my-context note "  Pace: 50% of points done"
   my-context note "  Status: ON TRACK"
   my-context note "  Risks: Database migration falling behind"
   my-context note "  Action: Pair database team with helper"
   ```

4. **Retrospective data collection**:
   ```bash
   # Collect throughout sprint
   my-context note "LEARNING: API-first approach worked well"
   my-context note "ISSUE: Acceptance criteria unclear on Story #42"

   # At retrospective time:
   my-context show  # Review all notes
   my-context export "Sprint 6 Retro" --to retro.md
   ```

---

## Product Manager / Product Owner

### Your Challenge
You define and track:
- Feature requirements
- User stories
- Acceptance criteria
- Priority and business value
- Release planning

### Daily Workflow

```bash
# Feature definition
my-context start "Feature: User preferences panel"

my-context note "BUSINESS VALUE:"
my-context note "  - User request: 15 customers asking for this"
my-context note "  - Impact: Reduce support tickets by 20%"
my-context note "  - Priority: High"

my-context note "REQUIREMENTS:"
my-context note "  - Allow users to change language"
my-context note "  - Allow timezone selection"
my-context note "  - Save preferences across sessions"
my-context note "  - Display in user menu"

my-context note "ACCEPTANCE CRITERIA:"
my-context note "  1. User can open preferences from account menu"
my-context note "  2. Language change applies immediately"
my-context note "  3. Timezone selection shows UTC offset"
my-context note "  4. Preferences persist across sessions"
my-context note "  5. Works on mobile (responsive)"

my-context note "DESIGN REVIEW: Required from Design team"
my-context note "ESTIMATION: 8 story points (planned for Sprint 6)"

my-context stop

# Stakeholder update prep
my-context start "Stakeholder Update: Q4 Progress"

my-context note "COMPLETED IN Q3:"
my-context note "  - User authentication revamp"
my-context note "  - Payment processor integration"
my-context note "  - Mobile app launch"

my-context note "IMPACT:"
my-context note "  - Reduced login failures by 60%"
my-context note "  - Increased transaction volume by 40%"
my-context note "  - 25k mobile users in first month"

my-context note "PLANNED FOR Q4:"
my-context note "  - Analytics dashboard"
my-context note "  - Advanced reporting"
my-context note "  - Team collaboration features"

my-context export "Q4 Stakeholder Update" --to stakeholder-update.md
my-context stop

# Release planning
my-context start "Release v3.0 Planning"

my-context note "TIMELINE: Dec 1, 2025"
my-context note "FEATURES PLANNED:"
my-context note "  1. User preferences (Sprint 6)"
my-context note "  2. Analytics dashboard (Sprint 7)"
my-context note "  3. Advanced reporting (Sprint 7-8)"
my-context note "  4. Team collaboration (Sprint 8)"

my-context note "DEPENDENCIES:"
my-context note "  - Design team: Must complete mockups by Nov 10"
my-context note "  - Backend: Database schema review by Nov 5"
my-context note "  - QA: Full regression testing by Nov 20"

my-context note "GO/NO-GO CRITERIA:"
my-context note "  - Test coverage > 85%"
my-context note "  - Performance: API < 200ms P95"
my-context note "  - No critical bugs"
my-context note "  - Design review approved"

my-context export "Release v3.0 Planning" --to release-plan.md
my-context stop
```

### Your Key Commands

```bash
# Feature definition
my-context start "Feature: Feature name"

# Structure your notes
my-context note "BUSINESS VALUE: Why we're building this"
my-context note "REQUIREMENTS: What needs to be built"
my-context note "ACCEPTANCE CRITERIA: How we know it's done"
my-context note "PRIORITY: High/Medium/Low"
my-context note "ESTIMATE: Story points"

# Export for team
my-context export "Feature name" --to feature-spec.md

# Release planning
my-context list --search "v3.0"  # See all v3.0 work
```

### Tips

1. **User story structure**:
   ```bash
   my-context start "Story: User can export data as CSV"
   my-context note "As a power user"
   my-context note "I want to export my data as CSV"
   my-context note "So that I can analyze it in Excel"
   my-context note ""
   my-context note "ACCEPTANCE CRITERIA:"
   my-context note "  1. Export button visible in settings"
   my-context note "  2. CSV includes all user data"
   my-context note "  3. File downloads to downloads folder"
   my-context note "  4. Works on Chrome, Safari, Firefox"
   ```

2. **Business value tracking**:
   ```bash
   my-context note "BUSINESS IMPACT:"
   my-context note "  - Addresses customer complaint #4823"
   my-context note "  - Competitive advantage over Competitor X"
   my-context note "  - Expected to increase engagement 15%"
   my-context note "  - Revenue impact: +$50k/year"
   ```

3. **Dependency tracking**:
   ```bash
   my-context note "DEPENDENCIES:"
   my-context note "  - Blocked: Waiting for API docs from backend (ETA: Oct 25)"
   my-context note "  - Blocks: Mobile app team can't start testing"
   my-context note "  - Needs: Design review from UX team"
   ```

4. **Release readiness**:
   ```bash
   my-context note "GO/NO-GO CHECKLIST:"
   my-context note "  ☑ Feature complete"
   my-context note "  ☑ Tested (QA approved)"
   my-context note "  ☐ Documentation complete (pending)"
   my-context note "  ☑ Design review approved"
   my-context note "  Status: READY FOR RELEASE"
   ```

---

## Quick Reference by Role

| Role | Start With | Key Practice | Main Commands |
|------|-----------|---|---|
| Backend Dev | Stage 1 → 2 | Document decisions | `start`, `note`, `file`, `export` |
| Frontend Dev | Stage 1 → 2 | Design/A11y notes | `start`, `note`, `file`, `export` |
| QA Engineer | Stage 1 → 3 | Test results/bugs | `start`, `note`, `file`, `export` |
| Tech Lead | Stage 3 → 5 | ADR decisions | `start`, `note`, `export`, `archive` |
| Scrum Master | Stage 2 → 4 | Sprint tracking | `start`, `note`, `list`, `export` |
| Product Manager | Stage 1 → 3 | Feature definition | `start`, `note`, `list`, `export` |

---

## Recommended Reading Order

**For your role:**
1. Read PROGRESSIVE-GUIDE.md (Stages 1-3)
2. Read your role-specific section above
3. Try the workflow examples
4. Read TRIGGERS-TUTORIAL.md if you want automation
5. Reference README.md for complete command syntax

**Common path:**
```
Week 1: Stages 1-2 (basic context tracking)
Week 2: Stage 3 (organizing work)
Week 3: Stage 4 (sharing with team)
Week 4+: Stage 5-6 (automation, team coordination)
```

Start with Stage 1 and take your time. You'll use different features at different times.

Happy tracking! 🚀
