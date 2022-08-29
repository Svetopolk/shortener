package service

import (
	"github.com/Svetopolk/shortener/internal/app/exceptions"
	"github.com/Svetopolk/shortener/internal/logging"
)

var _ ShortService = &MockShortService{}

type MockShortService struct{}

func NewMockShortService() *MockShortService {
	return &MockShortService{}
}

func (s *MockShortService) Get(hash string) (string, bool) {
	logging.Enter()
	defer logging.Exit()
	if hash == "12345" {
		return "https://ya.ru", true
	}
	return "", false
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
	values := make([]string, 0, len(hashes))
	for i := range hashes {
		values = append(values, hashes[i])
	}
	return values, nil
}
