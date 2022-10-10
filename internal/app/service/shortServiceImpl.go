package service

import (
	"errors"
	"log"

	"github.com/Svetopolk/shortener/internal/app/exceptions"
	"github.com/Svetopolk/shortener/internal/app/storage"
	"github.com/Svetopolk/shortener/internal/app/util"
	_ "github.com/Svetopolk/shortener/internal/logging"
)

var _ ShortService = &ShortServiceImpl{}

type ShortServiceImpl struct {
	storage           storage.Storage
	asyncStorage      storage.Storage
	initialHashLength int
}

func NewShortService(st storage.Storage) *ShortServiceImpl {
	return &ShortServiceImpl{storage: st, asyncStorage: storage.NewAsyncStorage(st), initialHashLength: 6}
}

func (s *ShortServiceImpl) Get(hash string) (string, error) {
	log.Printf("ShortServiceImpl: get key %s\n", hash)
	return s.storage.Get(hash)
}

func (s *ShortServiceImpl) GetAll() (map[string]string, error) {
	return s.storage.GetAll()
}

func (s *ShortServiceImpl) Save(url string) (string, error) {
	log.Printf("ShortServiceImpl: save url %s\n", url)
	hash := s.generateHash()
	hash, err := s.storage.Save(hash, url)
	return hash, err
}

func (s *ShortServiceImpl) SaveBatch(hashes []string, urls []string) ([]string, error) {
	return s.storage.SaveBatch(hashes, urls)
}

func (s *ShortServiceImpl) generateHash() string {
	return s.checkOrChange(util.RandomString(s.initialHashLength))
}

func (s *ShortServiceImpl) checkOrChange(hash string) string {
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

func (s *ShortServiceImpl) Delete(hash string) error {
	return s.storage.Delete(hash)
}

func (s *ShortServiceImpl) BatchDelete(hashes []string) error {
	return s.asyncStorage.BatchDelete(hashes)
}

func (s *ShortServiceImpl) Shutdown() {
	s.asyncStorage.Shutdown()
}
