package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// HandleGetMessages ist ein HTTP-Handler, der Nachrichten aus der Datenbank abruft und als JSON zurückgibt.
func HandleGetMessages(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)

	// Prüft, ob die Anfrage die richtige Methode (GET) verwendet
	if r.Method != http.MethodGet {
		log.Printf("Ungültige Methode: %s, erwartet GET", r.Method)
		http.Error(w, "Nur GET-Anfragen sind erlaubt", http.StatusMethodNotAllowed)
		return
	}

	// Nachrichten aus der Datenbank abrufen
	log.Println("Rufe Nachrichten aus der Datenbank ab...")
	messages, err := GetMessages() // GetMessages: Funktion, die Nachrichten aus der Datenbank holt
	if err != nil {
		log.Printf("Fehler beim Abrufen der Nachrichten: %v", err)
		http.Error(w, "Fehler beim Abrufen der Nachrichten", http.StatusInternalServerError)
		return
	}

	// Nachrichten als JSON zurückgeben
	log.Println("Gebe Nachrichten als JSON zurück")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		log.Printf("Fehler beim Kodieren der Nachrichten: %v", err)
		http.Error(w, "Fehler beim Kodieren der Nachrichten", http.StatusInternalServerError)
	}
}

// HandleEvents ist ein HTTP-Handler für Server-Sent Events (SSE), der kontinuierlich neue Nachrichten an Clients sendet.
func HandleEvents(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleEvents wurde aufgerufen")
	log.Printf("Bearbeite Events-Anfrage: %s %s", r.Method, r.URL.Path)

	// Nur GET-Anfragen sind erlaubt
	if r.Method != http.MethodGet {
		log.Println("Falsche Methode verwendet, Verbindung wird abgelehnt")
		http.Error(w, "Nur GET-Anfragen sind erlaubt", http.StatusMethodNotAllowed)
		return
	}

	// HTTP-Header für Server-Sent Events setzen
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	notify := r.Context().Done() // Kanal, um zu erkennen, wenn der Client die Verbindung schließt
	log.Println("Verbindung für Server-Sent Events geöffnet")
	defer log.Println("Verbindung für Server-Sent Events geschlossen")

	var lastTimestamp time.Time // Zeitstempel der letzten gesendeten Nachricht

	for {
		select {
		case <-notify:
			log.Println("Client hat die Verbindung geschlossen")
			return
		default:
			// Neue Nachrichten seit dem letzten Zeitstempel abrufen
			newMessages, err := GetMessagesSince(lastTimestamp) // Funktion, die neue Nachrichten abruft
			if err != nil {
				log.Printf("Fehler beim Abrufen neuer Nachrichten: %v", err)
				http.Error(w, "Fehler beim Abrufen neuer Nachrichten", http.StatusInternalServerError)
				return
			}

			// Wenn neue Nachrichten vorhanden sind, wird sie an den Client gesendet
			if len(newMessages) > 0 {
				lastTimestamp = newMessages[len(newMessages)-1].Timestamp
				for _, msg := range newMessages {
					messageData := map[string]interface{}{
						"username":  GetUsernameByID(msg.SenderID), // Funktion, die den Benutzernamen anhand der ID ermittelt
						"content":   msg.Content,
						"timestamp": msg.Timestamp.Format(time.RFC3339),
					}

					messageJSON, _ := json.Marshal(messageData)
					log.Printf("Sende neue Nachricht: %s", messageJSON)
					fmt.Fprintf(w, "data: %s\n\n", messageJSON)

					if f, ok := w.(http.Flusher); ok {
						f.Flush() // Daten direkt an den Client senden
					}
				}
			}

			time.Sleep(1 * time.Second)
		}
	}
}

// HandleCreateUser verarbeitet Anfragen zum Erstellen eines neuen Benutzers.
func HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)

	// Nur POST-Anfragen sind erlaubt
	if r.Method != http.MethodPost {
		log.Printf("Ungültige Methode: %s, erwartet POST", r.Method)
		http.Error(w, "Nur POST-Anfragen sind erlaubt", http.StatusMethodNotAllowed)
		return
	}

	var user User // Erwartetes JSON-Format: { "username": "name", "password": "passwort" }
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil || user.Username == "" {
		log.Printf("Ungültige Anfragedaten: %v", err)
		http.Error(w, "Ungültige Anfragedaten", http.StatusBadRequest)
		return
	}

	// Benutzer in der Datenbank speichern
	log.Printf("Speichere Benutzer in der Datenbank: %s", user.Username)
	err := InsertUser(user) // Funktion, die den Benutzer in die Datenbank schreibt
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

// HandleCreateMessage verarbeitet Anfragen zum Erstellen einer neuen Nachricht.
func HandleCreateMessage(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)

	// Nur POST-Anfragen sind erlaubt
	if r.Method != http.MethodPost {
		log.Printf("Ungültige Methode: %s, erwartet POST", r.Method)
		http.Error(w, "Nur POST-Anfragen sind erlaubt", http.StatusMethodNotAllowed)
		return
	}

	// Erwartetes JSON-Format: { "username": "name", "content": "text" }
	var msg struct {
		Username string `json:"username"`
		Content  string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil || msg.Username == "" || msg.Content == "" {
		log.Printf("Ungültige Anfragedaten: %v", err)
		http.Error(w, "Ungültige Anfragedaten", http.StatusBadRequest)
		return
	}

	// Benutzer-ID aus der Datenbank abrufen
	log.Printf("Hole Benutzer-ID für Benutzername: %s", msg.Username)
	var senderID int
	err := db.QueryRow(`SELECT id FROM users WHERE username = $1`, msg.Username).Scan(&senderID)
	if err != nil {
		log.Printf("Fehler beim Abrufen der Benutzer-ID: %v", err)
		http.Error(w, "Benutzer nicht gefunden", http.StatusNotFound)
		return
	}

	// Nachricht in der Datenbank speichern
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

// HandleLogin verarbeitet Login-Anfragen von Benutzern.
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)

	// Nur POST-Anfragen sind erlaubt
	if r.Method != http.MethodPost {
		log.Printf("Ungültige Methode: %s, erwartet POST", r.Method)
		http.Error(w, "Nur POST-Anfragen sind erlaubt", http.StatusMethodNotAllowed)
		return
	}

	// Erwartetes JSON-Format: { "username": "name", "password": "passwort" }
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil || credentials.Username == "" || credentials.Password == "" {
		log.Printf("Ungültige Anmeldedaten: %v", err)
		http.Error(w, "Ungültige Anfragedaten", http.StatusBadRequest)
		return
	}

	// Passwort-Hash aus der Datenbank abrufen
	log.Printf("Rufe gehashtes Passwort für Benutzername: %s ab", credentials.Username)
	var hashedPassword string
	query := `SELECT password FROM users WHERE username = $1`
	err := db.QueryRow(query, credentials.Username).Scan(&hashedPassword)
	if err != nil {
		log.Printf("Benutzername nicht gefunden oder Passwort falsch: %v", err)
		http.Error(w, "Benutzername oder Passwort ist falsch", http.StatusUnauthorized)
		return
	}

	// Passwort validieren
	log.Println("Vergleiche angegebenes Passwort mit gespeichertem Hash")
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(credentials.Password)); err != nil {
		log.Printf("Passwort stimmt nicht überein für Benutzername: %s", credentials.Username)
		http.Error(w, "Benutzername oder Passwort ist falsch", http.StatusUnauthorized)
		return
	}

	// Login erfolgreich
	log.Printf("Benutzer %s erfolgreich eingeloggt", credentials.Username)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Login erfolgreich",
		"username": credentials.Username,
	})
}
