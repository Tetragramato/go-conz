package internal

import (
	"encoding/json"
	"log"
	"net/http"
)

func Serve() {
	http.HandleFunc("/sensors", returnSensors)
	log.Fatal(http.ListenAndServe(Config.HttpPort, nil))
}

func returnSensors(w http.ResponseWriter, r *http.Request) {
	log.Println("Handle sensors request")
	csvModel, err := LoadModelFromCsv(Config.CsvPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(csvModel)
	if err != nil {
		log.Fatal(err)
		return
	}
}
