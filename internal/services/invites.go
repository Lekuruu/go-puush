package services

import (
	"time"

	"github.com/Lekuruu/go-puush/internal/authentication"
	"github.com/Lekuruu/go-puush/internal/database"
	"github.com/Lekuruu/go-puush/internal/state"
	"gorm.io/gorm"
)

const invitationKeyLength = 16

func CreateInvitationKey(expiry time.Duration, state *state.State) (*database.InvitationKey, error) {
	creationTime := time.Now()
	expiryTime := creationTime.Add(expiry)
	key, err := authentication.GenerateToken(invitationKeyLength)
	if err != nil {
		return nil, err
	}

	invitationKey := &database.InvitationKey{
		Key:       key,
		CreatedAt: creationTime,
		ExpiresAt: &expiryTime,
	}

	result := state.Database.Create(invitationKey)
	if result.Error != nil {
		return nil, result.Error
	}

	return invitationKey, nil
}

func IsValidInvitationKey(key string, state *state.State) (bool, error) {
	invitationKey := &database.InvitationKey{}
	result := state.Database.Where("key = ?", key).First(invitationKey)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, result.Error
	}

	if invitationKey.IsExpired() {
		return false, nil
	}

	return true, nil
}

func DeleteInvitationKey(key string, state *state.State) error {
	result := state.Database.Where("key = ?", key).Delete(&database.InvitationKey{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}
