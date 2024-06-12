package utils

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/o1egl/paseto"
)

func GenerateToken(userID int32, userRole string, privateKey ed25519.PrivateKey) (string, error) {
	// Create a new PASETO maker with version 2
	maker := paseto.NewV2()

	// Set the token claims
	claims := map[string]interface{}{
		"user_id":   userID,
		"user_role": userRole,
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	// Generate and return the token
	token, err := maker.Sign(privateKey, claims, nil)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return token, nil
}

func GeneratePasetoKeyPair() (ed25519.PrivateKey, ed25519.PublicKey, error) {
	// Generate an Ed25519 key pair
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate Ed25519 key pair: %v", err)
	}

	return privateKey, publicKey, nil
}