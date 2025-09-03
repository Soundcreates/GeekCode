package ws 

import(
	"github.com/gorilla/websocket"
	// "gorm.io/gorm"
	"net/http"
	"github.com/gin-gonic/gin"
	"fmt"
	// "geekCode/internal/models"
	"sync"
	"log"
	"encoding/json"
	"github.com/google/uuid"
)


//defining structs here

type Message struct {
	Type string `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Client struct {
	Conn *websocket.Conn
	Room *Room
	Send chan []byte
	UserID string
}

type Room struct {
	ID string
	Clients map[*Client]bool
	Broadcast chan []byte
	Register chan *Client
	Unregister chan *Client
	Mutex sync.Mutex
}


//now our global state is herre
var upgrader = websocket.Upgrader {
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}


var rooms = make(map[string]*Room)
var roomsMutex sync.Mutex

//room methods are now here

func NewRoom(id string) *Room {
	return  &Room{
		ID: id,
		Clients: make(map[*Client]bool),
		Broadcast: make(chan []byte),
		Register: make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (room *Room) Run() {
	for {
		select{
		case client := <-room.Register:
				room.registerClient(client)
		
	case client := <-room.Unregister:
		room.unregisterClient(client)

	case message:= <-room.Broadcast:
		room.broadcastMessage(message)
	}
}
}


func (room *Room) registerClient(client *Client) {
	room.Mutex.Lock()
	defer room.Mutex.Unlock()
	room.Clients[client] = true
	log.Printf("Client %s joined room %s", client.UserID, room.ID)
}

func (room *Room) unregisterClient(client *Client) {
	room.Mutex.Lock() 
	defer room.Mutex.Unlock()
	if _, ok := room.Clients[client]; ok {
		delete(room.Clients, client)
		log.Printf("Client %s has left room %s", client.UserID, room.ID)
		room.notifyUserListChanged()
	}

	if len(room.Clients) == 0{
		log.Printf("Room %s is empty, deleting.", room.ID)
		roomsMutex.Lock()
		delete(rooms, room.ID)
		roomsMutex.Unlock()
		return
	}
}

func (room *Room) broadcastMessage(message []byte) {
	room.Mutex.Lock()
	defer room.Mutex.Unlock()
	for client := range room.Clients {
		select{ 
		case client.Send <- message:
		default:
			close(client.Send)
			delete(room.Clients, client)
		}
	}
}


func (room *Room) notifyUserListChanged() {
	var userIDs []string
	for client:= range room.Clients {
		userIDs= append(userIDs, client.UserID)
	}

	payload, _ := json.Marshal(userIDs)
	msg := Message {
		Type: "user_list_update",
		Payload : payload,
	}

	msgBytes, _ := json.Marshal(msg)
	room.broadcastMessage(msgBytes)
}

//now client methods

func (c *Client) ReadPump() {
	defer func() {
		c.Room.Unregister <- c
		c.Conn.Close()
	}() 
	for {
		_, messageBytes, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("error reading message: %v", err)
			break
		}
		
		var msg Message
		if err := json.Unmarshal(messageBytes, &msg); err != nil {
			log.Printf("error unmarshalling message: %v", err)
			continue
		}
		
		var payloadMap map[string]interface{}
		json.Unmarshal(msg.Payload, &payloadMap)
		payloadMap["senderId"] = c.UserID

		newPayload , _ := json.Marshal(payloadMap)
		msg.Payload = newPayload

		finalMessage, _ := json.Marshal(msg)
		c.Room.Broadcast <- finalMessage
	}
}

func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()
	for message := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("error writing message: %v", err)
			break
		}
	}
}

//now handling the gin part

func HandleWebSocket(c *gin.Context) {
	fmt.Println("Handling WebSocket connection")
	roomId := c.Param("roomId")
	//i will  be following some steps
	//step1 - create or find a room
	roomsMutex.Lock()
	room , ok := rooms[roomId] //checking if room exists or not
	if !ok {
		log.Printf("Creating new room: %s", roomId)
		room = NewRoom(roomId) //just making a new room if it does not exist
		rooms[roomId] = room //assinging that room to that index

		go room.Run() //running the room in a go routine
	}

	//step2- upgrade the connection to websocket
	conn, err :=upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil{
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	//step3-to register the new client with rooms hub
	client := &Client {
		Conn: conn,
		Room: room,
		Send: make(chan []byte, 256), //buffering the channel for outgoing messages
		UserID: uuid.New().String(),
	}

	//step4- basically just send the client to the clients room and then register in it
	client.Room.Register <- client

	//step5 - start the clients read and write pummps in seperate goroutines
	//this lets the handler to return while the connection remains up and active

	go client.WritePump()
	go client.ReadPump()
}


