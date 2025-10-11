# Workflow Methodology: Release Train Model

**Version**: 2.0.0
**Date**: 2025-10-11
**Supersedes**: Original 1-spec-per-sprint model
**Applies To**: All projects (my-context, deb-sanity, ps-cli, etc.)

---

## Executive Summary

**Problem**: Original "1 spec = 1 sprint, sequential execution" model doesn't scale
- Bottlenecks on constitution updates
- Can't start spec N+1 until spec N fully complete
- Arbitrary time constraints
- Constant roadmap thrashing

**Solution**: **Release Train Model**
- Specs have independent lifecycles (can be at any stage simultaneously)
- Sprint = integration event when specs ready (not calendar-based)
- Continuous dev branch (always deployable)
- Release train (version when value accumulated)
- Constitution validates throughout (not end-gate)

**Result**: Scalable, parallel development, rapid iteration, clear process

---

## Core Definitions

### Spec
**Unit of work** defining a cohesive feature or improvement
- Has independent lifecycle (draft → merged)
- Tracked in specs/spec-NNN-feature-name/ directory
- Can be at any stage (draft, implementing, testing, etc.)
- Multiple specs can exist in parallel

### Sprint
**Integration event** (not time period) triggered when 1+ specs ready to merge
- Merges completed specs to dev branch
- Runs integration tests
- Validates against constitution
- Decides version bump
- Duration: Hours to days (based on scope, not calendar)

### Iteration
**Synonym for Spec** in development context
- "Iteration 006" = spec-006-signaling
- Emphasizes incremental delivery
- Better than "Sprint" (which means integration event)

### Release Train
**Pattern** of regular versioned releases when value accumulated
- Not time-based (release when features ready)
- Not arbitrary (release when coherent value delivered)
- Tags: v2.0.0, v2.1.0, v2.2.0, v3.0.0

---

## Branch Strategy

### Permanent Branches (2 - Never Delete)

**main**
- Purpose: Public releases only
- Tags: v1.0.0, v2.0.0, v3.0.0, etc.
- Updates: Only from dev branch during releases
- Audience: Public users, production deployments

**dev**
- Purpose: Continuous integration branch
- Content: Merged specs, always deployable
- Updates: Receives merged spec branches continuously
- Audience: Development team, internal testing

**Workflow**:
```
spec branches → merge to master → push to origin/dev → accumulate
When ready: dev → main + version tag → release
```

---

### Temporary Branches (Delete After Merge)

**spec-NNN-feature-name**
- Purpose: Develop one spec/feature
- Lifetime: Created → implemented → tested → merged → DELETED
- Example: spec-006-signaling, spec-007-integration

**hotfix-NNN-issue**
- Purpose: Urgent production fixes
- Lifetime: Created → fixed → merged → DELETED quickly
- Example: hotfix-001-crash-fix

**release-prep-vX.Y.Z** (Optional)
- Purpose: Stage major releases for final validation
- Lifetime: Created → validated → merged → DELETED
- Example: release-prep-v3.0.0

**Max Simultaneous**: 2-5 spec branches (more = need better prioritization)

---

## Spec Lifecycle: 10 Stages

### Stage 0: draft
**Entry**: Feature idea or requirement identified
**Activities**:
- Write spec using /specify command (or manually)
- Identify requirements, user stories
- Mark unclear areas with [NEEDS CLARIFICATION]

**Ceremony**: `/specify <feature description>`
**Journal**: `{project}: spec-NNN__planning` (STANDARD level, 20-30 notes)
**Quality Gate**: 
- Spec template followed
- Mandatory sections present (User Stories, Requirements, Success Criteria)
- Independent testability for user stories

**Artifacts**:
- specs/spec-NNN-feature-name/spec.md
- Branch: spec-NNN-feature-name created from master

**Exit Criteria**: Spec written, ready for clarification

---

### Stage 1: clarified
**Entry**: Spec has open questions

**Activities**:
- Answer all [NEEDS CLARIFICATION] markers
- Document decisions with rationale
- Resolve ambiguities

**Ceremony**: `/clarify` (spec kit command)
**Journal**: Continue spec-NNN__planning context
**Quality Gate**:
- Zero [NEEDS CLARIFICATION] remaining
- All decisions documented
- Rationale provided for each decision

**Artifacts**:
- specs/spec-NNN-feature-name/clarify.md

**Exit Criteria**: No ambiguities, ready to plan implementation

---

### Stage 2: planned
**Entry**: Requirements clear, ready for design

**Activities**:
- Create implementation plan (architecture, tech stack, phases)
- Generate design artifacts (data-model.md, contracts/, etc.)
- Identify dependencies and risks

**Ceremony**: `/plan` (spec kit command)
**Journal**: Continue spec-NNN__planning context
**Quality Gate**:
- Tech stack defined
- Architecture documented
- Phases identified
- Risks addressed

**Artifacts**:
- specs/spec-NNN-feature-name/plan.md
- specs/spec-NNN-feature-name/data-model.md (if applicable)
- specs/spec-NNN-feature-name/contracts/ (if applicable)

**Exit Criteria**: Implementation approach clear, ready for task breakdown

---

### Stage 3: tasked
**Entry**: Plan complete, ready for task generation

**Activities**:
- Generate numbered task list (T001, T002, ...)
- Specify file paths for each task
- Mark parallel tasks [P]
- Document dependencies

**Ceremony**: `/tasks` (spec kit command)
**Journal**: Continue spec-NNN__planning context, consider stopping
**Quality Gate**:
- All tasks numbered sequentially
- File paths specified
- Dependencies clear
- Parallel opportunities marked

**Artifacts**:
- specs/spec-NNN-feature-name/tasks.md

**Exit Criteria**: Task list complete, ready for implementation

---

### Stage 4: checklist-validated
**Entry**: Tasks ready, need requirements quality check

**Activities**:
- Generate requirements quality checklist
- Validate spec completeness, clarity, consistency
- Check acceptance criteria quality

**Ceremony**: `/checklist` (Cursor spec kit command)
**Journal**: Note checklist results
**Quality Gate**:
- All checklist items passing OR
- Documented exceptions for incomplete items

**Artifacts**:
- specs/spec-NNN-feature-name/checklists/requirements.md

**Exit Criteria**: Requirements quality validated, implementation can begin

---

### Stage 5: implementing
**Entry**: Ready to write code

**Activities**:
- Execute tasks in order
- Follow TDD approach (tests first)
- Track progress with checkpoints
- Signal heartbeats, phase boundaries

**Ceremony**: `/implement` (spec kit command) OR manual if supervised
**Journal**: `{project}: spec-NNN__implementation-phaseN` (VERBOSE level)
**Journaling Protocol**:
- Start: "🚀 Starting Phase N: {name} (Tasks {X}-{Y})"
- Heartbeat: `my-context touch` every 10 minutes
- Checkpoint: Every 10-15 tasks completed
- File tracking: Significant changes (>10 lines)
- Phase end: "✅ Phase N complete: {summary} ({done}/{total} tasks)"
- Target: 50 notes per phase, split if approaching 75

**Quality Gate**:
- All tasks marked [X] complete
- Tests passing
- No regressions in existing functionality
- Code builds successfully

**Artifacts**:
- Source code changes
- Test files
- Updated tasks.md (all [X])

**Exit Criteria**: All tasks complete, tests passing, ready for UAT

---

### Stage 6: testing
**Entry**: Implementation complete, needs validation

**Activities**:
- User acceptance testing
- Performance validation
- Integration testing
- Bug discovery and fixing

**Ceremony**: Manual testing, feedback collection
**Journal**: `{project}: spec-NNN__validation` (DETAILED level, 30-50 notes)
**Journaling Protocol**:
- Test results for each scenario
- Bugs discovered
- Performance metrics
- User feedback

**Quality Gate**:
- All acceptance scenarios validated
- Critical bugs fixed
- Success criteria met
- Performance targets achieved

**Artifacts**:
- Test results documentation
- Bug reports and fixes
- Validation evidence

**Exit Criteria**: UAT passed, ready for review

---

### Stage 7: review
**Entry**: Testing complete, needs approval to merge

**Activities**:
- Code review
- Spec review
- Constitution validation check
- Integration impact assessment

**Ceremony**: Manual review or automated checks
**Journal**: Note review feedback
**Quality Gate**:
- No blocking issues
- Constitution principles validated
- Integration conflicts identified (if any)
- Documentation complete

**Artifacts**:
- Review feedback
- Constitution validation notes

**Exit Criteria**: Approved for merge, no blockers

**TRIGGERS SPRINT**: Spec at stage 7 is ready for integration event

---

### Stage 8: merged
**Entry**: Spec approved, sprint integration occurring

**Activities**:
- Merge spec branch to master
- Resolve conflicts (if any)
- Push to origin/dev
- Run integration tests with other merged specs

**Ceremony**: Git merge (part of sprint integration event)
**Journal**: `{project}: sprint-integration-vX.Y.Z` (STANDARD level)
**Quality Gate**:
- Merge successful
- No conflicts (or resolved)
- Integration tests passing
- Dev branch deployable

**Artifacts**:
- Merged code in master/dev
- Integration test results

**Exit Criteria**: Integrated to dev branch successfully

---

### Stage 9: released
**Entry**: Sprint complete, version tagged

**Activities**:
- Version tag applied
- CHANGELOG updated
- Main branch updated (dev → main)
- Release published (GitHub, docs, etc.)

**Ceremony**: Release preparation and tagging
**Journal**: Continue sprint-integration context
**Quality Gate**:
- Version number decided
- CHANGELOG complete
- Documentation updated
- Release notes written

**Artifacts**:
- Git tag (vX.Y.Z)
- Updated CHANGELOG
- Release notes

**Exit Criteria**: Version released, publicly available

---

### Stage 10: closed
**Entry**: Spec released, cleanup needed

**Activities**:
- Delete spec branch (local + GitHub)
- Archive all spec-NNN contexts
- Update project status
- Close related issues/PRs

**Ceremony**: Branch and context cleanup
**Journal**: Final note in sprint-integration context
**Quality Gate**:
- No orphaned branches
- No orphaned contexts
- Documentation references updated

**Artifacts**:
- Clean branch state
- Archived contexts

**Exit Criteria**: Complete cleanup, ready for next specs

---

## Sprint as Integration Event

### Trigger Conditions

**Sprint occurs when**: 1 or more specs reach Stage 7 (review)

**Frequency**: Variable (could be daily, weekly, or monthly based on development pace)

**Not triggered by**:
- Calendar dates (not "every 2 weeks")
- Arbitrary milestones
- External schedules

**Triggered by**: Work readiness

---

### Sprint Composition

**Small Sprint** (1-2 specs, quick fixes):
```
Sprint for v2.2.1:
└── spec-012-bugfix (1 hour of work)

Duration: 2-3 hours (merge + test + release)
```

**Medium Sprint** (2-4 specs, mixed sizes):
```
Sprint for v2.3.0:
├── spec-006-signaling (large, 1 week work)
├── spec-011-docs-update (small, 2 hours)
└── spec-013-minor-improvement (medium, 1 day)

Duration: 4-8 hours (integration + testing + release)
```

**Large Sprint** (major version):
```
Sprint for v3.0.0:
├── spec-008-multi-channel (2 weeks)
├── spec-009-branching (1 week)
├── spec-014-breaking-change (3 days)

Duration: 1-2 days (integration + extensive testing + migration docs)
```

---

### Sprint Activities (Ceremony)

**Phase 1: Pre-Integration** (30min)
```
Journal: my-context start "{project}: sprint-integration-v2.2.0"

Activities:
- List all specs at stage 7 (ready to merge)
- Review for conflicts (same files modified?)
- Decide merge order
- Plan integration testing

Journal notes:
- "🚀 Sprint integration starting: v2.2.0"
- "📋 Specs ready: spec-006, spec-011, spec-013 (3 total)"
- "⚠️ Potential conflict: spec-006 and spec-013 both modify state.go"
```

**Phase 2: Integration** (1-4 hours)
```
For each spec (in order):
  1. Merge spec-NNN to master
  2. Resolve conflicts if any
  3. Run spec's tests
  4. my-context note "✅ spec-NNN integrated"
  5. my-context file <changed-files>

After all merged:
  1. Push master → origin/dev
  2. Run integration test suite
  3. my-context note "✅ All specs integrated: {count} specs merged"
```

**Phase 3: Constitution Validation** (30min-1h)
```
Review merged work:
- Check clarity (is code understandable?)
- Check completeness (all features working?)
- Check quality (tests passing, no hacks?)
- Check value (delivers user benefit?)

Journal:
- my-context note "📖 Constitution check: {principle} validated"
- Document any constitution learnings
- Note if constitution needs updates (rare)
```

**Phase 4: Version Decision** (15min)
```
Evaluate changes:
- Patch (v2.2.1): Bug fixes only, no new features
- Minor (v2.3.0): New features, backward compatible
- Major (v3.0.0): Breaking changes, new capabilities

Journal:
- my-context note "📦 Version decision: v{X.Y.Z} ({rationale})"
```

**Phase 5: Release Preparation** (30min-1h)
```
Activities:
1. Update CHANGELOG.md
2. Update README.md (if needed)
3. Update version in code (if hardcoded)
4. Build and test final binary
5. Run smoke tests

Journal:
- my-context file CHANGELOG.md
- my-context note "✅ CHANGELOG updated: {N} entries"
- my-context note "✅ Binary built and tested: v{X.Y.Z}"
```

**Phase 6: Release** (15min)
```
Activities:
1. Merge dev → main (or fast-forward)
2. Tag version: git tag -a vX.Y.Z -m "{message}"
3. Push tag: git push origin vX.Y.Z
4. GitHub Actions builds (if configured)
5. Publish release notes

Journal:
- my-context note "🎉 v{X.Y.Z} released to origin/main"
- my-context note "🔗 Release: {GitHub URL}"
```

**Phase 7: Cleanup** (15min)
```
Activities:
1. Delete merged spec branches (local + GitHub)
2. Archive sprint integration context
3. Archive individual spec contexts
4. Update project status documents

Journal:
- my-context note "🧹 Cleanup: {N} branches deleted"
- my-context stop
- my-context archive "{project}: sprint-integration-v{X.Y.Z}"
```

**Total Sprint Duration**: 3-8 hours (flexible based on scope)

---

## Journaling Policies

### By Lifecycle Stage

| Stage | Context Pattern | Level | Target Notes | Heartbeat |
|-------|----------------|-------|--------------|-----------|
| 0-3 (Planning) | spec-NNN__planning | STANDARD | 20-30 | No |
| 4 (Checklist) | spec-NNN__planning | MINIMAL | 5-10 | No |
| 5 (Implementing) | spec-NNN__impl-phaseN | VERBOSE | 50/phase | Yes (10min) |
| 6 (Testing) | spec-NNN__validation | DETAILED | 30-50 | No |
| 7 (Review) | spec-NNN__review | MINIMAL | 5-10 | No |
| 8-10 (Sprint) | sprint-integration-vX.Y.Z | STANDARD | 10-20 | No |

### Journaling Levels Defined

**MINIMAL** (5-10 notes):
- Major milestones only
- Entry/exit of stage
- Critical decisions
- Use for: Quick stages (checklist, review)

**STANDARD** (20-30 notes):
- Significant activities
- Key decisions
- Important observations
- Use for: Planning, sprint integration

**DETAILED** (30-50 notes):
- Comprehensive documentation
- All test results
- Bug tracking
- Use for: Testing/validation

**VERBOSE** (50+ notes per phase):
- Frequent checkpoints (every 10-15 tasks)
- Heartbeats (touch every 10min)
- All significant events
- Phase boundaries
- Use for: Implementation

### Heartbeat Protocol (Implementation Stage)

**Purpose**: Prove agent/human is making progress, enable supervisor monitoring

**Pattern**:
```bash
# Every 10 minutes during active implementation
my-context touch

# Every 10-15 tasks
my-context note "📊 Checkpoint: Tasks T010-T025 complete - {what was built}"

# At phase boundaries
my-context note "🚀 Starting Phase 2: {name} (Tasks T026-T050)"
my-context note "✅ Phase 1 complete: {summary} (25/25 tasks)"

# When blocked
my-context note "🚧 BLOCKED: {issue} - investigating"
my-context note "✅ UNBLOCKED: {resolution} - continuing"
```

**Supervisor Monitoring**:
```bash
# Manual check
my-context show spec-NNN__impl-phase1 | tail -20

# Automated (until Iteration 006 watch command)
watch -n 60 'my-context show | tail -20'

# After Iteration 006
my-context watch spec-NNN --new-notes --exec="notify-progress.sh"
```

---

## Quality Gates

### Spec Kit Command Gates

**Each command is a ceremony with quality gate**:

| Command | Stage Transition | Gate Check |
|---------|-----------------|------------|
| /specify | → draft | Template followed, sections complete |
| /clarify | draft → clarified | No [NEEDS CLARIFICATION] remaining |
| /plan | clarified → planned | Architecture documented, tech stack defined |
| /tasks | planned → tasked | All tasks numbered, dependencies clear |
| /checklist | tasked → validated | Requirements quality checked |
| /implement | validated → implementing | Checklists passing, ready to code |

**Gates prevent**:
- Under-specified features
- Missing requirements
- Unclear architecture
- Incomplete task lists

**Gates allow**: Fast-tracking (can skip for hotfixes, document exception)

---

### Test Gates

**Unit Tests**: Must pass before marking tasks complete
**Integration Tests**: Must pass before stage 6 (testing)
**UAT**: Must pass before stage 7 (review)
**Regression Tests**: Must pass during sprint integration

**Gate Policy**: Red tests block progression (fix before continuing)

---

### Constitution Gates (Continuous, Not Blocking)

**Validation Points** (Throughout):
- Stage 0: Spec aligns with principles?
- Stage 2: Architecture follows guidelines?
- Stage 5: Code meets quality standards?
- Stage 6: Delivers user value?
- Stage 8: Integration maintains coherence?

**Not Blocking**: Constitution informs, doesn't block
**Updates Constitution**: Significant learnings captured, not every spec
**Frequency**: When principles evolve (not routine)

---

## Parallel Spec Development

### How Multiple Specs Coexist

**Example State** (All Independent):
```
specs/ directory:
├── spec-003-summary/ (Stage 1: clarified, paused for data)
├── spec-006-signaling/ (Stage 5: implementing, Phase 2 of 3)
├── spec-007-integration/ (Stage 3: tasked, ready to implement)
├── spec-008-channels/ (Stage 2: planned, generating tasks)
├── spec-011-bugfix/ (Stage 7: review, ready to merge NOW)

Active Contexts:
├── my-context: spec-006__impl-phase2 (45 notes, implementing)
├── my-context: spec-007__planning (12 notes, stopped - ready to resume)
├── my-context: spec-008__planning (8 notes, active)

Next Sprint: Merge spec-011 (ready now)
Following Sprint: spec-006 + spec-007 when both reach stage 7
```

**No blocking** - each spec progresses at its own pace

---

### Prioritization

**High Priority Specs**: More resources, parallel agents
**Low Priority Specs**: Progress as time allows
**Paused Specs**: Stop at any stage, resume when ready

**Resource Allocation**:
- 1 spec with Agent A (full-time implementation)
- 2 specs with Agent B (part-time planning)
- 1 spec paused (waiting for data/decision)

---

## Rapid Deployment

### Hotfix Pattern (Fast Track)

**Trigger**: Critical production bug

**Compressed Workflow** (1-2 hours):
```
1. Create branch: hotfix-NNN-issue from master (2min)
2. Minimal spec: Problem + solution only (10min)
3. Implement fix (30min)
4. Test fix (15min)
5. Merge to master → dev (5min)
6. IMMEDIATE SPRINT: Release as patch version (30min)
   - Skip accumulation
   - Skip constitution review (document exception)
   - Fast-track to production

Journal:
- Single context: hotfix-NNN (rapid notes, ~10 total)
- Compressed ceremony (all stages in one context)

Total: 1.5 hours bug → production
```

**Gates Can Skip**: Planning depth, extensive testing (document why)
**Gates Cannot Skip**: Tests must pass, no regressions

---

### Feature Fast Track (Quick Wins)

**For small features** (2-4 hours total work):
```
1. Lightweight spec (1 page, minimal clarification)
2. Quick implementation (1-2 hours)
3. Basic testing (30min)
4. Merge and release

Can combine multiple quick wins into single sprint
Example: Sprint for v2.1.1 might have 3-4 small fixes/improvements
```

---

## Continuous Dev Branch

### Purpose

**dev branch** = Integration branch (always deployable)

**Characteristics**:
- Receives merged specs continuously
- Always builds successfully
- All tests passing
- Can deploy any commit
- Accumulates features between releases

**NOT**:
- Not unstable/experimental (that's feature branches)
- Not "staging" (it's THE integration point)
- Not optional (it's required for release train)

---

### Dev Branch Workflow

**Continuous Flow**:
```
Day 1: Merge spec-011 → master → push to dev (v2.2.0-dev)
Day 3: Merge spec-006 → master → push to dev (v2.2.0-dev+signaling)
Day 5: Merge spec-013 → master → push to dev (v2.2.0-dev+3specs)
Day 6: SPRINT → dev → main + tag v2.2.0 → release

Dev branch accumulated 3 specs over 6 days, released together
```

**Benefits**:
- Features integrate as they complete (not batched)
- Always have deployable build (any dev commit works)
- Can cherry-pick urgent fixes if needed
- Release when ready (not calendar-driven)

---

## Release Train

### Version Strategy

**Semantic Versioning**:
- **Patch** (v2.2.1): Bug fixes, no new features
- **Minor** (v2.3.0): New features, backward compatible
- **Major** (v3.0.0): Breaking changes, new capabilities

**Release Frequency**: When value accumulated
- Small patch: 1-3 days after bugfix
- Minor release: 1-2 weeks after features complete
- Major release: 1-3 months for significant capabilities

**Not Time-Based**: Release when coherent value delivered, not arbitrary dates

---

### Release Decision Framework

**Trigger Sprint When**:
- 1+ specs at stage 7 (ready to merge)
- Urgent fix needed (hotfix pattern)
- Significant value accumulated on dev (3+ features)
- Major milestone reached (v3.0.0 planning)

**Version Decision**:
```
Review merged specs:
- Any breaking changes? → Major version
- All backward compatible features? → Minor version
- Only bug fixes? → Patch version
- Mix? → Minor version (document breaking changes in CHANGELOG)
```

---

## Cross-Project Template

### Apply to All Projects

**Directory Structure** (Standard):
```
project/
├── .specify/
│   ├── config/
│   │   ├── journaling-policy.yml
│   │   └── workflow-config.yml
│   ├── scripts/
│   │   └── context-tracking.sh
│   └── templates/
│       └── (spec kit templates)
├── specs/
│   ├── spec-001-{name}/
│   │   ├── .status (file: "implementing")
│   │   ├── spec.md
│   │   ├── clarify.md
│   │   ├── plan.md
│   │   └── tasks.md
│   └── spec-NNN-{name}/
├── docs/
│   ├── WORKFLOW-METHODOLOGY.md (this document)
│   ├── ROADMAP.md
│   └── CHANGELOG.md
├── .github/workflows/ (CI/CD)
└── (source code)
```

**All Projects Use**:
- Same spec lifecycle (10 stages)
- Same journaling policies
- Same quality gates
- Same release train
- Same branch strategy

---

## Implementation Guide

### Starting New Spec

```bash
# 1. Create spec branch from master
git checkout master
git pull origin dev  # Get latest
git checkout -b spec-NNN-feature-name

# 2. Start planning context
my-context start "{project}: spec-NNN__planning"
my-context note "Starting spec for: {feature description}"

# 3. Run spec kit
/specify "{feature description}"
# Generates: specs/spec-NNN-feature-name/spec.md

# 4. Progress through stages
/clarify
/plan
/tasks
/checklist  # If using Cursor

# 5. Stop planning, start implementation
my-context stop
my-context start "{project}: spec-NNN__impl-phase1"

# 6. Implement
/implement  # Or manual with supervision

# 7. Test
my-context start "{project}: spec-NNN__validation"
# Run UAT, document results

# 8. Mark ready
# Spec is now at stage 7 (review) - ready for next sprint
```

---

### Running a Sprint (Integration Event)

```bash
# Pre-requisite: 1+ specs at stage 7

# 1. Start sprint context
my-context start "{project}: sprint-integration-v{X.Y.Z}"
my-context note "🚀 Sprint starting: {N} specs ready"

# 2. Merge each spec
for spec in {ready-specs}; do
  git merge $spec --no-ff -m "Merge $spec"
  my-context note "✅ $spec merged"
done

# 3. Push to dev
git push origin master:dev
my-context note "✅ Pushed to origin/dev: {N} specs integrated"

# 4. Integration testing
run-integration-tests.sh
my-context note "✅ Integration tests: {passing}/{total}"

# 5. Constitution check
review-against-principles.sh
my-context note "📖 Constitution validated: {principles checked}"

# 6. Version and release
decide-version.sh  # Determines patch/minor/major
git tag -a v{X.Y.Z} -m "{release notes}"
git push origin v{X.Y.Z}
git push origin master:main  # Update public branch
my-context note "🎉 v{X.Y.Z} released"

# 7. Cleanup
for spec in {merged-specs}; do
  git branch -d $spec
  git push origin --delete $spec
done
my-context note "🧹 {N} branches cleaned up"
my-context stop
my-context archive "{project}: sprint-integration-v{X.Y.Z}"
```

---

## Adoption Plan (All Projects)

### Phase 1: Document (1 day)
1. Create WORKFLOW-METHODOLOGY.md (this document)
2. Update SDLC.md to reference new model
3. Add .specify/config/journaling-policy.yml template
4. Document in each project README

### Phase 2: Migrate Current State (2-4 hours)
1. Close 005-ux-polish-and (absorbed)
2. Decide on 003-daily-summary (pause or close)
3. Rename existing "Sprint" references to "Iteration" in docs
4. Add .status files to existing specs

### Phase 3: Apply to Next Spec (Iteration 006)
1. Use new workflow for spec-006-signaling
2. Validate methodology works
3. Refine based on learnings

### Phase 4: Rollout to Other Projects (1 week)
1. Apply to deb-sanity (spec-007 using new model)
2. Apply to ps-cli (if active development)
3. Train team on new process

---

## Key Benefits

### Scalability
✅ Multiple specs in parallel (no sequential bottleneck)
✅ Each progresses independently
✅ Resources allocated flexibly

### Flexibility
✅ Sprint when ready (not calendar-based)
✅ Specs can pause/resume at any stage
✅ Rapid deployment available (hotfix pattern)

### Quality
✅ Gates at each stage (prevent low-quality work)
✅ Continuous constitution validation
✅ Integration testing before release

### Clarity
✅ Clear stages (know where each spec is)
✅ Clear ceremonies (spec kit commands)
✅ Clear branch strategy (2 permanent + temp features)

### Auditability
✅ Journaling policies (appropriate to stage)
✅ Complete context trail
✅ Version history clear

---

## Migration from Old Model

**Old Approach** (Sequential):
```
003-summary: Needs clarification, BLOCKS everything
↓ (can't start 004 until 003 done)
004-lifecycle: Waiting...
```

**New Approach** (Parallel):
```
003-summary: Stage 1 (paused - waiting for data)
004-lifecycle: Stage 9 (released as v2.0.0) ✅
005-ux-polish: Stage 0 (absorbed into 007)
006-signaling: Stage 0 (ready to start) ← NEXT
007-integration: Stage 0 (can start anytime)
```

**Impact**: 004 shipped while 003 still paused - no blocking!

---

## Example: Next 2 Weeks

**Week 1** (Iteration 006):
```
Day 1: spec-006 stages 0-3 (specify/clarify/plan/tasks)
Day 2-3: spec-006 stage 5 (implementing phase 1-2)
Day 4: spec-006 stage 5 (implementing phase 3)
Day 5: spec-006 stage 6 (testing)

Parallel: spec-007 stages 0-2 (planning while 006 implements)
```

**Week 2** (Sprint + Iteration 007):
```
Day 6: spec-006 stage 7 (review)
Day 6: SPRINT - Integrate spec-006 → v2.2.0 release (4 hours)
Day 7-8: spec-007 stage 5 (implementing)
Day 9: spec-007 stage 6 (testing)
Day 10: spec-007 stage 7 (review)

Parallel: spec-008 stages 0-2 (planning)
```

**Week 3** (Sprint + Continue):
```
Day 11: SPRINT - Integrate spec-007 → v2.3.0 release (3 hours)
Day 12-15: spec-008 stage 5 (implementing)

Parallel: spec-009 stages 0-3 (planning)
```

**Result**: 3 releases in 3 weeks, 4 specs completed, continuous progress

---

## Summary: Release Train Workflow

**Specs**: Independent lifecycles, parallel development
**Sprints**: Integration events when specs ready
**Dev Branch**: Continuous integration, always deployable
**Release Train**: Version when value accumulated
**Journaling**: Policies per stage, complete auditability
**Quality**: Gates at transitions, continuous validation
**Constitution**: Validates throughout, not end-gate
**Scalable**: Works for 1 spec or 10 specs in parallel

**This model scales to multiple projects and teams!**

---

**Document Version**: 2.0.0
**Status**: Proposed methodology for adoption
**Next**: Apply to Iteration 006 (Signaling) as validation
