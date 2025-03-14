package utils

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/o1egl/paseto"
)

var (
	publicKey  ed25519.PublicKey
	privateKey ed25519.PrivateKey
	blacklist  sync.Map
)

func InitializeTokenConfig() error {
	publicKeyStr := os.Getenv("TOKEN_PUBLIC_KEY")
	privateKeyStr := os.Getenv("TOKEN_PRIVATE_KEY")

	if publicKeyStr == "" || privateKeyStr == "" {
		pub, priv, err := GeneratePasetoKeyPair()
		if err != nil {
			return fmt.Errorf("failed to generate PASETO key pair: %v", err)
		}

		publicKey = pub
		privateKey = priv

		os.Setenv("TOKEN_PUBLIC_KEY", base64.StdEncoding.EncodeToString(pub))
		os.Setenv("TOKEN_PRIVATE_KEY", base64.StdEncoding.EncodeToString(priv))
	} else {
		pub, err := base64.StdEncoding.DecodeString(publicKeyStr)
		if err != nil {
			return fmt.Errorf("failed to decode public key: %v", err)
		}
		publicKey = ed25519.PublicKey(pub)

		priv, err := base64.StdEncoding.DecodeString(privateKeyStr)
		if err != nil {
			return fmt.Errorf("failed to decode private key: %v", err)
		}
		privateKey = ed25519.PrivateKey(priv)
	}

	return nil
}

func GenerateToken(userID int32, userName, userRole, userEmail string) (string, error) {
	maker := paseto.NewV2()

	claims := map[string]interface{}{
		"user_id":    float64(userID), // Convert to float64 to ensure consistency
		"user_name":  userName,
		"user_role":  userRole,
		"user_email": userEmail,
		"exp":        float64(time.Now().Add(time.Hour * 168).Unix()), // Convert to float64
	}

	token, err := maker.Sign(privateKey, claims, nil)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return token, nil
}

func GeneratePasetoKeyPair() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate Ed25519 key pair: %v", err)
	}

	return publicKey, privateKey, nil
}

func InvalidateToken(token string) {
	blacklist.Store(token, time.Now().Add(168*time.Hour)) // Invalidate for 7 days
}

func IsTokenInvalid(token string) bool {
	expiry, exists := blacklist.Load(token)
	if !exists {
		return false
	}
	if time.Now().After(expiry.(time.Time)) {
		blacklist.Delete(token)
		return false
	}
	return true
}

func ValidateToken(token string) (map[string]interface{}, error) {
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

func CleanupExpiredTokens() {
	now := time.Now()
	blacklist.Range(func(key, value interface{}) bool {
		expiry := value.(time.Time)
		if now.After(expiry) {
			blacklist.Delete(key)
		}
		return true
	})
}

func StartTokenCleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			CleanupExpiredTokens()
		}
	}()
}
