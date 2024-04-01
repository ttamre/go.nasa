package api

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GetPage(w, r)
	case "POST":
		GetImage(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Serve the HTML page
func GetPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/index.html")
}

func GetImage(w http.ResponseWriter, r *http.Request) {
	/* QUERY PARAMS
	- lat			FLOAT	Latitude of the image center.
	- lon			FLOAT	Longitude of the image center.
	- dim			FLOAT	(default 0.025) Width and height of image in degrees.
	- date			STRING	(default: today) Date of image in format YYYY-MM-DD
	- cloud_score	BOOL 	(default: false; *NOT AVAILABLE ATM*) calculate the percentage of the image covered by clouds
	- api_key		STRING	(default: DEMO_KEY) NASA API key.
	*/

	var lat string
	var lon string
	var dim string
	var date string
	var cloud_score string
	var api_key = GetAPIKey()

	if r.Method == "POST" {
		// Process form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusInternalServerError)
			log.Fatal(err)
		}

		// Get the query parameters from form data
		log.Printf("FORM: %s", r.Form)
		lat = r.FormValue("lat")
		lon = r.FormValue("lon")
		dim = r.FormValue("dim")
		date = r.FormValue("date")
		cloud_score = r.FormValue("cloud_score")
	}

	// Create a new URL
	url := url.URL{Scheme: "https", Host: "api.nasa.gov", Path: "/planetary/earth/imagery"}
	values := url.Query()

	// required query parameters
	values.Add("lat", lat)
	values.Add("lon", lon)
	values.Add("api_key", api_key)

	// optional query parameters
	if dim != "" {
		values.Add("dim", dim)
	}
	if date != "" {
		values.Add("date", date)
	}
	if cloud_score != "" {
		values.Add("cloud_score", cloud_score)
	}

	url.RawQuery = values.Encode()
	log.Printf("URL: %s", url.String())

	// Make a GET request to the NASA API
	resp, err := http.Get(url.String())
	if err != nil {
		http.Error(w, "Failed to get image", http.StatusInternalServerError)
		log.Printf("Failed to get image: %s", resp.Status)
	}

	// Check the response status code. If not 200, log error and redirect to main page
	if resp.StatusCode != http.StatusOK {
		// http.Error(w, "Bad response from NASA API", http.StatusInternalServerError)
		log.Printf("Bad response from NASA API: %s", resp.Status)
		GetPage(w, r)
		return

	}

	// Read the response body
	defer resp.Body.Close()
	image, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read image", http.StatusInternalServerError)
		log.Printf("Failed to read image: %s", err)
	}

	// Write the image to the response
	w.Header().Set("Content-Type", "image/png")
	_, err = w.Write(image)
	if err != nil {
		http.Error(w, "Failed to write image", http.StatusInternalServerError)
		log.Printf("Failed to write image: %s", err)
	}

	log.Printf("RESPONSE: %s", resp.Status)
	log.Printf("IMAGE: %d bytes", len(image))
}

// Returns the NASA API key from the .env file
func GetAPIKey() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("NASA_API_KEY")
}
