// Package mcp provides a Model Context Protocol server for my-context.
// This allows AI agents (Claude, etc.) to programmatically manage work contexts.
package mcp

import (
	"context"
	"fmt"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server wraps the MCP server with my-context functionality.
type Server struct {
	server  *mcp.Server
	version string
}

// NewServer creates a new MCP server for my-context.
func NewServer(version string) *Server {
	s := &Server{
		version: version,
	}

	// Create MCP server with implementation info
	s.server = mcp.NewServer(&mcp.Implementation{
		Name:    "my-context",
		Version: version,
	}, &mcp.ServerOptions{
		Instructions: "my-context is a CLI tool for managing developer work contexts. " +
			"Use these tools to start/stop contexts, add notes, and track your work sessions.",
	})

	// Register tools
	s.registerTools()

	return s
}

// registerTools adds all available tools to the MCP server.
func (s *Server) registerTools() {
	// start_context - Start a new work context
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "start_context",
		Description: "Start a new work context. Automatically stops any active context.",
	}, handleStartContext)

	// stop_context - Stop the active context
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "stop_context",
		Description: "Stop the currently active context.",
	}, handleStopContext)

	// add_note - Add a timestamped note to the active context
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "add_note",
		Description: "Add a timestamped note to the currently active context.",
	}, handleAddNote)

	// get_active_context - Get information about the active context
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "get_active_context",
		Description: "Get information about the currently active context including notes and files.",
	}, handleGetActiveContext)

	// list_contexts - List all contexts
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "list_contexts",
		Description: "List all contexts with optional filtering.",
	}, handleListContexts)

	// add_file - Associate a file with the active context
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "add_file",
		Description: "Associate a file with the currently active context for tracking.",
	}, handleAddFile)

	// list_files - List files associated with a context
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "list_files",
		Description: "List files associated with a context (defaults to active context).",
	}, handleListFiles)

	// export_context - Export a context to markdown or JSON
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "export_context",
		Description: "Export a context to markdown or JSON format for documentation or backup.",
	}, handleExportContext)

	// archive_context - Archive a stopped context
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "archive_context",
		Description: "Archive a stopped context. Cannot archive active contexts.",
	}, handleArchiveContext)

	// search_contexts - Search contexts by name and optionally notes
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "search_contexts",
		Description: "Search contexts by name and optionally search within note content.",
	}, handleSearchContexts)
}

// Run starts the MCP server with stdio transport.
// This blocks until the client disconnects or context is cancelled.
func (s *Server) Run(ctx context.Context) error {
	fmt.Fprintln(os.Stderr, "my-context MCP server started")
	return s.server.Run(ctx, &mcp.StdioTransport{})
}
