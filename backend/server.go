package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var port int

func main() {
	flag.IntVar(&port, "port", 8080, "port to listen on")
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir("../frontend")))

	log.Printf("Web Server listening on http://localhost:%d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal(err)
	}
}
