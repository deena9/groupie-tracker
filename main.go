package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var (
	homeTmpl   *template.Template
	artistTmpl *template.Template
	artists    []Artist
)

// Artist represents the structure of the data for each artist


func homeHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/" && r.Method == "GET" {

		makeArtists(w)

		// Execute the template
		if err := homeTmpl.Execute(w, artists); err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			log.Println("Error rendering template:", err)
			return
		}
		return
	}

	if r.URL.Path != "/" && r.Method == "GET" {
		http.Error(w, "404: Not found", http.StatusNotFound)
		return
	}

	http.Error(w, "405: Not allowed", http.StatusMethodNotAllowed)

}
func artistHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Path[len("/artist/"):]
	artistID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "400: Bad Request", http.StatusBadRequest)
		return
	}

	makeArtists(w)
	thisArtist := artists[artistID-1]

	// Execute the template
	if err := artistTmpl.Execute(w, thisArtist); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Println("Error rendering template:", err)
		return
	}

}

func main() {

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/artist/", artistHandler)
	http.HandleFunc("/locations/", artistHandler)
	

	var err error
	homeTmpl, err = template.ParseFiles("index.html")
	if err != nil {
		log.Fatal("Error parsing index template", err.Error())
		return
	}

	artistTmpl, err = template.ParseFiles("artist.html")
	if err != nil {
		log.Fatal("Error parsing artist template", err.Error())
		return
	}

	// Start the server
	log.Println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server error:", err)
	}
}
