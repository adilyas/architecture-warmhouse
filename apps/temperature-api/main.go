package main

import (
	"encoding/json"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"time"
)

// TemperatureResponse represents the response from the temperature API
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
	http.HandleFunc("/temperature/", temperatureHandlerWithID)
	http.HandleFunc("/temperature", temperatureHandler)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	port := getEnv("PORT", ":8081")
	log.Printf("Temperature API starting on %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}

func temperatureHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	location := r.URL.Query().Get("location")

	var sensorID string
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

	temperature := rand.Float64() * 30.0

	status := "active"
	description := "Sensor is properly working"

	response := TemperatureResponse{
		Value:       temperature,
		Unit:        "Celsius",
		Timestamp:   time.Now(),
		Location:    location,
		Status:      status,
		SensorID:    sensorID,
		SensorType:  "temperature",
		Description: description,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Temperature request: location=%s, sensorID=%s, value=%.2f\n", location, sensorID, temperature)
}

func temperatureHandlerWithID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sensorID := r.URL.Path[len("/temperature/"):]
	if sensorID == "" {
		http.Error(w, "Sensor ID is required", http.StatusBadRequest)
		return
	}

	var location string
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

	temperature := rand.Float64() * 30.0

	status := "active"
	description := "Sensor is properly working"

	response := TemperatureResponse{
		Value:       temperature,
		Unit:        "Celsius",
		Timestamp:   time.Now(),
		Location:    location,
		Status:      status,
		SensorID:    sensorID,
		SensorType:  "temperature",
		Description: description,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Temperature request by ID: sensorID=%s, location=%s, value=%.2f\n", sensorID, location, temperature)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
