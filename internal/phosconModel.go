package internal

import (
	"github.com/mitchellh/mapstructure"
)

const (
	PressureType    = "ZHAPressure"
	TemperatureType = "ZHATemperature"
	HumidityType    = "ZHAHumidity"
)

type APIKey []struct {
	Success struct {
		Username string `json:"username"`
	} `json:"success"`
}

type Gateway []struct {
	Internalipaddress string `json:"internalipaddress"`
	Internalport      int    `json:"internalport"`
}

type Sensor struct {
	Config struct {
		On        bool
		Battery   int
		Reachable bool
	}
	Ep               int
	Etag             string
	Lastseen         string
	Manufacturername string
	Modelid          string
	Name             string
	State            struct {
		Lastupdated string
		Temperature int
		Humidity    int
		Pressure    int
	}
	Swversion string
	Type      string
	Uniqueid  string
}

type CsvSensor struct {
	Uniqueid    string
	Etag        string
	Name        string
	Lastupdated string
	Temperature int
	Humidity    int
	Pressure    int
}

// GetSensorsByEtag Group sensors by Etag and convert to Sensor struct
func GetSensorsByEtag(parsedSensors map[string]interface{}) (map[string][]*Sensor, error) {
	sensorsByEtag := make(map[string][]*Sensor, len(parsedSensors))
	for _, s := range parsedSensors {
		sensor := &Sensor{}
		err := mapstructure.Decode(s, sensor)
		if err != nil {
			return nil, err
		}
		sensorsByEtag[sensor.Etag] = append(sensorsByEtag[sensor.Etag], sensor)
	}
	return sensorsByEtag, nil
}
