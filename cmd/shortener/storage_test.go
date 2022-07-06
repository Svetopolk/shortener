package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorage_save(t *testing.T) {
	s := NewMemStorage()

	hash := s.save("https://ya.ru")
	url := s.get(hash)
	assert.Equal(t, "https://ya.ru", url)
	assert.Equal(t, "", s.get("hashDoesNotExist"))
}
