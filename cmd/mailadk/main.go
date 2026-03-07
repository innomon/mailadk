package main

import (
	"fmt"
	"log"

	"github.com/innomon/mailadk/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Println("Config Loaded Successfully")
	fmt.Printf("IMAP Host: %s\n", cfg.IMAP.Host)
	fmt.Printf("IMAP User: %s\n", cfg.IMAP.User)
	fmt.Printf("Poll Interval: %v\n", cfg.App.PollInterval)
}
