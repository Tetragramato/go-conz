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
	GetAll() ([]*SensorsList, error)
	SaveAll([]*SensorsList) error
	Save(sensor *SensorsList) error
}

func (repo *sensorRepository) GetAll() ([]*SensorsList, error) {
	items, err := repo.operation.GetAll()
	if err != nil {
		return nil, err
	}
	var listOfSensorsList []*SensorsList
	for _, item := range items {
		if string(item.Key()) != DbApiKey {
			val, err := getSensorsList(item)
			if err != nil {
				return nil, err
			}
			listOfSensorsList = append(listOfSensorsList, val)
		}
	}
	return listOfSensorsList, nil
}

func (repo *sensorRepository) Save(sensorsList *SensorsList) error {
	return repo.operation.InsertOrUpdate(sensorsList, sensorsList.Etag)
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

//TODO peu peut être mieux faire pour eviter l'import de badger
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
