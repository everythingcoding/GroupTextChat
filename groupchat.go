package main

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type connsType map[string]*websocket.Conn

var (
	defaultGroupChat = &GroupChat{
		operationChan: make(chan func(connsMap connsType)),
	}
)

type GroupChat struct {
	operationChan chan func(connsMap connsType)
}

// keeps on receiving operations in the form of functions from channel
// and executes them
func (g *GroupChat) executeOperations() {
	m := make(connsType)
	for op := range g.operationChan {
		op(m)
	}
}

func (g *GroupChat) addConnection(conn *websocket.Conn) {
	g.operationChan <- func(connsMap connsType) {
		connsMap[conn.RemoteAddr().String()] = conn
	}
}

func (g *GroupChat) removeConnection(remoteAddr string) {
	g.operationChan <- func(connsMap connsType) {
		delete(connsMap, remoteAddr)
	}
}

func (g *GroupChat) broadcast(senderConn *websocket.Conn, message string) {
	g.operationChan <- func(connsMap connsType) {
		for remoteAddr, conn := range connsMap {
			// do not send to the self
			if remoteAddr == senderConn.RemoteAddr().String() {
				continue
			}
			err := conn.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Errorf("broadcast: send error: sender: %v, recipient: %v", senderConn.RemoteAddr().String(), remoteAddr)
				continue
			}
			log.Infof("broadcast: send success: sender: %v, recipient: %v, message: %v", senderConn.RemoteAddr().String(), remoteAddr, message)
		}
	}
}
