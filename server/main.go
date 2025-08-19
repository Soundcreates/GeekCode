package main

import (
    "fmt"
    "net/http"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true // allow all for dev
    },
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan []byte)

func main() {
    // Serve WebSocket endpoint
    http.HandleFunc("/ws", handleWebSocket)

    // Serve static files (your HTML client) from "static" folder
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/", fs)

    // Start broadcaster
    go handleMessages()

    fmt.Println("Server started at http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println("Upgrade error:", err)
        return
    }
    defer conn.Close()

    clients[conn] = true
    fmt.Println("New client connected")

    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            fmt.Println("Client disconnected:", err)
            delete(clients, conn)
            break
        }
        broadcast <- msg
    }
}

func handleMessages() {
    for {
        msg := <-broadcast
        fmt.Println("Broadcasting:", string(msg)) // 👀 debug log
        for client := range clients {
            err := client.WriteMessage(websocket.TextMessage, msg)
            if err != nil {
                fmt.Println("Write error:", err)
                client.Close()
                delete(clients, client)
            }
        }
    }
}

