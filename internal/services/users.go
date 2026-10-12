package services

import (
	"time"

	"github.com/Lekuruu/go-puush/internal/authentication"
	"github.com/Lekuruu/go-puush/internal/database"
	"github.com/Lekuruu/go-puush/internal/state"
	"gorm.io/gorm"
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
	user.Active = !state.Config.Service.RequireActivation
	if err := updateRecord(state.Database, user, "active"); err != nil {
		return nil, err
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
	if err := updateRecord(state.Database, user, "default_pool_id"); err != nil {
		return nil, err
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
	apiKey, err := authentication.GenerateApiKey()
	if err != nil {
		return "", err
	}
	updates := &database.User{Id: userId, ApiKey: apiKey}

	if err := updateRecord(state.Database, updates, "api_key"); err != nil {
		return "", err
	}
	return apiKey, nil
}

func UpdateUserDiskUsage(userId int, size int64, state *state.State) error {
	result := state.Database.Model(&database.User{}).
		Where("id = ?", userId).
		UpdateColumn("disk_usage", gorm.Expr("disk_usage + ?", size))
	return requireUpdatedRecord(result)
}

func UpdateUserLatestActivity(userId int, state *state.State) error {
	updates := &database.User{Id: userId, LatestActivity: time.Now()}
	return updateRecord(state.Database, updates, "latest_activity")
}

func UpdateUserDefaultPool(userId int, poolId int, state *state.State) error {
	updates := &database.User{Id: userId, DefaultPoolId: poolId}
	return updateRecord(state.Database, updates, "default_pool_id")
}

func UpdateUserPassword(userId int, passwordHash string, state *state.State) error {
	updates := &database.User{Id: userId, Password: passwordHash}
	return updateRecord(state.Database, updates, "password")
}

func ActivateUser(userId int, state *state.State) error {
	updates := &database.User{Id: userId, Active: true}
	return updateRecord(state.Database, updates, "active")
}

func UpdateUserViewType(userId int, viewType database.ViewType, state *state.State) error {
	updates := &database.User{Id: userId, ViewType: viewType}
	return updateRecord(state.Database, updates, "view_type")
}

func ClaimUsername(userId int, username string, state *state.State) error {
	updates := &database.User{
		Id:                    userId,
		Name:                  username,
		UsernameSetupReminder: false,
	}
	return updateRecord(
		state.Database.Where("name = ?", ""),
		updates,
		"name",
		"username_setup_reminder",
	)
}

func UpdateUsernameSetupReminder(userId int, enabled bool, state *state.State) error {
	updates := &database.User{Id: userId, UsernameSetupReminder: enabled}
	return updateRecord(state.Database, updates, "username_setup_reminder")
}
