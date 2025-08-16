package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var storedCreds Credentials

func init() {
	// Load static credentials from file
	file, err := os.Open("credentials.json")
	if err != nil {
		panic("Error opening credentials file: " + err.Error())
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&storedCreds); err != nil {
		panic("Error decoding credentials: " + err.Error())
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var input Credentials
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if input.Username == storedCreds.Username && input.Password == storedCreds.Password {
		// Set session cookie
		http.SetCookie(w, &http.Cookie{
			Name:  "session",
			Value: "authenticated",
			Path:  "/",
			// Secure: true,   // enable if HTTPS
			// HttpOnly: true, // JS can't access
		})

		json.NewEncoder(w).Encode(map[string]bool{"success": true})
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"success": false})
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	// Invalidate the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0), // Expire immediately
	})

	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
