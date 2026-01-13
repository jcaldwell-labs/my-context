package commands

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jefferycaldwell/my-context-copilot/internal/mcp"
	"github.com/spf13/cobra"
)

// NewMCPServerCmd creates the mcp-server command.
func NewMCPServerCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp-server",
		Short: "Run as MCP (Model Context Protocol) server",
		Long: `Start my-context as an MCP server for AI agent integration.

This allows AI assistants like Claude to programmatically manage work contexts
through the Model Context Protocol (MCP).

The server communicates via stdio (stdin/stdout) using JSON-RPC.

Example Claude Desktop configuration (claude_desktop_config.json):
  {
    "mcpServers": {
      "my-context": {
        "command": "my-context",
        "args": ["mcp-server"]
      }
    }
  }

Available tools:
  - start_context: Start a new work context
  - stop_context: Stop the active context
  - add_note: Add a timestamped note
  - get_active_context: Get active context info
  - list_contexts: List all contexts`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMCPServer(version)
		},
	}

	return cmd
}

func runMCPServer(version string) error {
	// Create MCP server
	server := mcp.NewServer(version)

	// Create context with signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		cancel()
	}()

	// Run the server (blocks until done)
	return server.Run(ctx)
}
