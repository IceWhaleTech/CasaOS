package sign

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"sync"
	"time"

	"github.com/IceWhaleTech/CasaOS/pkg/sign"
)

var once sync.Once
var instance sign.Sign

func Sign(data string) string {

	return NotExpired(data)

}

func WithDuration(data string, d time.Duration) string {
	once.Do(Instance)
	return instance.Sign(data, time.Now().Add(d).Unix())
}

func NotExpired(data string) string {
	once.Do(Instance)
	return instance.Sign(data, 0)
}

func Verify(data string, sign string) error {
	once.Do(Instance)
	return instance.Verify(data, sign)
}

func Instance() {
	// Try to load key from environment variable first
	keyEnv := os.Getenv("CASAOS_HMAC_SECRET")
	if len(keyEnv) > 0 {
		instance = sign.NewHMACSign([]byte(keyEnv))
		return
	}

	// Try to load key from file
	keyFile := "/var/lib/casaos/hmac_secret.key"
	if data, err := os.ReadFile(keyFile); err == nil && len(data) > 0 {
		instance = sign.NewHMACSign(data)
		return
	}

	// Generate a new random key and persist it
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		// Fallback: use a random key per startup (less ideal but better than hardcoded)
		instance = sign.NewHMACSign(key)
		return
	}
	os.WriteFile(keyFile, key, 0o600)
	instance = sign.NewHMACSign(key)
}

// GenerateKey generates a random hex-encoded key for use in configuration
func GenerateKey() string {
	key := make([]byte, 32)
	rand.Read(key)
	return hex.EncodeToString(key)
}
