# deb-sanity ↔ my-context Integration Analysis

**Date**: 2025-10-10  
**Author**: AI Agent (Cursor)  
**Purpose**: Evaluate synergies, integration opportunities, and feature requests

---

## Executive Summary

**deb-sanity** and **my-context** are **complementary tools** with minimal overlap and strong potential for integration. They operate at different layers of the developer workflow:

- **deb-sanity**: Environment-level (tools, projects, services, worktrees)
- **my-context**: Session-level (work contexts, notes, decisions, time tracking)

**Key Finding**: These tools already work together organically (my-context is registered in deb-sanity's tool registry, and deb-sanity work sessions are tracked in my-context).

---

## Tool Comparison Matrix

| Feature | deb-sanity | my-context | Overlap? |
|---------|------------|------------|----------|
| **Environment reporting** | ✅ Comprehensive | ❌ N/A | None |
| **Tool registration** | ✅ Registry-based | ❌ N/A | None |
| **Project discovery** | ✅ Filesystem scan | 🔶 Manual project filter | Minimal |
| **Context management** | ❌ N/A | ✅ Session-based | None |
| **Note-taking** | ❌ N/A | ✅ Timestamped | None |
| **File tracking** | ❌ N/A | ✅ Associations | None |
| **Git worktree management** | ✅ Full suite | ❌ N/A | None |
| **Service orchestration** | ✅ Docker/services | ❌ N/A | None |
| **Path translation (WSL ↔ Win)** | ✅ Built-in | ❌ N/A | None |
| **Time tracking** | 🔶 Logs only | ✅ Duration tracking | Minimal |
| **Export/Archive** | ❌ N/A | ✅ Markdown/JSON | None |

**Overlap Score**: ~5% (very low - excellent complementarity!)

---

## Current Integration Points (Already Working!)

### 1. Tool Registry Integration ✅
```bash
$ deb-sanity --tools | grep my-context
my-context
Description: Context management tool (.exe via WSL interop)
Status: FOUND
Location: /home/be-dev-agent/.local/bin/my-context
```

**Status**: deb-sanity already knows about and tracks my-context!

### 2. Context-Aware Work Sessions ✅
```bash
$ my-context list | grep deb-sanity
○ deb-sanity.sh version 1.6.0 uat (stopped)
  Duration: 46m
```

**Status**: my-context already tracks deb-sanity development work!

---

## Integration Opportunities

### 🎯 High-Value Integrations

#### 1. **Auto-Context Creation from Project Detection**
**Concept**: When deb-sanity discovers a new project, offer to create a my-context for it.

```bash
# deb-sanity --projects (discovers new project)
New project detected: /home/user/projects/new-api

# Prompt integration
Would you like to create a my-context for this project? (y/n)
> y
Created context: "new-api: Initial setup"
```

**Implementation**:
- Add `--create-context` flag to deb-sanity
- Call `my-context start "<project-name>: Initial setup"`
- Auto-associate project directory with context

**Value**: Seamless workflow transition from project discovery → context creation

---

#### 2. **Environment Snapshot in Context Notes** ⚠️ REVISED
**Concept**: Capture environment data when starting a context (via shared data).

```bash
# No flag needed - automatic if data available
$ my-context start "Bug Fix #123"

# If ~/.dev-tools-data/environment.json exists:
# Auto-creates notes:
- [timestamp] Environment: WSL 2 / Debian trixie
- [timestamp] Tools available: go (1.21), node (18.x), docker (running)
- [timestamp] Project: /home/user/projects/api-service
- [timestamp] Active worktree: feature/bug-123

# If file doesn't exist: proceeds normally (no error)
```

**Implementation** (Revised - No Direct Dependency):
```go
// On context start
func CreateContext(name string) {
    // ... normal context creation ...
    
    // Try to read shared environment data (OPTIONAL)
    envData, err := readEnvironmentData()  // reads ~/.dev-tools-data/environment.json
    if err == nil {
        // Add environment as initial notes
        addNote(fmt.Sprintf("Environment: %s", envData.OS))
        addNote(fmt.Sprintf("Tools: %s", strings.Join(envData.Tools, ", ")))
    }
    // If error: silently continue (no warning, no failure)
}
```

**Shared Data Contract** (written by deb-sanity):
```json
{
  "updated_at": "2025-10-10T20:30:00Z",
  "os": "WSL 2 / Debian trixie",
  "tools": ["go:1.21", "docker:24.0", "node:18.x"],
  "services": [{"name": "mysql", "status": "running"}],
  "active_worktree": "feature/bug-123"
}
```

**Value**: Context includes complete environment state for reproducibility
**Risk Mitigation**: No runtime dependency - works with or without deb-sanity

---

#### 3. **Project-Context Linkage**
**Concept**: Associate my-context sessions with deb-sanity project index.

**Data Structure** (add to my-context meta.json):
```json
{
  "name": "ps-cli: Phase 2",
  "project_metadata": {
    "deb_sanity_project": "ps-cli",
    "project_path": "/mnt/c/Users/.../ps-cli-dev",
    "git_remote": "github.com/user/ps-cli",
    "worktree": "feature/enhancement"
  }
}
```

**Use Cases**:
- `my-context list --project ps-cli` uses deb-sanity's project index
- Export includes full project metadata
- Retroactively link contexts to projects

**Value**: Unified project tracking across both tools

---

#### 4. **Service State Recording**
**Concept**: Record which services were running during a context.

```bash
$ my-context start "Integration Testing"
# Auto-detects via deb-sanity
Services active: mysql, redis, docker (api-gateway container)
Note added: "Environment: 3 services running"
```

**Implementation**:
- Check `deb-sanity --services --status` on context start
- Add as initial note
- On context stop, compare state changes

**Value**: Know exactly what environment the context was in

---

### 🔧 Medium-Value Integrations

#### 5. **Worktree Context Sync**
**Concept**: Auto-switch my-context when switching git worktrees.

```bash
# In deb-sanity worktree operations
$ deb-sanity --worktree-switch feature/new-login

# Integration point:
Switching worktree to: feature/new-login
Found my-context: "ps-cli: Login Redesign"
Switch to this context? (y/n)
```

**Value**: Context follows your worktree automatically

---

#### 6. **Unified Time Tracking**
**Concept**: Aggregate time across tools.

```bash
$ my-context report --time-by-project

ps-cli project (via deb-sanity integration):
  - 18 contexts tracked
  - Total time: 47h 23m
  - Most recent: "ps-cli: Phase 3" (2h 15m)
  - Worktrees used: 5 (feature/*, hotfix/*)
```

**Value**: Complete project time analytics

---

#### 7. **Export Enhancement: Full Environment Snapshot**
**Concept**: Enrich my-context exports with deb-sanity data.

```markdown
# Context Export: API Bug Fix

## Environment
- **OS**: WSL 2 / Debian trixie x86_64
- **Tools**: go 1.21, docker 24.0, mysql-client 8.0
- **Services**: mysql (running), redis (stopped)
- **Worktree**: /home/user/projects/api/feature/bug-fix
- **Git Remote**: github.com/company/api-service

## Notes
[... rest of export ...]
```

**Value**: Exports become comprehensive documentation

---

## Feature Requests

### For **deb-sanity**

#### FR-DS-001: Context Management Awareness
**Description**: Add `--active-context` flag to show current my-context.

```bash
$ deb-sanity --active-context
Active my-context: "ps-cli: Phase 3" (started 2h ago)
Notes: 12 | Files: 5 | Touches: 3
```

**Rationale**: Quick status check for developers using both tools.

---

#### FR-DS-002: Time-Aware Project Scanning
**Description**: Show last modification time when scanning projects.

```bash
$ deb-sanity --projects --recent

ps-cli (last modified: 2h ago)
  - 3 active worktrees
  - my-context: "ps-cli: Phase 3" (active)
  
api-service (last modified: 2 days ago)
  - No active contexts
```

**Rationale**: Help identify which projects are currently being worked on.

---

#### FR-DS-003: Context History Export
**Description**: Export project activity including my-context data.

```bash
$ deb-sanity --export-activity ps-cli --format json

{
  "project": "ps-cli",
  "contexts": [
    {"name": "ps-cli: Phase 3", "duration": "2h15m", "notes": 12},
    {"name": "ps-cli: Phase 2", "duration": "5h45m", "notes": 34}
  ],
  "worktrees": [...],
  "total_time": "47h23m"
}
```

**Rationale**: Unified project reporting across both tools.

---

### For **my-context**

#### FR-MC-001: Environment Tags
**Description**: Tag contexts with environment metadata.

```bash
$ my-context start "Bug Fix" --env wsl --tools go,docker

# Later filter:
$ my-context list --env wsl --has-tool docker
```

**Rationale**: Filter contexts by environment requirements.

---

#### FR-MC-002: Project Path Association
**Description**: Link contexts to filesystem paths (not just names).

```bash
$ my-context start "API Work" --path /home/user/projects/api-service

# Later:
$ cd /home/user/projects/api-service
$ my-context list --here
# Shows contexts associated with current directory
```

**Rationale**: Better integration with project-based workflows.

---

#### FR-MC-003: Tool Dependency Tracking
**Description**: Record which tools were used during a context.

```bash
$ my-context file src/auth.go
$ my-context tool-used go
$ my-context tool-used docker

# Export shows:
## Tools Used
- go (10 files)
- docker (mentioned in 3 notes)
```

**Rationale**: Know what tools are needed to reproduce work.

---

#### FR-MC-004: Auto-Project Detection
**Description**: Infer project from current directory on start.

```bash
$ cd /home/user/projects/api-service
$ my-context start "Bug fix"
# Auto-detects project: api-service
# Creates: "api-service: Bug fix"

# Or explicit:
$ my-context start "Bug fix" --auto-project
```

**Rationale**: Reduce manual project name entry, leverage deb-sanity detection.

---

#### FR-MC-005: Rich Export with Environment
**Description**: Include environment snapshot in exports (integrate with deb-sanity).

```bash
$ my-context export "API Work" --include-environment

# Calls deb-sanity internally for environment data
# Export includes: OS, tools, services, worktree info
```

**Rationale**: Make exports self-documenting with full context.

---

## Recommended Integration Roadmap

### Phase 1: Foundation (Sprint 3)
- [ ] **FR-MC-002**: Add project path association to my-context
- [ ] **FR-DS-001**: Add `--active-context` to deb-sanity
- [ ] Create shared data format specification

### Phase 2: Automation (Sprint 4)
- [ ] **Integration #2**: Environment snapshot on context start
- [ ] **Integration #4**: Service state recording
- [ ] **FR-MC-004**: Auto-project detection

### Phase 3: Deep Integration (Sprint 5)
- [ ] **Integration #1**: Auto-context creation from project detection
- [ ] **Integration #5**: Worktree context sync
- [ ] **Integration #7**: Enhanced exports with environment

### Phase 4: Analytics (Sprint 6)
- [ ] **Integration #6**: Unified time tracking
- [ ] **FR-DS-003**: Context history export
- [ ] Dashboard/reporting features

---

## Data Sharing Specification

### Proposed: `~/.dev-tools-data/` Directory

**Structure**:
```
~/.dev-tools-data/
├── registry.json          # Shared tool/project registry
├── active-context.json    # Current my-context metadata
├── projects/
│   └── api-service.json   # Per-project metadata
└── sessions/
    └── 2025-10-10.json    # Daily session data
```

**Format** (registry.json):
```json
{
  "version": "1.0",
  "tools": {
    "my-context": {
      "path": "/home/user/.local/bin/my-context",
      "version": "e972750",
      "platform": "both",
      "source": "deb-sanity"
    }
  },
  "projects": {
    "api-service": {
      "path": "/home/user/projects/api-service",
      "remote": "github.com/company/api",
      "active_context": "api-service: Bug Fix",
      "discovered_by": "deb-sanity",
      "last_context_start": "2025-10-10T20:00:00Z"
    }
  }
}
```

**Benefits**:
- Both tools read/write shared data
- No tight coupling (tools still work independently)
- Future tools can join ecosystem

---

## Risk Assessment & Mitigation

| Risk | Severity | Mitigation Strategy |
|------|----------|---------------------|
| **Tight coupling** | **HIGH** | ⚠️ **CRITICAL**: Use shared data files ONLY. Never call deb-sanity commands directly from my-context |
| **Surface area growth** | **MEDIUM** | Use smart defaults, not flags. Auto-detect instead of --auto-project |
| **Performance overhead** | Low | Cache shared data, read-only access, async updates |
| **Data conflicts** | Low | File locking (flock), append-only logs, atomic writes |
| **Backward compatibility** | Medium | Schema versioning ("schema_version": 2), optional fields with `omitempty` |
| **Cross-platform gaps** | **MEDIUM** | Graceful degradation on Windows, test without deb-sanity |
| **Error handling** | Medium | Fail silently if shared data unavailable, log warnings only |

### Design Principles for Integration

1. **Loose Coupling via Shared Data** (not direct calls)
   ```bash
   # ❌ BAD: Direct dependency
   my-context calls deb-sanity --tools
   
   # ✅ GOOD: Shared data
   deb-sanity --scan → writes ~/.dev-tools-data/environment.json
   my-context reads environment.json (if exists)
   ```

2. **Smart Defaults over Flags** (Constitution: Minimal Surface Area)
   ```bash
   # ❌ BAD: Too many flags
   --capture-env, --auto-project, --include-environment
   
   # ✅ GOOD: Intelligent behavior
   cd /path/to/project && my-context start "Work"
   # Auto-detects project from pwd (no flag)
   # Auto-includes environment.json if exists (no flag)
   ```

3. **Graceful Degradation** (works without deb-sanity)
   ```bash
   # If ~/.dev-tools-data/ doesn't exist:
   $ my-context start "Work"
   Context started: "Work"  # No warnings, just works
   
   # If environment.json exists:
   $ my-context start "Work"
   Context started: "Work"
   Environment captured: WSL/Debian, 5 tools detected
   ```

---

## Cross-Platform Compatibility

### Windows Support Strategy

**Current State**:
- **my-context**: ✅ Full Windows support (cmd.exe, PowerShell, Git Bash)
- **deb-sanity**: ❓ Primarily WSL/Linux focused

**Integration Behavior on Windows**:

```powershell
# Windows PowerShell (without deb-sanity)
PS> my-context start "Work"
Context started: "Work"
# No environment capture (no error, silent skip)

# Windows with WSL + deb-sanity
PS> my-context start "Work"
# Could read \\wsl$\debian\home\user\.dev-tools-data\environment.json
# Requires path translation logic
```

**Implementation Requirements**:
1. **Platform detection** in my-context:
   ```go
   func tryReadEnvironmentData() (*EnvData, error) {
       if runtime.GOOS == "windows" {
           // Try WSL path: \\wsl$\<distro>\home\<user>\.dev-tools-data\
           // If not found: return nil (silent skip)
       }
       // Unix: ~/.dev-tools-data/environment.json
   }
   ```

2. **Graceful degradation**:
   - Windows native (no WSL): integration disabled, no warnings
   - Windows + WSL: try to read via WSL path, fail silently if unavailable
   - Linux/macOS: full integration available

---

## Error Handling Specification

### Failure Modes & Responses

| Scenario | my-context Behavior | User Experience |
|----------|---------------------|-----------------|
| deb-sanity not installed | ✅ Works normally | No warnings, no errors |
| ~/.dev-tools-data/ missing | ✅ Works normally | Silent skip of integration features |
| environment.json corrupted | ✅ Works normally | Logs warning to stderr, continues |
| Shared data stale (>24h old) | ⚠️ Use data anyway | Note shows "(stale data)" in export |
| File lock timeout (>1s) | ✅ Skip read | Fails silently, no blocking |

### Error Logging (Not User-Facing)

```go
// Log to ~/.my-context/integration.log (optional debug file)
if err != nil {
    logIntegrationError("Failed to read environment data: %v", err)
    // Continue execution - do NOT surface error to user
}
```

**Philosophy**: Integration features are **enhancements**, not requirements. Any failure should be invisible to users.

---

## Testing Strategy

### Unit Tests (No External Dependencies)

```go
// Test with mock shared data
func TestEnvironmentCapture_WithSharedData(t *testing.T) {
    // Create fake ~/.dev-tools-data/environment.json
    tmpDir := t.TempDir()
    writeJSON(tmpDir+"/environment.json", mockEnvData)
    
    ctx := CreateContext("Test")
    
    // Verify environment notes added
    assert.Contains(ctx.Notes, "Environment: WSL 2")
}

func TestEnvironmentCapture_WithoutSharedData(t *testing.T) {
    // No shared data file
    ctx := CreateContext("Test")
    
    // Verify context still created successfully
    assert.NotNil(ctx)
    assert.Equal("Test", ctx.Name)
}
```

### Integration Tests (Decoupled)

```bash
#!/bin/bash
# test-integration.sh (no deb-sanity dependency)

# Test 1: Standalone (no shared data)
my-context start "Test1"
my-context stop
assert_success "Standalone mode works"

# Test 2: With shared data
mkdir -p ~/.dev-tools-data
echo '{"os":"Linux","tools":["go"]}' > ~/.dev-tools-data/environment.json
my-context start "Test2"
my-context show | grep "Environment: Linux"
assert_success "Integration mode works"

# Test 3: With corrupt data
echo 'invalid json' > ~/.dev-tools-data/environment.json
my-context start "Test3"
assert_success "Graceful degradation works"
```

**Key Principle**: Tests should **never** require deb-sanity to be installed or running.

---

## User Stories & Success Metrics

### User Story 1: Reproducible Environment
**As a** developer debugging production issues  
**I want** my-context to capture my local environment automatically  
**So that** I can reproduce the exact toolchain later

**Acceptance Criteria**:
- Context exports show OS, tool versions, services
- No manual flag required (`--capture-env` removed)
- Works without deb-sanity (graceful degradation)

**Success Metric**: 80% of exported contexts include environment data (when deb-sanity available)

---

### User Story 2: Project-Aware Contexts
**As a** developer working on 5+ projects simultaneously  
**I want** my-context to auto-detect the current project  
**So that** I don't manually type project names repeatedly

**Acceptance Criteria**:
```bash
cd /home/user/projects/api-service
my-context start "Bug fix"
# Auto-creates: "api-service: Bug fix"

my-context list --here
# Shows only contexts for api-service
```

**Success Metric**: 50% reduction in manual project name entry

---

### User Story 3: Context Discovery
**As a** developer returning to an old project after months  
**I want** to find relevant contexts by project path  
**So that** I can review my previous work and decisions

**Acceptance Criteria**:
```bash
cd /old/project/path
my-context list --here
# Shows 12 historical contexts for this project

my-context export --all --here --to project-history.md
# Creates documentation of all work done here
```

**Success Metric**: 90% of developers use `--here` filter within first month

---

## Rollback & Escape Hatches

### Disabling Integration Features

**Environment Variable** (global disable):
```bash
export MY_CONTEXT_DISABLE_INTEGRATION=1
my-context start "Work"
# Skips all shared data reads, even if available
```

**Per-Context Disable**:
```bash
my-context start "Work" --no-integration
# One-time disable for this context
```

**Configuration File** (~/.my-context/config.json):
```json
{
  "integration": {
    "enabled": false,
    "features": {
      "environment_capture": false,
      "auto_project": true
    }
  }
}
```

### Migration Path (Rollback to v2.0)

If integration causes issues:
```bash
# Downgrade to Sprint 2 version
cd ~/.my-context/
mv config.json config.json.backup
# Reinstall v2.0 binary
# Data remains compatible (new fields ignored)
```

**Backward Compatibility Guarantee**:
- Sprint 2 binary can read Sprint 3+ data (extra fields ignored)
- Sprint 3+ binary can read Sprint 2 data (missing fields = defaults)
- No migration required

---

## Revised Feature Flag Summary

### Removed Flags (Use Smart Defaults)
- ~~`--capture-env`~~ → Automatic if environment.json exists
- ~~`--auto-project`~~ → Always enabled (infer from pwd)
- ~~`--include-environment`~~ → Automatic in exports if data present

### New Flags (Minimal Addition)
- `--path <dir>` (explicit project path)
- `--here` (filter by current directory)
- `--no-integration` (opt-out escape hatch)

**Net Change**: +3 flags (down from +9 in original proposal)

**Constitution Compliance**: ✅ Maintains minimal surface area

---

## Conclusion

**Verdict**: **Strong synergy with controlled risk**

### Key Strengths
1. ✅ **Complementary**, not competitive
2. ✅ **Already working together** organically
3. ✅ **Clear integration points** identified
4. ✅ **Both use similar patterns** (plain text, JSON, Unix philosophy)

### Recommended Actions (Revised Based on Review)

**Sprint 3 (Immediate - Foundation)**:
1. ✅ **FR-MC-002**: Project path association (`--path`, `--here`)
   - Implementation: Add `project_path` to meta.json
   - Auto-detect from `pwd` (no flag required)
   - Filter: `my-context list --here`
   - **Risk**: Low, no external dependencies

2. 📝 **Shared Data Spec**: Document `~/.dev-tools-data/` contract
   - JSON schema definition
   - File locking rules
   - Versioning strategy
   - **Don't implement yet** - just specification

**Sprint 4 (Quick Wins)**:
3. ⚡ **Environment Capture**: Read shared data (if available)
   - Read `~/.dev-tools-data/environment.json`
   - Add as initial notes
   - Graceful degradation (no deb-sanity required)
   - **Risk**: Low, loose coupling via shared files

4. 🎯 **Smart Defaults**: Auto-project detection
   - Infer project from `pwd` automatically
   - No new flags required
   - **Risk**: None, pure enhancement

**Sprint 5+ (Hold for Feedback)**:
- Worktree sync (requires deb-sanity changes)
- Unified time tracking (complex reconciliation)
- Auto-context creation (tight coupling risk)

### Expected Benefits (Quantified)
- **50% reduction** in manual project name entry (User Story 2)
- **80% adoption** of environment capture (when deb-sanity present)
- **90% usage** of `--here` filter within first month
- **Zero breakage** for users without deb-sanity (graceful degradation)

### Constitution Compliance (Final Check)
| Principle | Status | Notes |
|-----------|--------|-------|
| Unix Philosophy | ✅ Pass | Loose coupling via shared files |
| Cross-Platform | ✅ Pass | Windows graceful degradation added |
| Stateful Context | ✅ Pass | No core model changes |
| Minimal Surface | ✅ Pass | +3 flags (was +9, reduced) |
| Data Portability | ✅ Pass | JSON, plain text, greppable |
| User-Driven | ✅ Pass | Based on observed workflows |

---

**Next Steps**: 
1. Implement FR-MC-002 in Sprint 3
2. Test without deb-sanity (ensure graceful degradation)
3. Coordinate shared data spec with deb-sanity maintainer
4. Collect user feedback before Sprint 5 features

