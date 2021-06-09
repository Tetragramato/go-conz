package internal

import (
	"github.com/gocarina/gocsv"
	"os"
)

func writeCsv(csvFile string, sensorsByEtag map[string][]*Sensor) (err error) {
	file, err := os.OpenFile(csvFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}

	defer func(file *os.File) {
		cerr := file.Close()
		if cerr != nil {
			err = cerr
		}
	}(file)

	fileStat, err := file.Stat()
	if err != nil {
		return
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
		return
	}
	return
}

func LoadModelFromCsv(csvFile string) (csvSensors []*CsvSensor, err error) {
	file, err := os.Open(csvFile)
	if err != nil {
		return
	}

	defer func(file *os.File) {
		cerr := file.Close()
		if cerr != nil {
			err = cerr
		}
	}(file)

	err = gocsv.UnmarshalFile(file, csvSensors)
	if err != nil {
		return
	}
	return
}

func PersistSensors(sensorsByEtag map[string][]*Sensor) error {
	err := writeCsv(Config.CsvPath, sensorsByEtag)
	if err != nil {
		return err
	}
	return nil
}
