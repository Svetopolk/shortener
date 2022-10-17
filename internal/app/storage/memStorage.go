package storage

import (
	"log"
	"sync"

	"github.com/Svetopolk/shortener/internal/app/exceptions"
)

type MemRecord struct {
	url     string
	deleted bool
}

type MemStorage struct {
	data map[string]MemRecord
	mtx  sync.RWMutex
}

var _ Storage = &MemStorage{}

func NewMemStorage() *MemStorage {
	return &MemStorage{data: make(map[string]MemRecord)}
}

func (s *MemStorage) Save(hash string, url string) (string, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	_, ok := s.data[hash]
	if ok {
		return hash, exceptions.ErrURLAlreadyExist
	}
	s.data[hash] = MemRecord{url, false}
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
		if value.deleted {
			return value.url, exceptions.ErrURLDeleted
		}
		return value.url, nil
	}
	return value.url, exceptions.ErrURLNotFound
}

func (s *MemStorage) GetAll() (map[string]string, error) {
	log.Print("MemStorage GetAll")
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	output := make(map[string]string)
	for hash, record := range s.data {
		if !record.deleted {
			output[hash] = record.url
		}
	}
	return output, nil
}

func (s *MemStorage) Delete(hash string) error {
	log.Print("MemStorage Delete ", hash)
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	record := s.data[hash]
	record.deleted = true
	s.data[hash] = record
	return nil
}

func (s *MemStorage) BatchDelete(hashes []string) error {
	log.Print("MemStorage BatchDelete ", hashes)
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	for _, hash := range hashes {
		record := s.data[hash]
		record.deleted = true
		s.data[hash] = record
	}
	return nil
}

func (s *MemStorage) Close() {
}
