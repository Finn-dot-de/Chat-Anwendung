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

	// Mux erstellen und Routen registrieren
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./frontend")))
	mux.HandleFunc("POST /api/create/user", HandleCreateUser)
	mux.HandleFunc("POST /api/new/message", HandleCreateMessage)
	mux.HandleFunc("GET /api/get/message", HandleGetMessages)
	mux.HandleFunc("POST /api/login", HandleLogin)
	mux.HandleFunc("GET /api/events", HandleEvents)

	// Logging-Middleware anwenden
	loggedMux := LoggerMiddleware(mux)

	log.Printf("Web Server listening on http://localhost:%d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), loggedMux); err != nil {
		log.Fatal(err)
	}
}
