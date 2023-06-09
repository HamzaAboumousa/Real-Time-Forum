package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

const (
	Init    = "init"
	Connect = "Connection"
	Typing  = "typing"
	Msg     = "msg"
)

type Messagedb struct {
	DB *sql.DB
}
type Commande struct {
	Type string          `json:"Type"`
	Data json.RawMessage `json:"Data"`
}
type Message struct {
	From    string `json:"From"`
	To      string `json:"To"`
	Date    string `json:"Date"`
	Content string `json:"Content"`
}
type Clientjs struct {
	Isconnected bool            `json:"Isconnected"`
	Pseudo      string          `json:"Pseudo"`
	Lastmessage json.RawMessage `json:"Lastmessage"`
	Allmessage  json.RawMessage `json:"Allmessage"`
}
type AllClients struct {
	ClientsList json.RawMessage `json:"ClientsList"`
	Visibility  bool            `json:"Visibility"`
}
type Typing_com struct {
	From string `json:"From"`
	To   string `json:"To"`
}

func bytetocommend(data []byte) Commande {
	var a Commande
	err := json.Unmarshal(data, &a)
	if err != nil {
		fmt.Println(err)
	}

	return a
}

func commandetobyte(c Commande) []byte {
	a, err := json.Marshal(c)
	if err != nil {
		fmt.Println(err)
	}
	return a
}
