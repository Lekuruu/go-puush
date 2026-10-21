package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	storageDirectoryMode = 0700
	storageFileMode      = 0600
)

type FileStorage struct {
	dataPath string
}

func NewFileStorage(dataPath string) Storage {
	return &FileStorage{dataPath: dataPath}
}

func (storage *FileStorage) Setup() error {
	folders := []string{
		storage.dataPath,
		filepath.Join(storage.dataPath, "uploads"),
		filepath.Join(storage.dataPath, "thumbnails"),
	}

	for _, folder := range folders {
		if _, err := os.Stat(folder); !os.IsNotExist(err) {
			continue
		}

		err := os.MkdirAll(folder, storageDirectoryMode)
		if err != nil {
			return fmt.Errorf("failed to create storage directory %s: %w", folder, err)
		}
	}

	// Ensure that all files & directories have the correct permissions
	return filepath.WalkDir(storage.dataPath, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("storage: refusing to follow symbolic link %s", path)
		}
		if entry.IsDir() {
			return os.Chmod(path, storageDirectoryMode)
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("storage: unsupported file type at %s", path)
		}
		return os.Chmod(path, storageFileMode)
	})
}

func (storage *FileStorage) Read(key string, folder string) ([]byte, error) {
	path := fmt.Sprintf("%s/%s/%s", storage.dataPath, folder, key)
	return os.ReadFile(path)
}

func (storage *FileStorage) ReadStream(key string, folder string) (io.ReadSeekCloser, error) {
	path := fmt.Sprintf("%s/%s/%s", storage.dataPath, folder, key)
	return os.Open(path)
}

func (storage *FileStorage) Save(key string, folder string, data []byte) error {
	path := fmt.Sprintf("%s/%s", storage.dataPath, folder)
	err := os.MkdirAll(path, storageDirectoryMode)
	if err != nil {
		return err
	}

	return os.WriteFile(fmt.Sprintf("%s/%s", path, key), data, storageFileMode)
}

func (storage *FileStorage) SaveStream(key string, folder string, stream io.Reader) error {
	path := fmt.Sprintf("%s/%s", storage.dataPath, folder)
	err := os.MkdirAll(path, storageDirectoryMode)
	if err != nil {
		return err
	}

	filePath := fmt.Sprintf("%s/%s", path, key)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Chmod(storageFileMode)
	defer file.Close()

	buffer := make([]byte, 256*1024)
	_, err = io.CopyBuffer(file, stream, buffer)
	return err
}

func (storage *FileStorage) Exists(key string, folder string) bool {
	path := fmt.Sprintf("%s/%s/%s", storage.dataPath, folder, key)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func (storage *FileStorage) Remove(key string, folder string) error {
	path := fmt.Sprintf("%s/%s/%s", storage.dataPath, folder, key)
	return os.Remove(path)
}

func (storage *FileStorage) SaveUpload(key string, data []byte) error {
	return storage.Save(key, "uploads", data)
}

func (storage *FileStorage) SaveUploadStream(key string, stream io.Reader) error {
	return storage.SaveStream(key, "uploads", stream)
}

func (storage *FileStorage) ReadUpload(key string) ([]byte, error) {
	return storage.Read(key, "uploads")
}

func (storage *FileStorage) ReadUploadStream(key string) (io.ReadSeekCloser, error) {
	return storage.ReadStream(key, "uploads")
}

func (storage *FileStorage) RemoveUpload(key string) error {
	return storage.Remove(key, "uploads")
}

func (storage *FileStorage) UploadExists(key string) bool {
	return storage.Exists(key, "uploads")
}

func (storage *FileStorage) SaveThumbnail(key string, data []byte) error {
	return storage.Save(key, "thumbnails", data)
}

func (storage *FileStorage) ReadThumbnail(key string) ([]byte, error) {
	return storage.Read(key, "thumbnails")
}

func (storage *FileStorage) ReadThumbnailStream(key string) (io.ReadSeekCloser, error) {
	return storage.ReadStream(key, "thumbnails")
}

func (storage *FileStorage) RemoveThumbnail(key string) error {
	return storage.Remove(key, "thumbnails")
}

func (storage *FileStorage) ThumbnailExists(key string) bool {
	return storage.Exists(key, "thumbnails")
}
