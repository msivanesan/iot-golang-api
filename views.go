package main

import (
	"html/template"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	tmpl.Execute(w, struct{}{})
}

func dashboard(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value != "authenticated" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	tmpl := template.Must(template.ParseFiles("templates/dashboard.html"))
	tmpl.Execute(w, struct{}{})
}

// Serve the contacts page
func handleContactsPage(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value != "authenticated" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	tmpl := template.Must(template.ParseFiles("templates/contacts.html"))
	tmpl.Execute(w, struct{}{})

}
