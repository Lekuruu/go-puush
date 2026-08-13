package services

import (
	"time"

	"github.com/Lekuruu/go-puush/internal/authentication"
	"github.com/Lekuruu/go-puush/internal/database"
	"github.com/Lekuruu/go-puush/internal/state"
)

func CreateUser(email string, password string, state *state.State) (*database.User, error) {
	passwordHash, err := authentication.CreatePasswordHash(password)
	if err != nil {
		return nil, err
	}
	apiKey, err := authentication.GenerateApiKey()
	if err != nil {
		return nil, err
	}

	user := &database.User{
		Name:           "",
		Email:          email,
		Password:       passwordHash,
		CreatedAt:      time.Now(),
		LatestActivity: time.Now(),
		Active:         !state.Config.Service.RequireActivation,
		ApiKey:         apiKey,
	}
	result := state.Database.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	// force the "active" field to be set correctly, since for some reason
	// gorm doesn't set boolean fields properly on creation
	result = state.Database.Model(user).Update("active", !state.Config.Service.RequireActivation)
	if result.Error != nil {
		return nil, result.Error
	}

	publicPoolIdentifier, err := GeneratePoolIdentifier(state)
	if err != nil {
		return nil, err
	}
	publicPool := &database.Pool{
		UserId:     user.Id,
		Name:       "Public",
		Identifier: publicPoolIdentifier,
		Type:       database.PoolTypePublic,
		CreatedAt:  time.Now(),
		LastUpload: time.Now(),
	}
	result = state.Database.Create(publicPool)
	if result.Error != nil {
		return nil, result.Error
	}

	privatePoolIdentifier, err := GeneratePoolIdentifier(state)
	if err != nil {
		return nil, err
	}
	privatePool := &database.Pool{
		UserId:     user.Id,
		Name:       "Private",
		Identifier: privatePoolIdentifier,
		Type:       database.PoolTypePrivate,
		CreatedAt:  time.Now(),
		LastUpload: time.Now(),
	}
	result = state.Database.Create(privatePool)
	if result.Error != nil {
		return nil, result.Error
	}

	galleryPoolIdentifier, err := GeneratePoolIdentifier(state)
	if err != nil {
		return nil, err
	}
	galleryPool := &database.Pool{
		UserId:     user.Id,
		Name:       "Gallery",
		Identifier: galleryPoolIdentifier,
		Type:       database.PoolTypeGallery,
		CreatedAt:  time.Now(),
		LastUpload: time.Now(),
	}

	result = state.Database.Create(galleryPool)
	if result.Error != nil {
		return nil, result.Error
	}

	user.DefaultPoolId = publicPool.Id
	result = state.Database.Save(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func FetchUserById(id int, state *state.State, preload ...string) (*database.User, error) {
	user := &database.User{}
	result := preloadQuery(state, preload).First(user, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func FetchUserByName(name string, state *state.State, preload ...string) (*database.User, error) {
	user := &database.User{}
	query := preloadQuery(state, preload).Where("name = ?", name)
	result := query.First(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func FetchUserByEmail(email string, state *state.State, preload ...string) (*database.User, error) {
	user := &database.User{}
	query := preloadQuery(state, preload).Where("email = ?", email)
	result := query.First(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func FetchUserByNameOrEmail(input string, state *state.State, preload ...string) (*database.User, error) {
	user := &database.User{}
	query := preloadQuery(state, preload).Where("name = ? OR email = ?", input, input)
	result := query.First(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func FetchUserByApiKey(apiKey string, state *state.State, preload ...string) (*database.User, error) {
	user := &database.User{}
	query := preloadQuery(state, preload).Where("api_key = ?", apiKey)
	result := query.First(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func RegenerateUserApiKey(userId int, state *state.State) (string, error) {
	user, err := FetchUserById(userId, state)
	if err != nil {
		return "", err
	}

	apiKey, err := authentication.GenerateApiKey()
	if err != nil {
		return "", err
	}
	user.ApiKey = apiKey
	result := state.Database.Save(user)

	if result.Error != nil {
		return "", result.Error
	}

	return user.ApiKey, nil
}

func UpdateUserDiskUsage(userId int, size int64, state *state.State) error {
	result := state.Database.Exec(
		"UPDATE users SET disk_usage = disk_usage + ? WHERE id = ?",
		size, userId,
	)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func UpdateUserLatestActivity(userId int, state *state.State) error {
	result := state.Database.Exec(
		"UPDATE users SET latest_activity = ? WHERE id = ?",
		time.Now(), userId,
	)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func UpdateUserDefaultPool(userId int, poolId int, state *state.State) error {
	result := state.Database.Exec(
		"UPDATE users SET default_pool_id = ? WHERE id = ?",
		poolId, userId,
	)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func UpdateUserPassword(userId int, passwordHash string, state *state.State) error {
	result := state.Database.Exec(
		"UPDATE users SET password = ? WHERE id = ?",
		passwordHash, userId,
	)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func ActivateUser(userId int, state *state.State) error {
	result := state.Database.Exec(
		"UPDATE users SET active = ? WHERE id = ?",
		true, userId,
	)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func UpdateUser(user *database.User, state *state.State) error {
	result := state.Database.Save(user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
