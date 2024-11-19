package main

import (
	"net/http"
	"time"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Message struct {
	SenderID  int       `json:"sender_id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// responseWriterWrapper ist eine Struktur, die http.ResponseWriter erweitert,
// um Statuscode und Antwortgröße mitzuprotokollieren.
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
	size       int
}
