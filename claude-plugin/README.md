# my-context Claude Code Plugin

Track developer work contexts with notes, files, and timestamps directly from Claude Code.

## Installation

```bash
# From the plugin directory
claude plugins add ./claude-plugin
```

Or install from GitHub:

```bash
claude plugins add github:jcaldwell-labs/my-context/claude-plugin
```

## Prerequisites

The `my-context` CLI must be installed and available in your PATH.

```bash
# Install my-context
go install github.com/jefferycaldwell/my-context-copilot/cmd/my-context@latest
```

## Commands

| Command              | Description                 |
| -------------------- | --------------------------- |
| `/my-context:start`  | Start a new work context    |
| `/my-context:stop`   | Stop the current context    |
| `/my-context:note`   | Add a timestamped note      |
| `/my-context:show`   | Show active context details |
| `/my-context:list`   | List all contexts           |
| `/my-context:export` | Export context to markdown  |

## Usage Examples

```
User: /my-context:start
Claude: Starting context "feature-implementation"...

User: /my-context:note "Completed Phase 1 of MCP server"
Claude: Added note to context...

User: /my-context:show
Claude: [Shows context with notes, files, duration]

User: /my-context:stop
Claude: Stopped context after 2h 15m with 12 notes captured
```

## Why Plugin vs MCP?

| Approach       | Best For                                      |
| -------------- | --------------------------------------------- |
| **Plugin**     | Claude Code - simpler, direct bash calls      |
| **MCP Server** | Claude Desktop, other AI clients without bash |

The plugin provides the same functionality as the MCP server but optimized for Claude Code's bash-native environment.

## Configuration

Set `MY_CONTEXT_HOME` to customize storage location:

- File-based: `export MY_CONTEXT_HOME=~/.my-context`
- Database: `export MY_CONTEXT_HOME=db`
- Partitioned: `export MY_CONTEXT_HOME=db:project-name`
