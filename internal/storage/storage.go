package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type State struct {
	LastFetchedTime time.Time `json:"last_fetched_time"`
}

type Manager struct {
	stateFile string
	tempDir   string
}

func NewManager(stateFile, tempDir string) (*Manager, error) {
	stateDir := filepath.Dir(stateFile)
	if err := os.MkdirAll(stateDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create state directory: %w", err)
	}

	if tempDir != "" {
		if err := os.MkdirAll(tempDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create temp directory: %w", err)
		}
	}

	return &Manager{
		stateFile: stateFile,
		tempDir:   tempDir,
	}, nil
}

func (m *Manager) LoadState(defaultStart time.Time) (time.Time, error) {
	data, err := os.ReadFile(m.stateFile)
	if err != nil {
		if os.IsNotExist(err) {
			return defaultStart, nil
		}
		return time.Time{}, err
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return time.Time{}, err
	}

	return state.LastFetchedTime, nil
}

func (m *Manager) SaveState(t time.Time) error {
	state := State{LastFetchedTime: t}
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return os.WriteFile(m.stateFile, data, 0644)
}

func (m *Manager) CreateTempFile(uid uint32, from string, subject string, body []byte) (string, error) {
	filename := fmt.Sprintf("email_%d.txt", uid)
	filePath := filepath.Join(m.tempDir, filename)

	content := fmt.Sprintf("From: %s\nSubject: %s\n\n%s", from, subject, string(body))
	
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write temp email file: %w", err)
	}

	return filePath, nil
}

func (m *Manager) CleanupTempFile(path string) error {
	return os.Remove(path)
}
