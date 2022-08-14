package storage

import "github.com/Svetopolk/shortener/internal/app/db"

type DBStorage struct {
	dbSource *db.Source
}

var _ Storage = &DBStorage{}

func NewDBStorage(db *db.Source) *DBStorage {
	db.InitTables()
	return &DBStorage{dbSource: db}
}

func (s *DBStorage) Save(hash string, url string) string {
	s.dbSource.Save(hash, url)
	return hash
}

func (s *DBStorage) Get(hash string) (string, bool) {
	url := s.dbSource.Get(hash)
	return url, true
}

func (s *DBStorage) GetAll() map[string]string {
	return map[string]string{}
}
