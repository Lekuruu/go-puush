package services

import (
	"time"

	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/database"
	"gorm.io/gorm"
)

const invitationKeyLength = 16

func CreateInvitationKey(expiry time.Duration, state *app.State) (*database.InvitationKey, error) {
	creationTime := time.Now()
	expiryTime := creationTime.Add(expiry)

	invitationKey := &database.InvitationKey{
		Key:       app.RandomString(invitationKeyLength),
		CreatedAt: creationTime,
		ExpiresAt: &expiryTime,
	}

	result := state.Database.Create(invitationKey)
	if result.Error != nil {
		return nil, result.Error
	}

	return invitationKey, nil
}

func IsValidInvitationKey(key string, state *app.State) (bool, error) {
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

func DeleteInvitationKey(key string, state *app.State) error {
	result := state.Database.Where("key = ?", key).Delete(&database.InvitationKey{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}
