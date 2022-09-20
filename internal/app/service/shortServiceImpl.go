package service

import (
	"errors"
	"log"

	"github.com/Svetopolk/shortener/internal/app/exceptions"
	"github.com/Svetopolk/shortener/internal/app/storage"
	"github.com/Svetopolk/shortener/internal/app/util"
	"github.com/Svetopolk/shortener/internal/logging"
)

var _ ShortService = &ShortServiceImpl{}

type ShortServiceImpl struct {
	storage           storage.Storage
	initialHashLength int
}

func NewShortService(storage storage.Storage) *ShortServiceImpl {
	logging.Enter()
	defer logging.Exit()

	return &ShortServiceImpl{storage: storage, initialHashLength: 6}
}

func (s *ShortServiceImpl) Get(hash string) (string, error) {
	logging.Enter()
	defer logging.Exit()

	log.Printf("ShortServiceImpl: get key %s\n", hash)
	return s.storage.Get(hash)
}

func (s *ShortServiceImpl) GetAll() (map[string]string, error) {
	logging.Enter()
	defer logging.Exit()

	return s.storage.GetAll()
}

func (s *ShortServiceImpl) Save(url string) (string, error) {
	logging.Enter()
	defer logging.Exit()

	log.Printf("ShortServiceImpl: save url %s\n", url)
	hash := s.generateHash()
	hash, err := s.storage.Save(hash, url)
	return hash, err
}

func (s *ShortServiceImpl) SaveBatch(hashes []string, urls []string) ([]string, error) {
	logging.Enter()
	defer logging.Exit()
	return s.storage.SaveBatch(hashes, urls)
}

func (s *ShortServiceImpl) generateHash() string {
	logging.Enter()
	defer logging.Exit()

	return s.checkOrChange(util.RandomString(s.initialHashLength))
}

func (s *ShortServiceImpl) checkOrChange(hash string) string {
	logging.Enter()
	defer logging.Exit()

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
	log.Print("delete", hash)
	return nil
}

func (s *ShortServiceImpl) BatchDelete(hashes []string) error {
	log.Print("delete", hashes)
	return nil
}
