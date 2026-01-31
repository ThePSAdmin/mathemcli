package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const (
	configDir  = ".mathemcli"
	configFile = "session.json"
)

// Session represents stored session data
type Session struct {
	SessionID string `json:"session_id"`
	CSRFToken string `json:"csrf_token"`
	Email     string `json:"email"`
}

// ConfigPath returns the path to the config directory
func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configDir), nil
}

// SessionPath returns the path to the session file
func SessionPath() (string, error) {
	configPath, err := ConfigPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(configPath, configFile), nil
}

// LoadSession loads the saved session from disk
func LoadSession() (*Session, error) {
	sessionPath, err := SessionPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(sessionPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil // No session saved
		}
		return nil, err
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

// SaveSession saves the session to disk
func SaveSession(session *Session) error {
	configPath, err := ConfigPath()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configPath, 0700); err != nil {
		return err
	}

	sessionPath, err := SessionPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(sessionPath, data, 0600)
}

// ClearSession removes the saved session
func ClearSession() error {
	sessionPath, err := SessionPath()
	if err != nil {
		return err
	}

	err = os.Remove(sessionPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
