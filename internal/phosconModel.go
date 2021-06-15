package internal

import (
	"github.com/mitchellh/mapstructure"
)

const (
	PressureType    = "ZHAPressure"
	TemperatureType = "ZHATemperature"
	HumidityType    = "ZHAHumidity"
)

func (sensor *Sensor) IsSensor() bool {
	return sensor.Type == TemperatureType || sensor.Type == HumidityType || sensor.Type == PressureType
}

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

// PersistedSensor TODO make a mapper model to model
type PersistedSensor struct {
	Uniqueid    string
	Etag        string
	Name        string
	Lastupdated string
	Temperature int
	Humidity    int
	Pressure    int
}

func GetSensors(parsedSensors map[string]interface{}) ([]*Sensor, error) {
	var sensors []*Sensor
	for _, s := range parsedSensors {
		sensor := &Sensor{}
		err := mapstructure.Decode(s, sensor)
		if err != nil {
			return nil, err
		}
		sensors = append(sensors, sensor)
	}
	return sensors, nil
}
