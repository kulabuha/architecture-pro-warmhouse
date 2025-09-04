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

type TelemetryPoint struct {
	DeviceID string    `json:"deviceId"`
	Metric   string    `json:"metric"`
	Value    float64   `json:"value"`
	TS       time.Time `json:"ts"`
}

type TemperatureResponse struct {
	Location    string  `json:"location"`
	SensorID    string  `json:"sensorId"`
	Temperature float64 `json:"temperature"`
	Unit        string  `json:"unit"`
	Timestamp   string  `json:"timestamp"`
}

var store = map[string][]TelemetryPoint{}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/healthz", healthHandler).Methods(http.MethodGet)
	r.HandleFunc("/ingest", ingestHandler).Methods(http.MethodPost)
	r.HandleFunc("/query", queryHandler).Methods(http.MethodGet)
	r.HandleFunc("/temperature", temperatureHandler).Methods(http.MethodGet)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("telemetry-service listening on :%s", port)
	log.Fatal(srv.ListenAndServe())
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("ok"))
}

func ingestHandler(w http.ResponseWriter, r *http.Request) {
	var p TelemetryPoint
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	store[p.DeviceID] = append(store[p.DeviceID], p)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"status": "ok", "count": len(store[p.DeviceID])})
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	deviceID := r.URL.Query().Get("deviceId")
	metric := r.URL.Query().Get("metric")
	if metric == "" {
		metric = "temperature"
	}
	points := store[deviceID]
	out := make([]TelemetryPoint, 0, len(points))
	for i := len(points) - 1; i >= 0 && len(out) < 100; i-- {
		if points[i].Metric == metric {
			out = append(out, points[i])
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"items": out, "count": len(out)})
}

func temperatureHandler(w http.ResponseWriter, r *http.Request) {
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
	temp := round(base+jitter, 1)
	now := time.Now().UTC()

	devID := "dev-" + sensorID
	store[devID] = append(store[devID], TelemetryPoint{
		DeviceID: devID, Metric: "temperature", Value: temp, TS: now,
	})

	resp := TemperatureResponse{
		Location:    location,
		SensorID:    sensorID,
		Temperature: temp,
		Unit:        "C",
		Timestamp:   now.Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func round(v float64, prec int) float64 {
	p := 1.0
	for i := 0; i < prec; i++ {
		p *= 10
	}
	return float64(int(v*p+0.5)) / p
}
