package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const API_KEY = "MY_SECRET_KEY" // Change this to your own secure key

type SensorData struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	SoundLevel  float64 `json:"sound_level"`
}

func handleSensorData(w http.ResponseWriter, r *http.Request) {
	// API key authentication
	clientKey := r.Header.Get("X-API-Key")
	if clientKey != API_KEY {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	var data SensorData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// File paths
	tempHumFile := "static/temp_humidity.csv"
	soundFile := "static/sound_level.csv"

	// ----- Store Temperature + Humidity -----
	storeCSV(tempHumFile, []string{"timestamp", "temperature", "humidity"},
		[]string{
			time.Now().Format(time.RFC3339),
			fmt.Sprintf("%.2f", data.Temperature),
			fmt.Sprintf("%.2f", data.Humidity),
		})

	// ----- Store Sound Level -----
	storeCSV(soundFile, []string{"timestamp", "sound_level"},
		[]string{
			time.Now().Format(time.RFC3339),
			fmt.Sprintf("%.2f", data.SoundLevel),
		})

	alert, _ := checkAndHandleAlert(data)

	if alert {
		w.Write([]byte("ALERT"))
	} else {
		w.Write([]byte("OK"))
	}
}

// Helper function to store data in CSV
func storeCSV(filePath string, headers []string, row []string) {
	fileExists := true
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fileExists = false
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	if !fileExists {
		writer.Write(headers)
	}

	writer.Write(row)
	writer.Flush()
}
