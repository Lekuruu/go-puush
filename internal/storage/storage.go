package storage

type Storage interface {
	Setup() error
	Save(key string, bucket string, data []byte) error
	Read(key string, bucket string) ([]byte, error)
	Remove(key string, bucket string) error
	Download(url string, key string, bucket string) error
}
