package storage

import "github.com/Svetopolk/shortener/internal/app/exceptions"

var _ Storage = AsyncStorage{}

type AsyncStorage struct {
	storage Storage
}

func NewAsyncStorage(st Storage) Storage {
	return &AsyncStorage{storage: st}
}

func (a AsyncStorage) Save(hash string, url string) (string, error) {
	return "", exceptions.ErrNotImplemented
}

func (a AsyncStorage) SaveBatch(hashes []string, urls []string) ([]string, error) {
	return nil, exceptions.ErrNotImplemented
}

func (a AsyncStorage) Get(hash string) (string, error) {
	return "", exceptions.ErrNotImplemented
}

func (a AsyncStorage) GetAll() (map[string]string, error) {
	return nil, exceptions.ErrNotImplemented
}

func (a AsyncStorage) Delete(hash string) error {
	return exceptions.ErrNotImplemented
}

func (a AsyncStorage) BatchDelete(hashes []string) error {
	go func() {
		err := a.storage.BatchDelete(hashes)
		if err != nil {
		}
	}()

	return nil
}
