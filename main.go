package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	http.HandleFunc("/F-MAP", FrontMapHandler)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func FrontMapHandler(w http.ResponseWriter, r *http.Request) {
	// Make sure the file exists and is reachable
	tmpl, err := template.ParseFiles("maps/map.html")
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := "DATA"
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Execute error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}