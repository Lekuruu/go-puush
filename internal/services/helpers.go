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

var ErrNoUpdateColumns = errors.New("update requires explicit columns")

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

func preloadQuery(state *state.State, preload []string) *gorm.DB {
	result := state.Database

	for _, p := range preload {
		result = result.Preload(p)
	}

	return result
}

func updateRecord[T any](db *gorm.DB, updates *T, columns ...string) error {
	if len(columns) == 0 {
		return ErrNoUpdateColumns
	}
	result := db.Model(updates).Select(columns).Updates(updates)
	return requireUpdatedRows(result.RowsAffected, result.Error)
}

func requireUpdatedRows(rowsAffected int64, err error) error {
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func requireUpdatedRecord(result *gorm.DB) error {
	return requireUpdatedRows(result.RowsAffected, result.Error)
}
