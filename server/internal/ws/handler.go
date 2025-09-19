package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	room string
	user string
	userID string
	joinedAt time.Time
}

type ClientInfo struct {
	User string `json:"user"`
	UserID string `json:"userId"`
	JoinedAt time.Time `json:"joinedAt"`
	isOnline bool `json:"isOnline"`
}
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var rooms = make(map[string]map[*Client]bool) // roomId to clients
var roomsMutex = &sync.Mutex{}

type Message struct {
	Action string          `json:"action"`
	Room   string          `json:"room,omitempty"`
	User   string          `json:"user,omitempty"`
	Change json.RawMessage `json:"change,omitempty"`
}

func HandleWebSocket(c *gin.Context) {
	roomId := c.Param("roomId")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error upgrading to websocket:", err)
		return
	}

	client := &Client{
		conn: conn,
		room: roomId,
	}

	registerClient(client)

	go client.readMessages()

	select {}
}

func registerClient(c *Client) {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	if rooms[c.room] == nil {
		rooms[c.room] = make(map[*Client]bool)
	}

	rooms[c.room][c] = true
	log.Printf("Client joined room %s, total clients: %d\n", c.room, len(rooms[c.room]))
}

func unregisterClient(c *Client) {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	clients, found := rooms[c.room]
	if !found {
		return
	}

	delete(clients, c)

	if len(clients) == 0 {
		delete(rooms, c.room)
	}

	log.Printf("Client left room %s, total clients: %d\n", c.room, len(clients))

	c.conn.Close()
}

func (c *Client) readMessages() {
	defer unregisterClient(c) // unregister and close connection

	for {
		_, msgBytes, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		var msg Message
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

		switch msg.Action {
		case "join":
			c.user = msg.User
			c.room = msg.Room
			registerClient(c)
			broadcastSystemMessage(c.room, c.user+" joined the room", c)

		case "edit":
			broadcastToRoom(msg.Room, msgBytes, c)

		case "leave":
			broadcastSystemMessage(c.room, c.user+" left the room", c)
			unregisterClient(c)
			return

		default:
			log.Printf("Unknown action: %s\n", msg.Action)
		}
	}
}

func broadcastToRoom(roomId string, msg []byte, sender *Client) {
	roomsMutex.Lock()
	clients, found := rooms[roomId]
	if !found {
		roomsMutex.Unlock()
		return
	}

	// Copy clients to avoid holding lock while writing
	targets := make([]*Client, 0, len(clients))
	for client := range clients {
		if client != sender {
			targets = append(targets, client)
		}
	}
	roomsMutex.Unlock()

	for _, client := range targets {
		if err := client.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("Broadcast write error:", err)
		}
	}
}

func broadcastSystemMessage(roomId, text string, exclude *Client) {
	sysMsg := Message{
		Action: "system",
		Room:   roomId,
		Change: json.RawMessage(`{"text":"` + text + `"}`),
	}
	msgBytes, _ := json.Marshal(sysMsg)

	broadcastToRoom(roomId, msgBytes, exclude)
}
