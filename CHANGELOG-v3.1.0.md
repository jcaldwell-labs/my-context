# Changelog - v3.1.0 (Sprint 3)

**Release Date:** 2025-11-18
**Codename:** "Organization & Hierarchy"

## 🎉 Major Features

### Context Tags/Labels System
- **NEW:** `my-context tag add <context> <tags...>` - Add tags to organize contexts
- **NEW:** `my-context tag remove <context> <tags...>` - Remove tags from contexts
- **NEW:** `my-context tag list` - View all tags with usage counts
- **NEW:** `my-context tag list <context>` - View tags for specific context
- **ENHANCED:** `my-context list --tag <tag>` - Filter contexts by tag
- **ENHANCED:** `my-context start --labels <tags>` - Add tags when creating context

**Use Cases:**
- Priority tracking (urgent, low-priority, blocked)
- Team organization (frontend, backend, qa)
- Feature categorization (auth, payments, notifications)
- Status tracking (in-review, needs-testing)

### Parent-Child Context Relationships
- **NEW:** `my-context link <child> <parent>` - Create hierarchical relationships
- **NEW:** `my-context unlink <context>` - Remove parent relationship
- **NEW:** `my-context up [context]` - Show parent context
- **NEW:** `my-context down [context]` - List child contexts
- **ENHANCED:** `my-context start --parent <parent>` - Set parent when creating

**Use Cases:**
- Sprint → Task → Subtask organization
- Project → Phase → Task structure
- Epic → Story → Task tracking
- Campaign → Initiative → Action planning

### Context Tree Visualization
- **NEW:** `my-context tree [context]` - Display hierarchical tree view
- **NEW:** ASCII tree rendering with proper branch characters
- **NEW:** Automatic circular dependency detection
- **NEW:** Root contexts view (contexts without parents)
- **NEW:** JSON tree output for programmatic use

**Tree Output Example:**
```
└─ Sprint 3
   ├─ API refactor
   │  ├─ Unit tests
   │  └─ Integration tests
   └─ Database migration
```

## ✨ Enhancements

### Enhanced List Command
- Added `--tag <tag>` flag for filtering by tags
- Tag filtering works with all existing filters (--project, --search, --limit, etc.)
- Case-insensitive tag matching

### Enhanced Start Command
- Added `--labels <comma-separated>` flag for initial tagging
- Added `--parent <context>` flag for immediate hierarchy creation
- Added `--created-by <user>` flag for attribution
- All metadata fields are optional

### Metadata System
- Contexts now support optional metadata structure:
  - `created_by`: User attribution
  - `parent`: Parent context name
  - `labels`: Array of tags/labels
- Metadata stored in `meta.json` alongside existing context data
- **100% Backward Compatible** - old contexts work seamlessly

## 📦 New Commands

| Command | Alias | Description |
|---------|-------|-------------|
| `tag` | `t` | Manage context tags |
| `link` | - | Link child to parent |
| `unlink` | - | Remove parent link |
| `tree` | - | Show context hierarchy |
| `up` | `parent` | Navigate to parent |
| `down` | `children` | List children |

## 🔧 Technical Changes

### Core Functions Added
- `AddTags()`, `RemoveTags()`, `GetContextTags()`, `GetAllTags()`
- `SetParent()`, `ClearParent()`, `GetChildren()`
- `GetContextTree()`, `GetRootContexts()`
- Tree building with cycle detection

### Data Model Updates
- New `pkg/models/context.go` with `ContextMetadata` struct
- `ContextWithMetadata` extends base context
- Validation for metadata fields (max 10 labels, 50 chars each)
- No whitespace allowed in tags

### Architecture Improvements
- Metadata gracefully degrades for old contexts
- File-based storage maintained (no database required)
- All operations remain atomic and file-system safe

## 🐛 Bug Fixes

- None (new features only in this release)

## 📚 Documentation

- **NEW:** `docs/SPRINT-3-FEATURES.md` - Comprehensive Sprint 3 feature documentation
- Complete API reference for all new commands
- Use case examples and tutorials
- Migration guide (spoiler: no migration needed!)
- JSON output examples
- Performance benchmarks

## 🧪 Testing

### Manual Testing Performed
- Tag operations (add, remove, list, filter)
- Parent-child linking with various hierarchies
- Tree visualization with multiple levels
- Circular dependency detection
- Backward compatibility with Sprint 2 contexts

### Test Scenarios Covered
- Single-level tagging
- Multi-level hierarchies (3+ levels deep)
- Root contexts without parents
- Contexts with both tags and parents
- Filtering combinations (--tag + --project + --search)

## 🚀 Performance

Tested with 1000 contexts:
- **Add tag:** <10ms
- **List tags:** <50ms
- **Filter by tag:** <100ms
- **Build tree:** <150ms
- **Link contexts:** <15ms

All operations remain fast and responsive.

## ⬆️ Upgrade Path

**From v2.x to v3.1.0:**

No migration required! Simply:
1. Replace binary with v3.1.0
2. Start using new features immediately
3. Existing contexts work without modification

**Metadata is created automatically when you:**
- Add tags to an existing context
- Link contexts
- Create new contexts with `--labels` or `--parent`

## 🔄 Backward Compatibility

✅ **100% Backward Compatible**
- All existing commands work unchanged
- Old contexts display and function normally
- Old contexts can be enhanced with new features
- File format extensions are additive only

## 🌟 Highlights

**Most Requested Features:**
1. ✅ Context tagging and categorization
2. ✅ Hierarchical project organization
3. ✅ Visual tree representation
4. ✅ Enhanced filtering capabilities

**Best New Workflows:**
- Tag-based context discovery
- Sprint/task hierarchy management
- Team coordination with labels
- Project phase organization

## 🎯 Use Case Examples

### Multi-Sprint Project
```bash
my-context start "Q4 2025" --labels quarterly,planning
my-context start "Sprint 3" --parent "Q4 2025" --labels sprint,active
my-context start "API refactor" --parent "Sprint 3" --labels backend
my-context tree "Q4 2025"
```

### Bug Tracking
```bash
my-context start "Login bug" --labels bug,urgent,frontend
my-context list --tag bug
my-context list --tag urgent --tag frontend
```

### Team Organization
```bash
my-context tag add "Feature X" alice frontend
my-context list --tag alice
my-context down "Sprint 3"
```

## 📋 Known Limitations

- Tree visualization limited to terminal width
- No multi-parent support (single parent only)
- Tag autocomplete not yet implemented
- No tag rename functionality (remove + add)

## 🔮 Coming Next (Sprint 3.2)

**Planned features:**
- Context templates (bug-fix, feature, meeting, research)
- Date range filtering (--from, --to)
- Full-text note search
- File association filtering
- Database backend option (SQLite, PostgreSQL)
- Performance optimizations

## 🙏 Acknowledgments

- Sprint 3 planning based on user feedback and observed workflows
- Inspired by Git branches, issue trackers, and project management tools
- Designed with Unix philosophy: composable, text-based, single-purpose

---

## Installation

### Binary Installation
```bash
# Linux/macOS
wget https://github.com/jefferycaldwell/my-context/releases/download/v3.1.0/my-context-linux-amd64
chmod +x my-context-linux-amd64
sudo mv my-context-linux-amd64 /usr/local/bin/my-context

# Windows
# Download my-context-windows-amd64.exe from releases page
```

### Build from Source
```bash
git clone https://github.com/jefferycaldwell/my-context.git
cd my-context
git checkout v3.1.0
go build -o my-context ./cmd/my-context/
```

---

**Full Documentation:** [docs/SPRINT-3-FEATURES.md](docs/SPRINT-3-FEATURES.md)
**Previous Release:** [v2.3.0](CHANGELOG-v2.3.0.md)
**Repository:** https://github.com/jefferycaldwell/my-context
**License:** MIT
