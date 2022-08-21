package storage

import (
	"testing"

	"github.com/Svetopolk/shortener/internal/app/db"
	"github.com/Svetopolk/shortener/internal/app/util"
	"github.com/stretchr/testify/assert"
)

func TestDBStorage(t *testing.T) {
	hash := util.RandomString(10)

	dbSource, err := db.NewDB("postgres://shortener:pass@localhost:5432/shortener")

	if err != nil {
		t.Error("error when NewDB", err)
	}

	err = dbSource.Ping()
	if err != nil {
		t.Skip("no db connection")
	}
	storage := NewDBStorage(dbSource)

	storage.Save(hash, "http://url")
	storage.Save(hash, "http://url2")

	url, _ := storage.Get(hash)
	assert.Equal(t, "http://url", url) //old value

	data := storage.GetAll()
	assert.GreaterOrEqual(t, len(data), 1)

	urlFromMap := data[hash]
	assert.Equal(t, "http://url", urlFromMap)
}
