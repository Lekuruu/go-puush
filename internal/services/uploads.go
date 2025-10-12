package services

import (
	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/database"
)

func CreateUpload(upload *database.Upload, state *app.State) error {
	result := state.Database.Create(upload)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func FetchUploadById(id int, state *app.State, preload ...string) (*database.Upload, error) {
	upload := &database.Upload{}
	result := preloadQuery(state, preload).First(upload, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return upload, nil
}

func FetchUploadByChecksum(checksum string, state *app.State, preload ...string) (*database.Upload, error) {
	upload := &database.Upload{}
	query := preloadQuery(state, preload).Where("checksum = ?", checksum)
	result := query.First(upload)

	if result.Error != nil {
		return nil, result.Error
	}

	return upload, nil
}

func FetchRecentUploadsByUser(user *database.User, state *app.State, limit int, preload ...string) ([]*database.Upload, error) {
	var uploads []*database.Upload
	query := preloadQuery(state, preload).Where("user_id = ?", user.Id).Order("created_at DESC").Limit(limit)
	result := query.Find(&uploads)

	if result.Error != nil {
		return nil, result.Error
	}

	return uploads, nil
}

func FetchLastPoolUpload(poolId int, state *app.State, preload ...string) (*database.Upload, error) {
	upload := &database.Upload{}
	query := preloadQuery(state, preload).Where("pool_id = ?", poolId).Order("created_at DESC")
	result := query.First(upload)

	if result.Error != nil {
		return nil, result.Error
	}

	return upload, nil
}

func FetchUploadByFilenameAndPool(filename string, poolId int, state *app.State, preload ...string) (*database.Upload, error) {
	upload := &database.Upload{}
	query := preloadQuery(state, preload).Where("filename = ? AND pool_id = ?", filename, poolId)
	result := query.First(upload)

	if result.Error != nil {
		return nil, result.Error
	}

	return upload, nil
}

func FetchUploadsByPool(poolId int, state *app.State, preload ...string) ([]*database.Upload, error) {
	var uploads []*database.Upload
	query := preloadQuery(state, preload).Where("pool_id = ?", poolId).Order("created_at DESC")
	result := query.Find(&uploads)

	if result.Error != nil {
		return nil, result.Error
	}

	return uploads, nil
}

func FetchPoolUploadCount(poolId int, state *app.State) (int64, error) {
	var count int64
	result := state.Database.Model(&database.Upload{}).Where("pool_id = ?", poolId).Count(&count)

	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func UpdateUpload(upload *database.Upload, state *app.State) error {
	result := state.Database.Save(upload)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func UpdateUploadPool(uploadId int, poolId int, state *app.State) error {
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

func UpdatePoolUploadCount(poolId int, state *app.State) error {
	count, err := FetchPoolUploadCount(poolId, state)
	if err != nil {
		return err
	}

	pool := &database.Pool{}
	result := state.Database.First(pool, poolId)
	if result.Error != nil {
		return result.Error
	}

	pool.UploadCount = int(count)
	err = UpdatePool(pool, state)
	if err != nil {
		return err
	}

	return nil
}

func UpdatePoolUploadCounts(user *database.User, state *app.State) error {
	var pools []*database.Pool
	result := state.Database.Where("user_id = ?", user.Id).Find(&pools)
	if result.Error != nil {
		return result.Error
	}

	for _, pool := range pools {
		count, err := FetchPoolUploadCount(pool.Id, state)
		if err != nil {
			return err
		}
		pool.UploadCount = int(count)
		err = UpdatePool(pool, state)
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteUpload(upload *database.Upload, state *app.State) error {
	result := state.Database.Delete(upload)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
