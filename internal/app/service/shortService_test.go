package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Svetopolk/shortener/internal/app/exceptions"
	"github.com/Svetopolk/shortener/internal/app/storage"
)

func TestShotServiceSaveGet(t *testing.T) {
	s := NewShortService(storage.NewMemStorage())

	hash, err := s.Save("https://ya.ru")
	assert.Nil(t, err)
	assert.Len(t, hash, 6)
	url, err2 := s.Get(hash)
	assert.Equal(t, "https://ya.ru", url)
	assert.Nil(t, err2)

	url2, err3 := s.Get("hashDoesNotExist")
	assert.Equal(t, "", url2)
	assert.NotNil(t, err3)
	assert.Equal(t, exceptions.ErrURLNotFound, err3)
}
