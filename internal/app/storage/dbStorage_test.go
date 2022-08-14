package storage

import (
	"testing"

	"github.com/Svetopolk/shortener/internal/app/db"
	"github.com/Svetopolk/shortener/internal/app/util"
	"github.com/stretchr/testify/assert"
)

func TestDBStorage(t *testing.T) {
	hash := util.RandomString(10)

	dbSource := db.NewDB("postgres://shortener:pass@localhost:5432/shortener")

	err := dbSource.Ping()
	if err != nil {
		t.Skip("no db connection")
	}
	storage := NewDBStorage(dbSource)

	storage.Save(hash, "url")
	storage.Save(hash, "url2")

	url, _ := storage.Get(hash)
	assert.Equal(t, "url", url) //old value

	data := storage.GetAll()
	assert.GreaterOrEqual(t, len(data), 1)

	urlFromMap := data[hash]
	assert.Equal(t, "url", urlFromMap)
}
