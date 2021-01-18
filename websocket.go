package main

import (
	"net"
	"net/http"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/google/uuid"
)

var wsOnce sync.Once
var wsRwm sync.RWMutex

// WebsocketServer struct
type WebsocketServer struct {
	connections map[[16]byte]*net.Conn
}

var websocketInstance *WebsocketServer

// AddConnection adds new conection from http handler
func (server WebsocketServer) AddConnection(w http.ResponseWriter, r *http.Request, f func([]byte)) ([16]byte, error) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)

	var idConnection [16]byte

	if err != nil {
		return idConnection, err
	}

	wsRwm.Lock()
	defer wsRwm.Unlock()

	idConnection = uuid.New()

	websocketInstance.connections[idConnection] = &conn

	go func() {
		for {
			msg, _, err := wsutil.ReadClientData(conn)
			if err == nil {
				f(msg)
			} else {
				wsRwm.Lock()
				defer wsRwm.Unlock()
				delete(websocketInstance.connections, idConnection)
				return
			}
		}
	}()

	return idConnection, nil
}

// InitWebsocketServer initialize the websocket server
func InitWebsocketServer() error {

	if websocketInstance == nil {
		wsOnce.Do(func() {
			websocketInstance = &WebsocketServer{
				connections: make(map[[16]byte]*net.Conn),
			}
		})
	}

	return nil
}

// GetWebsocketServerInstance returns websocket server instance
func GetWebsocketServerInstance() *WebsocketServer {
	return websocketInstance
}

// BroadcastMessage broadcast message to all connections
func (server WebsocketServer) BroadcastMessage(body []byte) {
	wsRwm.Lock()
	defer wsRwm.Unlock()

	for _, v := range server.connections {
		wsutil.WriteServerMessage(*v, ws.OpText, body)
	}
}

// SendMessage send message to individual connection
func (server WebsocketServer) SendMessage(body []byte, idConnection [16]byte) {
	wsRwm.Lock()
	defer wsRwm.Unlock()

	connection, exists := server.connections[idConnection]

	if exists {
		wsutil.WriteServerMessage(*connection, ws.OpText, body)
	}
}
