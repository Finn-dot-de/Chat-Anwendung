package main

import (
	"log"
	"net/http"
	"time"
)

// Definiert ein einheitliches Zeitformat als Konstante
const timeFormat = time.RFC1123

// LoggerMiddleware protokolliert die Details jeder eingehenden HTTP-Anfrage und -Antwort,
// einschließlich Startzeit, Methode, URL, Statuscode und Dauer der Anfrage.
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Startzeit wird aufgezeichnet
		start := time.Now()
		log.Printf("Startzeit: %s | Methode: %s | URL: %s | RemoteAddr: %s",
			start.Format(timeFormat), r.Method, r.RequestURI, r.RemoteAddr)

		// Wrap den ResponseWriter, um Status und Größe zu erfassen
		wrappedWriter := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

		// Der nächste Handler in der Kette wird aufgerufen
		next.ServeHTTP(wrappedWriter, r)

		// Endzeit und Dauer werden protokolliert
		end := time.Now()
		duration := end.Sub(start)
		log.Printf("Endzeit: %s | Dauer: %s | Statuscode: %d | Antwortgröße: %d Bytes",
			end.Format(timeFormat), duration, wrappedWriter.statusCode, wrappedWriter.size)
	})
}
