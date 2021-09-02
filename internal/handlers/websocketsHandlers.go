package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/gorilla/websocket"
)

var wsChan = make(chan WSPayload)                  // this channel only accepts payloads from the client
var clients = make(map[WebSocketConnection]string) // all of our connected clients

// use this variable to upgrade incoming requests to websocket connection
var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// the type of information we will be sending back to the client from the websockets server
type WSJSONResponse struct {
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	MessageType    string   `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
}

type WebSocketConnection struct {
	*websocket.Conn
}

// the type of information we are sending from the client to the websockets server
type WSPayload struct {
	Action   string              `json:"action"`
	Username string              `json:"username"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json:"_"`
}

// upgrades connection to websockets and sends back a JSON response
// used for both submitting new messages to and receivng messages from the chat service
func WsEndPoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil) // we aren't going to worry about the response header right now
	if err != nil {
		log.Println(err)
	}

	log.Println("Client connected to endpoint")

	// send back a response back to the client using JSON so it is easy to parse
	var response WSJSONResponse
	response.Message = `<em><small>Connected to Server</small></em>`

	// when someone new joins the chatroom, add them to the clients map
	conn := WebSocketConnection{Conn: ws}
	clients[conn] = "" // initialize with 0 clients connected

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err)
	}

	// take everyone who is in our clients map and put them in a new goroutine
	// goroutines always start with the go keyword
	go ListenForWS(&conn)
}

// goroutine 
func ListenForWS(conn *WebSocketConnection) {
	// if the goroutine stops for any reason, start again
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()

	var payload WSPayload

	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			// do nothing, there is no payload
		} else {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}

// goroutine
func ListenToWSChannel() {
	var response WSJSONResponse // what we are sending back the client
	for {
		e := <-wsChan

		switch e.Action {
		case "username": // if the Action on the event is "username"
			// get a list of user and send it back via the broadcast function
			clients[e.Conn] = e.Username
			users := getUserList()
			response.Action = "list_users"
			response.ConnectedUsers = users
			broadcaseToAll(response)
		case "left":
			// send an action to the users to tell that someone has left
			// remove that user from the list of users
			response.Action = "list_users"
			delete(clients, e.Conn)
			users := getUserList()
			response.ConnectedUsers = users
			broadcaseToAll(response)
		case "broadcast":
			response.Action = "broadcast"
			response.Message = fmt.Sprintf("<strong>%s</strong>: %s", e.Username, e.Message)
			broadcaseToAll(response)
		}
	}
}

func getUserList() []string {
	var userList []string

	for _, x := range clients {
		if x != "" {
			userList = append(userList, x)
		}
	}
	sort.Strings(userList)

	return userList
}

// broadcase response to all of the connected users
func broadcaseToAll(response WSJSONResponse) {
	for client := range clients { // loop through all of the clients
		err := client.WriteJSON(response) // write the JSON response to all of the connected clients (users)
		if err != nil {
			log.Println("websocket err")
			_ = client.Close()      // close that connection for that client
			delete(clients, client) // remove them from the map
		}
	}
}

/*
	How does our websockets server handle requests from the client?
	After someone connects, throw them off to a goroutine that will run forever, that determines what we will do with a particular request from the client
	When we get a payload coming in, we will determine what to do depending on the action,
	and hand that off to another go routine that listens to a specific channel and does different things depending on the content of the payload
*/
