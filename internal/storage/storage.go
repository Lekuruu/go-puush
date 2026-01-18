package storage

import "io"

type Storage interface {
	Setup() error
	Save(key string, bucket string, data []byte) error
	SaveStream(key string, bucket string, stream io.Reader) error
	Read(key string, bucket string) ([]byte, error)
	ReadStream(key string, bucket string) (io.ReadSeekCloser, error)
	Remove(key string, bucket string) error
	Exists(key string, bucket string) bool

	SaveUpload(key string, data []byte) error
	SaveUploadStream(key string, stream io.Reader) error
	ReadUpload(key string) ([]byte, error)
	ReadUploadStream(key string) (io.ReadSeekCloser, error)
	RemoveUpload(key string) error
	UploadExists(key string) bool

	SaveThumbnail(key string, data []byte) error
	ReadThumbnail(key string) ([]byte, error)
	ReadThumbnailStream(key string) (io.ReadSeekCloser, error)
	RemoveThumbnail(key string) error
	ThumbnailExists(key string) bool
}
