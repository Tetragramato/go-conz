package internal

import (
	"bytes"
	"encoding/gob"
	"github.com/dgraph-io/badger/v3"
	"log"
)

type database struct {
	instance *badger.DB
}

type Operable interface {
	InsertOrUpdate(value interface{}, key string) error
	GetAll() ([]*badger.Item, error)
	Get(key string) ([]byte, error)
}

func NewDB() Operable {
	opts := badger.DefaultOptions(Config.DatabasePath)
	opts.Logger = nil

	db, err := badger.Open(opts)
	if err != nil {
		log.Fatalln(err)
	}
	return &database{db}
}

func (database database) InsertOrUpdate(val interface{}, key string) error {
	var dataAsBytes []byte
	switch v := val.(type) {
	case string:
		dataAsBytes = []byte(v)
	default:
		var buffer bytes.Buffer
		encoder := gob.NewEncoder(&buffer)
		err := encoder.Encode(v)
		if err != nil {
			return err
		}

		dataAsBytes = buffer.Bytes()
	}

	return database.instance.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), dataAsBytes)
	})
}

func (database database) Get(key string) ([]byte, error) {
	var valCopy []byte
	err := database.instance.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		valCopy, err = item.ValueCopy(valCopy)
		if err != nil {
			return err
		}
		return nil
	})
	return valCopy, err
}

// GetAll TODO moyen de faire mieux avec la réflexion pour éviter le badger.Item et plutôt passer par interface {}? (Type, Value....)
func (database database) GetAll() ([]*badger.Item, error) {
	var items []*badger.Item
	err := database.instance.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			items = append(items, it.Item())
		}
		return nil
	})
	return items, err
}
