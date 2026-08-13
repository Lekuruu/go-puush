package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/Lekuruu/go-puush/internal/database"
	"github.com/Lekuruu/go-puush/internal/state"
)

func CreateSession(userId int, duration time.Duration, state *state.State) (*database.Session, error) {
	token, err := generateSessionToken()
	if err != nil {
		return nil, err
	}

	session := &database.Session{
		UserId:    userId,
		Token:     token,
		ExpiresAt: time.Now().Add(duration),
	}

	if err := state.Database.Create(session).Error; err != nil {
		return nil, err
	}
	return session, nil
}

func ValidateSession(token string, state *state.State, preload ...string) (*database.Session, error) {
	session := &database.Session{}
	query := preloadQuery(state, preload).Where("token = ?", token)
	result := query.First(session)

	if result.Error != nil {
		return nil, result.Error
	}

	if session.IsExpired() {
		DeleteSession(token, state)
		return nil, errors.New("session expired")
	}

	return session, nil
}

func DeleteSession(token string, state *state.State) error {
	result := state.Database.Where("token = ?", token).Delete(&database.Session{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("session not found")
	}

	return nil
}

func DeleteAllSessionsForUser(userId int, state *state.State) error {
	result := state.Database.Where("user_id = ?", userId).Delete(&database.Session{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func ExtendSession(session *database.Session, duration time.Duration, state *state.State) error {
	session.ExpiresAt = time.Now().Add(duration)
	result := state.Database.Save(session)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func CleanupExpiredSessions(state *state.State) error {
	result := state.Database.Where("expires_at < ?", time.Now()).Delete(&database.Session{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func generateSessionToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
