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
			oldHash, err2 := s.dbSource.GetHashByURL(url)
			if err2 != nil {
				return "", err2 // somebody deleted url form db?
			}
			return oldHash, exceptions.ErrURLAlreadyExist
		}
		return "", err // unexpected error
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
	err := s.dbSource.SaveBatch(hashes, urls)
	if err != nil {
		return nil, err
	}
	return hashes, nil
}

func (s *DBStorage) Get(hash string) (string, error) {
	return s.dbSource.Get(hash)
}

func (s *DBStorage) GetAll() (map[string]string, error) {
	return s.dbSource.GetAll()
}

func (s *DBStorage) Delete(hash string) error {
	return s.dbSource.Delete(hash)
}

func (s *DBStorage) BatchDelete(hashes []string) error {
	return s.dbSource.BatchDelete(hashes)
}

func (s *DBStorage) Close() {
	s.dbSource.Close()
}
