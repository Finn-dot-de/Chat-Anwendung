package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// HTTP-Handler für neue Nachrichten
func HandleNewMessage(w http.ResponseWriter, r *http.Request) {
	// Überprüfen, ob die Anfrage ein POST-Request ist
	if r.Method != http.MethodPost {
		http.Error(w, "Nur POST-Anfragen sind erlaubt", http.StatusMethodNotAllowed)
		return
	}

	// Nachricht aus dem Request-Body
	var msg Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "Ungültige Anfragedaten", http.StatusBadRequest)
		return
	}

	// Nachricht in die Datenbank einfügen
	err = InsertMessage(msg)
	if err != nil {
		http.Error(w, "Fehler beim Einfügen der Nachricht in die Datenbank", http.StatusInternalServerError)
		return
	}

	// Erfolgreiche Antwort senden
	w.WriteHeader(http.StatusCreated)
	log.Panicln(w, "Nachricht erfolgreich gespeichert")
}


func HandleGetMessages(w http.ResponseWriter, r *http.Request) {
	// Überprüfen, ob die Anfrage ein GET-Request ist
	if r.Method != http.MethodGet {
		http.Error(w, "Nur GET-Anfragen sind erlaubt", http.StatusMethodNotAllowed)
		return
	}

	// Nachrichten aus der Datenbank abrufen
	messages, err := GetMessages()
	if err != nil {
		http.Error(w, "Fehler beim Abrufen der Nachrichten aus der Datenbank", http.StatusInternalServerError)
		return
	}

	// Nachrichten als JSON zurückgeben
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func HandleEvents(w http.ResponseWriter, r *http.Request) {
	// Überprüfen, ob die Anfrage ein GET-Request ist
	if r.Method != http.MethodGet {
		http.Error(w, "Nur GET-Anfragen sind erlaubt", http.StatusMethodNotAllowed)
		return
	}

	// Nachrichten aus der Datenbank abrufen
	messages, err := GetMessages()
	if err != nil {
		http.Error(w, "Fehler beim Abrufen der Nachrichten aus der Datenbank", http.StatusInternalServerError)
		return
	}

	// Nachrichten als JSON zurückgeben
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for _, msg := range messages {
		fmt.Fprintf(w, "data: %s\n\n", msg.Content)
		w.(http.Flusher).Flush()
	}
}