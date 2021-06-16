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

func (sensor *PhosconSensor) ToSensor() *Sensor {
	return &Sensor{
		Uniqueid:    sensor.Uniqueid,
		Etag:        sensor.Etag,
		Name:        sensor.Name,
		Lastupdated: sensor.State.Lastupdated,
		Temperature: sensor.State.Temperature,
		Humidity:    sensor.State.Humidity,
		Pressure:    sensor.State.Pressure,
	}
}

func (sensor *PhosconSensor) IsSensor() bool {
	return sensor.Type == TemperatureType || sensor.Type == HumidityType || sensor.Type == PressureType
}

type Sensor struct {
	Uniqueid    string `json:"uniqueId"`
	Etag        string `json:"etag"`
	Name        string `json:"name"`
	Lastupdated string `json:"lastUpdated"`
	Temperature int    `json:"temperature"`
	Humidity    int    `json:"humidity"`
	Pressure    int    `json:"pressure"`
}

type SensorsList struct {
	Etag    string    `json:"etag"`
	Sensors []*Sensor `json:"sensors"`
}

func GetListOfSensorsList(parsedSensors map[string]interface{}) ([]*SensorsList, error) {

	sensorsByEtag := make(map[string][]*PhosconSensor, len(parsedSensors))
	var etags []string
	for _, s := range parsedSensors {
		phosconSensor := &PhosconSensor{}
		err := mapstructure.Decode(s, phosconSensor)
		if err != nil {
			return nil, err
		}
		if phosconSensor.IsSensor() {
			if !Contains(etags, phosconSensor.Etag) {
				etags = append(etags, phosconSensor.Etag)
			}
			sensorsByEtag[phosconSensor.Etag] = append(sensorsByEtag[phosconSensor.Etag], phosconSensor)
		}
	}

	var listOfSensorsList []*SensorsList
	for _, etag := range etags {
		sensorsList := &SensorsList{Etag: etag}
		var sensor []*Sensor
		for _, sensorVal := range sensorsByEtag[etag] {
			sensor = append(sensor, sensorVal.ToSensor())
		}
		sensorsList.Sensors = sensor
		listOfSensorsList = append(listOfSensorsList, sensorsList)
	}

	return listOfSensorsList, nil
}
