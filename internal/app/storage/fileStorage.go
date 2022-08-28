package storage

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/Svetopolk/shortener/internal/logging"
)

type FileStorage struct {
	data            map[string]string
	fileStoragePath string
	producer        *producer
	mtx             sync.RWMutex
}

var _ Storage = &FileStorage{}

func NewFileStorage(fileStoragePath string) *FileStorage {
	logging.Enter()
	defer logging.Exit()

	checkDirExistOrCreate(fileStoragePath)
	mapStore := readFromFileIntoMap(fileStoragePath)

	fileProducer, err := NewProducer(fileStoragePath)
	if err != nil {
		log.Println("can not create NewFileStorage", err)
	}
	return &FileStorage{data: mapStore, fileStoragePath: fileStoragePath, producer: fileProducer}
}

func (s *FileStorage) Save(hash string, url string) (string, error) {
	logging.Enter()
	defer logging.Exit()

	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.data[hash] = url
	s.writeToFile(hash, url)
	return hash, nil
}

func (s *FileStorage) SaveBatch(hashes []string, urls []string) ([]string, error) {
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

func (s *FileStorage) Get(hash string) (string, bool) {
	logging.Enter()
	defer logging.Exit()

	s.mtx.RLock()
	defer s.mtx.RUnlock()
	value, ok := s.data[hash]
	return value, ok
}

func (s *FileStorage) GetAll() map[string]string {
	logging.Enter()
	defer logging.Exit()

	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.data
}

func checkDirExistOrCreate(fileStoragePath string) {
	logging.Enter()
	defer logging.Exit()

	dir, _ := filepath.Split(fileStoragePath)
	if dir == "" {
		return
	}
	if _, err := os.Stat(fileStoragePath); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0700)
		if err != nil {
			log.Println("error accessing file system:", err)
		}
	}
}

func readFromFileIntoMap(fileStoragePath string) map[string]string {
	logging.Enter()
	defer logging.Exit()

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
	logging.Enter()
	defer logging.Exit()

	record := Record{hash, url}
	err := s.producer.WriteRecord(&record)
	if err != nil {
		log.Println("error writing to file:", err)
	}
}
