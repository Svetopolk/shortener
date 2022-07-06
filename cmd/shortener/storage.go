package main

import (
	"log"
)

type MemStorage struct {
	mapStore map[string]string
}

type Storage interface {
	save(url string) string
	get(hash string) string
}

func NewMemStorage() Storage {
	return &MemStorage{mapStore: make(map[string]string)}
}

func (s *MemStorage) save(url string) string {
	log.Printf("storage: save url %s\n", url)
	hash := RandStringRunes(5)
	s.mapStore[hash] = url
	return hash
}

func (s *MemStorage) get(hash string) string {
	log.Printf("storage: get key %s\n", hash)
	value := s.mapStore[hash]

	return value
}
