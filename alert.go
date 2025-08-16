package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var lastAlertState = false // false = OK, true = Alert active

func loadThresholds() (Thresholds, error) {
	file, err := os.Open(thresholdFile)
	if err != nil {
		return Thresholds{}, err
	}
	defer file.Close()

	var t Thresholds
	err = json.NewDecoder(file).Decode(&t)
	if err != nil {
		return Thresholds{}, err
	}

	return t, nil
}

// sendNotification sends the alert to authorised persons
func sendNotification(message string) {
	// TODO: integrate with Twilio/SMTP/Telegram
	// Example: pass the CSV file path directly
	filePath := "static/contacts.csv"
	emails, phones, err := ReadContactsEmailsPhones(filePath)
	if err != nil {
		fmt.Println("Error reading contacts:", err)
		return
	}
	fmt.Println("Emails:", emails)
	fmt.Println("Phones:", phones)
	fmt.Println("NOTIFICATION:", message)
	//send mail
	sendEmail(emails, "NOTIFICATION:", message)
	//send sms
	// sendSMSAsync(phones, message)
}

// checkAndHandleAlert checks thresholds and sends alert only on state change
func checkAndHandleAlert(data SensorData) (bool, string) {
	// Load thresholds from JSON file
	thresholds, err := loadThresholds()
	if err != nil {
		fmt.Println("Error loading thresholds:", err)
		// fallback to default if file error
		thresholds = Thresholds{Temperature: 30.0, Humidity: 70.0, Sound: 80.0}
	}

	alert := false
	alertMsg := ""

	if data.Temperature > thresholds.Temperature {
		alert = true
		alertMsg += fmt.Sprintf("High Temp: %.2f°C\n", data.Temperature)
	}
	if data.Humidity > thresholds.Humidity {
		alert = true
		alertMsg += fmt.Sprintf("High Humidity: %.2f%%\n", data.Humidity)
	}
	if int(data.SoundLevel) > int(thresholds.Sound) {
		alert = true
		alertMsg += fmt.Sprintf("High Sound: %.2f\n", data.SoundLevel)
	}

	// Send alert only when state changes
	if alert && !lastAlertState {
		sendNotification("🚨 ALERT TRIGGERED:\n" + alertMsg)
		lastAlertState = true
	}
	if !alert && lastAlertState {
		sendNotification("✅ All readings back to normal")
		lastAlertState = false
	}
	return alert, alertMsg
}
