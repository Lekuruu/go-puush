package services

import (
	"time"

	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/database"
)

func CreateUser(email string, password string, state *app.State) (*database.User, error) {
	passwordHash, err := app.CreatePasswordHash(password)
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
		ApiKey:         app.GenerateApiKey(),
	}
	result := state.Database.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	publicPool := &database.Pool{
		UserId:     user.Id,
		Name:       "Public",
		Identifier: app.GeneratePoolIdentifier(),
		Type:       database.PoolTypePublic,
		CreatedAt:  time.Now(),
		LastUpload: time.Now(),
	}
	privatePool := &database.Pool{
		UserId:     user.Id,
		Name:       "Private",
		Identifier: app.GeneratePoolIdentifier(),
		Type:       database.PoolTypePrivate,
		CreatedAt:  time.Now(),
		LastUpload: time.Now(),
	}
	galleryPool := &database.Pool{
		UserId:     user.Id,
		Name:       "Gallery",
		Identifier: app.GeneratePoolIdentifier(),
		Type:       database.PoolTypeGallery,
		CreatedAt:  time.Now(),
		LastUpload: time.Now(),
	}

	result = state.Database.Create(publicPool)
	if result.Error != nil {
		return nil, result.Error
	}

	result = state.Database.Create(privatePool)
	if result.Error != nil {
		return nil, result.Error
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

func FetchUserById(id int, state *app.State, preload ...string) (*database.User, error) {
	user := &database.User{}
	result := preloadQuery(state, preload).First(user, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func FetchUserByName(name string, state *app.State, preload ...string) (*database.User, error) {
	user := &database.User{}
	query := preloadQuery(state, preload).Where("name = ?", name)
	result := query.First(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func FetchUserByEmail(email string, state *app.State, preload ...string) (*database.User, error) {
	user := &database.User{}
	query := preloadQuery(state, preload).Where("email = ?", email)
	result := query.First(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func FetchUserByNameOrEmail(input string, state *app.State, preload ...string) (*database.User, error) {
	user := &database.User{}
	query := preloadQuery(state, preload).Where("name = ? OR email = ?", input, input)
	result := query.First(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func FetchUserByApiKey(apiKey string, state *app.State, preload ...string) (*database.User, error) {
	user := &database.User{}
	query := preloadQuery(state, preload).Where("api_key = ?", apiKey)
	result := query.First(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func RegenerateUserApiKey(userId int, state *app.State) (string, error) {
	user, err := FetchUserById(userId, state)
	if err != nil {
		return "", err
	}

	user.ApiKey = app.GenerateApiKey()
	result := state.Database.Save(user)

	if result.Error != nil {
		return "", result.Error
	}

	return user.ApiKey, nil
}

func UpdateUserDiskUsage(userId int, size int64, state *app.State) error {
	result := state.Database.Exec(
		"UPDATE users SET disk_usage = disk_usage + ? WHERE id = ?",
		size, userId,
	)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func UpdateUserLatestActivity(userId int, state *app.State) error {
	result := state.Database.Exec(
		"UPDATE users SET latest_activity = ? WHERE id = ?",
		time.Now(), userId,
	)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func UpdateUserDefaultPool(userId int, poolId int, state *app.State) error {
	result := state.Database.Exec(
		"UPDATE users SET default_pool_id = ? WHERE id = ?",
		poolId, userId,
	)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func UpdateUserPassword(userId int, passwordHash string, state *app.State) error {
	result := state.Database.Exec(
		"UPDATE users SET password = ? WHERE id = ?",
		passwordHash, userId,
	)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func UpdateUser(user *database.User, state *app.State) error {
	result := state.Database.Save(user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
