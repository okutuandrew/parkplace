package WbSocks 

import (
    "log"
    "net/http"
    "github.com/gorilla/websocket"
)

///WEBSOCKET UPGRADER 
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

	///////////////
	err = conn.WriteMessage(websocket.TextMessage,[]byte("Welcome to Park Place 🚗"),)
if err != nil {
    log.Println(err)
    return
}

	///////////////

    for {
        _, message, err := conn.ReadMessage()
        if err != nil {
            log.Println("Disconnected")
            break
        }

        log.Printf("Received: %s\n", message)
    }
}