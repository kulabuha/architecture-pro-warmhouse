package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Device struct {
	ID       string `json:"id"`
	Location string `json:"location"`
	SensorID string `json:"sensorId"`
	Type     string `json:"type"`
	Status   string `json:"status"`
}

var (
	devices = map[string]*Device{}
	seq     = 1
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/healthz", healthHandler).Methods(http.MethodGet)
	r.HandleFunc("/devices", devicesHandler).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/map", mapHandler).Methods(http.MethodGet)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("devices-service listening on :%s", port)
	log.Fatal(srv.ListenAndServe())
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("ok"))
}

func devicesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		list := make([]*Device, 0, len(devices))
		for _, d := range devices {
			list = append(list, d)
		}
		json.NewEncoder(w).Encode(map[string]any{"items": list})
	case http.MethodPost:
		var body struct{ Location string }
		_ = json.NewDecoder(r.Body).Decode(&body)
		id := randID()
		sid := nextSensorID()
		loc := strings.TrimSpace(body.Location)
		if loc == "" {
			loc = "Unknown"
		}
		d := &Device{ID: id, Location: loc, SensorID: sid, Type: "Sensor", Status: "Active"}
		devices[id] = d
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(d)
	}
}

func mapHandler(w http.ResponseWriter, r *http.Request) {
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
	json.NewEncoder(w).Encode(map[string]string{
		"location": location,
		"sensorId": sensorID,
	})
}

func nextSensorID() string { seq++; return string('0' + (seq % 10)) }
func randID() string       { return time.Now().Format("20060102150405") }
