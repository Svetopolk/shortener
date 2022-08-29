package service

import (
	"log"

	"github.com/Svetopolk/shortener/internal/logging"

	"github.com/Svetopolk/shortener/internal/app/storage"
	"github.com/Svetopolk/shortener/internal/app/util"
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

func (s *ShortServiceImpl) Get(hash string) (string, bool) {
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

	if _, ok := s.storage.Get(hash); !ok {
		return hash
	}

	newHash := hash + util.RandomString(1)
	log.Printf("hash %s already exist, generate a new one %s", hash, newHash)
	return s.checkOrChange(newHash)
}
