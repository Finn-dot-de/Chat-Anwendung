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

	db, err = ConnectToDB()
	if err != nil {
		log.Fatalf("Fehler beim Verbinden mit der Datenbank: %v", err)
	}

	// Schließt die Datenbankverbindung bei Programmende
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalln("Fehler beim Schließen der DB:", err)
		}
	}()

	flag.IntVar(&port, "port", 8080, "port to listen on")
	flag.Parse()

	// Routen definieren
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("POST /api/create/user", HandleCreateUser)
	http.HandleFunc("POST /api/new/message", HandleCreateMessage)
	http.HandleFunc("POST /api/login", HandleLogin)
	http.HandleFunc("GET /api/events", HandleEvents)

	// Server starten
	log.Printf("Web Server listening on http://localhost:%d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal(err)
	}
}
