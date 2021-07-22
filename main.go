package main

import (
	"github.com/tetragramato/go-conz/internal"
	"log"
	"time"
)

func init() {
	internal.InitConfig()
}

func main() {
	log.Println("Start GO-CONZ...")
	//Init db/repo/httpClient
	db := internal.NewDB()
	httpClient := internal.NewHttpClient()
	// Get Gateway specs
	gatewayResp, err := httpClient.GetGateway()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Getting and setting API Key...")
	apiKey, err := internal.NewApiKeyConfig(db).RegisterApiKey(gatewayResp)
	if err != nil {
		log.Fatal(err)
	}
	internal.Config.ApiKey = apiKey

	sensorRepo := internal.NewSensorRepository(db)
	internal.Parallelize(
		func() {
			if !internal.Config.ReadOnly {
				internal.RunNewPoller(
					time.Second*internal.Config.DelayInSecond,
					func() {
						listOfSensors, err := httpClient.GetAndParseSensors(gatewayResp)
						if err != nil {
							log.Printf("Error while retrieving sensors: %v", err.Error())
						} else {
							err = sensorRepo.SaveAll(listOfSensors)
							if err != nil {
								log.Printf("Error while saving sensors: %v", err.Error())
							}
						}
					})
			}
		},
		func() {
			log.Fatal(internal.Serve(sensorRepo))
		},
	)
	log.Println("Close GO-CONZ")
}
