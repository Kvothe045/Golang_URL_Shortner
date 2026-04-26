// url_shortner/store/store.go
package store

import (
	"sync"
)

type DB struct {
	EncodeDB map[string]string
	DecodeDB map[string]string
	mu       sync.RWMutex
}

func NewDB() *DB {
	return &DB{
		EncodeDB: make(map[string]string),
		DecodeDB: make(map[string]string),
	}
}

func (db *DB) Save(url string, code string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.EncodeDB[url] = code
	db.DecodeDB[code] = url
}

func (db *DB) GetEncodedURL(url string) (string, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	val, exists := db.EncodeDB[url]

	return val, exists
}

func (db *DB) GetOriginalURL(code string) (string, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	val, exists := db.DecodeDB[code]
	return val, exists
}
