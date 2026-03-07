package mcp

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Client struct {
	command []string
	model   string
	client  *mcp.Client
	session *mcp.ClientSession
}

func NewClient(commandStr, model string) (*Client, error) {
	cmdParts := strings.Fields(commandStr)
	if len(cmdParts) == 0 {
		return nil, fmt.Errorf("empty MCP command")
	}

	return &Client{
		command: cmdParts,
		model:   model,
	}, nil
}

func (c *Client) Start(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, c.command[0], c.command[1:]...)
	
	transport := &mcp.CommandTransport{
		Command: cmd,
	}
	
	c.client = mcp.NewClient(
		&mcp.Implementation{
			Name:    "mailadk-puller",
			Version: "1.0.0",
		},
		nil,
	)

	session, err := c.client.Connect(ctx, transport, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to MCP server: %w", err)
	}

	c.session = session
	return nil
}

func (c *Client) Ingest(ctx context.Context, path string) error {
	args := map[string]interface{}{
		"model": c.model,
		"path":  path,
	}

	_, err := c.session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "ingest_asset",
		Arguments: args,
	})
	
	if err != nil {
		return fmt.Errorf("failed to call ingest_asset tool: %w", err)
	}

	return nil
}

func (c *Client) Close() error {
	if c.session != nil {
		return c.session.Close()
	}
	return nil
}
