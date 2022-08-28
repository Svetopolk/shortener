package storage

type Storage interface {
	Save(hash string, url string) (string, error)
	SaveBatch(hashes []string, urls []string) ([]string, error)
	Get(hash string) (string, bool)
	GetAll() map[string]string
}
