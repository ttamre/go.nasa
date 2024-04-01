package main

import (
	"log"
	"net/http"

	"github.com/ttamre/go.nasa/api"
)

const NASA_URL = "https://api.nasa.gov/planetary/earth/imagery"
const PORT = 8080

func main() {

	http.HandleFunc("/", api.ImageHandler)
	log.Printf("Server started on: http://localhost:%d", PORT)
	http.ListenAndServe(":8080", nil)
}
