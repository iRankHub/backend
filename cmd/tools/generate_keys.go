package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func main() {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Printf("Error generating key pair: %v\n", err)
		return
	}

	publicKeyEncoded := base64.StdEncoding.EncodeToString(publicKey)
	privateKeyEncoded := base64.StdEncoding.EncodeToString(privateKey)

	fmt.Printf("Public Key: %s\n", publicKeyEncoded)
	fmt.Printf("Private Key: %s\n", privateKeyEncoded)
}
