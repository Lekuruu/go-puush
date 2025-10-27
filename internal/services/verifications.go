package services

import (
	"errors"
	"time"

	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/database"
	"gorm.io/gorm"
)

const verificationKeyLength = 32

func CreateEmailVerification(action database.EmailVerificationAction, expiry time.Duration, state *app.State) (*database.EmailVerification, error) {
	creationTime := time.Now()
	var expiresAt *time.Time

	if expiry > 0 {
		expiration := creationTime.Add(expiry)
		expiresAt = &expiration
	}

	verification := &database.EmailVerification{
		Key:       app.RandomString(verificationKeyLength),
		Action:    action,
		CreatedAt: creationTime,
		ExpiresAt: expiresAt,
	}

	if err := state.Database.Create(verification).Error; err != nil {
		return nil, err
	}

	return verification, nil
}

func FetchEmailVerificationByKey(key string, state *app.State, preload ...string) (*database.EmailVerification, error) {
	verification := &database.EmailVerification{}
	query := preloadQuery(state, preload).Where("key = ?", key)
	result := query.First(verification)

	if result.Error != nil {
		return nil, result.Error
	}

	return verification, nil
}

func ValidateEmailVerification(key string, action database.EmailVerificationAction, state *app.State, preload ...string) (*database.EmailVerification, error) {
	verification := &database.EmailVerification{}
	query := preloadQuery(state, preload).Where("key = ? AND action = ?", key, action)
	result := query.First(verification)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	if verification.IsExpired() {
		DeleteEmailVerificationByID(verification.Id, state)
		return nil, nil
	}

	return verification, nil
}

func DeleteEmailVerificationByID(id int, state *app.State) error {
	result := state.Database.Delete(&database.EmailVerification{}, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func DeleteEmailVerificationByKey(key string, state *app.State) error {
	result := state.Database.Where("key = ?", key).Delete(&database.EmailVerification{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func CleanupExpiredEmailVerifications(state *app.State) error {
	result := state.Database.Where("expires_at IS NOT NULL AND expires_at < ?", time.Now()).Delete(&database.EmailVerification{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}
