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
    conn     *websocket.Conn
    room     string
    user     string
    userID   string
    joinedAt time.Time
}

type ClientInfo struct {
    User     string    `json:"user"`
    UserID   string    `json:"userId"`
    JoinedAt time.Time `json:"joinedAt"`
    IsOnline bool      `json:"isOnline"`
}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

var rooms = make(map[string]map[*Client]bool)
var roomsMutex = &sync.Mutex{}

type Message struct {
    Action      string          `json:"action"`
    Room        string          `json:"room,omitempty"`
    User        string          `json:"user,omitempty"`
    UserID      string          `json:"userId,omitempty"`
    Change      json.RawMessage `json:"change,omitempty"`
    Clients     []ClientInfo    `json:"clients,omitempty"`
    ClientCount int             `json:"clientCount,omitempty"`
    Timestamp   time.Time       `json:"timestamp,omitempty"`
}

func HandleWebSocket(c *gin.Context) {
    roomId := c.Param("roomId")
    log.Printf("WebSocket connection request for room: %s", roomId)

    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Printf("Error upgrading to websocket: %v", err)
        return
    }

    client := &Client{
        conn:     conn,
        room:     roomId,
        joinedAt: time.Now(),
    }

    log.Printf("WebSocket connection established for room: %s", roomId)
    go client.readMessages()
    
    // Keep the connection alive
    select {}
}

func registerClient(c *Client) {
    roomsMutex.Lock()
    defer roomsMutex.Unlock()

    if rooms[c.room] == nil {
        rooms[c.room] = make(map[*Client]bool)
    }

    rooms[c.room][c] = true
    log.Printf("Client %s joined room %s, total clients: %d", c.user, c.room, len(rooms[c.room]))

    // Broadcasting updated client count and list
    broadcastRoomUpdate(c.room)
}

func unregisterClient(c *Client) {
    //handling the room mutex
    roomsMutex.Lock()
    defer roomsMutex.Unlock()

    clients, found := rooms[c.room]
    if !found {
        return
    }

    delete(clients, c)
    clientCount := len(clients)

    if clientCount == 0 {
        delete(rooms, c.room)
    }

    log.Printf("Client %s left room %s, total clients: %d", c.user, c.room, clientCount)

    c.conn.Close()

    if clientCount > 0 {
        broadcastRoomUpdate(c.room)
    }
}

func (c *Client) readMessages() {
    defer unregisterClient(c)

    for {
        _, msgBytes, err := c.conn.ReadMessage()
        if err != nil {
            log.Printf("Error reading message from %s: %v", c.user, err)
            break
        }

        var msg Message
        if err := json.Unmarshal(msgBytes, &msg); err != nil {
            log.Printf("Error unmarshalling message: %v", err)
            continue
        }

        log.Printf("Received message: %s from user: %s", msg.Action, msg.User)

        switch msg.Action {
        case "join":
            c.user = msg.User
            c.userID = msg.UserID
            c.room = msg.Room
            registerClient(c)
            broadcastSystemMessage(c.room, c.user+" joined the room", c)

        case "edit":
            log.Printf("Broadcasting edit message in room: %s", msg.Room)
            broadcastToRoom(msg.Room, msgBytes, c)

        case "get_room_info":
            sendRoomInfo(c)

        case "leave":
            broadcastSystemMessage(c.room, c.user+" left the room", c)
						unregisterClient(c)
            return

        default:
            log.Printf("Unknown action: %s", msg.Action)
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

    targets := make([]*Client, 0, len(clients))
    for client := range clients {
        if client != sender {
            targets = append(targets, client)
        }
    }
    roomsMutex.Unlock()

    log.Printf("Broadcasting to %d clients in room %s", len(targets), roomId)
    for _, client := range targets {
        if err := client.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
            log.Printf("Broadcast write error to %s: %v", client.user, err)
        }
    }
}

func broadcastSystemMessage(roomId, text string, exclude *Client) {
    sysMsg := Message{
        Action:    "system",
        Room:      roomId,
        Change:    json.RawMessage(`{"text":"` + text + `"}`),
        Timestamp: time.Now(),
    }
    msgBytes, _ := json.Marshal(sysMsg)

    broadcastToRoom(roomId, msgBytes, exclude)
}

func getRoomInfo(roomId string) ([]ClientInfo, int) {
    roomsMutex.Lock()
    defer roomsMutex.Unlock()

    clients, found := rooms[roomId]
    if !found {
			log.Printf("No clients found in room: %s", roomId)
        return []ClientInfo{}, 0
    }

    clientList := make([]ClientInfo, 0, len(clients))
    for client := range clients {
        clientList = append(clientList, ClientInfo{
            User:     client.user,
            UserID:   client.userID,
            JoinedAt: client.joinedAt,
            IsOnline: true,
        })
    }
    return clientList, len(clientList)
}

func broadcastRoomUpdate(roomId string) {
    clientList, clientCount := getRoomInfo(roomId)

    updateMsg := Message{
        Action:      "room_update",
        Room:        roomId,
        Clients:     clientList,
        ClientCount: clientCount,
        Timestamp:   time.Now(),
    }

    msgBytes, _ := json.Marshal(updateMsg)
    broadcastToRoom(roomId, msgBytes, nil)
}

func sendRoomInfo(client *Client) {
    clientList, clientCount := getRoomInfo(client.room)

    msg := Message{
        Action:      "room_info",
        Room:        client.room,
        Clients:     clientList,
        ClientCount: clientCount,
        Timestamp:   time.Now(),
    }

    msgBytes, _ := json.Marshal(msg)
    if err := client.conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
        log.Printf("Error sending room info to %s: %v", client.user, err)
    }
}