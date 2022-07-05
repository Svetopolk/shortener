package main

import (
	"log"
)

type Storage struct {
	mapStore map[string]string
}

func NewStorage() Storage {
	return Storage{mapStore: make(map[string]string)}
}

func (s *Storage) save(url string) string {
	log.Printf("storage: save url %s\n", url)
	hash := RandStringRunes(5)
	s.mapStore[hash] = url
	return hash
}

func (s *Storage) get(hash string) string {
	log.Printf("storage: get key %s\n", hash)
	value := s.mapStore[hash]

	return value
}
