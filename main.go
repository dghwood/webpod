package main

import (
	"encoding/json"
	api "github.com/dghwood/webpod/api"
	"log"
	"net/http"
	"os"
)

// TODO: Move this to api?
func HandleURL2Pod(w http.ResponseWriter, r *http.Request) {
	request := api.URL2PodRequest{}
	json.NewDecoder(r.Body).Decode(&request)
	response, _ := api.URL2Pod(request)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ERROR: Failed to encode info response, %s", err)
	}
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile("static/index.html")
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(content)
}

func main() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/home/dgh_wood/key.json")

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	// TODO: This is probably insecure
	/*
		http.Handle(
			"/static/",
			http.StripPrefix("/static/",
				http.FileServer(http.Dir("./static"))))
	*/

	//http.HandleFunc("/", HandleIndex)
	http.HandleFunc("/api/url2pod", HandleURL2Pod)

	log.Printf("Running at http://0.0.0.0:%s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
