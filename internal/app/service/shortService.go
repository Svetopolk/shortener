package service

type ShortService interface {
	Get(hash string) (string, bool)
	GetAll() (map[string]string, error)
	Save(url string) (string, error)
	SaveBatch(hashes []string, urls []string) ([]string, error)
}
