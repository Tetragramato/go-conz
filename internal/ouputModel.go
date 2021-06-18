package internal

type OutputSensors struct {
	Uniqueid string         `json:"uniqueId"`
	Name     string         `json:"name"`
	Type     string         `json:"type"`
	Events   []*SensorEvent `json:"events"`
}

type pairIdName struct {
	uniqueid   string
	name       string
	sensorType string
}

type SensorEvent struct {
	Etag        string `json:"etag"`
	Lastupdated string `json:"lastUpdated"`
	Temperature int    `json:"temperature"`
	Humidity    int    `json:"humidity"`
	Pressure    int    `json:"pressure"`
}

func GetOutputSensors(listSensorsList []*InputSensors) []*OutputSensors {
	var flatMapSensors []*InputSensor
	for _, sensorList := range listSensorsList {
		flatMapSensors = append(flatMapSensors, sensorList.Sensors...)
	}

	mapSensor := make(map[pairIdName][]*SensorEvent)
	for _, sensor := range flatMapSensors {
		pairIdName := pairIdName{
			uniqueid:   sensor.Uniqueid,
			name:       sensor.Name,
			sensorType: sensor.Type,
		}
		mapSensor[pairIdName] = append(mapSensor[pairIdName], sensor.toSensorEvent())
	}

	var outputSensors []*OutputSensors
	for key, sensors := range mapSensor {
		outputSensors = append(outputSensors, &OutputSensors{
			Uniqueid: key.uniqueid,
			Name:     key.name,
			Type:     key.sensorType,
			Events:   sensors,
		})
	}
	return outputSensors
}
