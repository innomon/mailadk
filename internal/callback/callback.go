package callback

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Sender struct {
	callbackURL string
	client      *http.Client
}

func NewSender(callbackURL string) *Sender {
	return &Sender{
		callbackURL: callbackURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *Sender) PostToken(token string) error {
	payload := map[string]string{"token": token}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := s.client.Post(s.callbackURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to POST token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("SPA callback returned non-2xx status code: %d", resp.StatusCode)
	}

	return nil
}
