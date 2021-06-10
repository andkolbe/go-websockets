package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// use this variable to upgrade to websocket connection
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
type WsPaylod struct {
	Action   string              `json:"action"`
	Username string              `json:"username"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json:"_"`
}

// upgrades connection to websockets and sends back a JSON response
func WsEndPoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil) // we aren't going to worry about the response header right now
	if err != nil {
		log.Println(err)
	}

	log.Println("Client connected to endpoint")

	// send back a response back to the client using JSON so it is easy to parse
	var response WSJSONResponse
	response.Message = `<em><small>Connected to Server</small></em>`

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err)
	}

}
