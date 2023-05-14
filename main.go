package main

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	// start http server
	http.HandleFunc("/ws", wsHandler)

	go defaultGroupChat.executeOperations()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Errorln("Failed to upgrade connection:", err)
		return
	}

	// make sure to close the connection before exit
	defer conn.Close()

	defaultGroupChat.addConnection(conn)
	// remove from conns map
	defer defaultGroupChat.removeConnection(conn.RemoteAddr().String())

	for {
		messageType, messageBytes, err := conn.ReadMessage()
		if err != nil {
			log.Errorln("Failed to read message:", err)
			break
		}

		conn.SetReadDeadline(time.Now().Add(time.Second * 300))
		conn.SetWriteDeadline(time.Now().Add(time.Second * 300))

		message := string(messageBytes)
		log.Infof("Received message type %d with contents: %s\n", messageType, message)
		defaultGroupChat.broadcast(conn, message)
	}
}
