package authentication

import (
	"crypto/rand"
	"errors"
	"math/big"
)

const tokenCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateToken returns a cryptographically secure alphanumeric token.
func GenerateToken(length int) (string, error) {
	if length < 0 {
		return "", errors.New("authentication: token length cannot be negative")
	}

	result := make([]byte, length)
	for i := range result {
		index, err := rand.Int(rand.Reader, new(big.Int).SetInt64(int64(len(tokenCharacters))))
		if err != nil {
			return "", err
		}
		result[i] = tokenCharacters[index.Int64()]
	}
	return string(result), nil
}

func GenerateApiKey() (string, error) {
	return GenerateToken(32)
}
