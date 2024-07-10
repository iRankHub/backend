package utils

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/o1egl/paseto"
)

var (
	publicKey ed25519.PublicKey
	blacklist = struct {
		sync.RWMutex
		tokens map[string]time.Time
	}{tokens: make(map[string]time.Time)}
)

func GenerateToken(userID int32, userRole string, privateKey ed25519.PrivateKey) (string, error) {
	maker := paseto.NewV2()

	claims := map[string]interface{}{
		"user_id":   userID,
		"user_role": userRole,
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	token, err := maker.Sign(privateKey, claims, nil)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return token, nil
}

func GeneratePasetoKeyPair() (ed25519.PrivateKey, ed25519.PublicKey, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate Ed25519 key pair: %v", err)
	}

	return privateKey, publicKey, nil
}

func SetPublicKey(key ed25519.PublicKey) {
	publicKey = key
}

func InvalidateToken(token string) {
	blacklist.Lock()
	defer blacklist.Unlock()
	blacklist.tokens[token] = time.Now().Add(24 * time.Hour) // Invalidate for 24 hours
}

func IsTokenInvalid(token string) bool {
	blacklist.RLock()
	defer blacklist.RUnlock()
	expiry, exists := blacklist.tokens[token]
	if !exists {
		return false
	}
	if time.Now().After(expiry) {
		delete(blacklist.tokens, token)
		return false
	}
	return true
}

func ValidateToken(token string) (map[string]interface{}, error) {
	if publicKey == nil {
		return nil, fmt.Errorf("public key not set")
	}

	if IsTokenInvalid(token) {
		return nil, fmt.Errorf("token has been invalidated")
	}

	maker := paseto.NewV2()
	var claims map[string]interface{}
	err := maker.Verify(token, publicKey, &claims, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %v", err)
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid expiration claim")
	}

	if time.Now().Unix() > int64(exp) {
		return nil, fmt.Errorf("token has expired")
	}

	return claims, nil
}

// New function to clean up expired tokens from the blacklist
func CleanupExpiredTokens() {
	blacklist.Lock()
	defer blacklist.Unlock()
	now := time.Now()
	for token, expiry := range blacklist.tokens {
		if now.After(expiry) {
			delete(blacklist.tokens, token)
		}
	}
}

// This function runs periodically to clean up expired tokens
func StartTokenCleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			CleanupExpiredTokens()
		}
	}()
}