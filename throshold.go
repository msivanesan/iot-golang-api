package main

import (
	"encoding/json"
	"net/http"
	"os"
)

type Thresholds struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Sound       float64 `json:"sound"`
}

const thresholdFile = "static/thresholds.json"

func handleGetThresholds(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	file, err := os.Open(thresholdFile)
	if err != nil {
		// Default thresholds if file doesn't exist
		json.NewEncoder(w).Encode(Thresholds{
			Temperature: 30.0,
			Humidity:    70.0,
			Sound:       80.0,
		})
		return
	}
	defer file.Close()

	var t Thresholds
	if err := json.NewDecoder(file).Decode(&t); err != nil {
		http.Error(w, "Error reading thresholds", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(t)
}

func handleUpdateThresholds(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var t Thresholds
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	file, err := os.Create(thresholdFile)
	if err != nil {
		http.Error(w, "Error saving thresholds", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	json.NewEncoder(file).Encode(t)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
