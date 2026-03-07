package main

import (
	"context"
	"log"
	"time"

	"github.com/innomon/mailadk/internal/config"
	"github.com/innomon/mailadk/internal/imap"
	"github.com/innomon/mailadk/internal/mcp"
	"github.com/innomon/mailadk/internal/storage"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. Initialize MCP Client
	mcpClient, err := mcp.NewClient(cfg.Puller.MCPCommand, cfg.Puller.MCPModel)
	if err != nil {
		log.Fatalf("Failed to initialize MCP client: %v", err)
	}
	ctx := context.Background()
	if err := mcpClient.Start(ctx); err != nil {
		log.Fatalf("Failed to start MCP client: %v", err)
	}
	defer mcpClient.Close()

	// 3. Initialize Storage Manager (for state and temp files)
	sm, err := storage.NewManager(cfg.Puller.StateFile, cfg.Puller.TempDir)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// 4. Parse Default Start Date
	defaultStart, err := time.Parse("2006-01-02", cfg.Puller.StartDate)
	if err != nil {
		log.Printf("Warning: failed to parse start_date %v, falling back to 24h ago: %v", cfg.Puller.StartDate, err)
		defaultStart = time.Now().Add(-24 * time.Hour)
	}

	// 5. Load last fetched time
	lastFetched, err := sm.LoadState(defaultStart)
	if err != nil {
		log.Fatalf("Failed to load state: %v", err)
	}

	// 6. Initialize IMAP Poller
	poller := imap.NewPoller(cfg.IMAP.Host, cfg.IMAP.User, cfg.IMAP.Pass, cfg.IMAP.Subject)

	// 7. Define callback for ingesting emails via MCP
	ingestCallback := func(uid uint32, from string, subject string, body []byte) error {
		log.Printf("Fetching email UID %v from %s", uid, from)
		
		// Save to temp file
		path, err := sm.CreateTempFile(uid, from, subject, body)
		if err != nil {
			return err
		}
		defer sm.CleanupTempFile(path)

		// Ingest via MCP
		log.Printf("Ingesting email UID %v into Dyna-SLM...", uid)
		if err := mcpClient.Ingest(ctx, path); err != nil {
			return err
		}

		return nil
	}

	// 8. Perform Fetch
	log.Printf("Starting fetch from %v...", lastFetched)
	latestTime, err := poller.FetchSince(ctx, lastFetched, ingestCallback)
	if err != nil {
		log.Fatalf("Fetch failed: %v", err)
	}

	// 9. Save State
	if err := sm.SaveState(latestTime); err != nil {
		log.Fatalf("Failed to save state: %v", err)
	}

	log.Printf("Fetch and Ingestion completed. Next sync will start from %v", latestTime)
}
