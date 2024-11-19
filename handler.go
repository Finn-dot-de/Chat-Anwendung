package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func HandleGetMessages(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodGet {
		http.Error(w, "Nur GET-Anfragen sind erlaubt", http.StatusMethodNotAllowed)
		return
	}

	// Nachrichten aus der Datenbank abrufen
	messages, err := GetMessages()
	if err != nil {
		http.Error(w, "Fehler beim Abrufen der Nachrichten", http.StatusInternalServerError)
		return
	}

	// Nachrichten als JSON zurückgeben
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		http.Error(w, "Fehler beim Kodieren der Nachrichten", http.StatusInternalServerError)
	}
}

func HandleEvents(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodGet {
		http.Error(w, "Nur GET-Anfragen sind erlaubt", http.StatusMethodNotAllowed)
		return
	}

	// Header für Server-Sent Events setzen
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Nachrichten abrufen und streamen
	messages, err := GetMessages()
	if err != nil {
		http.Error(w, "Fehler beim Abrufen der Nachrichten", http.StatusInternalServerError)
		return
	}

	for _, msg := range messages {
		fmt.Fprintf(w, "data: %s\n\n", msg.Content)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}
}

func HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		http.Error(w, "Nur POST-Anfragen sind erlaubt", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil || user.Username == "" {
		http.Error(w, "Ungültige Anfragedaten", http.StatusBadRequest)
		return
	}

	// Benutzer speichern
	err := InsertUser(user)
	if err != nil {
		http.Error(w, "Fehler beim Speichern des Benutzers", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Fehler beim Kodieren des Benutzers", http.StatusInternalServerError)
	}
}

func HandleCreateMessage(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		http.Error(w, "Nur POST-Anfragen sind erlaubt", http.StatusMethodNotAllowed)
		return
	}

	var msg struct {
		Username string `json:"username"`
		Content  string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil || msg.Username == "" || msg.Content == "" {
		http.Error(w, "Ungültige Anfragedaten", http.StatusBadRequest)
		return
	}

	// Benutzer-ID abrufen
	var senderID int
	err := db.QueryRow(`SELECT user_id FROM users WHERE username = $1`, msg.Username).Scan(&senderID)
	if err != nil {
		http.Error(w, "Benutzer nicht gefunden", http.StatusNotFound)
		return
	}

	// Nachricht speichern
	err = InsertMessage(Message{
		SenderID:  senderID,
		Content:   msg.Content,
		Timestamp: time.Now(),
	})
	if err != nil {
		http.Error(w, "Fehler beim Speichern der Nachricht", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		http.Error(w, "Nur POST-Anfragen sind erlaubt", http.StatusMethodNotAllowed)
		return
	}

	// Benutzername und Passwort aus der Anfrage lesen
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil || credentials.Username == "" || credentials.Password == "" {
		http.Error(w, "Ungültige Anfragedaten", http.StatusBadRequest)
		return
	}

	// Benutzername in der Datenbank suchen und Passwort abrufen
	var hashedPassword string
	query := `SELECT password FROM users WHERE username = $1`
	err := db.QueryRow(query, credentials.Username).Scan(&hashedPassword)
	if err != nil {
		http.Error(w, "Benutzername oder Passwort ist falsch", http.StatusUnauthorized)
		return
	}

	// Passwort validieren
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(credentials.Password)); err != nil {
		http.Error(w, "Benutzername oder Passwort ist falsch", http.StatusUnauthorized)
		return
	}

	// Erfolgreiche Authentifizierung
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Login erfolgreich",
		"username": credentials.Username,
	})
}
