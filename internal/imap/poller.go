package imap

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
)

type Message struct {
	From    string
	AuthJWT string
}

type Poller struct {
	host    string
	user    string
	pass    string
	subject string
}

func NewPoller(host, user, pass, subject string) *Poller {
	return &Poller{
		host:    host,
		user:    user,
		pass:    pass,
		subject: subject,
	}
}

func (p *Poller) Poll(ctx context.Context, callback func(msg *Message) error) error {
	c, err := imapclient.DialTLS(p.host, &imapclient.Options{
		TLSConfig: &tls.Config{InsecureSkipVerify: false},
	})
	if err != nil {
		return fmt.Errorf("failed to dial IMAP: %w", err)
	}
	defer c.Logout()

	if err := c.Login(p.user, p.pass).Wait(); err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	// Select INBOX
	if _, err := c.Select("INBOX", nil).Wait(); err != nil {
		return fmt.Errorf("failed to select INBOX: %w", err)
	}

	// Search for messages with the given subject that are UNSEEN
	// and arrived since yesterday (as a safety window).
	searchData, err := c.Search(&imap.SearchCriteria{
		Header:  []imap.SearchCriteriaHeaderField{{Key: "Subject", Value: p.subject}},
		NotFlag: []imap.Flag{imap.FlagSeen},
		Since:   time.Now().Add(-24 * time.Hour),
	}, nil).Wait()
	if err != nil {
		return fmt.Errorf("failed to search: %w", err)
	}

	if searchData.All == nil {
		return nil
	}

	// Fetch messages
	fetchOptions := &imap.FetchOptions{
		Envelope: true,
		UID:      true,
		BodySection: []*imap.FetchItemBodySection{
			{Specifier: imap.PartSpecifierText}, // Fetch the text part
		},
	}

	messages, err := c.Fetch(searchData.All, fetchOptions).Collect()
	if err != nil {
		return fmt.Errorf("failed to fetch messages: %w", err)
	}

	for _, metadata := range messages {
		// Extract From address
		from := ""
		if metadata.Envelope != nil && len(metadata.Envelope.From) > 0 {
			f := metadata.Envelope.From[0]
			from = fmt.Sprintf("%s@%s", f.Mailbox, f.Host)
		}

		// Extract body and find auth JSON
		var authJWT string
		for _, body := range metadata.BodySection {
			bodyData := body.Bytes

			// Parse JSON: {"auth": "..."}
			var payload struct {
				Auth string `json:"auth"`
			}
			// The body might contain more than just JSON (e.g. if it's multipart)
			if err := json.Unmarshal(bodyData, &payload); err == nil {
				authJWT = payload.Auth
			} else {
				// Try a more robust approach if it's not pure JSON
				start := strings.Index(string(bodyData), "{")
				end := strings.LastIndex(string(bodyData), "}")
				if start != -1 && end != -1 && end > start {
					if err := json.Unmarshal(bodyData[start:end+1], &payload); err == nil {
						authJWT = payload.Auth
					}
				}
			}

			if authJWT != "" {
				break
			}
		}

		if from != "" && authJWT != "" {
			if err := callback(&Message{From: from, AuthJWT: authJWT}); err == nil {
				// Mark as seen/deleted if callback succeeded
				storeFlags := &imap.StoreFlags{
					Op:    imap.StoreFlagsAdd,
					Flags: []imap.Flag{imap.FlagSeen, imap.FlagDeleted},
				}
				if err := c.Store(imap.UIDSetNum(metadata.UID), storeFlags, nil).Close(); err != nil {
					log.Printf("failed to mark message %v as seen/deleted: %v", metadata.UID, err)
				}
			}
		}
	}

	// Expunge if we marked as Deleted
	if err := c.Expunge().Close(); err != nil {
		log.Printf("failed to expunge: %v", err)
	}

	return nil
}

func (p *Poller) FetchSince(ctx context.Context, since time.Time, callback func(uid uint32, from string, subject string, body []byte) error) (time.Time, error) {
	c, err := imapclient.DialTLS(p.host, &imapclient.Options{
		TLSConfig: &tls.Config{InsecureSkipVerify: false},
	})
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to dial IMAP: %w", err)
	}
	defer c.Logout()

	if err := c.Login(p.user, p.pass).Wait(); err != nil {
		return time.Time{}, fmt.Errorf("failed to login: %w", err)
	}

	if _, err := c.Select("INBOX", nil).Wait(); err != nil {
		return time.Time{}, fmt.Errorf("failed to select INBOX: %w", err)
	}

	searchData, err := c.Search(&imap.SearchCriteria{
		Since: since,
	}, nil).Wait()
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to search: %w", err)
	}

	if searchData.All == nil {
		return since, nil
	}

	fetchOptions := &imap.FetchOptions{
		Envelope: true,
		UID:      true,
		BodySection: []*imap.FetchItemBodySection{
			{Specifier: imap.PartSpecifierText},
		},
	}

	messages, err := c.Fetch(searchData.All, fetchOptions).Collect()
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to fetch: %w", err)
	}

	latest := since
	for _, m := range messages {
		from := ""
		if m.Envelope != nil && len(m.Envelope.From) > 0 {
			f := m.Envelope.From[0]
			from = fmt.Sprintf("%s@%s", f.Mailbox, f.Host)
		}

		body := []byte("")
		if len(m.BodySection) > 0 {
			body = m.BodySection[0].Bytes
		}

		if err := callback(uint32(m.UID), from, m.Envelope.Subject, body); err != nil {
			log.Printf("Callback error for UID %v: %v", m.UID, err)
			continue
		}

		// Update latest timestamp based on Envelope date if available, or just use now
		if m.Envelope != nil && m.Envelope.Date.After(latest) {
			latest = m.Envelope.Date
		}
	}

	return latest, nil
}

func (p *Poller) Run(ctx context.Context, interval time.Duration, callback func(msg *Message) error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		log.Printf("Polling IMAP...")
		if err := p.Poll(ctx, callback); err != nil {
			log.Printf("IMAP poll error: %v", err)
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}
