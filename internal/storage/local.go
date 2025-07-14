package storage

import (
	"fmt"
	"io"
	"net/http"
	"os"
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
	}

	for _, folder := range folders {
		if _, err := os.Stat(folder); !os.IsNotExist(err) {
			continue
		}

		err := os.MkdirAll(folder, 0755)
		if err != nil {
			return fmt.Errorf("failed to create storage directory %s: %w", folder, err)
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

func (storage *FileStorage) Remove(key string, folder string) error {
	path := fmt.Sprintf("%s/%s/%s", storage.dataPath, folder, key)
	return os.Remove(path)
}

func (storage *FileStorage) Download(url string, key string, folder string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Download failed: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return storage.Save(key, folder, data)
}
