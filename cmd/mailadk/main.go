package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/innomon/mailadk/internal/callback"
	"github.com/innomon/mailadk/internal/config"
	"github.com/innomon/mailadk/internal/imap"
	"github.com/innomon/mailadk/internal/token"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Println("Configuration loaded.")

	// 2. Initialize Token Verifier/Signer
	verifier, err := token.NewVerifier(cfg.Auth.SPAPublicKey, cfg.Auth.VerifierPrivateKey)
	if err != nil {
		log.Fatalf("Failed to initialize token verifier: %v", err)
	}

	// 3. Initialize Callback Sender
	sender := callback.NewSender(cfg.App.SPACallbackURL)

	// 4. Initialize IMAP Poller
	poller := imap.NewPoller(cfg.IMAP.Host, cfg.IMAP.User, cfg.IMAP.Pass, cfg.IMAP.Subject)

	// 5. Setup Context for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 6. Define the message processing callback
	processMessage := func(msg *imap.Message) error {
		log.Printf("Processing message from: %s", msg.From)

		// a. Decode Base32
		tokenStr, err := verifier.DecodeBase32(msg.AuthJWT)
		if err != nil {
			log.Printf("Failed to decode base32 JWT: %v", err)
			return err
		}

		// b. Verify SPA Token
		spaClaims, err := verifier.VerifySPAToken(tokenStr)
		if err != nil {
			log.Printf("Failed to verify SPA token: %v", err)
			return err
		}

		// c. Sign Final Token
		finalToken, err := verifier.SignFinalToken(spaClaims.LipID, msg.From)
		if err != nil {
			log.Printf("Failed to sign final token: %v", err)
			return err
		}

		// d. Post Final Token to SPA callback
		if err := sender.PostToken(finalToken); err != nil {
			log.Printf("Failed to post final token to callback: %v", err)
			return err
		}

		log.Printf("Successfully processed and posted token for %s", msg.From)
		return nil
	}

	// 7. Start Polling
	log.Printf("Starting mailadk poller with interval: %v", cfg.App.PollInterval)
	poller.Run(ctx, cfg.App.PollInterval, processMessage)

	log.Println("mailadk shutting down...")
}
