package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type data struct {
	Street string  // exported so template can see it
}

func main() {
	var details data

	// Pass details (or a pointer to it) into the handler
	// by wrapping it in a closure
	http.HandleFunc("/F-MAP", func(w http.ResponseWriter, r *http.Request) {
		FrontMapHandler(w, r, &details)
	})

	http.HandleFunc("/BOOKPARKING", Bookparking)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func FrontMapHandler(w http.ResponseWriter, r *http.Request, d *data) {
	// d is *data, so you can modify it
	d.Street = "KIMATHI STREET"

	tmpl, err := template.ParseFiles("maps/map.html")
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, d)
	if err != nil {
		http.Error(w, "Execute error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func Bookparking(w http.ResponseWriter, r *http.Request){

fmt.Fprintf(w, "Hello, this is your Go HTTP server!")
}