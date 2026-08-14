package WbSocks

import (
	"log"
	"net/http"
	"os"
	"parkplace/static"
	"time"

	"github.com/gorilla/websocket"
)

/// WEBSOCKET UPGRADER
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // OK for development
	},
}

func DriverMapWbScock(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	log.Println("Client connected")

	// Send initial greeting message
	err = conn.WriteMessage(websocket.TextMessage, []byte("Welcome to Park Place 🚗"))
	if err != nil {
		log.Println(err)
		return
	}

	// 1. START BACKGROUND WORKER (Goroutine)
	// This monitors the file and pushes live updates independently of the reading loop below.
	go func() {
		var lastModified time.Time

		for {
			fileInfo, err := os.Stat("static/parkingData.json")
			if err != nil {
				log.Println("Error reading parkingData.json status:", err)
				time.Sleep(2 * time.Second)
				continue
			}

			// If the file modification time has changed, process data and send it
			if fileInfo.ModTime().After(lastModified) {
				lastModified = fileInfo.ModTime()

				// Fetch calculated tallies from your function
				tallyData := static.ParkingData()

				// Push data to frontend as JSON over the open socket connection
				err = conn.WriteJSON(tallyData)
				if err != nil {
					// Connection likely dropped; exiting this goroutine safely
					log.Println("WriteJSON failed, closing stream:", err)
					return
				}
				log.Println("⚡ Live metrics pushed automatically to dashboard!")
			}

			// Check file state every 1 second
			time.Sleep(1 * time.Second)
		}
	}()

	// 2. EXISTING READ LOOP
	// This remains active, listening for incoming data strings sent from your web browser client.
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Disconnected")
			break
		}

		log.Printf("Received: %s\n", message)
	}
}
