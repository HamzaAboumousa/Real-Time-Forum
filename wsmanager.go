package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var websocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Manager struct {
	Client Clientliste
	sync.RWMutex
}

func newManager() *Manager {
	return &Manager{
		Client: make(Clientliste),
	}
}

func (m *Manager) Servws(w http.ResponseWriter, r *http.Request, s *Session) {
	// if r.Header["Sec-Fetch-Mode"][0] != "websocket" {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	w.Write([]byte("Bad request"))
	// 	return
	// }
	fmt.Println("new connection from:", strings.Split(r.RemoteAddr, ":")[0])
	// upgrade to ws
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	client := Newclient(conn, m, s)

	m.addClient(client)
	client.initco()
	// read and write
	go client.readMessage()
	go client.writeMessage()
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	m.Client[client] = true
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	for v := range client.manager.Client {
		if v.session.Username != client.session.Username {
			var conn Commande
			var cl Clientjs
			var err error
			var data []byte
			conn.Type = Connect
			cl.Pseudo = client.session.Username
			cl.Isconnected = false
			conn.Data, err = json.Marshal(cl)
			if err != nil {
				fmt.Println(err)
			}
			data, err = json.Marshal(conn)
			if err != nil {
				fmt.Println(err)
			}
			v.connection.WriteMessage(websocket.TextMessage, data)
		}
	}
	if _, ok := m.Client[client]; ok {
		client.connection.Close()
		delete(m.Client, client)
	}
}
