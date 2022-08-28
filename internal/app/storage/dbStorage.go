package storage

import (
	"strings"

	"github.com/Svetopolk/shortener/internal/app/db"
	"github.com/Svetopolk/shortener/internal/app/exceptions"
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

func (s *DBStorage) Save(hash string, url string) (string, error) {
	err := s.dbSource.Save(hash, url)
	if err != nil {
		if isHashUniqueViolation(err) {
			return hash, exceptions.ErrHashAlreadyExist
		}
		if isURLUniqueViolation(err) {
			oldHash, ok := s.dbSource.GetHashByURL(url)
			if ok {
				return oldHash, exceptions.ErrURLAlreadyExist
			}
			return "", err // somebody deleted url form db?
		}
		return "", err //unexpected error
	}
	return hash, nil
}

func isHashUniqueViolation(err error) bool {
	return strings.Contains(err.Error(), "urls_pk")
}

func isURLUniqueViolation(err error) bool {
	return strings.Contains(err.Error(), "urls_url_uindex")
}

func (s *DBStorage) SaveBatch(hashes []string, urls []string) ([]string, error) {
	// TODO there is no collision hashed check
	s.dbSource.SaveBatch(hashes, urls)
	return hashes, nil
}

func (s *DBStorage) Get(hash string) (string, bool) {
	return s.dbSource.Get(hash)

}

func (s *DBStorage) GetAll() map[string]string {
	urls := s.dbSource.GetAll()
	return urls
}
