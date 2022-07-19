package storage

type MemStorage struct {
	mapStore map[string]string
}

var _ Storage = &MemStorage{}

func NewMemStorage() *MemStorage {
	return &MemStorage{mapStore: make(map[string]string)}
}

func (s *MemStorage) Save(hash string, url string) string {
	s.mapStore[hash] = url
	return hash
}

func (s *MemStorage) Get(hash string) (string, bool) {
	value, ok := s.mapStore[hash]
	return value, ok
}
