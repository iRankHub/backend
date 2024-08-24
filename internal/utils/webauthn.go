package utils

import (
	"crypto/rand"
	"encoding/binary"
)

// GenerateWebAuthnUserID generates a random byte array to use as WebAuthn user ID
func GenerateWebAuthnUserID() ([]byte, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Uint64ToBytes converts a uint64 to a byte slice
func Uint64ToBytes(i uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return b
}

// BytesToUint64 converts a byte slice to uint64
func BytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

// WebAuthnIDToInt32 converts a WebAuthn ID ([]byte) to int32
func WebAuthnIDToInt32(id []byte) int32 {
	// Ensure the id is at least 4 bytes long
	if len(id) < 4 {
		return 0
	}
	return int32(binary.BigEndian.Uint32(id[:4]))
}

// Int32ToWebAuthnID converts an int32 to a WebAuthn ID ([]byte)
func Int32ToWebAuthnID(id int32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(id))
	return b
}
