package internal

import (
	"encoding/json"
	"log"
	"net/http"
)

func Serve(repo SensorRepository) {
	http.HandleFunc("/sensors", func(w http.ResponseWriter, r *http.Request) {
		returnSensors(repo, w, r)
	})
	log.Fatal(http.ListenAndServe(Config.HttpPort, nil))
}

func returnSensors(repo SensorRepository, w http.ResponseWriter, r *http.Request) {
	log.Println("Handle sensors request")
	sensors, err := repo.GetAll()
	if err != nil {
		log.Fatal(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	err = json.NewEncoder(w).Encode(&sensors)
	if err != nil {
		log.Fatal(err)
		return
	}
}
