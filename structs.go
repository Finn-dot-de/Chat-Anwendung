package main

// Nachricht ist die Struktur für eine neue Nachricht
type Message struct {
	SenderID int    `json:"sender_id"`
	Content  string `json:"content"`
}
