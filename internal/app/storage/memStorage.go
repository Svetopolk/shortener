package storage

import (
	"sync"

	"github.com/Svetopolk/shortener/internal/logging"
)

type MemStorage struct {
	data map[string]string
	mtx  sync.RWMutex
}

var _ Storage = &MemStorage{}

func NewMemStorage() *MemStorage {
	logging.Enter()
	defer logging.Exit()

	return &MemStorage{data: make(map[string]string)}
}

func (s *MemStorage) Save(hash string, url string) string {
	logging.Enter()
	defer logging.Exit()

	s.mtx.RLock()
	defer s.mtx.RUnlock()
	s.data[hash] = url
	return hash
}

func (s *MemStorage) SaveBatch(hashes []string, urls []string) []string {
	values := make([]string, 0, len(hashes))
	for i := range hashes {
		values = append(values, s.Save(hashes[i], urls[i]))
	}
	return values
}

func (s *MemStorage) Get(hash string) (string, bool) {
	logging.Enter()
	defer logging.Exit()

	s.mtx.Lock()
	defer s.mtx.Unlock()
	value, ok := s.data[hash]
	return value, ok
}

func (s *MemStorage) GetAll() map[string]string {
	logging.Enter()
	defer logging.Exit()

	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.data
}
