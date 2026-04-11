package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"
)

type PageData struct {
    Street string
}

func main() {
    // Static files handler - MUST come BEFORE the specific routes
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

    // Your page route
    http.HandleFunc("/F-MAP", FrontMapHandler)

    // Optional route
    http.HandleFunc("/BOOKPARKING", Bookparking)

    fmt.Println("🚀 Server running on http://localhost:8080")
    fmt.Println("   Open → http://localhost:8080/F-MAP")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func FrontMapHandler(w http.ResponseWriter, r *http.Request) {
    data := PageData{Street: "KIMATHI STREET"}

    tmpl, err := template.ParseFiles("maps/map.html")
    if err != nil {
        http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    log.Println("✅ User accessed F-MAP page")
    err = tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, "Execute error: "+err.Error(), http.StatusInternalServerError)
    }
}

func Bookparking(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, this is your Go HTTP server!")
}