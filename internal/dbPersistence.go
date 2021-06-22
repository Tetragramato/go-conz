package internal

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"github.com/dgraph-io/badger/v3"
	"log"
	"time"
)

type database struct {
	instance   *badger.DB
	keysLookup map[string][]byte
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
	database := &database{
		instance:   db,
		keysLookup: make(map[string][]byte),
	}
	err = database.preloadKeys()
	if err != nil {
		log.Fatalln(err)
	}
	return database
}

func u64ToBytes(i uint64) (b []byte) {
	b = make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return
}

// GetConcreteKey TODO voir si on peu virer ca pour rester encapsulé
func GetConcreteKey(dbKey []byte) string {
	lenKey := len(dbKey)
	return string(dbKey[8:lenKey])
}

func (database *database) preloadKeys() error {
	err := database.instance.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			tmpKey := GetConcreteKey(it.Item().Key())
			database.keysLookup[tmpKey] = it.Item().Key()
		}
		return nil
	})
	return err
}

func (database *database) getPrefixedKey(key string) []byte {
	prefixedVal, ok := database.keysLookup[key]
	if !ok {
		timeKey := u64ToBytes(uint64(time.Now().UnixNano()))
		bytesKey := append(timeKey, []byte(key)...)
		database.keysLookup[key] = bytesKey
		return bytesKey
	}
	return prefixedVal
}

func (database *database) InsertOrUpdate(val interface{}, key string) error {
	verify, err := database.verifyByKey(key)
	if err != nil {
		return err
	}
	if !verify {
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
			return txn.Set(database.getPrefixedKey(key), dataAsBytes)
		})
	}
	return nil
}

func (database *database) Get(key string) ([]byte, error) {
	var valCopy []byte
	err := database.instance.View(func(txn *badger.Txn) error {
		item, err := txn.Get(database.getPrefixedKey(key))
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

func (database *database) verifyByKey(key string) (bool, error) {
	verifBool := false
	err := database.instance.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			if bytes.Equal(it.Item().Key(), database.getPrefixedKey(key)) {
				verifBool = true
				return nil
			}
		}
		return nil
	})
	return verifBool, err
}

// GetAll TODO moyen de faire mieux avec la réflexion pour éviter le badger.Item et plutôt passer par interface {}? (Type, Value....)
func (database *database) GetAll() ([]*badger.Item, error) {
	var items []*badger.Item
	err := database.instance.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			items = append(items, it.Item())
		}
		return nil
	})
	return items, err
}
