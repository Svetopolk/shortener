package storage

type Storage interface {
	Save(hash string, url string) string
	Get(hash string) (string, bool)
	GetAll() map[string]string
}
