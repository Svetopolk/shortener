package service

import (
	"log"

	"github.com/Svetopolk/shortener/internal/app/storage"
	"github.com/Svetopolk/shortener/internal/app/util"
)

type ShortService struct {
	storage           storage.Storage
	initialHashLength int
}

func NewShortService(storage storage.Storage) *ShortService {
	return &ShortService{storage: storage, initialHashLength: 6}
}

func (s *ShortService) Get(hash string) string {
	log.Printf("ShortService: get key %s\n", hash)
	value, _ := s.storage.Get(hash)
	return value
}

func (s *ShortService) Save(url string) string {
	log.Printf("ShortService: save url %s\n", url)
	return s.storage.Save(s.generateHash(), url)
}

func (s *ShortService) generateHash() string {
	return s.checkOrChange(util.RandomString(s.initialHashLength))
}

func (s *ShortService) checkOrChange(hash string) string {
	if _, ok := s.storage.Get(hash); !ok {
		return hash
	}
	hash = hash + util.RandomString(1)
	return s.checkOrChange(hash)
}
