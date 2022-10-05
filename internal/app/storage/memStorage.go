package storage

import (
	"log"
	"sync"

	"github.com/Svetopolk/shortener/internal/app/exceptions"
)

type MemStorage struct {
	data map[string]string
	mtx  sync.RWMutex
}

var _ Storage = &MemStorage{}

func NewMemStorage() *MemStorage {
	return &MemStorage{data: make(map[string]string)}
}

func (s *MemStorage) Save(hash string, url string) (string, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	s.data[hash] = url
	return hash, nil
}

func (s *MemStorage) SaveBatch(hashes []string, urls []string) ([]string, error) {
	values := make([]string, 0, len(hashes))
	for i := range hashes {
		hash, _ := s.Save(hashes[i], urls[i])
		values = append(values, hash)
	}
	return values, nil
}

func (s *MemStorage) Get(hash string) (string, error) {
	log.Print("MemStorage Get ", hash)
	s.mtx.Lock()
	defer s.mtx.Unlock()
	value, ok := s.data[hash]
	if ok {
		return value, nil
	}
	return value, exceptions.ErrURLNotFound
}

func (s *MemStorage) GetAll() (map[string]string, error) {
	log.Print("MemStorage GetAll")
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.data, nil
}

func (s *MemStorage) Delete(hash string) error {
	log.Print("MemStorage Delete ", hash)
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	delete(s.data, hash)
	return nil
}

func (s *MemStorage) BatchDelete(hashes []string) error {
	log.Print("MemStorage BatchDelete ", hashes)
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	for _, hash := range hashes {
		delete(s.data, hash)
	}
	return nil
}
