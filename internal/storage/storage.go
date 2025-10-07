package storage

import "io"

type Storage interface {
	Setup() error
	Save(key string, bucket string, data []byte) error
	Read(key string, bucket string) ([]byte, error)
	ReadStream(key string, bucket string) (io.ReadCloser, error)
	Remove(key string, bucket string) error
	Exists(key string, bucket string) bool

	SaveUpload(key string, data []byte) error
	ReadUpload(key string) ([]byte, error)
	ReadUploadStream(key string) (io.ReadCloser, error)
	RemoveUpload(key string) error
	UploadExists(key string) bool

	SaveThumbnail(key string, data []byte) error
	ReadThumbnail(key string) ([]byte, error)
	ReadThumbnailStream(key string) (io.ReadCloser, error)
	RemoveThumbnail(key string) error
	ThumbnailExists(key string) bool
}
