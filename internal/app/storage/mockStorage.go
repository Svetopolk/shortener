package storage

import (
	"log"

	"github.com/Svetopolk/shortener/internal/app/exceptions"
)

type MockStorage struct {
	requestCount int
}

var _ Storage = &MockStorage{}

func NewMockStorage() *MockStorage {
	return &MockStorage{}
}

func (m *MockStorage) Save(hash string, url string) (string, error) {
	if url == "https://ya.ru" {
		return "12345", nil
	}
	if url == "https://already.exist" {
		return "urlAlreadyExistHash", exceptions.ErrURLAlreadyExist
	}
	return "67890", nil
}

func (m *MockStorage) SaveBatch(hashes []string, urls []string) ([]string, error) {
	_ = urls
	values := make([]string, 0, len(hashes))
	for i := range hashes {
		values = append(values, hashes[i])
	}
	return values, nil
}

func (m *MockStorage) Get(hash string) (string, error) {
	if hash == "12345" {
		return "https://ya.ru", nil
	}
	if hash == "_deleted_" {
		return "", exceptions.ErrURLDeleted
	}
	return "", exceptions.ErrURLNotFound
}

func (m *MockStorage) GetAll() (map[string]string, error) {
	data := make(map[string]string)
	data["12345"] = "https://ya.ru"
	return data, nil
}

func (m *MockStorage) Delete(hash string) error {
	log.Print("delete ", hash)
	return nil
}

func (m *MockStorage) BatchDelete(hashes []string) error {
	if len(hashes) > 0 {
		log.Print("delete ", hashes)
	}
	return nil
}

func (m *MockStorage) Close() {
}
