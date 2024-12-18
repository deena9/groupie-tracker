package main

import (
	"net/http"
	"log"
	"encoding/json"
)

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}



func makeArtists(w http.ResponseWriter) {
	// Fetch data from the API
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		http.Error(w, "Failed to fetch artist data", http.StatusInternalServerError)
		log.Println("Error fetching API data:", err)
		return
	}

	defer resp.Body.Close()
	// Parse JSON data into a slice of Artist
	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		http.Error(w, "Failed to parse artist data", http.StatusInternalServerError)
		log.Println("Error decoding JSON:", err)
		return
	}
}