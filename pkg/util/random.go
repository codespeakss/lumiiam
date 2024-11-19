package util

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"fmt"
)

func GenerateRefreshToken() (string, error) {
	size := 32
	randomBytes := make([]byte, size)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}

	refreshToken := base64.URLEncoding.EncodeToString(randomBytes)
	return refreshToken, nil
}

func GenerateToken() (string, error) {
	size := 32 // Desired length in bytes
	randomBytes := make([]byte, size)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}

	// Encode to base32 to avoid special characters
	// Base32 output is uppercase; you can use strings.ToLower if lowercase is preferred
	token := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	return token, nil
}
