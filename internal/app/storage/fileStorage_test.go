package storage

import (
	"testing"

	"github.com/Svetopolk/shortener/internal/app/util"
	"github.com/stretchr/testify/assert"
)

func TestReadFromFileWhenCreated(t *testing.T) {
	hash := util.RandomString(10)

	fileStorage := NewFileStorage("file.log")
	fileStorage.writeToFile(hash, "url")

	url, _ := fileStorage.Get(hash)
	assert.Equal(t, "", url)

	anotherStorage := NewFileStorage("file.log")
	url, _ = anotherStorage.Get(hash)
	assert.Equal(t, "url", url)
}

func TestDoubleSave(t *testing.T) {
	hash := util.RandomString(10)

	fileStorage := NewFileStorage("file.log")
	fileStorage.Save(hash, "url")
	fileStorage.Save(hash, "url2")

	url, _ := fileStorage.Get(hash)
	assert.Equal(t, "url2", url)

	anotherStorage := NewFileStorage("file.log")
	url, _ = anotherStorage.Get(hash)
	assert.Equal(t, "url2", url)
}

func Test(t *testing.T) {
	assert.NotEmpty(t, NewFileStorage("/tmp/shortener/shortener.log"))
}

func TestDirNotExist(t *testing.T) {
	dir := util.RandomString(10)
	assert.NotEmpty(t, NewFileStorage("/tmp/shortener/"+dir+"/log.file"))
}
