package api

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

func GetAPIKey() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("NASA_API_KEY")
}
func ImageHandler(w http.ResponseWriter, r *http.Request) {
	/* QUERY PARAMS
	- lat			FLOAT	Latitude of the image center.
	- lon			FLOAT	Longitude of the image center.
	- dim			FLOAT	(default 0.025) Width and height of image in degrees.
	- date			STRING	(default: today) Date of image in format YYYY-MM-DD
	- cloud_score	BOOL 	(default: false; *NOT AVAILABLE ATM*) calculate the percentage of the image covered by clouds
	- api_key		STRING	(default: DEMO_KEY) NASA API key.
	*/

	var lat = "1.5"
	var lon = "100.75"
	var dim = "0.1"
	var date = "2024-01-01"
	// var cloud_score String

	if r.Method == "POST" {
		// Process form data
		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}

		// Get the query parameters from form data
		lat = r.Form.Get("lat")
		lon = r.Form.Get("lon")
		dim = r.Form.Get("dim")
		date = r.Form.Get("date")
		// cloud_score := r.Form.Get("cloud_score")
	}

	// Create a new URL
	url := url.URL{Scheme: "https", Host: "api.nasa.gov", Path: "/planetary/earth/imagery"}
	values := url.Query()
	values.Add("lat", lat)
	values.Add("lon", lon)
	values.Add("dim", dim)
	values.Add("date", date)
	// values.Add("cloud_score", "true")
	values.Add("api_key", "DEMO_KEY")
	url.RawQuery = values.Encode()

	// Make a GET request to the NASA API
	resp, err := http.Get(url.String())
	if err != nil {
		log.Fatal(err)
	}

	// Read the response body
	defer resp.Body.Close()
	image, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Write the image to the response
	w.Header().Set("Content-Type", "image/png")
	_, err = w.Write(image)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("LAT: %s, LON: %s, DIM: %s, DATE: %s", lat, lon, dim, date)
}
