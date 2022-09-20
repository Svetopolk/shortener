package service

type ShortService interface {
	Get(hash string) (string, error)
	GetAll() (map[string]string, error)
	Save(url string) (string, error)
	SaveBatch(hashes []string, urls []string) ([]string, error)
	Delete(url string) error
	BatchDelete(hashes []string) error
}
