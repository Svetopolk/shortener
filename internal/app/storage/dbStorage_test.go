package storage

import (
	"testing"

	"github.com/Svetopolk/shortener/internal/app/db"
	"github.com/Svetopolk/shortener/internal/app/util"
	"github.com/stretchr/testify/assert"
)

func TestDBStorageDoubleSave(t *testing.T) {
	hash := util.RandomString(10)

	dbService := db.NewDB("postgres://shortener:pass@localhost:5432/shortener")

	err := dbService.Ping()
	if err != nil {
		t.Skip("no db connection")
	}
	storage := NewDBStorage(dbService)

	storage.Save(hash, "url")
	storage.Save(hash, "url2")

	url, _ := storage.Get(hash)
	assert.Equal(t, "url", url) //old value
}
