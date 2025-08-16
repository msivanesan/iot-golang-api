package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT") // Render provides this
	if port == "" {
		port = "8000" // fallback for local run
	}
	// Ensure static folder exists
	if _, err := os.Stat("static"); os.IsNotExist(err) {
		os.Mkdir("static", 0755)
	}
	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Routes
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", handleLogout)
	http.HandleFunc("/sensor-data", handleSensorData)
	http.HandleFunc("/", home)
	http.HandleFunc("/get-thresholds", handleGetThresholds)
	http.HandleFunc("/update-thresholds", handleUpdateThresholds)
	http.HandleFunc("/contacts", handleContactsPage)
	http.HandleFunc("/contacts/add", handleAddContact)
	http.HandleFunc("/contacts/list", handleGetContacts)
	http.HandleFunc("/dashboard", dashboard)
	http.HandleFunc("/contacts/delete", handleDeleteContact)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}
