package utils

import (
	"crypto/rand"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+"

func HashPassword(password string) (string, error) {
    hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedBytes), nil
}

func ComparePasswords(hashedPassword, plainPassword string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

func GenerateRandomPassword() string {
    length := 12 // You can adjust the length as needed
    password := make([]byte, length)
    charsetLength := big.NewInt(int64(len(charset)))

    for i := range password {
        randomIndex, err := rand.Int(rand.Reader, charsetLength)
        if err != nil {
            // If there's an error, fall back to a default character
            password[i] = 'x'
        } else {
            password[i] = charset[randomIndex.Int64()]
        }
    }

    return string(password)
}