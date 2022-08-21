package storage

import (
	"github.com/Svetopolk/shortener/internal/app/db"
	"github.com/Svetopolk/shortener/internal/logging"
)

type DBStorage struct {
	dbSource *db.Source
}

var _ Storage = &DBStorage{}

func NewDBStorage(db *db.Source) *DBStorage {
	logging.Enter()
	defer logging.Exit()

	db.InitTables()
	return &DBStorage{dbSource: db}
}

func (s *DBStorage) Save(hash string, url string) string {
	urlFromDB := s.dbSource.Get(hash)
	if urlFromDB != "" {
		return hash
	}
	s.dbSource.Save(hash, url)
	return hash
}

func (s *DBStorage) Get(hash string) (string, bool) {
	url := s.dbSource.Get(hash)
	return url, true
}

func (s *DBStorage) GetAll() map[string]string {
	urls := s.dbSource.GetAll()
	return urls
}
