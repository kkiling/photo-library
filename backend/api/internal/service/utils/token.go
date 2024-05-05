package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateCode() (string, error) {
	token := make([]byte, 8)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}

func GenerateAPIToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}
