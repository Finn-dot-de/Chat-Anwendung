package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// ConnectToDB stellt eine Verbindung zur Datenbank her und gibt sie zurück.
func ConnectToDB() (*sql.DB, error) {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	// Verbindungszeichenkette erstellen
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("Erfolgreich mit der Datenbank verbunden!")
	return db, nil
}

// InsertMessage fügt eine neue Nachricht in die Datenbank ein.
func InsertMessage(msg Message) error {
	log.Println("Hier!!")
	query := `INSERT INTO messages (sender_id, content, timestamp) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, msg.SenderID, msg.Content, time.Now())
	return err
}

// GetMessages holt Nachrichten aus der Datenbank.
func GetMessages() ([]Message, error) {
	query := `SELECT sender_id, content, timestamp FROM messages ORDER BY timestamp DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.SenderID, &msg.Content, &msg.Timestamp); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func GetMessagesSince(timestamp time.Time) ([]Message, error) {
	query := `SELECT sender_id, content, timestamp FROM messages WHERE timestamp > $1 ORDER BY timestamp ASC`
	rows, err := db.Query(query, timestamp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.SenderID, &msg.Content, &msg.Timestamp); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}


func InsertUser(user User) error {
	// Passwort hashen
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("fehler beim Hashen des Passworts: %v", err)
	}

	query := `INSERT INTO users (username, password) VALUES ($1, $2)`
	_, err = db.Exec(query, user.Username, string(hashedPassword))
	return err
}

// Funktion zum Überprüfen eines Passworts
func ValidateUser(username, password string) (bool, error) {
	var hashedPassword string
	query := `SELECT password FROM users WHERE username = $1`
	err := db.QueryRow(query, username).Scan(&hashedPassword)
	if err != nil {
		return false, err
	}

	// Passwort validieren
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false, nil
	}
	return true, nil
}

func GetUsernameByID(userID int) string {
	var username string
	err := db.QueryRow(`SELECT username FROM users WHERE id = $1`, userID).Scan(&username)
	if err != nil {
		log.Printf("Fehler beim Abrufen des Benutzernamens für ID %d: %v", userID, err)
		return "Unbekannt"
	}
	return username
}
