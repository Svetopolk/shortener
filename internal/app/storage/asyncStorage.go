package storage

import (
	"log"

	"github.com/Svetopolk/shortener/internal/app/exceptions"
)

var _ Storage = &AsyncStorage{}

type AsyncStorage struct {
	storage     Storage
	deleteQueue chan string
	batchSize   int
}

func NewAsyncStorage(st Storage, bufferSize int) *AsyncStorage {
	asyncStorage := &AsyncStorage{
		storage:     st,
		deleteQueue: make(chan string, bufferSize),
		batchSize:   10,
	}
	asyncStorage.startAsyncDeleter()
	return asyncStorage
}

func (a *AsyncStorage) Save(hash string, url string) (string, error) {
	return "", exceptions.ErrNotImplemented
}

func (a *AsyncStorage) SaveBatch(hashes []string, urls []string) ([]string, error) {
	return nil, exceptions.ErrNotImplemented
}

func (a *AsyncStorage) Get(hash string) (string, error) {
	return "", exceptions.ErrNotImplemented
}

func (a *AsyncStorage) GetAll() (map[string]string, error) {
	return nil, exceptions.ErrNotImplemented
}

func (a *AsyncStorage) Delete(hash string) error {
	return exceptions.ErrNotImplemented
}

func (a *AsyncStorage) BatchDelete(hashes []string) error {
	log.Println("BatchDelete() hashes=", hashes)
	for _, hash := range hashes {
		a.deleteQueue <- hash
	}
	return nil
}

func (a *AsyncStorage) Close() {
	log.Println("close deleted queue")
	close(a.deleteQueue)
}

func (a *AsyncStorage) startAsyncDeleter() {
	go func() {
		for hash := range a.deleteQueue {
			err := a.storage.Delete(hash)
			if err != nil {
				log.Println("delete error for hash=", hash, err)
			}
		}
		log.Println("deleteQueue closed, startAsyncDeleter closed")
	}()
}
