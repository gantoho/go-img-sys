package auth

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

// APIKey represents an API key with expiration
type APIKey struct {
	Key       string
	ExpiresAt time.Time
	Active    bool
	CreatedAt time.Time
}

// KeyManager manages API keys
type KeyManager struct {
	mu   sync.RWMutex
	keys map[string]*APIKey
}

var keyManager *KeyManager

func init() {
	keyManager = &KeyManager{
		keys: make(map[string]*APIKey),
	}
}

// GetManager returns the singleton key manager
func GetManager() *KeyManager {
	return keyManager
}

// GenerateKey generates a new API key hash
func GenerateKey(plainKey string) string {
	hash := sha256.Sum256([]byte(plainKey))
	return fmt.Sprintf("%x", hash)
}

// CreateKey creates a new API key with expiration
func (km *KeyManager) CreateKey(plainKey string, expireAfterDays int) string {
	km.mu.Lock()
	defer km.mu.Unlock()

	hashedKey := GenerateKey(plainKey)
	expiresAt := time.Now().AddDate(0, 0, expireAfterDays)

	km.keys[hashedKey] = &APIKey{
		Key:       hashedKey,
		ExpiresAt: expiresAt,
		Active:    true,
		CreatedAt: time.Now(),
	}

	return plainKey // Return plain key only once for user to save
}

// ValidateKey validates an API key
func (km *KeyManager) ValidateKey(plainKey string) bool {
	km.mu.RLock()
	defer km.mu.RUnlock()

	hashedKey := GenerateKey(plainKey)
	apiKey, exists := km.keys[hashedKey]

	if !exists {
		return false
	}

	// Check if key is active
	if !apiKey.Active {
		return false
	}

	// Check if key has expired
	if time.Now().After(apiKey.ExpiresAt) {
		return false
	}

	return true
}

// RevokeKey revokes an API key
func (km *KeyManager) RevokeKey(plainKey string) bool {
	km.mu.Lock()
	defer km.mu.Unlock()

	hashedKey := GenerateKey(plainKey)
	apiKey, exists := km.keys[hashedKey]

	if !exists {
		return false
	}

	apiKey.Active = false
	return true
}

// GetKeyInfo returns information about a key
func (km *KeyManager) GetKeyInfo(plainKey string) *APIKey {
	km.mu.RLock()
	defer km.mu.RUnlock()

	hashedKey := GenerateKey(plainKey)
	return km.keys[hashedKey]
}

// ListKeys returns all keys info (without the actual key hash for security)
func (km *KeyManager) ListKeys() []map[string]interface{} {
	km.mu.RLock()
	defer km.mu.RUnlock()

	result := make([]map[string]interface{}, 0, len(km.keys))
	for _, apiKey := range km.keys {
		result = append(result, map[string]interface{}{
			"created_at": apiKey.CreatedAt.Unix(),
			"expires_at": apiKey.ExpiresAt.Unix(),
			"active":     apiKey.Active,
			"is_expired": time.Now().After(apiKey.ExpiresAt),
			"key_hash":   apiKey.Key[:16] + "...", // Show partial hash only
		})
	}
	return result
}

// CleanupExpiredKeys removes expired keys periodically
func (km *KeyManager) CleanupExpiredKeys() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		km.mu.Lock()
		now := time.Now()
		for key, apiKey := range km.keys {
			if now.After(apiKey.ExpiresAt) && !apiKey.Active {
				delete(km.keys, key)
			}
		}
		km.mu.Unlock()
	}
}

// InitDefaultKeys initializes some default keys for development
func (km *KeyManager) InitDefaultKeys() {
	km.CreateKey("demo-key-12345", 30) // Valid for 30 days
	km.CreateKey("test-key-67890", 7)  // Valid for 7 days
}
