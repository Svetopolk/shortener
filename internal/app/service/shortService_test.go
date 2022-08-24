package service

import (
	"log"
	"testing"

	"github.com/Svetopolk/shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
)

func TestShotServiceSaveGet(t *testing.T) {
	s := NewShortService(storage.NewMemStorage())

	hash := s.Save("https://ya.ru")
	assert.Len(t, hash, 6)
	url := s.Get(hash)
	assert.Equal(t, "https://ya.ru", url)
	assert.Equal(t, "", s.Get("hashDoesNotExist"))
}

func TestShotServiceHashCollision(t *testing.T) {
	s := NewShortService(&MockStorage{})
	hash := s.Save("https://ya.ru")
	assert.Len(t, hash, 7)
}

type MockStorage struct {
	requestCount int
}

var _ storage.Storage = &MockStorage{}

func (m *MockStorage) Save(hash string, _ string) string {
	return hash
}

func (m *MockStorage) SaveBatch(hashes []string, urls []string) []string {
	values := make([]string, 0, len(hashes))
	for i := range hashes {
		values = append(values, m.Save(hashes[i], urls[i]))
	}
	return values
}

func (m *MockStorage) Get(hash string) (string, bool) {
	log.Default().Println("mock storage get with hash: ", hash)
	if m.requestCount > 0 {
		return "", false
	}
	m.requestCount++
	return "hashExists", true
}

func (m *MockStorage) GetAll() map[string]string {
	return make(map[string]string)
}
