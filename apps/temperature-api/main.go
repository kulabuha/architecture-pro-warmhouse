package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type TemperatureResponse struct {
	Location    string  `json:"location"`
	SensorID    string  `json:"sensorId"`
	Temperature float64 `json:"temperature"`
	Unit        string  `json:"unit"`
	Timestamp   string  `json:"timestamp"`
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/temperature", temperatureHandler).Methods("GET")

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","endpoints":["/temperature?location=Living Room","/temperature?sensorId=1"]}`))
	}).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("temperature-api listening on :%s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}

func temperatureHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	location := r.URL.Query().Get("location")
	sensorID := r.URL.Query().Get("sensorId")

	if location == "" {
		switch sensorID {
		case "1":
			location = "Living Room"
		case "2":
			location = "Bedroom"
		case "3":
			location = "Kitchen"
		default:
			location = "Unknown"
		}
	}

	if sensorID == "" {
		switch location {
		case "Living Room":
			sensorID = "1"
		case "Bedroom":
			sensorID = "2"
		case "Kitchen":
			sensorID = "3"
		default:
			sensorID = "0"
		}
	}

	base := 18.0 + rand.Float64()*10.0
	sid, _ := strconv.Atoi(sensorID)
	jitter := float64((sid%5)-2) * 0.1
	temp := base + jitter

	resp := TemperatureResponse{
		Location:    location,
		SensorID:    sensorID,
		Temperature: round(temp, 1),
		Unit:        "C",
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(resp)
}

func round(v float64, prec int) float64 {
	pow := 1.0
	for i := 0; i < prec; i++ {
		pow *= 10
	}
	return float64(int(v*pow+0.5)) / pow
}
