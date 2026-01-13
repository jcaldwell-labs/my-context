# MCP Server Implementation Research

**Date**: 2026-01-13
**Issue**: #10 - Implement MCP (Model Context Protocol) server

## Go SDK Options

### 1. Official SDK (Recommended)

**Repository**: [github.com/modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk)

| Attribute  | Details                                        |
| ---------- | ---------------------------------------------- |
| Maintainer | Anthropic + Google collaboration               |
| Version    | v1.2.0 (Dec 2025)                              |
| MCP Spec   | 2025-11-25, 2025-06-18, 2025-03-26, 2024-11-05 |
| Transports | Stdio, Command, Custom                         |
| Maturity   | Production-ready, uses gopls JSON-RPC          |

**Basic Server Example:**

```go
server := mcp.NewServer(&mcp.Implementation{Name: "my-context", Version: "v3.0.0"}, nil)
mcp.AddTool(server, &mcp.Tool{Name: "add_note", Description: "Add note to context"}, AddNoteHandler)
if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
    log.Fatal(err)
}
```

**Pros:**

- Official support from Anthropic
- Long-term maintenance guaranteed
- Battle-tested JSON-RPC from gopls
- Direct integration path with Claude

**Cons:**

- More verbose API
- Less documentation/examples than community SDK

### 2. Community SDK (mcp-go)

**Repository**: [github.com/mark3labs/mcp-go](https://github.com/mark3labs/mcp-go)

| Attribute  | Details                    |
| ---------- | -------------------------- |
| Maintainer | mark3labs (Ed Zynda)       |
| Adoption   | 1,307+ dependent packages  |
| Transports | Stdio, SSE, HTTP           |
| API Style  | Builder pattern, ergonomic |

**Basic Server Example:**

```go
s := server.NewMCPServer("my-context", "3.0.0")
tool := mcp.NewTool("add_note", mcp.WithDescription("Add note"), mcp.WithString("text", mcp.Required()))
s.AddTool(tool, addNoteHandler)
server.ServeStdio(s)
```

**Pros:**

- More ergonomic API
- Extensive documentation
- SSE transport for web integrations

**Cons:**

- Community maintained (risk of abandonment)
- May diverge from official spec

## Recommendation

**Use the Official SDK** (`github.com/modelcontextprotocol/go-sdk`) because:

1. **Official support** - Maintained by Anthropic, the creators of MCP
2. **Google collaboration** - Additional enterprise backing
3. **Long-term viability** - Unlikely to be abandoned
4. **Claude integration** - Direct path for Claude Desktop support

## Implementation Plan

### Phase 1: Core Server

1. Add `cmd/mcp-server/main.go` entry point
2. Implement stdio transport for Claude Desktop
3. Define basic tools: `start_context`, `stop_context`, `add_note`

### Phase 2: Resources

1. Implement `context://active` resource
2. Implement `context://list` resource
3. Implement `context://{name}` dynamic resource

### Phase 3: Full Tools

1. `add_file` - Associate file with context
2. `export_context` - Export to markdown/JSON
3. `search_contexts` - Search by keyword

### Phase 4: Prompts

1. `summarize_session` - Generate session summary
2. `list_decisions` - Extract decision notes
3. `handoff_context` - Create handoff document

## Claude Desktop Configuration

```json
{
  "mcpServers": {
    "my-context": {
      "command": "my-context",
      "args": ["mcp-server"],
      "env": {
        "MY_CONTEXT_HOME": "db"
      }
    }
  }
}
```

## Dependencies to Add

```bash
go get github.com/modelcontextprotocol/go-sdk@v1.2.0
```

## Estimated Effort

| Phase                | Effort   | Priority |
| -------------------- | -------- | -------- |
| Phase 1: Core Server | 2-3 days | High     |
| Phase 2: Resources   | 1-2 days | High     |
| Phase 3: Full Tools  | 2-3 days | Medium   |
| Phase 4: Prompts     | 1-2 days | Low      |

**Total**: ~7-10 days for full implementation

## References

- [MCP Specification](https://modelcontextprotocol.io/specification/2025-11-25/)
- [Official Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- [mcp-go Community SDK](https://github.com/mark3labs/mcp-go)
- [Building MCP Server in Go](https://navendu.me/posts/mcp-server-go/)
