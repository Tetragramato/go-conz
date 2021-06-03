package internal

import (
	"github.com/gocarina/gocsv"
	"log"
	"os"
)

func WriteCsv(csvFile string, sensorsByEtag map[string][]*Sensor) error {
	file, err := os.OpenFile(csvFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
			return
		}
	}(file)

	fileStat, err := file.Stat()
	if err != nil {
		return err
	}

	var csvSensors []*CsvSensor
	for key, value := range sensorsByEtag {
		for _, valSensor := range value {
			if valSensor.Type == TemperatureType || valSensor.Type == HumidityType || valSensor.Type == PressureType {
				csvSensors = append(
					csvSensors,
					&CsvSensor{
						key,
						valSensor.Name,
						valSensor.State.Lastupdated,
						valSensor.State.Temperature,
						valSensor.State.Humidity,
						valSensor.State.Pressure,
					},
				)
			}
		}
	}

	if fileStat.Size() == 0 {
		err = gocsv.MarshalFile(&csvSensors, file)
	} else {
		err = gocsv.MarshalWithoutHeaders(&csvSensors, file)
	}
	if err != nil {
		return err
	}
	return nil
}

func LoadModelFromCsv(csvFile string) ([]*CsvSensor, error) {
	file, err := os.Open(csvFile)
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
			return
		}
	}(file)

	var csvSensors []*CsvSensor
	err = gocsv.UnmarshalFile(file, &csvSensors)
	if err != nil {
		return nil, err
	}
	return csvSensors, nil
}
