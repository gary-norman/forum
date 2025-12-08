package handlers

import (
	"fmt"
	"github.com/gary-norman/forum/internal/app"
	"github.com/gorilla/websocket"
	"net/http"
)

type WebsocketHandler struct {
	App      *app.App
	User     *UserHandler
	Post     *PostHandler
	Comment  *CommentHandler
	Reaction *ReactionHandler
	Channel  *ChannelHandler
	Mod      *ModHandler
	// Notification *NotificationHandler
	// Membership *MembershipHandler
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (ws *WebsocketHandler) HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	//Handle Websocket messages here
	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Received message: %s\n", data)

		// Echo the message back to the client
		if err := conn.WriteMessage(messageType, data); err != nil {
			fmt.Println(err)
			return
		}
	}
}
