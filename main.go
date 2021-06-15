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

	httpClient := internal.NewHttpClient()
	db := internal.NewDB()
	sensorRepo := internal.NewSensorRepository(db)
	// Get Gateway specs
	gatewayResp, err := httpClient.GetGateway()
	if err != nil {
		log.Fatal(err)
	}
	internal.Parallelize(
		func() {
			for true {
				sensors, err := httpClient.GetAndParseSensors(gatewayResp)
				if err != nil {
					log.Fatal(err)
				}
				err = sensorRepo.SaveAll(sensors)
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
