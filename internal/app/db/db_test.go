package db

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Svetopolk/shortener/internal/app/exceptions"
	"github.com/Svetopolk/shortener/internal/app/util"
)

// launch docker before run tests otherwise most tests will be ignored
// docker run --name postgresql -e POSTGRES_USER=shortener -e POSTGRES_PASSWORD=pass -p 5432:5432 -d postgres

func TestPingWrongDBPort(t *testing.T) {
	db, err := NewDB("postgres://shortener:pass@localhost:5433/shortener")
	require.NoError(t, err)
	err = db.Ping()
	require.Error(t, err)
}

func TestPingOk(t *testing.T) {
	db := initDB(t)

	err := db.Ping()
	require.NoError(t, err)
}

func TestSaveGet(t *testing.T) {
	db := initDB(t)

	hash := util.RandomString(5)
	url := "https://" + util.RandomString(5)

	assert.Nil(t, db.Save(hash, url))

	urlFromDB, err := db.Get(hash)
	assert.Nil(t, err)
	assert.Equal(t, url, urlFromDB)
}

func TestGetEmpty(t *testing.T) {
	db := initDB(t)

	hash := util.RandomString(5)

	urlFromDB, err := db.Get(hash)
	assert.NotNil(t, err)
	assert.Equal(t, exceptions.ErrURLNotFound, err)
	assert.Equal(t, "", urlFromDB)
}

func TestSaveSameHash(t *testing.T) {
	db := initDB(t)

	hash := util.RandomString(5)
	url1 := "https://" + util.RandomString(5)
	url2 := "https://" + util.RandomString(5)

	err1 := db.Save(hash, url1)
	assert.Nil(t, err1)

	err2 := db.Save(hash, url2)
	assert.NotNil(t, err2)
	assert.Contains(t, err2.Error(), "urls_pk")
}

func TestSaveSameUrl(t *testing.T) {
	db := initDB(t)

	hash := util.RandomString(5)
	url1 := "https://" + util.RandomString(5)

	err1 := db.Save(hash, url1)
	assert.Nil(t, err1)

	hash2 := util.RandomString(5)

	err2 := db.Save(hash2, url1)
	assert.NotNil(t, err2)
	assert.Contains(t, err2.Error(), "urls_url_uindex")
}

func TestSave_GetHashByURL(t *testing.T) {
	db := initDB(t)

	hash1 := util.RandomString(5)
	url := "https://" + util.RandomString(5)

	err1 := db.Save(hash1, url)
	assert.Nil(t, err1)

	hash2, err2 := db.GetHashByURL(url)
	assert.Nil(t, err2)
	assert.Equal(t, hash1, hash2)
}

func TestSave_GetHashByURL_Empty(t *testing.T) {
	db := initDB(t)

	url := "https://" + util.RandomString(5)

	hash, err := db.GetHashByURL(url)
	assert.NotNil(t, err)
	assert.Equal(t, "", hash)
}

func TestGetAll(t *testing.T) {
	db := initDB(t)

	hash1 := util.RandomString(5)
	hash2 := util.RandomString(5)
	url1 := "https://" + hash1
	url2 := "https://" + hash2

	assert.Nil(t, db.Save(hash1, url1))
	assert.Nil(t, db.Save(hash2, url2))

	data, err := db.GetAll()
	assert.Nil(t, err)
	assert.Equal(t, data[hash1], url1)
	assert.Equal(t, data[hash2], url2)
}

func TestDelete(t *testing.T) {
	db := initDB(t)

	hash := util.RandomString(5)
	url := "https://" + util.RandomString(5)

	err1 := db.Save(hash, url)
	assert.Nil(t, err1)

	err2 := db.Delete(hash)
	assert.Nil(t, err2)

	url3, err3 := db.Get(hash)
	assert.NotNil(t, err3)
	assert.Equal(t, url, url3)
	assert.Contains(t, err3.Error(), "url is deleted")
}

func TestBatchDelete(t *testing.T) {
	db := initDB(t)

	hash1 := util.RandomString(5)
	hash2 := util.RandomString(5)
	hash3 := util.RandomString(5)

	_ = db.Save(hash1, "https://"+hash1)
	_ = db.Save(hash2, "https://"+hash2)
	_ = db.Save(hash3, "https://"+hash3)

	hashes := []string{hash1, hash2}
	err := db.BatchDelete(hashes)
	assert.Nil(t, err)

	_, err1 := db.Get(hash1)
	assert.NotNil(t, err1)
	assert.Contains(t, err1.Error(), "url is deleted")

	_, err2 := db.Get(hash2)
	assert.NotNil(t, err2)
	assert.Contains(t, err2.Error(), "url is deleted")

	_, err3 := db.Get(hash3)
	assert.Nil(t, err3)
}

func initDB(t *testing.T) *Source {
	db, err := NewDB("postgres://shortener:pass@localhost:5432/shortener")
	if err != nil {
		t.Skip("no db connection")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err = db.db.PingContext(ctx)

	if err != nil {
		log.Println("exceptions while ping DB:", err)
		t.Skip("no db connection")
	}
	db.InitTables()
	return db
}
