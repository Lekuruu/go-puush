package services

import (
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/Lekuruu/go-puush/internal/state"
	"gorm.io/gorm"
)

const identifierCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const identifierCharactersLength = int64(len(identifierCharacters))

func preloadQuery(state *state.State, preload []string) *gorm.DB {
	result := state.Database

	for _, p := range preload {
		result = result.Preload(p)
	}

	return result
}

func randomIdentifier(length int) (string, error) {
	if length < 0 {
		return "", errors.New("services: identifier length cannot be negative")
	}

	result := make([]byte, length)
	limit := big.NewInt(identifierCharactersLength)

	for i := range result {
		index, err := rand.Int(rand.Reader, limit)
		if err != nil {
			return "", err
		}
		result[i] = identifierCharacters[index.Int64()]
	}
	return string(result), nil
}
