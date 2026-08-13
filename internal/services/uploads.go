package services

import (
	"errors"

	"github.com/Lekuruu/go-puush/internal/database"
	"github.com/Lekuruu/go-puush/internal/state"
	"gorm.io/gorm"
)

func CreateUpload(upload *database.Upload, state *state.State) error {
	result := state.Database.Create(upload)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func FetchUploadById(id int, state *state.State, preload ...string) (*database.Upload, error) {
	upload := &database.Upload{}
	result := preloadQuery(state, preload).First(upload, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return upload, nil
}

func FetchUploadByChecksumForUser(userId int, checksum string, state *state.State, preload ...string) (*database.Upload, error) {
	if checksum == "" {
		return nil, errors.New("checksum is empty")
	}

	upload := &database.Upload{}
	query := preloadQuery(state, preload).Where("user_id = ? AND checksum = ?", userId, checksum)
	result := query.First(upload)

	if result.Error != nil {
		return nil, result.Error
	}

	return upload, nil
}

func FetchRecentUploadsByUser(user *database.User, state *state.State, limit int, preload ...string) ([]*database.Upload, error) {
	var uploads []*database.Upload
	query := preloadQuery(state, preload).Where("user_id = ?", user.Id).Order("created_at DESC, id DESC").Limit(limit)
	result := query.Find(&uploads)

	if result.Error != nil {
		return nil, result.Error
	}

	return uploads, nil
}

func FetchLastPoolUpload(poolId int, state *state.State, preload ...string) (*database.Upload, error) {
	upload := &database.Upload{}
	query := preloadQuery(state, preload).Where("pool_id = ?", poolId).Order("created_at DESC, id DESC")
	result := query.First(upload)

	if result.Error != nil {
		return nil, result.Error
	}

	return upload, nil
}

func FetchUploadByFilenameAndPool(filename string, poolId int, state *state.State, preload ...string) (*database.Upload, error) {
	upload := &database.Upload{}
	query := preloadQuery(state, preload).Where("filename = ? AND pool_id = ?", filename, poolId)
	result := query.First(upload)

	if result.Error != nil {
		return nil, result.Error
	}

	return upload, nil
}

func FetchUploadsByPool(poolId int, offset int, limit int, state *state.State, preload ...string) ([]*database.Upload, error) {
	var uploads []*database.Upload
	query := preloadQuery(state, preload).
		Where("pool_id = ?", poolId).
		Order("created_at DESC, id DESC").
		Offset(offset).
		Limit(limit)
	result := query.Find(&uploads)

	if result.Error != nil {
		return nil, result.Error
	}

	return uploads, nil
}

func FetchPoolUploadCount(poolId int, state *state.State) (int64, error) {
	var count int64
	result := state.Database.Model(&database.Upload{}).Where("pool_id = ?", poolId).Count(&count)

	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func SearchUploadsFromPool(queryStr string, poolId int, offset int, limit int, state *state.State, preload ...string) ([]*database.Upload, error) {
	var uploads []*database.Upload
	query := preloadQuery(state, preload).Where("pool_id = ?", poolId).
		Where("filename LIKE ?", "%"+queryStr+"%").
		Order("created_at DESC, id DESC").
		Offset(offset).
		Limit(limit)
	result := query.Find(&uploads)

	if result.Error != nil {
		return nil, result.Error
	}

	// TODO: Full-text search / fuzzy search
	return uploads, nil
}

func UpdateUpload(upload *database.Upload, state *state.State) error {
	result := state.Database.Save(upload)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func UpdateUploadPool(uploadId int, poolId int, state *state.State) error {
	upload := &database.Upload{}
	result := state.Database.First(upload, uploadId)
	if result.Error != nil {
		return result.Error
	}

	upload.PoolId = poolId
	result = state.Database.Save(upload)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func UpdateUploadChecksum(uploadId int, checksum string, state *state.State) error {
	result := state.Database.Exec(
		"UPDATE uploads SET checksum = ? WHERE id = ?",
		checksum, uploadId,
	)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdatePoolUploadCount(poolId int, state *state.State) error {
	result := state.Database.Model(&database.Pool{}).
		Where("id = ?", poolId).
		UpdateColumn("upload_count", gorm.Expr("(SELECT COUNT(*) FROM uploads WHERE uploads.pool_id = pools.id)"))
	return result.Error
}

func UpdatePoolUploadCounts(user *database.User, state *state.State) error {
	result := state.Database.Model(&database.Pool{}).
		Where("user_id = ?", user.Id).
		UpdateColumn("upload_count", gorm.Expr("(SELECT COUNT(*) FROM uploads WHERE uploads.pool_id = pools.id)"))
	return result.Error
}

func DeleteUpload(upload *database.Upload, state *state.State) error {
	result := state.Database.Delete(upload)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

const minimumLinkIdentifierLength = 5
const maximumLinkIdentifierLength = 16

func CreateUploadIdentifier(uploadId int, state *state.State) (string, error) {
	identifier, err := GenerateUploadIdentifier(minimumLinkIdentifierLength, state)
	if err != nil {
		return "", err
	}

	// Update the upload with the shortlink identifier
	result := state.Database.Model(&database.Upload{}).
		Where("id = ?", uploadId).
		Update("identifier", identifier)

	if result.Error != nil {
		return "", result.Error
	}

	return identifier, nil
}

func FetchUploadByIdentifier(identifier string, state *state.State, preload ...string) (*database.Upload, error) {
	upload := &database.Upload{}
	query := preloadQuery(state, preload).Where("identifier = ?", identifier)
	result := query.First(upload)

	if result.Error != nil {
		return nil, result.Error
	}

	return upload, nil
}

func FetchManyUploadsByIdentifiers(identifiers []string, state *state.State, preload ...string) ([]*database.Upload, error) {
	var uploads []*database.Upload
	query := preloadQuery(state, preload).Where("identifier IN ?", identifiers)
	result := query.Find(&uploads)

	if result.Error != nil {
		return nil, result.Error
	}

	return uploads, nil
}

func UploadIdentifierExists(identifier string, state *state.State) (bool, error) {
	var count int64
	result := state.Database.Model(&database.Upload{}).Where("identifier = ?", identifier).Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

func GenerateUploadIdentifier(length int, state *state.State) (string, error) {
	if length < minimumLinkIdentifierLength {
		length = minimumLinkIdentifierLength
	}
	if length > maximumLinkIdentifierLength {
		length = maximumLinkIdentifierLength
	}

	for i := length; i <= maximumLinkIdentifierLength; i++ {
		identifier, err := randomIdentifier(i)
		if err != nil {
			return "", err
		}
		exists, err := UploadIdentifierExists(identifier, state)
		if err != nil {
			return "", err
		}
		if !exists {
			return identifier, nil
		}
	}
	return "", errors.New("could not generate unique identifier")
}
