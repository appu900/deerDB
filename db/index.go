package db

import (
	"log"

	"github.com/cockroachdb/pebble"
)

type Document struct {
	ID   string      `json:"id"`
	Data interface{} `json:"data"`
}

func Testdatabase() {
	db, err := pebble.Open("testDB", &pebble.Options{})
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	key := []byte("id1")
	if err := db.Set(key, []byte("value1"), pebble.Sync); err != nil {
		log.Fatalf("failed to set key: %v", err)
	}

	value, closer, err := db.Get(key)
	if err != nil {
		log.Fatalf("failed to get key: %v", err)
	}
	log.Printf("key: %s, value: %s", key, value)
	if err := closer.Close(); err != nil {
		log.Fatalf("failed to close closer: %v", err)
	}
}
