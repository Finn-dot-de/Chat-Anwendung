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
		log.Printf("Ungültige Methode: %s, erwartet GET", r.Method)
		http.Error(w, "Nur GET-Anfragen sind erlaubt", http.StatusMethodNotAllowed)
		return
	}

	log.Println("Rufe Nachrichten aus der Datenbank ab...")
	messages, err := GetMessages()
	if err != nil {
		log.Printf("Fehler beim Abrufen der Nachrichten: %v", err)
		http.Error(w, "Fehler beim Abrufen der Nachrichten", http.StatusInternalServerError)
		return
	}

	log.Println("Gebe Nachrichten als JSON zurück")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		log.Printf("Fehler beim Kodieren der Nachrichten: %v", err)
		http.Error(w, "Fehler beim Kodieren der Nachrichten", http.StatusInternalServerError)
	}
}

func HandleEvents(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleEvents wurde aufgerufen")
	log.Printf("Bearbeite Events-Anfrage: %s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodGet {
		log.Println("Falsche Methode verwendet, Verbindung wird abgelehnt")
		http.Error(w, "Nur GET-Anfragen sind erlaubt", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	notify := r.Context().Done()
	log.Println("Verbindung für Server-Sent Events geöffnet")
	defer log.Println("Verbindung für Server-Sent Events geschlossen")

	// Initialer Zeitstempel setzen
	var lastTimestamp time.Time

	for {
		select {
		case <-notify:
			log.Println("Client hat die Verbindung geschlossen")
			return
		default:
			// Neue Nachrichten abrufen
			newMessages, err := GetMessagesSince(lastTimestamp)
			if err != nil {
				log.Printf("Fehler beim Abrufen neuer Nachrichten: %v", err)
				http.Error(w, "Fehler beim Abrufen neuer Nachrichten", http.StatusInternalServerError)
				return
			}

			if len(newMessages) > 0 {
				lastTimestamp = newMessages[len(newMessages)-1].Timestamp
				for _, msg := range newMessages {
					// Nachricht als JSON-String senden
					messageData := map[string]interface{}{
						"username":  GetUsernameByID(msg.SenderID), // Funktion, um den Username anhand der ID zu holen
						"content":   msg.Content,
						"timestamp": msg.Timestamp.Format(time.RFC3339),
					}

					messageJSON, _ := json.Marshal(messageData)
					log.Printf("Sende neue Nachricht: %s", messageJSON)
					fmt.Fprintf(w, "data: %s\n\n", messageJSON)

					if f, ok := w.(http.Flusher); ok {
						f.Flush()
					}
				}
			}

			time.Sleep(1 * time.Second)
		}
	}
}

func HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		log.Printf("Ungültige Methode: %s, erwartet POST", r.Method)
		http.Error(w, "Nur POST-Anfragen sind erlaubt", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil || user.Username == "" {
		log.Printf("Ungültige Anfragedaten: %v", err)
		http.Error(w, "Ungültige Anfragedaten", http.StatusBadRequest)
		return
	}

	log.Printf("Speichere Benutzer in der Datenbank: %s", user.Username)
	err := InsertUser(user)
	if err != nil {
		log.Printf("Fehler beim Speichern des Benutzers: %v", err)
		http.Error(w, "Fehler beim Speichern des Benutzers", http.StatusInternalServerError)
		return
	}

	log.Printf("Benutzer %s erfolgreich erstellt", user.Username)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("Fehler beim Kodieren der Benutzerdaten: %v", err)
		http.Error(w, "Fehler beim Kodieren des Benutzers", http.StatusInternalServerError)
	}
}

func HandleCreateMessage(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		log.Printf("Ungültige Methode: %s, erwartet POST", r.Method)
		http.Error(w, "Nur POST-Anfragen sind erlaubt", http.StatusMethodNotAllowed)
		return
	}

	var msg struct {
		Username string `json:"username"`
		Content  string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil || msg.Username == "" || msg.Content == "" {
		log.Printf("Ungültige Anfragedaten: %v", err)
		http.Error(w, "Ungültige Anfragedaten", http.StatusBadRequest)
		return
	}

	log.Printf("Hole Benutzer-ID für Benutzername: %s", msg.Username)
	var senderID int
	err := db.QueryRow(`SELECT id FROM users WHERE username = $1`, msg.Username).Scan(&senderID)
	if err != nil {
		log.Printf("Fehler beim Abrufen der Benutzer-ID: %v", err)
		http.Error(w, "Benutzer nicht gefunden", http.StatusNotFound)
		return
	}

	log.Printf("Speichere Nachricht für Benutzer-ID %d", senderID)
	err = InsertMessage(Message{
		SenderID:  senderID,
		Content:   msg.Content,
		Timestamp: time.Now(),
	})
	if err != nil {
		log.Printf("Fehler beim Speichern der Nachricht: %v", err)
		http.Error(w, "Fehler beim Speichern der Nachricht", http.StatusInternalServerError)
		return
	}

	log.Printf("Nachricht erfolgreich erstellt von Benutzer %s", msg.Username)
	w.WriteHeader(http.StatusCreated)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		log.Printf("Ungültige Methode: %s, erwartet POST", r.Method)
		http.Error(w, "Nur POST-Anfragen sind erlaubt", http.StatusMethodNotAllowed)
		return
	}

	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil || credentials.Username == "" || credentials.Password == "" {
		log.Printf("Ungültige Anmeldedaten: %v", err)
		http.Error(w, "Ungültige Anfragedaten", http.StatusBadRequest)
		return
	}

	log.Printf("Rufe gehashtes Passwort für Benutzername: %s ab", credentials.Username)
	var hashedPassword string
	query := `SELECT password FROM users WHERE username = $1`
	err := db.QueryRow(query, credentials.Username).Scan(&hashedPassword)
	if err != nil {
		log.Printf("Benutzername nicht gefunden oder Passwort falsch: %v", err)
		http.Error(w, "Benutzername oder Passwort ist falsch", http.StatusUnauthorized)
		return
	}

	log.Println("Vergleiche angegebenes Passwort mit gespeichertem Hash")
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(credentials.Password)); err != nil {
		log.Printf("Passwort stimmt nicht überein für Benutzername: %s", credentials.Username)
		http.Error(w, "Benutzername oder Passwort ist falsch", http.StatusUnauthorized)
		return
	}

	log.Printf("Benutzer %s erfolgreich eingeloggt", credentials.Username)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Login erfolgreich",
		"username": credentials.Username,
	})
}
