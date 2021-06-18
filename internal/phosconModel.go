package internal

import (
	"github.com/mitchellh/mapstructure"
)

const (
	PressureType    = "ZHAPressure"
	TemperatureType = "ZHATemperature"
	HumidityType    = "ZHAHumidity"
)

type APIKey struct {
	Success struct {
		Username string `json:"username"`
	} `json:"success"`
}

func GetApiKey(parsedKey []interface{}) (*APIKey, error) {
	apiKey := &APIKey{}
	for _, s := range parsedKey {
		err := mapstructure.Decode(s, apiKey)
		if err != nil {
			return nil, err
		}
	}
	return apiKey, nil
}

type Gateway []struct {
	Internalipaddress string `json:"internalipaddress"`
	Internalport      int    `json:"internalport"`
}

type PhosconSensor struct {
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

func (sensor *PhosconSensor) ToInputSensor() *InputSensor {
	return &InputSensor{
		Uniqueid:    sensor.Uniqueid,
		Etag:        sensor.Etag,
		Name:        sensor.Name,
		Lastupdated: sensor.State.Lastupdated,
		Type:        sensor.Type,
		Temperature: sensor.State.Temperature,
		Humidity:    sensor.State.Humidity,
		Pressure:    sensor.State.Pressure,
	}
}

func (sensor *PhosconSensor) IsSensor() bool {
	return sensor.Type == TemperatureType || sensor.Type == HumidityType || sensor.Type == PressureType
}
