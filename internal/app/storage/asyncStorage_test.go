package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Svetopolk/shortener/internal/app/exceptions"
	"github.com/Svetopolk/shortener/internal/app/util"
)

func TestAsyncStorage(t *testing.T) {
	memStorage := NewMemStorage()
	asyncStorage := NewAsyncStorage(memStorage, 10)

	hashes := make([]string, 0)
	for i := 0; i < 15; i++ {
		hashes = append(hashes, util.RandomString(5))
		_, _ = memStorage.Save(hashes[i], "https://"+hashes[i])
	}
	_ = asyncStorage.BatchDelete(hashes[0:5])

	_, err := memStorage.Get(hashes[0])
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	_, err = memStorage.Get(hashes[0])
	assert.NotNil(t, err)
	assert.Equal(t, exceptions.ErrURLDeleted, err)
}
