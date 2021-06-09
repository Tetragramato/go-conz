package main

import (
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

	httpClient := &internal.HttpClient{Client: resty.New()}
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
				err = internal.PersistSensors(sensors)
				if err != nil {
					log.Fatal(err)
				}
				time.Sleep(internal.Config.DelayInSecond * time.Second)
			}
		},
		internal.Serve,
	)
	log.Println("Close GO-CONZ")
}
