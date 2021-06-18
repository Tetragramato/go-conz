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
	apiKey := internal.NewApiKeyConfig(db, httpClient).RegisterApiKey(gatewayResp)
	internal.Config.ApiKey = apiKey

	sensorRepo := internal.NewSensorRepository(db)
	internal.Parallelize(
		func() {
			for {
				listOfSensors, err := httpClient.GetAndParseSensors(gatewayResp)
				if err != nil {
					log.Fatal(err)
				}
				err = sensorRepo.SaveAll(listOfSensors)
				if err != nil {
					log.Fatal(err)
				}
				time.Sleep(internal.Config.DelayInSecond * time.Second)
			}
		},
		func() {
			internal.Serve(sensorRepo)
		},
	)
	log.Println("Close GO-CONZ")
}
