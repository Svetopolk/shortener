package storage

import (
	"log"
)

type FileStorage struct {
	mapStore        map[string]string
	fileStoragePath string
	producer        *producer
}

var _ Storage = &FileStorage{}

func NewFileStorage(fileStoragePath string) *FileStorage {
	mapStore := make(map[string]string)
	readFromFileIntoMap(fileStoragePath, mapStore)

	fileProducer, err := NewProducer(fileStoragePath)
	if err != nil {
		log.Fatal(err)
	}
	return &FileStorage{mapStore, fileStoragePath, fileProducer}
}

func readFromFileIntoMap(fileStoragePath string, mapStore map[string]string) {
	consumer, err := NewConsumer(fileStoragePath)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; ; i++ {
		record, err := consumer.ReadRecord()
		if err != nil {
			break
		}
		mapStore[record.Hash] = record.URL
	}
	consumer.Close()
}

func (s *FileStorage) Save(hash string, url string) string {
	s.mapStore[hash] = url
	s.writeToFile(hash, url)
	return hash
}

func (s *FileStorage) Get(hash string) (string, bool) {
	value, ok := s.mapStore[hash]
	return value, ok
}

func (s *FileStorage) writeToFile(hash string, url string) {
	record := Record{hash, url}
	err := s.producer.WriteRecord(&record)
	if err != nil {
		log.Fatal(err)
	}
}
