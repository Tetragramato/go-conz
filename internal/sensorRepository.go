package internal

import (
	"bytes"
	"encoding/gob"
	"github.com/dgraph-io/badger/v3"
)

type sensorRepository struct {
	operation Operable
}

func NewSensorRepository(operation Operable) SensorRepository {
	return &sensorRepository{operation}
}

type SensorRepository interface {
	GetAll() ([]*InputSensors, error)
	SaveAll([]*InputSensors) error
	Save(sensor *InputSensors) error
}

func (repo *sensorRepository) GetAll() ([]*InputSensors, error) {
	items, err := repo.operation.GetAll()
	if err != nil {
		return nil, err
	}
	var listOfSensors []*InputSensors
	for _, item := range items {
		if string(item.Key()) != DbApiKey {
			val, err := getSensorsList(item)
			if err != nil {
				return nil, err
			}
			listOfSensors = append(listOfSensors, val)
		}
	}
	return listOfSensors, nil
}

func (repo *sensorRepository) Save(sensorsList *InputSensors) error {
	return repo.operation.InsertOrUpdate(sensorsList, sensorsList.Etag)
}

func (repo *sensorRepository) SaveAll(listOfSensors []*InputSensors) error {
	for _, value := range listOfSensors {
		err := repo.Save(value)
		if err != nil {
			return err
		}
	}
	return nil
}

//TODO peu peut être mieux faire pour eviter l'import de badger
func getSensorsList(item *badger.Item) (*InputSensors, error) {
	var sensorsList InputSensors
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
