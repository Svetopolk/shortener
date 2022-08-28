package service

import (
	"testing"

	"github.com/Svetopolk/shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
)

func TestShotServiceSaveGet(t *testing.T) {
	s := NewShortService(storage.NewMemStorage())

	hash, err := s.Save("https://ya.ru")
	assert.Nil(t, err)
	assert.Len(t, hash, 6)
	url, ok := s.Get(hash)
	assert.Equal(t, "https://ya.ru", url)
	assert.True(t, ok)

	url2, ok := s.Get("hashDoesNotExist")
	assert.Equal(t, "", url2)
	assert.False(t, ok)
}
