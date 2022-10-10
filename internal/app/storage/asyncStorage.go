package storage

import (
	"log"
	"sync"
	"time"

	"github.com/Svetopolk/shortener/internal/app/exceptions"
)

var _ Storage = &AsyncStorage{}

type AsyncStorage struct {
	storage     Storage
	deleteQueue []string
	batchSize   int
	mtx         sync.RWMutex
}

func NewAsyncStorage(st Storage) *AsyncStorage {
	asyncStorage := &AsyncStorage{
		storage:     st,
		deleteQueue: make([]string, 0, 100),
		batchSize:   10,
	}
	asyncStorage.startPeriodicDelete(1 * time.Second)
	return asyncStorage
}

func (a *AsyncStorage) Save(hash string, url string) (string, error) {
	return "", exceptions.ErrNotImplemented
}

func (a *AsyncStorage) SaveBatch(hashes []string, urls []string) ([]string, error) {
	return nil, exceptions.ErrNotImplemented
}

func (a *AsyncStorage) Get(hash string) (string, error) {
	return "", exceptions.ErrNotImplemented
}

func (a *AsyncStorage) GetAll() (map[string]string, error) {
	return nil, exceptions.ErrNotImplemented
}

func (a *AsyncStorage) Delete(hash string) error {
	return exceptions.ErrNotImplemented
}

func (a *AsyncStorage) BatchDelete(hashes []string) error {
	log.Println("BatchDelete() hashes=", hashes)
	a.mtx.RLock()
	defer a.mtx.RUnlock()
	a.deleteQueue = append(a.deleteQueue, hashes...)
	return nil
}

func (a *AsyncStorage) Shutdown() {
	a.mtx.RLock()
	defer a.mtx.RUnlock()
	err := a.storage.BatchDelete(a.deleteQueue)
	if err != nil {
		log.Println("error when shutdown deleteQueue flash")
	}
}

func (a *AsyncStorage) BatchDeleteFromQueue() {
	a.mtx.RLock()
	defer a.mtx.RUnlock()
	if len(a.deleteQueue) == 0 {
		return
	}
	size := min(len(a.deleteQueue), a.batchSize)
	batch := a.deleteQueue[:size]
	a.deleteQueue = a.deleteQueue[size:]
	log.Println("BatchDeleteFromQueue(); batch size=", size, "batch=", batch)
	err := a.storage.BatchDelete(batch)
	if err != nil {
		log.Println("error when BatchDeleteFromQueue", err)
	}
}

func (a *AsyncStorage) startPeriodicDelete(period time.Duration) {
	ticker := time.NewTicker(period)
	go func() {
		for {
			<-ticker.C
			a.BatchDeleteFromQueue()
		}
	}()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
