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

// all the fields we will be sending back to the client from the websockets server
type WSJSONResponse struct {
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	MessageType    string   `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
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



