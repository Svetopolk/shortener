package storage

import "sync"

type MemStorage struct {
	data map[string]string
	mtx  sync.RWMutex
}

var _ Storage = &MemStorage{}

func NewMemStorage() *MemStorage {
	return &MemStorage{data: make(map[string]string)}
}

func (s *MemStorage) Save(hash string, url string) string {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	s.data[hash] = url
	return hash
}

func (s *MemStorage) Get(hash string) (string, bool) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	value, ok := s.data[hash]
	return value, ok
}

func (s *MemStorage) GetAll() map[string]string {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.data
}
