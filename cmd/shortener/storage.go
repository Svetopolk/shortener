package main

import (
	"log"
)

type Storage struct {
}

func (*Storage) save(url string) string {
	log.Printf("storage: save url %s\n", url)
	return "dfdsavarevw"
}

func (*Storage) get(hash string) string {
	log.Printf("storage: get key %s\n", hash)
	return "https://example.com/"
}
