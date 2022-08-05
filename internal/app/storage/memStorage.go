package storage

import "sync"

type MemStorage struct {
	mapStore map[string]string
	mtx      sync.RWMutex
}

var _ Storage = &MemStorage{}

func NewMemStorage() *MemStorage {
	return &MemStorage{mapStore: make(map[string]string)}
}

func (s *MemStorage) Save(hash string, url string) string {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	s.mapStore[hash] = url
	return hash
}

func (s *MemStorage) Get(hash string) (string, bool) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	value, ok := s.mapStore[hash]
	return value, ok
}

func (s *MemStorage) GetAll() map[string]string {
	//TODO implement me
	panic("implement me")
}
