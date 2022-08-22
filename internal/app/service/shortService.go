package service

import (
	"log"

	"github.com/Svetopolk/shortener/internal/logging"

	"github.com/Svetopolk/shortener/internal/app/storage"
	"github.com/Svetopolk/shortener/internal/app/util"
)

type ShortService struct {
	storage           storage.Storage
	initialHashLength int
}

func NewShortService(storage storage.Storage) *ShortService {
	logging.Enter()
	defer logging.Exit()

	return &ShortService{storage: storage, initialHashLength: 6}
}

func (s *ShortService) Get(hash string) string {
	logging.Enter()
	defer logging.Exit()

	log.Printf("ShortService: get key %s\n", hash)
	value, _ := s.storage.Get(hash)
	return value
}

func (s *ShortService) GetAll() map[string]string {
	logging.Enter()
	defer logging.Exit()

	values := s.storage.GetAll()
	return values
}

func (s *ShortService) Save(url string) string {
	logging.Enter()
	defer logging.Exit()

	log.Printf("ShortService: save url %s\n", url)
	return s.storage.Save(s.generateHash(), url)
}

func (s *ShortService) generateHash() string {
	logging.Enter()
	defer logging.Exit()

	return s.checkOrChange(util.RandomString(s.initialHashLength))
}

func (s *ShortService) checkOrChange(hash string) string {
	logging.Enter()
	defer logging.Exit()

	if _, ok := s.storage.Get(hash); !ok {
		return hash
	}

	newHash := hash + util.RandomString(1)
	log.Printf("hash %s already exist, generate a new one %s", hash, newHash)
	return s.checkOrChange(newHash)
}
