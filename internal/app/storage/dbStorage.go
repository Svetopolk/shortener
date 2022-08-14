package storage

import "github.com/Svetopolk/shortener/internal/app/db"

type DBStorage struct {
	db db.Source
}

var _ Storage = &DBStorage{}

func NewDBStorage() *MemStorage {
	return &MemStorage{}
}

func (s *DBStorage) Save(hash string, url string) string {
	return hash
}

func (s *DBStorage) Get(hash string) (string, bool) {
	return "", false
}

func (s *DBStorage) GetAll() map[string]string {
	return map[string]string{}
}
