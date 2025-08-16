package main

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"os"
)

// CSV file path
const contactsFile = "static/contacts.csv"

// Contact structure
type Contact struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// Add contact (POST)
func handleAddContact(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value != "authenticated" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var c Contact
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Read existing contacts to check duplicates
	var existingContacts [][]string
	if _, err := os.Stat(contactsFile); err == nil {
		file, err := os.Open(contactsFile)
		if err == nil {
			reader := csv.NewReader(file)
			existingContacts, _ = reader.ReadAll()
			file.Close()
		}
	}

	// Check for duplicates
	for i, row := range existingContacts {
		if i == 0 {
			continue // skip header
		}
		if len(row) >= 3 && row[1] == c.Email || row[2] == c.Phone {
			json.NewEncoder(w).Encode(map[string]string{"success": "false", "message": "Contact already exists"})
			return
		}
	}

	// Append to CSV
	fileExists := true
	if _, err := os.Stat(contactsFile); os.IsNotExist(err) {
		fileExists = false
	}

	file, err := os.OpenFile(contactsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, "Cannot open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if !fileExists {
		writer.Write([]string{"Name", "Email", "Phone"})
	}
	writer.Write([]string{c.Name, c.Email, c.Phone})
	writer.Flush()

	json.NewEncoder(w).Encode(map[string]string{"success": "true"})
}

// Get all contacts (GET)
func handleGetContacts(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value != "authenticated" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	file, err := os.Open(contactsFile)
	if err != nil {
		json.NewEncoder(w).Encode([]Contact{})
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, _ := reader.ReadAll()
	var contacts []Contact
	for i, row := range rows {
		if i == 0 {
			continue // skip header
		}
		if len(row) >= 3 {
			contacts = append(contacts, Contact{
				Name:  row[0],
				Email: row[1],
				Phone: row[2],
			})
		}
	}

	json.NewEncoder(w).Encode(contacts)
}

func handleDeleteContact(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value != "authenticated" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var c Contact
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	file, err := os.Open(contactsFile)
	if err != nil {
		http.Error(w, "Cannot open file", http.StatusInternalServerError)
		return
	}
	reader := csv.NewReader(file)
	rows, _ := reader.ReadAll()
	file.Close()

	var newRows [][]string
	newRows = append(newRows, []string{"Name", "Email", "Phone"}) // keep header

	for i, row := range rows {
		if i == 0 {
			continue // skip header
		}
		if len(row) >= 3 && !(row[0] == c.Name && row[1] == c.Email && row[2] == c.Phone) {
			newRows = append(newRows, row)
		}
	}

	file, err = os.Create(contactsFile)
	if err != nil {
		http.Error(w, "Cannot save file", http.StatusInternalServerError)
		return
	}
	writer := csv.NewWriter(file)
	writer.WriteAll(newRows)
	writer.Flush()
	file.Close()

	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
