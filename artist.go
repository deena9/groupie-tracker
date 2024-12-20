package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"text/template"
)

const (
	apiURL            = "https://groupietrackers.herokuapp.com/api"
	artistsEndpoint   = apiURL + "/artists"
	locationsEndpoint = apiURL + "/locations"
	datesEndpoint     = apiURL + "/dates"
	relationsEndpoint = apiURL + "/relation"
)

// Artist data structure remains the same
type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

// Artist data structure remains the same
type ArtistPageData struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    []string
	Dates        []string
	Relations    []string
}

type LocationAPIResponse struct {
	Index []Location `json:"index"`
}

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

type DatesAPIResponse struct {
	Index []Dates `json:"index"`
}

type Dates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

/* type Relation struct {
	Index []struct {
		ID        int      `json:"id"`
		Dates     []string `json:"dates"`
		Locations []string `json:"locations"`
	} `json:"index"`
} */

type Relation struct {
	Index []struct {
		ID             int                 `json:"id"`
		DatesLocations map[string][]string `json:"datesLocations"`
	} `json:"index"`
}

var (
	homeTmpl   *template.Template
	artistTmpl *template.Template
	artists    []Artist
	locations  []Location
	dates      []Dates
	relations  Relation
)

func fetchArtists() error {
	resp, err := http.Get(artistsEndpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &artists); err != nil {
		return err
	}
	log.Println("Artists fetched successfully")
	return nil
}

func fetchLocations() error {
	resp, err := http.Get(locationsEndpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var apiResponse LocationAPIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return err
	}
	locations = apiResponse.Index
	log.Println("Locations fetched successfully")
	return nil
}

func fetchDates() error {
	resp, err := http.Get(datesEndpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Use the wrapper struct to unmarshal the data
	var apiResponse DatesAPIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return err
	}

	// Assign the fetched dates
	dates = apiResponse.Index
	log.Println("Dates fetched successfully")
	return nil
}

func fetchRelations() error {
	resp, err := http.Get(relationsEndpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &relations); err != nil {
		return err
	}
	log.Println("Relations fetched successfully")
	return nil
}

// Fetch All Data at Once
func fetchAllData() {
	if len(artists) == 0 {
		if err := fetchArtists(); err != nil {
			log.Fatal("Error fetching artists:", err)
		}
		if err := fetchLocations(); err != nil {
			log.Fatal("Error fetching locations:", err)
		}
		if err := fetchDates(); err != nil {
			log.Fatal("Error fetching dates:", err)
		}
		if err := fetchRelations(); err != nil {
			log.Fatal("Error fetching relations:", err)
		}
	}
}

func fetchArtistData(id int) (Artist, []string, []string, map[string][]string, error) {
	// Fetch the artist by ID
	var selectedArtist Artist
	for _, artist := range artists {
		if artist.ID == id {
			selectedArtist = artist
			break
		}
	}

	// Fetch associated data (locations, dates, and relations)
	var associatedLocations []string
	var associatedDates []string
	//var associatedRelations []string
	var associatedRelations map[string][]string

	for _, location := range locations {
		if location.ID == id {
			associatedLocations = location.Locations
		}
	}

	for _, date := range dates {
		if date.ID == id {
			associatedDates = date.Dates
		}
	}

	for _, relation := range relations.Index {
		if relation.ID == id {
			//fmt.Println(relation)
			//associatedRelations = append(associatedRelations, fmt.Sprintf("Locations: %v, Dates: %v", relation.Locations, relation.Dates))
			associatedRelations = relation.DatesLocations
		}
	}

	return selectedArtist, associatedLocations, associatedDates, associatedRelations, nil
}

var tpl *template.Template

func renderError(w http.ResponseWriter, status int, errorTemplate string) {
	w.WriteHeader(status)
	err := tpl.ExecuteTemplate(w, errorTemplate, nil)
	if err != nil {
		http.Error(w, http.StatusText(status), status)
	}
}



func artistHandler(w http.ResponseWriter, r *http.Request) {
	// Get the artist ID from URL
	// e.g. /artist/1

	fetchAllData()

	artistID := r.URL.Path[len("/artist/"):]
	if artistID == "" {
		http.Error(w, "Artist ID is required", http.StatusBadRequest)
		return
	}

	id := 0
	// Convert the artist ID from string to int
	fmt.Sscanf(artistID, "%d", &id)

	// Fetch the artist and related data
	artist, locations, dates, relations, err := fetchArtistData(id)
	if err != nil {
		http.Error(w, "Failed to fetch artist data", http.StatusInternalServerError)
		return
	}

	stringRelations := []string{}
	for location, dates := range relations {
		for _, date := range dates {
			stringRelations = append(stringRelations, location+" "+date)
		}
	}

	APD := ArtistPageData{
		ID:           artist.ID,
		Image:        artist.Image,
		Name:         artist.Name,
		Members:      artist.Members,
		CreationDate: artist.CreationDate,
		FirstAlbum:   artist.FirstAlbum,
		Locations:    locations,
		Dates:        dates,
		Relations:    stringRelations,
	}


	if err := artistTmpl.Execute(w, APD); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Println("Template render error:", err)
	}

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fetchAllData()

	if err := homeTmpl.Execute(w, artists); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Println("Template render error:", err)
	}
}
