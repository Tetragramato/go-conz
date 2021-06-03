package main

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/tetragramato/go-conz/internal"
	"log"
	"time"
)

func main() {
	log.Println("Start GO-CONZ...")
	internal.Parallelize(
		func() {
			for true {
				getAndPersistSensors()
				time.Sleep(internal.DelayInSecond * time.Second)
			}
		},
		internal.Serve,
	)
	log.Println("Close GO_CONZ")
}

func getAndPersistSensors() {
	client := resty.New()

	//////////////////////
	// Get Gateway specs
	gatewayResp, err := internal.GetGateway(client)
	if err != nil {
		log.Fatal(err)
	}

	//////////////////////
	// Get sensors from Gateway
	sensors, err := internal.GetSensors(client, gatewayResp, internal.ApiKey)
	if err != nil {
		log.Fatal(err)
	}

	/////////////////////
	// Parse sensors response
	var parsed map[string]interface{}
	err = json.Unmarshal(sensors.Body(), &parsed)
	if err != nil {
		log.Fatal(err)
	}

	/////////////////////
	sensorsByEtag, err := internal.MakeSensorsGrouped(parsed)
	if err != nil {
		log.Fatal(err)
	}

	/////////////////////
	// Send structured and grouped sensors to persistence
	err = internal.WriteCsv(internal.CsvPath, sensorsByEtag)
	if err != nil {
		log.Fatal(err)
	}
}
