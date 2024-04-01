package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ttamre/go.nasa/api"
)

const PORT = 8080

func main() {
	// Create a file server to serve static files
	fileServer := http.FileServer(http.Dir("./web"))
	http.Handle("/web/", http.StripPrefix("/web/", fileServer))

	http.HandleFunc("/", api.ImageHandler)
	log.Printf("Server started on: http://localhost:%d", PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}
