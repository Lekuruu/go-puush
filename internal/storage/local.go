package storage

import (
	"fmt"
	"os"
	"strings"
)

type FileStorage struct {
	dataPath string
}

func NewFileStorage(dataPath string) Storage {
	return &FileStorage{dataPath: dataPath}
}

func (storage *FileStorage) Setup() error {
	var folders = []string{
		fmt.Sprintf("%s/uploads", storage.dataPath),
		fmt.Sprintf("%s/thumbnails", storage.dataPath),
		fmt.Sprintf("%s/update", storage.dataPath),
	}

	for _, folder := range folders {
		if _, err := os.Stat(folder); !os.IsNotExist(err) {
			continue
		}

		err := os.MkdirAll(folder, 0755)
		if err != nil {
			return fmt.Errorf("failed to create storage directory %s: %w", folder, err)
		}

		if !strings.HasSuffix(folder, "update") {
			continue
		}

		// Copy required files for updates
		if _, err := os.Stat(fmt.Sprintf("%s/puush-rss.xml", folder)); os.IsNotExist(err) {
			data, err := os.ReadFile("./.github/puush-rss.xml")
			if err != nil {
				return fmt.Errorf("failed to write puush-rss.xml: %w", err)
			}

			err = os.WriteFile(fmt.Sprintf("%s/puush-rss.xml", folder), data, 0644)
			if err != nil {
				return fmt.Errorf("failed to write puush-rss.xml: %w", err)
			}
		}

		if _, err := os.Stat(fmt.Sprintf("%s/puush-win.txt", folder)); os.IsNotExist(err) {
			data, err := os.ReadFile("./.github/puush-win.txt")
			if err != nil {
				return fmt.Errorf("failed to write puush-win.txt: %w", err)
			}

			err = os.WriteFile(fmt.Sprintf("%s/puush-win.txt", folder), data, 0644)
			if err != nil {
				return fmt.Errorf("failed to write puush-win.txt: %w", err)
			}
		}
	}

	return nil
}

func (storage *FileStorage) Read(key string, folder string) ([]byte, error) {
	path := fmt.Sprintf("%s/%s/%s", storage.dataPath, folder, key)
	return os.ReadFile(path)
}

func (storage *FileStorage) Save(key string, folder string, data []byte) error {
	path := fmt.Sprintf("%s/%s", storage.dataPath, folder)
	err := os.MkdirAll(path, 0755)

	if err != nil {
		return err
	}

	os.WriteFile(fmt.Sprintf("%s/%s", path, key), data, os.ModePerm)
	return nil
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

func (storage *FileStorage) ReadUpdateConfigurationWindows() ([]byte, error) {
	return storage.Read("puush-win.txt", "update")
}

func (storage *FileStorage) ReadUpdateConfigurationMacOS() ([]byte, error) {
	return storage.Read("puush-rss.xml", "update")
}

func (storage *FileStorage) SaveUpload(key string, data []byte) error {
	return storage.Save(key, "uploads", data)
}

func (storage *FileStorage) ReadUpload(key string) ([]byte, error) {
	return storage.Read(key, "uploads")
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

func (storage *FileStorage) RemoveThumbnail(key string) error {
	return storage.Remove(key, "thumbnails")
}

func (storage *FileStorage) ThumbnailExists(key string) bool {
	return storage.Exists(key, "thumbnails")
}
