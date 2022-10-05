package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Svetopolk/shortener/internal/app/exceptions"
	"github.com/Svetopolk/shortener/internal/app/util"
)

func TestMemStorage(t *testing.T) {
	storage := NewMemStorage()

	generatedHash1 := util.RandomString(10)
	url1 := "https://" + generatedHash1
	savedHash1, err1 := storage.Save(generatedHash1, url1)
	assert.Nil(t, err1)
	assert.Equal(t, generatedHash1, savedHash1)

	// save new url with same hash
	url2 := "https://" + util.RandomString(10)
	savedHash2, err2 := storage.Save(generatedHash1, url2)
	assert.NotNil(t, err2)
	assert.Error(t, exceptions.ErrHashAlreadyExist, err2)
	assert.Equal(t, generatedHash1, savedHash2)

	savedURL1, err := storage.Get(savedHash1)
	assert.Nil(t, err)
	assert.Equal(t, url1, savedURL1) // old value

	data, err4 := storage.GetAll()
	assert.Nil(t, err4)
	assert.GreaterOrEqual(t, len(data), 1)

	urlFromMap := data[generatedHash1]
	assert.Equal(t, url1, urlFromMap)
}

func TestMemStorageSaveBatch(t *testing.T) {
	storage := NewMemStorage()

	hash1 := util.RandomString(10)
	hash2 := util.RandomString(10)
	url1 := "https://" + hash1
	url2 := "https://" + hash2
	hashes := []string{hash1, hash2}
	urls := []string{url1, url2}

	savedHashes, err := storage.SaveBatch(hashes, urls)
	assert.Nil(t, err)
	assert.Equal(t, hashes, savedHashes)

	savedURL1, err1 := storage.Get(hash1)
	assert.Nil(t, err1)
	assert.Equal(t, url1, savedURL1)

	savedURL2, err2 := storage.Get(hash2)
	assert.Nil(t, err2)
	assert.Equal(t, url2, savedURL2)
}

func TestMemStorageDelete(t *testing.T) {
	storage := NewMemStorage()

	hash1 := util.RandomString(10)
	hash2 := util.RandomString(10)
	url1 := "https://" + hash1
	url2 := "https://" + hash2
	hashes := []string{hash1, hash2}
	urls := []string{url1, url2}

	savedHashes, err := storage.SaveBatch(hashes, urls)
	assert.Nil(t, err)
	assert.Equal(t, hashes, savedHashes)

	err = storage.Delete(hash1)
	assert.Nil(t, err)

	savedURL1, err1 := storage.Get(hash1)
	assert.NotNil(t, err1)
	assert.Contains(t, err1.Error(), "url is deleted")
	assert.Equal(t, url1, savedURL1)

	savedURL2, err2 := storage.Get(hash2)
	assert.Nil(t, err2)
	assert.Equal(t, url2, savedURL2)
}

func TestMemStorageBatchDelete(t *testing.T) {
	storage := NewMemStorage()

	hash1 := util.RandomString(10)
	hash2 := util.RandomString(10)
	hash3 := util.RandomString(10)
	url1 := "https://" + hash1
	url2 := "https://" + hash2
	url3 := "https://" + hash3
	hashes := []string{hash1, hash2, hash3}
	hashesToDelete := []string{hash1, hash2}
	urls := []string{url1, url2, url3}

	_, err := storage.SaveBatch(hashes, urls)
	assert.Nil(t, err)

	err = storage.BatchDelete(hashesToDelete)
	assert.Nil(t, err)

	savedURL1, err1 := storage.Get(hash1)
	assert.NotNil(t, err1)
	assert.Contains(t, err1.Error(), "url is deleted")
	assert.Equal(t, url1, savedURL1)

	savedURL2, err2 := storage.Get(hash2)
	assert.NotNil(t, err2)
	assert.Contains(t, err2.Error(), "url is deleted")
	assert.Equal(t, url2, savedURL2)

	savedURL3, err3 := storage.Get(hash3)
	assert.Nil(t, err3)
	assert.Equal(t, url3, savedURL3)
}
