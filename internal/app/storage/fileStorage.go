package storage

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/Svetopolk/shortener/internal/app/exceptions"
)

type FileStorage struct {
	data            map[string]string
	fileStoragePath string
	producer        *producer
	mtx             sync.RWMutex
}

var _ Storage = &FileStorage{}

func NewFileStorage(fileStoragePath string) *FileStorage {
	log.Print("FileStorage NewFileStorage")

	checkDirExistOrCreate(fileStoragePath)
	mapStore := readFromFileIntoMap(fileStoragePath)

	fileProducer, err := NewProducer(fileStoragePath)
	if err != nil {
		log.Println("can not create NewFileStorage", err)
	}
	return &FileStorage{data: mapStore, fileStoragePath: fileStoragePath, producer: fileProducer}
}

func (s *FileStorage) Save(hash string, url string) (string, error) {
	log.Print("FileStorage Save	", hash)

	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.data[hash] = url
	s.writeToFile(hash, url)
	return hash, nil
}

func (s *FileStorage) SaveBatch(hashes []string, urls []string) ([]string, error) {
	log.Print("FileStorage SaveBatch ", hashes)

	values := make([]string, 0, len(hashes))

	for i := range hashes {
		save, err := s.Save(hashes[i], urls[i])
		if err != nil {
			return values, err
		}
		values = append(values, save)
	}
	return values, nil
}

func (s *FileStorage) Get(hash string) (string, error) {
	log.Print("FileStorage Get", hash)
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	value, ok := s.data[hash]
	if ok {
		return value, nil
	}
	return value, exceptions.ErrURLNotFound
}

func (s *FileStorage) GetAll() (map[string]string, error) {
	log.Print("FileStorage GetAll")
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.data, nil
}

func (s *FileStorage) Close() {
}

func checkDirExistOrCreate(fileStoragePath string) {
	dir, _ := filepath.Split(fileStoragePath)
	if dir == "" {
		return
	}
	if _, err := os.Stat(fileStoragePath); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0o700)
		if err != nil {
			log.Println("error accessing file system:", err)
		}
	}
}

func readFromFileIntoMap(fileStoragePath string) map[string]string {
	consumer, err := NewConsumer(fileStoragePath)
	if err != nil {
		log.Println("error reading file from disk:", err)
	}
	defer consumer.Close()

	mapStore := make(map[string]string)
	for i := 0; ; i++ {
		record, err := consumer.ReadRecord()
		if err != nil {
			break
		}
		mapStore[record.Hash] = record.URL
	}
	return mapStore
}

func (s *FileStorage) writeToFile(hash string, url string) {
	record := Record{hash, url}
	err := s.producer.WriteRecord(&record)
	if err != nil {
		log.Println("error writing to file:", err)
	}
}

func (s *FileStorage) Delete(hash string) error {
	return exceptions.ErrNotImplemented
}

func (s *FileStorage) BatchDelete(hashes []string) error {
	return exceptions.ErrNotImplemented
}
