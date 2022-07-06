package main

import (
	"log"
	"net/http"
)

func main() {
	m := RequestHandler{NewMemStorage()}
	http.Handle("/", &m)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
