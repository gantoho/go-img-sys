package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gantoho/go-img-sys/pkg/utils"
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

func DefaultStorePath() string {
	return "configs/api_keys.json"
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

func (km *KeyManager) CreateRandomKey(expireAfterDays int) string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	plainKey := hex.EncodeToString(b)
	return km.CreateKey(plainKey, expireAfterDays)
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

func (km *KeyManager) SaveToFile(path string) error {
	km.mu.RLock()
	defer km.mu.RUnlock()
	dir := filepath.Dir(path)
	if err := utils.EnsureDir(dir); err != nil {
		return err
	}
	data, err := json.MarshalIndent(km.keys, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (km *KeyManager) LoadFromFile(path string) error {
	km.mu.Lock()
	defer km.mu.Unlock()
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var m map[string]*APIKey
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	km.keys = m
	return nil
}

// GetKeyInfo returns information about a key
func (km *KeyManager) GetKeyInfo(plainKey string) *APIKey {
	km.mu.RLock()
	defer km.mu.RUnlock()

	hashedKey := GenerateKey(plainKey)
	return km.keys[hashedKey]
}

func (km *KeyManager) GetKeyInfoByHash(hashedKey string) *APIKey {
	km.mu.RLock()
	defer km.mu.RUnlock()
	return km.keys[hashedKey]
}

func (km *KeyManager) ValidateKeyHash(hashedKey string) bool {
	km.mu.RLock()
	defer km.mu.RUnlock()
	apiKey, exists := km.keys[hashedKey]
	if !exists {
		return false
	}
	if !apiKey.Active {
		return false
	}
	if time.Now().After(apiKey.ExpiresAt) {
		return false
	}
	return true
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
			"key_hash":   apiKey.Key[:16] + "...",
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
