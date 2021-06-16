package internal

import (
	"bytes"
	"encoding/gob"
	"github.com/dgraph-io/badger/v3"
)

type sensorRepository struct {
	db *Database
}

func NewSensorRepository(db *Database) SensorRepository {
	return &sensorRepository{db}
}

type SensorRepository interface {
	GetAll() ([]*SensorsList, error)
	SaveAll([]*SensorsList) error
	Save(sensor *SensorsList) error
}

func (repo *sensorRepository) GetAll() ([]*SensorsList, error) {
	var sensorsList []*SensorsList

	err := repo.db.instance.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			sensor, err := getSensorsList(it.Item())
			if err != nil {
				return err
			}
			sensorsList = append(sensorsList, sensor)
		}
		return nil
	})
	return sensorsList, err
}

func (repo *sensorRepository) Save(sensorsList *SensorsList) error {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(sensorsList)
	if err != nil {
		return err
	}

	return repo.db.instance.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(sensorsList.Etag), buffer.Bytes())
	})
}

func (repo *sensorRepository) SaveAll(listOfSensors []*SensorsList) error {
	for _, value := range listOfSensors {
		err := repo.Save(value)
		if err != nil {
			return err
		}
	}
	return nil
}

func getSensorsList(item *badger.Item) (*SensorsList, error) {
	var sensorsList SensorsList
	var buffer bytes.Buffer

	err := item.Value(func(val []byte) error {
		_, err := buffer.Write(val)
		return err
	})
	if err != nil {
		return nil, err
	}

	decoder := gob.NewDecoder(&buffer)
	err = decoder.Decode(&sensorsList)

	return &sensorsList, err
}
