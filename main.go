package main

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/tetragramato/go-conz/internal"
	"log"
	"time"
)

func init() {
	internal.InitConfig()
}

func main() {
	log.Println("Start GO-CONZ...")

	client := resty.New()
	// Get Gateway specs
	gatewayResp, err := internal.GetGateway(client)
	if err != nil {
		log.Fatal(err)
	}
	internal.Parallelize(
		func() {
			for true {
				sensors, err := getSensors(client, gatewayResp)
				if err != nil {
					log.Fatal(err)
				}
				err = persistSensors(sensors)
				if err != nil {
					log.Fatal(err)
				}
				time.Sleep(internal.Config.DelayInSecond * time.Second)
			}
		},
		internal.Serve,
	)
	log.Println("Close GO_CONZ")
}

func getSensors(client *resty.Client, gatewayResp *internal.Gateway) (map[string]interface{}, error) {
	// Get sensors from Gateway
	sensors, err := internal.GetSensors(client, gatewayResp, internal.Config.ApiKey)
	if err != nil {
		return nil, err
	}

	var parsed map[string]interface{}
	err = json.Unmarshal(sensors.Body(), &parsed)
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

func persistSensors(sensors map[string]interface{}) error {
	sensorsByEtag, err := internal.GetSensorsByEtag(sensors)
	if err != nil {
		return err
	}

	// Send structured and grouped sensors to persistence
	err = internal.WriteCsv(internal.Config.CsvPath, sensorsByEtag)
	if err != nil {
		return err
	}
	return nil
}
