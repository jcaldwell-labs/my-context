# My-Context Roadmap

> Feature roadmap for my-context development.

## Current Version

**v3.3.1** — Database Backend with Partitions

## Completed Sprints

### Sprint 4 (v3.2.x - v3.3.x)
- PostgreSQL database backend (`MY_CONTEXT_HOME=db`)
- Partition support for project isolation (`db:project-name`)
- Cross-partition queries (`--all-partitions`)
- `partitions` command for partition management
- `which` command shows backend info
- Shell completions (bash/zsh/fish/powershell)
- MCP server support (`mcp-server` command)
- Stats command with date filters (`--today`, `--week`, `--month`, `--since`, `--until`)
- Record command for clipboard monitoring

### Sprint 3 (v3.1.0)
- Context tags and filtering (`--tag`, `--tags`)
- Parent-child context hierarchy (`--parent`)
- Enhanced metadata (created-by, labels)
- Signal coordination (`signal`, `wait`)
- Context watching (`watch`)

### Sprint 2 (v2.x)
- Project grouping (`--project` flag)
- Export to markdown/JSON
- Archive and delete commands
- List filtering (`--search`, `--limit`, `--archived`)
- Multi-platform builds

### Sprint 1 (v1.x)
- Core 8 commands (start, stop, note, file, touch, show, list, history)
- Plain-text storage
- Cross-platform support
- JSON output mode

## Planned Features

### Near-term (v3.4.x)

**Workflow Improvements**
- Resume by index from list output (#41)
- Focus command: stop + resume in one (#42)
- Date filters for history command (#37)
- Stale context detection and warnings (#96)
- Cascade-stop for parent/child hierarchies (#97)

**Developer Experience**
- Pre-release validation script (#43)
- Dependabot for security updates (#46)

### Medium-term (v3.5.x+)

**Code Quality**
- Unit test coverage improvements (#89)
- Consolidate duplicated filter logic (#30)
- golangci-lint v2 migration (#67)

**Features**
- Import markdown notes (#29)
- Timetracker integration docs (#95)

### Future Exploration

**Integration Ecosystem**
- VS Code extension
- JetBrains plugin
- Git hooks library
- Slack notifications

**Team Features**
- Shared context visibility
- Context templates
- Team activity dashboards

**Advanced Workflows**
- AI-powered context summaries
- Clawdbot/agent integration

## Contributing

Feature requests welcome! Open an issue at:
https://github.com/jcaldwell-labs/my-context/issues
