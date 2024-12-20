package main

import (
	"log"
	"net/http"
	"text/template"
)

func main() {

	var err error
	homeTmpl, err = template.ParseFiles("static/index.html")
	if err != nil {
		log.Fatal("Error parsing template:", err)
	}

	artistTmpl, err = template.ParseFiles("static/artist.html")
	if err != nil {
		log.Fatal("Error parsing template:", err)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/artist/", artistHandler)

	log.Println("Server started on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server error:", err)
	}
}

