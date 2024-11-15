package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

var (
	port int
	db   *sql.DB
)

func main() {
	// Laden der .env-Datei
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Fehler beim Laden der .env-Datei: %v", err)
	}

	// Verbindung zur Datenbank herstellen
	db, err := ConnectToDB()
	if err != nil {
		log.Fatalf("Fehler beim Verbinden mit der Datenbank: %v", err)
	}

	// Schließt die Datenbankverbindung bei Programmende
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalln("Fehler beim Schließen der DB:", err)
		}
	}(db)

	flag.IntVar(&port, "port", 8080, "port to listen on")
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("POST /api/newmessage", HandleNewMessage)
	http.HandleFunc("POST /api/newmessage", HandleNewMessage)

	log.Printf("Web Server listening on http://localhost:%d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal(err)
	}
}
