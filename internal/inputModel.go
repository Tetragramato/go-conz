package internal

import "github.com/mitchellh/mapstructure"

type InputSensor struct {
	Uniqueid    string
	Etag        string
	Name        string
	Lastupdated string
	Type        string
	Temperature int
	Humidity    int
	Pressure    int
}

func (inputSensor *InputSensor) toSensorEvent() *SensorEvent {
	return &SensorEvent{
		Etag:        inputSensor.Etag,
		Lastupdated: inputSensor.Lastupdated,
		Temperature: inputSensor.Temperature,
		Humidity:    inputSensor.Humidity,
		Pressure:    inputSensor.Pressure,
	}
}

type InputSensors struct {
	Etag    string
	Sensors []*InputSensor
}

// GetInputSensors TODO peut être simplifié ?
func GetInputSensors(parsedSensors map[string]interface{}) ([]*InputSensors, error) {
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

	var listOfSensors []*InputSensors
	for _, etag := range etags {
		sensorsList := &InputSensors{Etag: etag}
		var sensor []*InputSensor
		for _, sensorVal := range sensorsByEtag[etag] {
			sensor = append(sensor, sensorVal.ToInputSensor())
		}
		sensorsList.Sensors = sensor
		listOfSensors = append(listOfSensors, sensorsList)
	}
	return listOfSensors, nil
}
