package internal

import (
	"github.com/dgraph-io/badger/v3"
	"log"
)

type Database struct {
	instance *badger.DB
}

func NewDB() *Database {
	opts := badger.DefaultOptions(Config.DatabasePath)
	opts.Logger = nil

	db, err := badger.Open(opts)
	if err != nil {
		log.Fatalln(err)
	}
	return &Database{db}
}
