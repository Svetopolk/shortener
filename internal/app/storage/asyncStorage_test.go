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
	asyncStorage := NewAsyncStorage(memStorage)

	hashes := make([]string, 0)
	for i := 0; i < 15; i++ {
		hashes = append(hashes, util.RandomString(5))
		_, _ = memStorage.Save(hashes[i], "https://"+hashes[i])
	}
	_ = asyncStorage.BatchDelete(hashes[0:12])

	_, err1 := memStorage.Get(hashes[0])
	assert.Nil(t, err1)

	asyncStorage.BatchDeleteFromQueue()

	_, err2 := memStorage.Get(hashes[0])
	assert.Equal(t, exceptions.ErrURLDeleted, err2)

	_, err3 := memStorage.Get(hashes[11])
	assert.Nil(t, err3)

	time.Sleep(2 * time.Second)
	_, err4 := memStorage.Get(hashes[11])
	assert.Equal(t, exceptions.ErrURLDeleted, err4)
}
