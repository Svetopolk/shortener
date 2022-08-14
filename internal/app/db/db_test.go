package db

import (
	"testing"

	"github.com/Svetopolk/shortener/internal/app/util"
	"github.com/stretchr/testify/assert"
)

func TestPingWrongPort(t *testing.T) {
	db := NewDB("postgres://shortener:pass@localhost:5433/shortener")
	err := db.Ping()
	assert.Contains(t, err.Error(), "failed to connect")
}

//before run this
//docker run --name postgresql -e POSTGRES_USER=shortener -e POSTGRES_PASSWORD=pass -p 5432:5432 -d postgres

func _TestPingOk(t *testing.T) {
	db := NewDB("postgres://shortener:pass@localhost:5432/shortener")
	err := db.Ping()
	assert.Equal(t, nil, err)
}

func _TestSaveGet(t *testing.T) {
	db := NewDB("postgres://shortener:pass@localhost:5432/shortener")
	hash := util.RandomString(5)

	db.Save(hash, "someUrl")

	urlFromDB := db.Get(hash)
	assert.Equal(t, "someUrl", urlFromDB)
}

func _TestGetEmpty(t *testing.T) {
	db := NewDB("postgres://shortener:pass@localhost:5432/shortener")
	hash := util.RandomString(5)

	urlFromDB := db.Get(hash)
	assert.Equal(t, "someUrl", urlFromDB)
}
