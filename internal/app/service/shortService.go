package service

import (
	"errors"
	"log"

	"github.com/Svetopolk/shortener/internal/app/exceptions"
	"github.com/Svetopolk/shortener/internal/app/storage"
	"github.com/Svetopolk/shortener/internal/app/util"
	_ "github.com/Svetopolk/shortener/internal/logging"
)

type ShortService struct {
	storage           storage.Storage
	asyncStorage      storage.Storage
	initialHashLength int
}

func NewShortService(st storage.Storage) *ShortService {
	return &ShortService{storage: st, asyncStorage: storage.NewAsyncStorage(st), initialHashLength: 6}
}

func (s *ShortService) Get(hash string) (string, error) {
	log.Printf("ShortService: get key %s\n", hash)
	return s.storage.Get(hash)
}

func (s *ShortService) GetAll() (map[string]string, error) {
	return s.storage.GetAll()
}

func (s *ShortService) Save(url string) (string, error) {
	log.Printf("ShortService: save url %s\n", url)
	hash := s.generateHash()
	hash, err := s.storage.Save(hash, url)
	return hash, err
}

func (s *ShortService) SaveBatch(hashes []string, urls []string) ([]string, error) {
	return s.storage.SaveBatch(hashes, urls)
}

func (s *ShortService) generateHash() string {
	return s.checkOrChange(util.RandomString(s.initialHashLength))
}

func (s *ShortService) checkOrChange(hash string) string {
	_, err := s.storage.Get(hash)
	if err != nil {
		if errors.Is(err, exceptions.ErrURLNotFound) {
			return hash
		} else {
			log.Fatal("unexpected error", err)
		}
	}
	newHash := hash + util.RandomString(1)
	log.Printf("hash %s already exist, generate a new one %s", hash, newHash)
	return s.checkOrChange(newHash)
}

func (s *ShortService) Delete(hash string) error {
	return s.storage.Delete(hash)
}

func (s *ShortService) BatchDelete(hashes []string) error {
	return s.asyncStorage.BatchDelete(hashes)
}

func (s *ShortService) Shutdown() {
	s.asyncStorage.Close()
}
