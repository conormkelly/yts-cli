// internal/config/keyring.go
package config

import (
	"fmt"

	"github.com/zalando/go-keyring"
)

const (
	KeyringService = "yts-cli"
)

// APIKeyManager handles secure storage and retrieval of API keys
type APIKeyManager struct {
	service string
}

// NewAPIKeyManager creates a new instance of APIKeyManager
func NewAPIKeyManager() *APIKeyManager {
	return &APIKeyManager{
		service: KeyringService,
	}
}

// SetAPIKey stores an API key in the system keyring
func (m *APIKeyManager) SetAPIKey(provider, apiKey string) error {
	if err := keyring.Set(m.service, provider, apiKey); err != nil {
		return fmt.Errorf("failed to store API key for %s: %w", provider, err)
	}
	return nil
}

func (m *APIKeyManager) HasAPIKey(provider string) bool {
	_, err := keyring.Get(m.service, provider)
	return err == nil
}

// GetAPIKey retrieves an API key from the system keyring
func (m *APIKeyManager) GetAPIKey(provider string) (string, error) {
	apiKey, err := keyring.Get(m.service, provider)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve API key for %s: %w", provider, err)
	}
	return apiKey, nil
}

// DeleteAPIKey removes an API key from the system keyring
func (m *APIKeyManager) DeleteAPIKey(provider string) error {
	if err := keyring.Delete(m.service, provider); err != nil {
		return fmt.Errorf("failed to delete API key for %s: %w", provider, err)
	}
	return nil
}
