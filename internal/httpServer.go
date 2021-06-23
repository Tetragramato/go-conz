package internal

import (
	"encoding/json"
	"log"
	"net/http"
)

func Serve(repo SensorRepository) error {
	http.HandleFunc(
		"/sensors",
		func(w http.ResponseWriter, r *http.Request) {
			handleSensors(repo, w, r)
		},
	)
	return http.ListenAndServe(Config.HttpPort, nil)
}

func handleSensors(repo SensorRepository, w http.ResponseWriter, r *http.Request) {
	log.Println("Handle sensors request")
	listOfSensors, err := repo.GetAll()
	if err != nil {
		log.Fatal(err)
		return
	}
	outputSensors := GetOutputSensors(listOfSensors)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	err = json.NewEncoder(w).Encode(&outputSensors)
	if err != nil {
		log.Fatal(err)
		return
	}
}
