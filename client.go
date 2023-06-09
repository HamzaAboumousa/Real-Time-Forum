package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Clientliste map[*Client]bool

type Client struct {
	connection *websocket.Conn
	manager    Manager
	session    *Session
	// egress is a chan to avoid concurrent writes on ws conn
	egress chan Commande
}

func Newclient(conn *websocket.Conn, manager *Manager, s *Session) *Client {
	return &Client{
		connection: conn,
		manager:    *manager,
		session:    s,
		egress:     make(chan Commande),
	}
}

func (c *Client) readMessage() {
	defer func() {
		// clean conn
		c.manager.removeClient(c)
		log.Println("connection closed: ")
	}()
	for {
		_, payload, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway) {
				log.Printf("Error reading message %v\n", err)
			}
			break
		}
		c.egress <- bytetocommend(payload)
	}
}

func (c *Client) writeMessage() {
	defer func() {
		c.manager.removeClient(c)
		log.Println("connection closed: ")
	}()

	for {
		select {
		case value, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed: ", err)
				}
				return
			}
			switch value.Type {
			case Msg:
				var message Message

				err := json.Unmarshal(value.Data, &message)
				if err != nil {
					fmt.Println(err)
				}
				if message.Content != "" {
					message.From = c.session.Username
					value.Data, _ = json.Marshal(message)
					Msgdb.Add(message)
					for v := range c.manager.Client {
						if v.session.Username == message.To {
							v.connection.WriteMessage(websocket.TextMessage, commandetobyte(value))
						}
					}
				}
				break
			case Typing:
				var Typ Typing_com

				err := json.Unmarshal(value.Data, &Typ)
				if err != nil {
					fmt.Println(err)
				}
				Typ.From = c.session.Username
				value.Data, _ = json.Marshal(Typ)
				for v := range c.manager.Client {
					if v.session.Username == Typ.To {
						v.connection.WriteMessage(websocket.TextMessage, commandetobyte(value))
					}
				}

				break
			}
		}
	}
}

func (c *Client) initco() {
	user := c.session.Username
	Allmessages := Msgdb.Get(user)
	usersitems := guys.Get()
	var All AllClients
	All.Visibility = true
	var rawclients []Clientjs

	for _, value := range usersitems {
		if value.Username != user {
			var client1 Clientjs
			var allmsg []Message
			client1.Isconnected = false
			for _, messag := range Allmessages {
				if (messag.From == user && messag.To == value.Username) || (messag.To == user && messag.From == value.Username) {
					allmsg = append(allmsg, messag)
				}
			}
			client1.Allmessage, _ = json.Marshal(allmsg)
			client1.Pseudo = value.Username
			for c, v := range c.manager.Client {
				if c.session.Username == value.Username {
					if v {
						client1.Isconnected = true
						break
					}
				}
			}
			if len(allmsg) > 0 {
				client1.Lastmessage, _ = json.Marshal(allmsg[len(allmsg)-1])
			}
			rawclients = append(rawclients, client1)
		}
	}
	var errr error
	All.ClientsList, errr = json.Marshal(rawclients)
	if errr != nil {
		fmt.Println(errr)
	}
	alll, err := json.Marshal(All)
	if err != nil {
		fmt.Println(err)
	}
	var a Commande
	a.Type = Init
	a.Data = alll
	Data, err := json.Marshal(a)
	if err != nil {
		fmt.Println(err)
	}
	c.connection.WriteMessage(websocket.TextMessage, Data)
	for v := range c.manager.Client {
		if v.session.Username != c.session.Username {
			var conn Commande
			var cl Clientjs
			var err error
			var data []byte
			conn.Type = Connect
			cl.Pseudo = c.session.Username
			cl.Isconnected = true
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
}
