package internal

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Serve() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/sensors", returnSensors)

	log.Fatal(http.ListenAndServe(HttpPort, myRouter))
}

func returnSensors(w http.ResponseWriter, r *http.Request) {
	log.Println("Handle sensors request")
	csvModel, err := LoadModelFromCsv(CsvPath)
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
