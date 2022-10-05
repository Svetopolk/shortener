package service

import (
	"log"

	"github.com/Svetopolk/shortener/internal/app/exceptions"
	"github.com/Svetopolk/shortener/internal/logging"
)

var _ ShortService = &MockShortService{}

type MockShortService struct{}

func NewMockShortService() *MockShortService {
	return &MockShortService{}
}

func (s *MockShortService) Get(hash string) (string, error) {
	logging.Enter()
	defer logging.Exit()
	if hash == "12345" {
		return "https://ya.ru", nil
	}
	if hash == "_deleted_" {
		return "", exceptions.ErrURLDeleted
	}
	return "", exceptions.ErrURLNotFound
}

func (s *MockShortService) GetAll() (map[string]string, error) {
	logging.Enter()
	defer logging.Exit()
	data := make(map[string]string)
	data["12345"] = "https://ya.ru"
	return data, nil
}

func (s *MockShortService) Save(url string) (string, error) {
	logging.Enter()
	defer logging.Exit()

	if url == "https://ya.ru" {
		return "12345", nil
	}
	if url == "https://already.exist" {
		return "urlAlreadyExistHash", exceptions.ErrURLAlreadyExist
	}
	return "67890", nil
}

func (s *MockShortService) SaveBatch(hashes []string, urls []string) ([]string, error) {
	_ = urls
	values := make([]string, 0, len(hashes))
	for i := range hashes {
		values = append(values, hashes[i])
	}
	return values, nil
}

func (s *MockShortService) Delete(hash string) error {
	log.Print("delete ", hash)
	return nil
}

func (s *MockShortService) BatchDelete(hashes []string) error {
	log.Print("delete ", hashes)
	return nil
}
