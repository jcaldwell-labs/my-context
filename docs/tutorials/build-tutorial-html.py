#!/usr/bin/env python3
"""
Build Tutorial HTML Pages

Generates complete tutorial HTML pages with embedded panel exports,
bash code examples, and step-by-step walkthroughs.
"""

from pathlib import Path
from datetime import datetime

TUTORIALS_DIR = Path(__file__).parent
CSS_FILE = TUTORIALS_DIR / "shared-assets" / "tutorial-theme.css"

# Load CSS
with open(CSS_FILE, 'r') as f:
    TUTORIAL_CSS = f.read()

def create_html_page(title, subtitle, content, prev_link=None, next_link=None):
    """Create complete HTML page with navigation"""

    nav_html = '<div class="tutorial-nav">'
    nav_html += '<a href="../INDEX.html">← All Tutorials</a>'
    nav_html += '</div>'

    footer_nav = '<div class="footer-nav">'
    if prev_link:
        footer_nav += f'<a href="{prev_link}" class="prev">Previous Tutorial</a>'
    else:
        footer_nav += '<span></span>'

    footer_nav += '<a href="../INDEX.html">Tutorial Index</a>'

    if next_link:
        footer_nav += f'<a href="{next_link}" class="next">Next Tutorial</a>'
    else:
        footer_nav += '<span></span>'
    footer_nav += '</div>'

    return f'''<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{title} - My-Context Visual Tutorial</title>
    <style>
{TUTORIAL_CSS}
    </style>
</head>
<body>
    <div class="tutorial-header">
        <div class="container">
            <h1>{title}</h1>
            <p class="subtitle">{subtitle}</p>
            {nav_html}
        </div>
    </div>

    <div class="container">
{content}
    </div>

    <div class="tutorial-footer">
        <div class="container">
            {footer_nav}
        </div>
    </div>
</body>
</html>'''

def embed_panel(panel_file, title):
    """Create HTML for embedding a panel export"""
    return f'''<div class="panel-embed">
    <h3>{title}</h3>
    <iframe src="{panel_file}" title="{title}"></iframe>
</div>'''

# =============================================================================
# Tutorial 1: Backend Developer Solo
# =============================================================================

def build_tutorial_01():
    """Build Tutorial 1: Backend Developer Solo"""
    print("Building Tutorial 1...")

    content = '''
<div class="section scenario">
    <p><strong>Role:</strong> <span class="character backend">Alice - Backend Developer</span></p>
    <p><strong>Scenario:</strong> Alice is implementing a payment retry feature with exponential backoff. She wants to track her decisions and remember which files she modified.</p>
    <p><strong>Challenge:</strong> How do I track my work so I can review it later or share with teammates?</p>
</div>

<div class="objectives">
    <h3>🎯 Learning Objectives</h3>
    <ul>
        <li>Start your first context for tracking work</li>
        <li>Add notes to document important decisions</li>
        <li>Associate files you modify with the context</li>
        <li>View your work progress</li>
        <li>Stop the context when done</li>
    </ul>
</div>

<div class="section">
    <h2>Step-by-Step Walkthrough</h2>

    <div class="step">
        <h4><span class="step-number">1</span> Start a Context</h4>
        <div class="step-description">
            <p>Alice begins working on the payment retry feature. She starts a context to track her work:</p>
        </div>
        <pre><code>my-context start "payment-retry-logic" --project payment-service</code></pre>
        <div class="output success">Context 'payment-service: payment-retry-logic' started</div>
        <p><strong>What happened:</strong> My-context created a new context and is now tracking Alice's work session.</p>
    </div>

    <div class="step">
        <h4><span class="step-number">2</span> Add Notes About Decisions</h4>
        <div class="step-description">
            <p>As Alice makes design decisions, she documents them with notes:</p>
        </div>
        <pre><code>my-context note "DECISION: Exponential backoff strategy - 1s, 2s, 4s, 8s"
my-context note "Using exponential backoff to handle transient failures"
my-context note "Max retries: 3 attempts before marking payment as failed"</code></pre>
        <div class="output success">Note added to context
Note added to context
Note added to context</div>
        <p><strong>Why this matters:</strong> These notes explain the "why" behind Alice's implementation choices. Future Alice (or teammates) will understand her reasoning.</p>
    </div>

    <div class="step">
        <h4><span class="step-number">3</span> Track Files You Modify</h4>
        <div class="step-description">
            <p>Alice associates the files she's working on with the context:</p>
        </div>
        <pre><code>my-context file internal/payments/retry.go
my-context file internal/payments/backoff.go
my-context file tests/payments/retry_test.go</code></pre>
        <div class="output success">File associated: internal/payments/retry.go
File associated: internal/payments/backoff.go
File associated: tests/payments/retry_test.go</div>
        <p><strong>Benefit:</strong> She can now see exactly which files were part of this feature work.</p>
    </div>

    <div class="step">
        <h4><span class="step-number">4</span> Check Your Progress</h4>
        <div class="step-description">
            <p>At any time, Alice can view her current work:</p>
        </div>
        <pre><code>my-context show</code></pre>
        <div class="output">Active Context: payment-service: payment-retry-logic
Status: Active (started 47 minutes ago)

Notes:
  [14:20:32] DECISION: Exponential backoff strategy - 1s, 2s, 4s, 8s
  [14:28:15] Using exponential backoff to handle transient failures
  [14:35:22] Max retries: 3 attempts before marking payment as failed

Files:
  internal/payments/retry.go
  internal/payments/backoff.go
  tests/payments/retry_test.go</div>
        <p><strong>Result:</strong> A complete record of decisions and files for this feature.</p>
    </div>

    <div class="step">
        <h4><span class="step-number">5</span> Stop the Context When Done</h4>
        <div class="step-description">
            <p>When Alice finishes her work session, she stops the context:</p>
        </div>
        <pre><code>my-context stop</code></pre>
        <div class="output success">Context stopped
Duration: 47 minutes</div>
        <p><strong>Why stop?</strong> Accurate time tracking and clear "what am I working on" state.</p>
    </div>
</div>

<div class="section">
    <h2>Visual Result: Context Explorer Panel</h2>
    <p>Here's what Alice's context looks like in the visual explorer:</p>
    ''' + embed_panel('tutorial-01-backend-solo_explorer.html', 'Context Hierarchy - Alice\'s Work') + '''
</div>

<div class="takeaways">
    <h3>✅ Key Takeaways</h3>
    <ul>
        <li><strong>Start context:</strong> <code>my-context start "name" --project project-name</code></li>
        <li><strong>Document decisions:</strong> <code>my-context note "DECISION: explanation"</code></li>
        <li><strong>Track files:</strong> <code>my-context file path/to/file</code></li>
        <li><strong>Check progress:</strong> <code>my-context show</code></li>
        <li><strong>Stop when done:</strong> <code>my-context stop</code></li>
    </ul>
    <div class="alert info">
        <strong>💡 Pro Tip:</strong> Use prefixes in notes like "DECISION:", "BUG:", "QUESTION:" to make them searchable later.
    </div>
</div>

<div class="section">
    <h2>Try It Yourself</h2>
    <p>Copy and run these commands to create your own first context:</p>
    <pre><code># Start a context for your current work
my-context start "your-task-name" --project your-project

# Document what you're doing
my-context note "Starting work on [describe your task]"

# Add files as you modify them
my-context file path/to/your/file

# Check your progress
my-context show

# Stop when done
my-context stop</code></pre>
</div>
'''

    html = create_html_page(
        title="Tutorial 1: Your First Context",
        subtitle="Backend Developer - Solo Work Pattern",
        content=content,
        prev_link=None,
        next_link="../tutorial-02/tutorial-02.html"
    )

    output_file = TUTORIALS_DIR / "tutorial-01" / "tutorial-01.html"
    output_file.write_text(html, encoding='utf-8')
    print(f"  ✅ Created: {output_file}")

# Continue with remaining tutorials...
# For now, let me create a condensed version for tutorials 2-8

def build_tutorial_02():
    """Build Tutorial 2: Frontend Developer Solo"""
    print("Building Tutorial 2...")

    content = '''
<div class="section scenario">
    <p><strong>Role:</strong> <span class="character frontend">Bob - Frontend Developer</span></p>
    <p><strong>Scenario:</strong> Bob is building a responsive checkout UI with accessibility features. He wants to track design decisions and A11y testing notes.</p>
</div>

<div class="objectives">
    <h3>🎯 Learning Objectives</h3>
    <ul>
        <li>Track frontend-specific work (components, styling, A11y)</li>
        <li>Document design decisions</li>
        <li>Record accessibility testing results</li>
    </ul>
</div>

<div class="section">
    <h2>Bob's Workflow</h2>

    <div class="step">
        <h4><span class="step-number">1</span> Start Work on Checkout UI</h4>
        <pre><code>my-context start "checkout-ui-responsive" --project web-app</code></pre>
    </div>

    <div class="step">
        <h4><span class="step-number">2</span> Document Design Decisions</h4>
        <pre><code>my-context note "DECISION: Using CSS Grid for 3-column layout"
my-context note "Mobile-first approach - 1 column on mobile, 3 on desktop"
my-context note "Color contrast ratio: 4.5:1 (WCAG AA compliant)"</code></pre>
    </div>

    <div class="step">
        <h4><span class="step-number">3</span> Track A11y Testing</h4>
        <pre><code>my-context note "A11Y: Added ARIA labels to all payment form fields"
my-context note "A11Y: Keyboard navigation fully supported (Tab order tested)"
my-context note "Using semantic HTML for better screen reader support"</code></pre>
    </div>

    <div class="step">
        <h4><span class="step-number">4</span> Associate Files</h4>
        <pre><code>my-context file src/components/Checkout.tsx
my-context file src/components/PaymentForm.tsx
my-context file src/styles/checkout.css
my-context file tests/components/Checkout.test.tsx</code></pre>
    </div>
</div>

<div class="section">
    <h2>Visual Result</h2>
    ''' + embed_panel('tutorial-02-frontend-solo_explorer.html', 'Bob\'s Frontend Work') + '''
</div>

<div class="takeaways">
    <h3>✅ Key Takeaways</h3>
    <ul>
        <li>Frontend developers benefit from tracking <strong>design decisions</strong> and <strong>A11y notes</strong></li>
        <li>Same commands work for any role - <code>start</code>, <code>note</code>, <code>file</code>, <code>stop</code></li>
        <li>Notes create a searchable history of "why we built it this way"</li>
    </ul>
</div>
'''

    html = create_html_page(
        title="Tutorial 2: Frontend Developer",
        subtitle="Responsive UI with Accessibility",
        content=content,
        prev_link="../tutorial-01/tutorial-01.html",
        next_link="../tutorial-03/tutorial-03.html"
    )

    output_file = TUTORIALS_DIR / "tutorial-02" / "tutorial-02.html"
    output_file.write_text(html, encoding='utf-8')
    print(f"  ✅ Created: {output_file}")

def build_tutorial_03():
    """Build Tutorial 3: QA Engineer Solo"""
    print("Building Tutorial 3...")

    content = '''
<div class="section scenario">
    <p><strong>Role:</strong> <span class="character qa">Carol - QA Engineer</span></p>
    <p><strong>Scenario:</strong> Carol is testing the payment flow across multiple browsers. She needs to track test results and document bugs she discovers.</p>
</div>

<div class="objectives">
    <h3>🎯 Learning Objectives</h3>
    <ul>
        <li>Track QA test sessions</li>
        <li>Document bug discoveries with details</li>
        <li>Record browser compatibility results</li>
    </ul>
</div>

<div class="section">
    <h2>Carol's QA Workflow</h2>
    <div class="step">
        <h4><span class="step-number">1</span> Start Test Session</h4>
        <pre><code>my-context start "payment-flow-testing" --project qa-suite</code></pre>
    </div>
    <div class="step">
        <h4><span class="step-number">2</span> Track Test Results</h4>
        <pre><code>my-context note "✅ Test passed: Chrome 120 - payment successful"
my-context note "❌ BUG: Safari 16.2 - card validation fails on submit"</code></pre>
    </div>
</div>

<div class="section">
    <h2>Visual Result</h2>
    ''' + embed_panel('tutorial-03-qa-solo_explorer.html', 'Carol\'s QA Testing Session') + '''
</div>

<div class="takeaways">
    <h3>✅ Key Takeaways</h3>
    <ul>
        <li>QA engineers track <strong>test results</strong> and <strong>bugs</strong></li>
        <li>Use emoji prefixes (✅ ❌) for quick visual status</li>
    </ul>
</div>
'''

    html = create_html_page(
        title="Tutorial 3: QA Engineer Workflow",
        subtitle="Cross-Browser Testing and Bug Tracking",
        content=content,
        prev_link="../tutorial-02/tutorial-02.html",
        next_link="../tutorial-04/tutorial-04.html"
    )

    output_file = TUTORIALS_DIR / "tutorial-03" / "tutorial-03.html"
    output_file.write_text(html, encoding='utf-8')
    print(f"  ✅ Created: {output_file}")

def build_tutorial_04():
    """Build Tutorial 4: Multi-Project Management"""
    print("Building Tutorial 4...")

    content = '''
<div class="section scenario">
    <p><strong>Role:</strong> <span class="character backend">Alice - Backend Consultant</span></p>
    <p><strong>Scenario:</strong> Alice is a consultant juggling work for three clients simultaneously: ACME Corp (e-commerce), TechCorp (B2B SaaS), and a fintech Startup. She uses the <code>--project</code> flag to keep each client's contexts organized.</p>
    <p><strong>Time Period:</strong> Over one week, Alice switches between client projects multiple times per day.</p>
</div>

<div class="objectives">
    <h3>🎯 Learning Objectives</h3>
    <ul>
        <li>Use <code>--project</code> flag to organize contexts by client</li>
        <li>Switch between multiple projects seamlessly</li>
        <li>View contexts filtered by project</li>
        <li>Track parallel work streams without confusion</li>
    </ul>
</div>

<div class="section">
    <h2>Alice's Multi-Client Workflow</h2>

    <div class="step">
        <h4><span class="step-number">1</span> Monday Morning: ACME Corp Work</h4>
        <p class="step-description">Alice starts her week optimizing ACME's API performance.</p>
        <pre><code># Start ACME context with project flag
my-context start "api-optimization" --project client-acme

# Document the work
my-context note "Client context: E-commerce platform, high traffic"
my-context note "Performance issue: API latency >500ms on product search"
my-context note "SOLUTION: Added Redis caching layer for product catalog"
my-context file acme/internal/cache/redis.go
my-context stop</code></pre>
    </div>

    <div class="step">
        <h4><span class="step-number">2</span> Monday Afternoon: TechCorp Database Migration</h4>
        <p class="step-description">Switching to TechCorp's PostgreSQL migration project.</p>
        <pre><code># Start TechCorp context
my-context start "database-migration" --project client-techcorp

# Document migration strategy
my-context note "Client context: B2B SaaS, migrating MySQL → PostgreSQL"
my-context note "Migration strategy: Blue-green deployment with dual-write period"
my-context note "Data validation: Compare row counts after migration"
my-context file techcorp/migrations/001_initial_schema.sql
my-context stop</code></pre>
    </div>

    <div class="step">
        <h4><span class="step-number">3</span> Tuesday: Startup Security Audit</h4>
        <p class="step-description">Working on SOC2 audit preparation for the fintech startup.</p>
        <pre><code># Start Startup context
my-context start "security-audit" --project client-startup

# Document security findings
my-context note "Client context: Fintech startup preparing for SOC2 audit"
my-context note "FINDING: API keys stored in plaintext in config files"
my-context note "RECOMMENDATION: Move to AWS Secrets Manager"
my-context note "FINDING: No rate limiting on public API endpoints"
my-context note "RECOMMENDATION: Implement token bucket algorithm"
my-context file startup/security/audit-report.md
my-context stop</code></pre>
    </div>

    <div class="step">
        <h4><span class="step-number">4</span> Viewing Contexts by Project</h4>
        <p class="step-description">Alice can filter contexts to focus on one client at a time.</p>
        <pre><code># View all ACME contexts
my-context list --project client-acme

# View all TechCorp contexts
my-context list --project client-techcorp

# View all Startup contexts
my-context list --project client-startup

# View all contexts (across all clients)
my-context list</code></pre>
    </div>

    <div class="step">
        <h4><span class="step-number">5</span> Later in the Week: More Client Switching</h4>
        <p class="step-description">Alice continues work across all three clients, creating more contexts.</p>
        <pre><code># ACME: Stripe integration
my-context start "acme-payment-integration" --project client-acme
my-context note "Integrating Stripe payment gateway"
my-context stop

# TechCorp: CI/CD setup
my-context start "techcorp-ci-pipeline" --project client-techcorp
my-context note "Setting up GitHub Actions for automated testing"
my-context stop

# Startup: Monitoring setup
my-context start "startup-monitoring" --project client-startup
my-context note "Setting up Datadog for application monitoring"
my-context stop</code></pre>
    </div>
</div>

<div class="section">
    <h2>Visual Result: All Projects Organized</h2>
    <p>After a week of work, Alice has 6 contexts organized across 3 client projects. The explorer panel shows the clear hierarchy:</p>
    ''' + embed_panel('tutorial-04-multi-project_explorer.html', 'Alice\'s Multi-Client Context Organization') + '''
</div>

<div class="takeaways">
    <h3>✅ Key Takeaways</h3>
    <ul>
        <li>The <code>--project</code> flag is essential for <strong>organizing work across multiple codebases or clients</strong></li>
        <li>Contexts are automatically <strong>grouped by project</strong> in the explorer view</li>
        <li>You can <strong>filter contexts</strong> using <code>my-context list --project [name]</code></li>
        <li>Each context maintains its own <strong>notes, files, and metadata</strong> regardless of project</li>
        <li>Perfect for consultants, contractors, or anyone managing <strong>parallel work streams</strong></li>
    </ul>
</div>

<div class="section">
    <h2>Try It Yourself</h2>
    <p>Practice organizing contexts by project:</p>
    <pre><code># Create contexts for different projects
my-context start "feature-x" --project project-a
my-context note "Working on project A"
my-context stop

my-context start "feature-y" --project project-b
my-context note "Working on project B"
my-context stop

# List all contexts grouped by project
my-context list

# Filter by specific project
my-context list --project project-a</code></pre>
</div>
'''

    html = create_html_page(
        title="Tutorial 4: Multi-Project Consultant",
        subtitle="Managing Multiple Clients with the --project Flag",
        content=content,
        prev_link="../tutorial-03/tutorial-03.html",
        next_link="../tutorial-05/tutorial-05.html"
    )

    output_file = TUTORIALS_DIR / "tutorial-04" / "tutorial-04.html"
    output_file.write_text(html, encoding='utf-8')
    print(f"  ✅ Created: {output_file}")

def build_tutorial_05():
    """Build Tutorial 5: Scrum Master Sprint Management"""
    print("Building Tutorial 5...")

    content = '''
<div class="section scenario">
    <p><strong>Role:</strong> <span class="character scrum">Dave - Scrum Master</span></p>
    <p><strong>Scenario:</strong> Dave is managing Sprint 5 for Team Alpha (5 developers). He uses my-context to track sprint planning, daily standups, retrospectives, and action items throughout the sprint lifecycle.</p>
    <p><strong>Team Context:</strong> Payment platform development, following 2-week sprints</p>
</div>

<div class="objectives">
    <h3>🎯 Learning Objectives</h3>
    <ul>
        <li>Track sprint planning sessions with capacity and velocity calculations</li>
        <li>Document daily standup notes and action items</li>
        <li>Record retrospective insights for continuous improvement</li>
        <li>Organize sprint contexts using the <code>--project</code> flag</li>
    </ul>
</div>

<div class="section">
    <h2>Dave's Sprint Workflow</h2>

    <div class="step">
        <h4><span class="step-number">1</span> Sprint Planning: Setting Up Sprint 5</h4>
        <p class="step-description">Dave starts Sprint 5 planning by documenting team capacity, velocity, and goals.</p>
        <pre><code># Start Sprint 5 planning context
my-context start "sprint-5-planning" --project team-alpha

# Document sprint goals
my-context note "Sprint 5 goals: Payment integration + responsive UI"

# Track team capacity
my-context note "Team capacity: 5 developers × 8 days = 40 dev-days"
my-context note "Team velocity: 42 story points (3-sprint average: 38)"
my-context note "Sprint commitment: 40 story points (conservative)"

# Identify blockers early
my-context note "BLOCKER: API dependency on Platform team (payment gateway)"
my-context note "Mitigation: Daily sync with Platform team lead"

# Associate planning artifacts
my-context file sprint-5/planning/sprint-backlog.md
my-context file sprint-5/planning/capacity-plan.xlsx

my-context stop</code></pre>
    </div>

    <div class="step">
        <h4><span class="step-number">2</span> Daily Standup: Day 3 Progress Tracking</h4>
        <p class="step-description">On Day 3 of the sprint, Dave tracks team updates and action items from standup.</p>
        <pre><code># Track daily standup
my-context start "sprint-5-day-3-standup" --project team-alpha

# Document each team member's update
my-context note "Alice: Payment retry logic - 80% complete"
my-context note "Bob: Checkout UI - blocked on design review"
my-context note "Carol: E2E tests - Safari bug found"

# Capture action items
my-context note "ACTION ITEM: Schedule design review for Bob today"

my-context stop</code></pre>
    </div>

    <div class="step">
        <h4><span class="step-number">3</span> Sprint 4 Retrospective: Lessons Learned</h4>
        <p class="step-description">After Sprint 4 completes, Dave documents the retrospective insights.</p>
        <pre><code># Document retrospective
my-context start "sprint-4-retrospective" --project team-alpha

# Track what went well
my-context note "Sprint 4: Completed 38 / 40 story points"
my-context note "What went well: Good collaboration, clear requirements"

# Track improvements for next sprint
my-context note "What to improve: Earlier QA involvement, better estimation"

my-context stop</code></pre>
    </div>

    <div class="step">
        <h4><span class="step-number">4</span> Viewing All Sprint Contexts</h4>
        <p class="step-description">Dave can view all contexts for Team Alpha to see the full sprint timeline.</p>
        <pre><code># View all Team Alpha contexts
my-context list --project team-alpha

# Output shows:
# - sprint-5-planning (current sprint)
# - sprint-5-day-3-standup (current sprint)
# - sprint-4-retrospective (completed sprint)
# - sprint-4-planning (completed sprint)</code></pre>
    </div>
</div>

<div class="section">
    <h2>Visual Result: Sprint Context Organization</h2>
    <p>Dave's context history shows the complete sprint lifecycle, from planning through execution to retrospectives:</p>
    ''' + embed_panel('tutorial-05-scrum-master_explorer.html', 'Dave\'s Sprint Management Contexts') + '''
</div>

<div class="takeaways">
    <h3>✅ Key Takeaways for Scrum Masters</h3>
    <ul>
        <li>Use separate contexts for <strong>each sprint ceremony</strong> (planning, standup, retrospective)</li>
        <li>Document <strong>blockers and action items</strong> in real-time during meetings</li>
        <li>Track <strong>team velocity and capacity</strong> for sprint-over-sprint trend analysis</li>
        <li>Use the <code>--project</code> flag to <strong>group contexts by team</strong></li>
        <li>Build a <strong>searchable history</strong> of sprint decisions and retrospective insights</li>
    </ul>
</div>

<div class="section">
    <h2>Advanced Patterns for Scrum Masters</h2>

    <h3>Pattern 1: Sprint Goal Tracking</h3>
    <pre><code># Create a context for overall sprint tracking
my-context start "sprint-6-overview" --project team-alpha
my-context note "Sprint goal: Achieve 95% payment success rate"
my-context note "Success metric: Monitor Datadog dashboard"
my-context note "Risk: Third-party API rate limiting"</code></pre>

    <h3>Pattern 2: Blocker Escalation Trail</h3>
    <pre><code># Document blocker progression
my-context note "BLOCKER: Database migration blocked (Day 2)"
my-context note "ESCALATION: Reached out to DBA team (Day 3)"
my-context note "RESOLVED: Migration script approved (Day 4)"</code></pre>

    <h3>Pattern 3: Multi-Team Coordination</h3>
    <pre><code># Track dependencies across teams
my-context start "cross-team-sync" --project team-alpha
my-context note "Dependency: Waiting for Platform team API v2 release"
my-context note "Platform ETA: End of Sprint 5 (Day 10)"
my-context note "Mitigation: Build against v1 API, migrate to v2 next sprint"</code></pre>
</div>

<div class="section">
    <h2>Try It Yourself</h2>
    <p>Practice tracking a sprint ceremony:</p>
    <pre><code># Start your own sprint planning context
my-context start "sprint-planning" --project your-team

# Document sprint details
my-context note "Sprint goals: [Your sprint goals]"
my-context note "Team capacity: [Number of devs] × [Days] = [Total dev-days]"
my-context note "Team velocity: [Story points from last sprint]"

# Track any blockers
my-context note "BLOCKER: [Describe blocker]"
my-context note "Mitigation: [How you'll address it]"

# Stop when planning is done
my-context stop

# Later, review all sprint contexts
my-context list --project your-team</code></pre>
</div>
'''

    html = create_html_page(
        title="Tutorial 5: Scrum Master Sprint Management",
        subtitle="Planning, Standups, and Retrospectives",
        content=content,
        prev_link="../tutorial-04/tutorial-04.html",
        next_link="../tutorial-06/tutorial-06.html"
    )

    output_file = TUTORIALS_DIR / "tutorial-05" / "tutorial-05.html"
    output_file.write_text(html, encoding='utf-8')
    print(f"  ✅ Created: {output_file}")

def build_tutorial_06():
    """Build Tutorial 6: Team Handoff"""
    print("Building Tutorial 6...")

    content = '''
<div class="section scenario">
    <p><strong>Team Members:</strong>
        <span class="character backend">Alice - Backend Developer</span>
        <span class="character frontend">Bob - Frontend Developer</span>
    </p>
    <p><strong>Scenario:</strong> Alice has built a new Payment API v2.0 endpoint. Bob needs to integrate it into the frontend UI. They use my-context to share implementation details asynchronously through documented contexts.</p>
    <p><strong>Collaboration Pattern:</strong> Async context sharing (documentation handoff)</p>
</div>

<div class="objectives">
    <h3>🎯 Learning Objectives</h3>
    <ul>
        <li>Share context between team members asynchronously</li>
        <li>Reference another team member's context as documentation</li>
        <li>Track implementation decisions that affect multiple teams</li>
        <li>Build a shared knowledge base across frontend and backend</li>
    </ul>
</div>

<div class="section">
    <h2>The Handoff Workflow</h2>

    <div class="step">
        <h4><span class="step-number">1</span> Alice: Backend API Implementation</h4>
        <p class="step-description">Alice creates a context documenting the new Payment API v2.0 endpoint.</p>
        <pre><code># Alice starts backend work
my-context start "payment-api-v2" --project backend-services

# Document the API specification
my-context note "API SPEC: POST /api/v2/payments/process"
my-context note "Request body: { amount, currency, payment_method, customer_id }"
my-context note "Response: { payment_id, status, transaction_id }"

# Document important decisions
my-context note "DECISION: Using idempotency keys to prevent duplicate charges"
my-context note "Idempotency key header: X-Idempotency-Key (UUID)"
my-context note "Error handling: Return 400 for validation, 402 for payment failure"
my-context note "Rate limiting: 100 requests/minute per customer"

# Track implementation files
my-context file internal/api/v2/payments.go
my-context file docs/api/payment-endpoint-spec.md

my-context stop</code></pre>
    </div>

    <div class="step">
        <h4><span class="step-number">2</span> Alice: Share Context with Bob</h4>
        <p class="step-description">Alice can export her context for Bob to review, or Bob can directly reference the shared documentation file.</p>
        <pre><code># Option 1: Export context for Bob
my-context export payment-api-v2 > handoff-to-bob.md

# Option 2: Bob reads the API spec file Alice tracked
# Bob sees: docs/api/payment-endpoint-spec.md in Alice's context</code></pre>
    </div>

    <div class="step">
        <h4><span class="step-number">3</span> Bob: Frontend Integration (References Alice's Work)</h4>
        <p class="step-description">Bob starts his frontend work and explicitly references Alice's backend context.</p>
        <pre><code># Bob starts frontend integration
my-context start "payment-ui-integration" --project web-app

# Reference Alice's work
my-context note "REF: Alice's payment API spec in docs/api/payment-endpoint-spec.md"
my-context note "Implementing UI integration with POST /api/v2/payments/process"

# Document frontend-specific decisions
my-context note "Added UUID generation for X-Idempotency-Key header"
my-context note "Error handling: Display user-friendly message for 400/402 errors"
my-context note "Loading state: Show spinner during payment processing"

# Track frontend files
my-context file src/services/payment-api.ts
my-context file src/components/PaymentProcessor.tsx

my-context stop</code></pre>
    </div>

    <div class="step">
        <h4><span class="step-number">4</span> Visual Result: Two Connected Contexts</h4>
        <p class="step-description">Each team member has their own context, but they reference each other's work.</p>
    </div>
</div>

<div class="section">
    <h2>Alice's Context (Backend)</h2>
    ''' + embed_panel('tutorial-06-team-alice_explorer.html', 'Alice\'s Backend API Context') + '''
</div>

<div class="section">
    <h2>Bob's Context (Frontend)</h2>
    ''' + embed_panel('tutorial-06-team-bob_explorer.html', 'Bob\'s Frontend Integration Context') + '''
</div>

<div class="takeaways">
    <h3>✅ Key Takeaways for Team Collaboration</h3>
    <ul>
        <li><strong>Async documentation sharing</strong>: Backend and frontend teams work independently but stay aligned</li>
        <li><strong>Explicit references</strong>: Use "REF:" notes to link contexts across team members</li>
        <li><strong>Shared files</strong>: API specs and documentation files tracked in contexts</li>
        <li><strong>Decision transparency</strong>: Important choices (idempotency, error codes) documented in context</li>
        <li><strong>Knowledge preservation</strong>: Future developers can understand the "why" behind implementations</li>
    </ul>
</div>

<div class="section">
    <h2>Advanced Team Patterns</h2>

    <h3>Pattern 1: Export Context for Code Reviews</h3>
    <pre><code># Export your context as part of PR description
my-context export my-feature-name > pr-context.md

# Reviewers see your decision trail, not just code changes</code></pre>

    <h3>Pattern 2: Context Chains Across Multiple Developers</h3>
    <pre><code># Developer 1 (Backend)
my-context note "Created user authentication endpoint"

# Developer 2 (Middleware)
my-context note "REF: Backend auth endpoint from Alice's context"
my-context note "Added rate limiting middleware"

# Developer 3 (Frontend)
my-context note "REF: Backend auth + middleware from Alice and Charlie"
my-context note "Integrated login UI with rate-limited auth endpoint"</code></pre>

    <h3>Pattern 3: Shared Context Home for Team Projects</h3>
    <pre><code># Team shares a context home directory
export MY_CONTEXT_HOME=/team/shared-contexts

# Everyone can see each other's contexts
my-context list

# Search across all team contexts
my-context search "payment API"</code></pre>
</div>

<div class="section">
    <h2>Try It Yourself</h2>
    <p>Practice async team collaboration:</p>
    <pre><code># Developer 1: Create a backend context
my-context start "api-endpoint" --project backend
my-context note "API SPEC: GET /api/users/:id"
my-context note "Returns: { id, name, email, created_at }"
my-context file backend/api/users.go
my-context stop

# Developer 2: Reference it in frontend work
my-context start "user-profile-ui" --project frontend
my-context note "REF: Backend API spec from Developer 1"
my-context note "Fetching user data from GET /api/users/:id"
my-context file frontend/components/UserProfile.tsx
my-context stop</code></pre>
</div>
'''

    html = create_html_page(
        title="Tutorial 6: Team Handoff",
        subtitle="Async Collaboration Between Backend and Frontend",
        content=content,
        prev_link="../tutorial-05/tutorial-05.html",
        next_link="../tutorial-07/tutorial-07.html"
    )

    output_file = TUTORIALS_DIR / "tutorial-06" / "tutorial-06.html"
    output_file.write_text(html, encoding='utf-8')
    print(f"  ✅ Created: {output_file}")

def build_tutorial_07():
    """Build Tutorial 7: Signal Coordination"""
    print("Building Tutorial 7...")

    content = '''
<div class="section scenario">
    <p><strong>Team Members:</strong>
        <span class="character backend">Alice - Backend</span>
        <span class="character frontend">Bob - Frontend</span>
        <span class="character qa">Carol - QA</span>
        <span class="character product">Eve - Product Owner</span>
    </p>
    <p><strong>Scenario:</strong> The team is coordinating a Payment API v2.0 release to staging. Each role depends on the previous one completing their work. They use signals for real-time coordination across the release pipeline.</p>
    <p><strong>Collaboration Pattern:</strong> Real-time signal coordination (workflow dependencies)</p>
</div>

<div class="objectives">
    <h3>🎯 Learning Objectives</h3>
    <ul>
        <li>Use signals to coordinate work across multiple team members</li>
        <li>Create event-driven workflows with context signals</li>
        <li>Track sequential dependencies in a release pipeline</li>
        <li>Enable real-time team awareness of progress</li>
    </ul>
</div>

<div class="section">
    <h2>The Signal Coordination Flow</h2>

    <div class="step">
        <h4><span class="step-number">1</span> Alice: Backend Deployment (Creates First Signal)</h4>
        <p class="step-description">Alice deploys the Payment API v2.0 to staging and signals it's ready.</p>
        <pre><code># Alice deploys backend
my-context start "payment-api-release" --project backend-services

my-context note "Release: Payment API v2.0 to staging"
my-context note "Endpoints ready: /process, /refund, /status"
my-context note "Database migrations applied successfully"
my-context note "Health check: ✅ All endpoints responding"

# Signal to the team that backend is ready
my-context signal create api-v2-staging-ready

my-context stop</code></pre>
    </div>

    <div class="step">
        <h4><span class="step-number">2</span> Bob: Frontend Integration (Watches for Alice's Signal)</h4>
        <p class="step-description">Bob waits for the backend signal, then integrates the frontend.</p>
        <pre><code># Bob waits for backend signal, then proceeds
my-context start "frontend-integration" --project web-app

my-context note "Waiting for api-v2-staging-ready signal..."
# (Signal received from Alice!)
my-context note "Signal received! Starting integration work"

my-context note "Integrated payment API v2.0 endpoints"
my-context note "Manual testing: Payment flow working end-to-end"
my-context note "Deployed to staging environment"

# Signal to QA that frontend is ready
my-context signal create frontend-staging-ready

my-context stop</code></pre>
    </div>

    <div class="step">
        <h4><span class="step-number">3</span> Carol: QA Testing (Watches for Bob's Signal)</h4>
        <p class="step-description">Carol waits for the frontend signal, then runs E2E tests.</p>
        <pre><code># Carol waits for frontend signal
my-context start "integration-testing" --project qa-suite

my-context note "Waiting for frontend-staging-ready signal..."
# (Signal received from Bob!)
my-context note "Signal received! Starting E2E tests"

my-context note "✅ Test suite: payment-flow-e2e (15 tests passed)"
my-context note "✅ Test: Process payment with valid card"
my-context note "✅ Test: Handle invalid card gracefully"
my-context note "✅ Test: Refund processed payment"
my-context note "✅ Test: Check payment status"
my-context note "All tests passed - staging environment approved"

# Signal to Product that QA approves
my-context signal create qa-approved-staging

my-context stop</code></pre>
    </div>

    <div class="step">
        <h4><span class="step-number">4</span> Eve: Release Coordination (Watches for Carol's Signal)</h4>
        <p class="step-description">Eve waits for QA approval, then schedules the production release.</p>
        <pre><code># Eve coordinates final release
my-context start "release-coordination" --project product

my-context note "Release: Payment v2.0 feature"
my-context note "Waiting for qa-approved-staging signal..."
# (Signal received from Carol!)
my-context note "QA approval received!"

my-context note "DECISION: Release window - Friday 2pm PST"
my-context note "Stakeholder notification sent"
my-context note "Marketing: Blog post scheduled for Monday"
my-context note "Support team: Training completed on new payment flow"

my-context stop</code></pre>
    </div>
</div>

<div class="section">
    <h2>Visual Result: Signal Chain Across 4 Roles</h2>
    <p>Each team member's context shows their part of the coordinated release:</p>

    <h3>Alice - Backend Developer</h3>
    ''' + embed_panel('tutorial-07-release-alice_explorer.html', 'Alice\'s Backend Release Context') + '''

    <h3>Bob - Frontend Developer</h3>
    ''' + embed_panel('tutorial-07-release-bob_explorer.html', 'Bob\'s Frontend Integration Context') + '''

    <h3>Carol - QA Engineer</h3>
    ''' + embed_panel('tutorial-07-release-carol_explorer.html', 'Carol\'s Testing Context') + '''

    <h3>Eve - Product Owner</h3>
    ''' + embed_panel('tutorial-07-release-eve_explorer.html', 'Eve\'s Release Coordination Context') + '''
</div>

<div class="takeaways">
    <h3>✅ Key Takeaways for Signal Coordination</h3>
    <ul>
        <li><strong>Event-driven workflows</strong>: Signals enable real-time coordination without constant checking</li>
        <li><strong>Clear dependencies</strong>: Each team member knows exactly what they're waiting for</li>
        <li><strong>Automated handoffs</strong>: No manual "I'm done" messages needed in Slack</li>
        <li><strong>Audit trail</strong>: Signal history shows when each stage completed</li>
        <li><strong>Scalable coordination</strong>: Works for 2 people or 20</li>
    </ul>
</div>

<div class="section">
    <h2>Advanced Signal Patterns</h2>

    <h3>Pattern 1: Watch Multiple Signals</h3>
    <pre><code># Wait for both backend AND database migration
my-context signal watch api-ready
my-context signal watch db-migration-complete

# Proceed only when both are ready</code></pre>

    <h3>Pattern 2: Broadcast Signal to Multiple Teams</h3>
    <pre><code># Backend creates a signal
my-context signal create api-breaking-change

# Multiple teams watch for it
# Frontend team watches
# Mobile team watches
# Partner integration team watches</code></pre>

    <h3>Pattern 3: Conditional Signals</h3>
    <pre><code># Create different signals based on test results
if [ "$ALL_TESTS_PASSED" = "true" ]; then
  my-context signal create qa-approved
else
  my-context signal create qa-blocked
  my-context note "BLOCKER: 3 tests failing, see report.html"
fi</code></pre>
</div>

<div class="section">
    <h2>Try It Yourself</h2>
    <p>Practice signal coordination with a simplified pipeline:</p>
    <pre><code># Developer 1: Backend
my-context start "backend-work"
my-context note "API endpoints deployed"
my-context signal create backend-ready
my-context stop

# Developer 2: Frontend (watches for backend)
my-context start "frontend-work"
my-context signal watch backend-ready
my-context note "Backend ready! Integrating now..."
my-context signal create frontend-ready
my-context stop

# Developer 3: QA (watches for frontend)
my-context start "qa-work"
my-context signal watch frontend-ready
my-context note "Frontend ready! Running tests..."
my-context note "✅ All tests passed"
my-context signal create qa-approved
my-context stop</code></pre>
</div>
'''

    html = create_html_page(
        title="Tutorial 7: Signal Coordination",
        subtitle="Real-Time Team Coordination with Signals",
        content=content,
        prev_link="../tutorial-06/tutorial-06.html",
        next_link="../tutorial-08/tutorial-08.html"
    )

    output_file = TUTORIALS_DIR / "tutorial-07" / "tutorial-07.html"
    output_file.write_text(html, encoding='utf-8')
    print(f"  ✅ Created: {output_file}")

def build_tutorial_08():
    """Build Tutorial 8: Agents as Team Members"""
    print("Building Tutorial 8...")

    content = '''
<div class="section scenario">
    <p><strong>Team Members:</strong>
        <span class="character backend">Alice - Human Developer</span>
        <span class="character agent">Claude - AI Coding Agent</span>
        <span class="character agent">CI/CD - Automation Agent</span>
        <span class="character agent">QA Bot - Testing Agent</span>
    </p>
    <p><strong>Scenario:</strong> Alice is implementing OAuth 2.0 integration. She works alongside three AI agents: Claude (code generation), CI/CD (build automation), and QA Bot (automated testing). Each agent creates its own context as a first-class team member.</p>
    <p><strong>Collaboration Pattern:</strong> Human-AI agent collaboration with signal coordination</p>
</div>

<div class="objectives">
    <h3>🎯 Learning Objectives</h3>
    <ul>
        <li>See how AI agents use my-context as first-class team members</li>
        <li>Understand agent-to-agent signal coordination</li>
        <li>Track automated workflows with context history</li>
        <li>Enable human oversight of agent work through contexts</li>
    </ul>
</div>

<div class="section">
    <h2>The Agent-Augmented Workflow</h2>

    <div class="step">
        <h4><span class="step-number">1</span> Alice: Start OAuth Feature</h4>
        <p class="step-description">Alice starts the feature work and hands off code generation to Claude agent.</p>
        <pre><code># Alice starts the feature
my-context start "oauth-integration" --project backend-services --labels feature,backend

my-context note "Implementing OAuth 2.0 client flow"
my-context note "Providers: Google, GitHub"
my-context note "Using authorization code flow with PKCE"

my-context stop

# Alice signals Claude agent to help with implementation
my-context signal create code-assistance-requested</code></pre>
    </div>

    <div class="step">
        <h4><span class="step-number">2</span> Claude Agent: Generate OAuth Code</h4>
        <p class="step-description">Claude watches for Alice's signal, then creates a context for code generation.</p>
        <pre><code># Claude agent watches for signal, then acts
my-context start "oauth-code-assistance" --project backend-services --created-by claude-agent

my-context note "Parent context: oauth-integration (Alice)"
my-context note "Generated OAuth client boilerplate code"
my-context note "DECISION: Using golang.org/x/oauth2 library (official Google package)"
my-context note "Implemented token storage with encryption"
my-context note "Added automatic token refresh logic"

my-context file internal/auth/oauth_client.go
my-context file internal/auth/token_storage.go

# Signal Alice that code is ready for review
my-context signal create code-ready-for-review

my-context stop</code></pre>
    </div>

    <div class="step">
        <h4><span class="step-number">3</span> Alice: Review and Enhance Agent Code</h4>
        <p class="step-description">Alice reviews Claude's code and adds enhancements.</p>
        <pre><code># Alice reviews Claude's work
my-context start "oauth-integration-review" --project backend-services

my-context note "Reviewed Claude agent's generated code"
my-context note "Code quality: Excellent, follows Go best practices"
my-context note "Added error handling for network failures"
my-context note "Added logging for OAuth flow debugging"

# Signal that feature is ready for CI
my-context signal create feature-ready-for-ci

my-context stop</code></pre>
    </div>

    <div class="step">
        <h4><span class="step-number">4</span> CI/CD Agent: Build and Test</h4>
        <p class="step-description">CI/CD agent watches for Alice's signal, then runs the build pipeline.</p>
        <pre><code># CI/CD agent triggered by signal
my-context start "build-oauth-feature" --project backend-services --created-by cicd-agent

my-context note "Parent context: oauth-integration (Alice)"
my-context note "Build #5598 triggered by commit abc123def456"
my-context note "✅ Stage 1: Lint - No issues found"
my-context note "✅ Stage 2: Unit tests - 127/127 passed (0.8s)"
my-context note "✅ Stage 3: Integration tests - 43/43 passed (12.3s)"
my-context note "✅ Stage 4: Code coverage - 94.2% (target: 90%)"
my-context note "✅ Stage 5: Security scan - No vulnerabilities"
my-context note "Build artifacts uploaded to S3"

# Signal QA bot that build passed
my-context signal create ci-build-passed

my-context stop</code></pre>
    </div>

    <div class="step">
        <h4><span class="step-number">5</span> QA Bot: E2E Testing</h4>
        <p class="step-description">QA bot watches for CI success, then runs end-to-end tests.</p>
        <pre><code># QA bot triggered by CI signal
my-context start "e2e-oauth-flow" --project qa-automation --created-by qa-bot

my-context note "Parent context: oauth-integration (Alice)"
my-context note "Test suite: OAuth end-to-end flows"
my-context note "✅ Test: Authorization code flow (Google)"
my-context note "✅ Test: Authorization code flow (GitHub)"
my-context note "✅ Test: Token refresh on expiry"
my-context note "✅ Test: Handle invalid authorization code"
my-context note "✅ Test: Handle expired refresh token"
my-context note "Test duration: 45.2 seconds"
my-context note "All tests passed - feature ready for production"

# Signal final approval
my-context signal create qa-automated-passed

my-context stop</code></pre>
    </div>

    <div class="step">
        <h4><span class="step-number">6</span> Alice: Final Review and Deployment</h4>
        <p class="step-description">Alice reviews all agent work and approves for production.</p>
        <pre><code># Alice final review
my-context start "oauth-feature-complete" --project backend-services

my-context note "All agents completed successfully:"
my-context note "  ✅ Claude agent: Code generation"
my-context note "  ✅ CI/CD agent: Build and tests"
my-context note "  ✅ QA bot: E2E validation"
my-context note "Feature ready for production deployment"
my-context note "Deployment scheduled for: Tomorrow 10am PST"

my-context stop</code></pre>
    </div>
</div>

<div class="section">
    <h2>Visual Result: Human + Agent Contexts</h2>
    <p>Each participant (human and AI) has their own context showing their contributions:</p>

    <h3>Alice - Human Developer</h3>
    ''' + embed_panel('tutorial-08-human-alice_explorer.html', 'Alice\'s Human Developer Contexts') + '''

    <h3>Claude - AI Coding Agent</h3>
    ''' + embed_panel('tutorial-08-agent-claude_explorer.html', 'Claude Agent\'s Code Generation Context') + '''

    <h3>CI/CD - Automation Agent</h3>
    ''' + embed_panel('tutorial-08-agent-cicd_explorer.html', 'CI/CD Agent\'s Build Context') + '''

    <h3>QA Bot - Testing Agent</h3>
    ''' + embed_panel('tutorial-08-agent-qa_explorer.html', 'QA Bot\'s Testing Context') + '''
</div>

<div class="takeaways">
    <h3>✅ Key Takeaways for Agent Workflows</h3>
    <ul>
        <li><strong>Agents as teammates</strong>: AI agents use contexts just like human developers</li>
        <li><strong>Audit trail</strong>: Every agent decision is tracked and reviewable</li>
        <li><strong>Human oversight</strong>: Alice can see exactly what each agent did</li>
        <li><strong>Signal coordination</strong>: Agents coordinate with each other via signals</li>
        <li><strong>Parent contexts</strong>: Agents reference the human's original context</li>
        <li><strong>Automated workflows</strong>: Entire pipelines run autonomously with full tracking</li>
    </ul>
</div>

<div class="section">
    <h2>The Future: Agents Everywhere</h2>

    <h3>Example 1: PR Review Agent</h3>
    <pre><code># PR Review agent creates context
my-context start "pr-review-1234" --created-by review-agent

my-context note "Reviewing PR #1234: Add user authentication"
my-context note "Code quality: 9/10 - Well structured, good tests"
my-context note "Security: No vulnerabilities found"
my-context note "SUGGESTION: Add rate limiting to login endpoint"
my-context note "SUGGESTION: Consider using bcrypt cost factor 12 (currently 10)"

my-context signal create review-complete</code></pre>

    <h3>Example 2: Documentation Agent</h3>
    <pre><code># Docs agent watches for feature completion
my-context start "docs-oauth-feature" --created-by docs-agent

my-context note "Parent: oauth-integration (Alice)"
my-context note "Generated API documentation from code comments"
my-context note "Updated OAuth setup guide"
my-context note "Added troubleshooting section based on common errors"

my-context file docs/api/oauth.md
my-context file docs/guides/oauth-setup.md

my-context signal create docs-ready</code></pre>

    <h3>Example 3: Deployment Agent</h3>
    <pre><code># Deployment agent coordinates production release
my-context start "prod-deploy-oauth" --created-by deploy-agent

my-context note "Deploying OAuth feature to production"
my-context note "✅ Pre-deployment checks passed"
my-context note "✅ Database migrations applied"
my-context note "✅ Service deployed to 10 regions"
my-context note "✅ Health checks: All green"
my-context note "✅ Rollback plan prepared"

my-context signal create deployment-complete</code></pre>
</div>

<div class="section">
    <h2>Try It Yourself</h2>
    <p>Simulate an agent workflow:</p>
    <pre><code># Human: Start feature
my-context start "new-feature" --project my-app
my-context note "Building user profile page"
my-context signal create code-needed
my-context stop

# Simulate AI agent response
my-context start "code-gen-profile" --created-by ai-agent
my-context note "Parent: new-feature"
my-context note "Generated React component for user profile"
my-context file src/components/UserProfile.tsx
my-context signal create code-ready
my-context stop

# Human: Review and approve
my-context start "review-profile" --project my-app
my-context note "Reviewed AI-generated code - looks good!"
my-context note "Added custom styling"
my-context signal create ready-for-testing
my-context stop</code></pre>
</div>
'''

    html = create_html_page(
        title="Tutorial 8: Agents as Team Members",
        subtitle="AI Agents as First-Class Context Users",
        content=content,
        prev_link="../tutorial-07/tutorial-07.html",
        next_link="../INDEX.html"
    )

    output_file = TUTORIALS_DIR / "tutorial-08" / "tutorial-08.html"
    output_file.write_text(html, encoding='utf-8')
    print(f"  ✅ Created: {output_file}")

def build_all_tutorials():
    """Build all tutorial HTML pages"""
    print("=" * 70)
    print("BUILDING ALL TUTORIAL HTML PAGES (1-8)")
    print("=" * 70)
    print()

    build_tutorial_01()
    build_tutorial_02()
    build_tutorial_03()
    build_tutorial_04()
    build_tutorial_05()
    build_tutorial_06()
    build_tutorial_07()
    build_tutorial_08()

    print()
    print("=" * 70)
    print("✅ TUTORIAL HTML BUILD COMPLETE")
    print("=" * 70)
    print()
    print(f"Total tutorials built: 8")
    print(f"Tutorial hub: INDEX.html")
    print()
    print("🎉 ALL TUTORIALS COMPLETE! 🎉")
    print("All 8 tutorials ready to view!")

if __name__ == "__main__":
    build_all_tutorials()
