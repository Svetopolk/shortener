package db

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
	db := initDBStorage(t)

	err := db.Ping()
	require.NoError(t, err)
}

func TestSaveGet(t *testing.T) {
	db := initDBStorage(t)

	hash := util.RandomString(5)
	url := "https://" + util.RandomString(5)

	assert.Nil(t, db.Save(hash, url))

	urlFromDB, ok := db.Get(hash)
	assert.True(t, ok)
	assert.Equal(t, url, urlFromDB)
}

func TestGetEmpty(t *testing.T) {
	db := initDBStorage(t)

	hash := util.RandomString(5)

	urlFromDB, ok := db.Get(hash)
	assert.False(t, ok)
	assert.Equal(t, "", urlFromDB)
}

func TestSaveSameHash(t *testing.T) {
	db := initDBStorage(t)

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
	db := initDBStorage(t)

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
	db := initDBStorage(t)

	hash1 := util.RandomString(5)
	url := "https://" + util.RandomString(5)

	err1 := db.Save(hash1, url)
	assert.Nil(t, err1)

	hash2, err2 := db.GetHashByURL(url)
	assert.Nil(t, err2)
	assert.Equal(t, hash1, hash2)
}

func TestSave_GetHashByURL_Empty(t *testing.T) {
	db := initDBStorage(t)

	url := "https://" + util.RandomString(5)

	hash, err := db.GetHashByURL(url)
	assert.NotNil(t, err)
	assert.Equal(t, "", hash)
}

func TestGetAll(t *testing.T) {
	db := initDBStorage(t)

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

func initDBStorage(t *testing.T) *Source {
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
	return db
}
