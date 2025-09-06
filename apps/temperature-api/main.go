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
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	Timestamp   time.Time `json:"timestamp"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
	SensorID    string    `json:"sensor_id"`
	SensorType  string    `json:"sensor_type"`
	Description string    `json:"description"`
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","endpoints":["/temperature?location=Living Room","/temperature?sensorId=1","/temperature/1"]}`))
	}).Methods(http.MethodGet)

	router.HandleFunc("/temperature", temperatureHandler).Methods(http.MethodGet)

	router.HandleFunc("/temperature/{sensorId}", temperatureByIDHandler).Methods(http.MethodGet)

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

	temp := genTemp(sensorID)

	resp := TemperatureResponse{
		Value:       round(temp, 1),
		Unit:        "C",
		Timestamp:   time.Now().UTC(),
		Location:    location,
		Status:      "ok",
		SensorID:    sensorID,
		SensorType:  "temperature",
		Description: "Randomly generated temperature",
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func temperatureByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	sensorID := vars["sensorId"]

	location := "Unknown"
	switch sensorID {
	case "1":
		location = "Living Room"
	case "2":
		location = "Bedroom"
	case "3":
		location = "Kitchen"
	default:
		sensorID = "0"
	}

	temp := genTemp(sensorID)

	resp := TemperatureResponse{
		Value:       round(temp, 1),
		Unit:        "C",
		Timestamp:   time.Now().UTC(),
		Location:    location,
		Status:      "ok",
		SensorID:    sensorID,
		SensorType:  "temperature",
		Description: "Randomly generated temperature",
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func genTemp(sensorID string) float64 {
	base := 18.0 + rand.Float64()*10.0
	sid, _ := strconv.Atoi(sensorID)
	jitter := float64((sid%5)-2) * 0.1
	return base + jitter
}

func round(v float64, prec int) float64 {
	pow := 1.0
	for i := 0; i < prec; i++ {
		pow *= 10
	}
	return float64(int(v*pow+0.5)) / pow
}
