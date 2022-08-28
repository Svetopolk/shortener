package storage

import (
	"testing"

	"github.com/Svetopolk/shortener/internal/app/db"
	"github.com/Svetopolk/shortener/internal/app/exceptions"
	"github.com/Svetopolk/shortener/internal/app/util"
	"github.com/stretchr/testify/assert"
)

func TestDBStorage(t *testing.T) {
	storage := initDbStorage(t)

	generatedHash1 := util.RandomString(10)
	url1 := "https://" + generatedHash1
	savedHash1, err1 := storage.Save(generatedHash1, url1)
	assert.Nil(t, err1)
	assert.Equal(t, generatedHash1, savedHash1)

	// save new url with same hash
	url2 := "https://" + util.RandomString(10)
	savedHash2, err2 := storage.Save(generatedHash1, url2)
	assert.NotNil(t, err2)
	assert.Error(t, exceptions.HashAlreadyExist, err2)
	assert.Equal(t, generatedHash1, savedHash2)

	// save same url with new hash
	generatedHash3 := util.RandomString(10)
	savedHash3, err3 := storage.Save(generatedHash3, url1)
	assert.NotNil(t, err3)
	assert.Error(t, exceptions.UrlAlreadyExist, err3)
	assert.NotEqual(t, generatedHash1, generatedHash3)
	assert.Equal(t, generatedHash1, savedHash3)

	savedUrl1, ok := storage.Get(savedHash1)
	assert.True(t, ok)
	assert.Equal(t, url1, savedUrl1) //old value

	data := storage.GetAll()
	assert.GreaterOrEqual(t, len(data), 1)

	urlFromMap := data[generatedHash1]
	assert.Equal(t, url1, urlFromMap)
}

func TestDBStorageSaveBatch(t *testing.T) {
	storage := initDbStorage(t)

	hash1 := util.RandomString(10)
	hash2 := util.RandomString(10)
	url1 := "https://" + hash1
	url2 := "https://" + hash2
	hashes := []string{hash1, hash2}
	urls := []string{url1, url2}

	storage.SaveBatch(hashes, urls)

	savedUrl1, ok := storage.Get(hash1)
	assert.True(t, ok)
	assert.Equal(t, url1, savedUrl1)

	savedUrl2, ok := storage.Get(hash2)
	assert.True(t, ok)
	assert.Equal(t, url2, savedUrl2)
}

func initDbStorage(t *testing.T) *DBStorage {
	dbSource, err := db.NewDB("postgres://shortener:pass@localhost:5432/shortener")

	if err != nil {
		t.Error("exceptions when NewDB", err)
	}

	err = dbSource.Ping()
	if err != nil {
		t.Skip("no db connection")
	}
	storage := NewDBStorage(dbSource)
	return storage
}
