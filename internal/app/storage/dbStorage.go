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
	_, ok := s.dbSource.Get(hash)
	if ok {
		return hash
	}
	s.dbSource.Save(hash, url)
	return hash
}

func (s *DBStorage) SaveBatch(hashes []string, urls []string) []string {
	// TODO there is no collision hashed check
	s.dbSource.SaveBatch(hashes, urls)
	return hashes
}

func (s *DBStorage) Get(hash string) (string, bool) {
	return s.dbSource.Get(hash)

}

func (s *DBStorage) GetAll() map[string]string {
	urls := s.dbSource.GetAll()
	return urls
}
