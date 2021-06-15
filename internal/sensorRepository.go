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
	GetAll() ([]*PersistedSensor, error)
	SaveAll([]*Sensor) error
	Save(sensor *PersistedSensor) error
}

func (repo *sensorRepository) GetAll() ([]*PersistedSensor, error) {
	var sensors []*PersistedSensor

	err := repo.db.instance.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			sensor, err := getSensor(it.Item())
			if err != nil {
				return err
			}
			sensors = append(sensors, sensor)
		}
		return nil
	})
	return sensors, err
}

func (repo *sensorRepository) Save(sensor *PersistedSensor) error {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(sensor)
	if err != nil {
		return err
	}

	return repo.db.instance.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(sensor.Etag), buffer.Bytes())
	})
}

func (repo *sensorRepository) SaveAll(sensors []*Sensor) error {
	for _, value := range sensors {
		if value.IsSensor() {
			val := &PersistedSensor{
				Uniqueid:    value.Uniqueid,
				Etag:        value.Etag,
				Name:        value.Name,
				Lastupdated: value.State.Lastupdated,
				Temperature: value.State.Temperature,
				Humidity:    value.State.Humidity,
				Pressure:    value.State.Pressure,
			}
			err := repo.Save(val)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getSensor(item *badger.Item) (*PersistedSensor, error) {
	var sensor PersistedSensor
	var buffer bytes.Buffer

	err := item.Value(func(val []byte) error {
		_, err := buffer.Write(val)
		return err
	})

	decoder := gob.NewDecoder(&buffer)
	err = decoder.Decode(&sensor)

	return &sensor, err
}
