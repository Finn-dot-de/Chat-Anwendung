package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// ConnectToDB stellt eine Verbindung zur Datenbank her und gibt diese zurück.
func ConnectToDB() (*sql.DB, error) {
	// Laden der Umgebungsvariablen.
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	// Erstellen der Verbindungszeichenkette.
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Versuch, eine Verbindung zur Datenbank herzustellen.
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	// Versuch, die Datenbank anzupingen, um die Verbindung zu testen.
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Wenn die Verbindung erfolgreich hergestellt wurde, wird eine Erfolgsmeldung gedruckt.
	fmt.Println("Successfully connected!")

	// Gibt die Datenbankverbindung und nil für den Fehler zurück.
	return db, nil
}

func InsertMessage(msg Message) error {
	query := `INSERT INTO messages (sender_id, content, timestamp) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, msg.SenderID, msg.Content, time.Now())
	return err
}

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
		err := rows.Scan(&msg.SenderID, &msg.Content, &msg.Timestamp)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func InsertUser(user User) error {
	query := `INSERT INTO users (username, password) VALUES ($1, $2)`
	_, err := db.Exec(query, user.Username, user.Password)
	return err
}
