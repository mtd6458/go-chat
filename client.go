package main

import (
	"github.com/gorilla/websocket"
	"log"
)

// Clientはチャットを行なっている1人のユーザを表します
type client struct {
	// socketはこのクライアントの為のWebSocketです。
	socket *websocket.Conn
	// sendはメッセージが送られるチャネルです。
	send chan []byte
	// roomはこのクライアントが参加しているチャットルームです。
	room *room
}

func (c *client) read() {
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			break
		}
		c.room.forward <- msg
	}
	err := c.socket.Close()
	if err != nil {
		log.Fatal("socket close:", err)
	}
}

func (c *client) write() {
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
	err := c.socket.Close()
	if err != nil {
		log.Fatal("socket close:", err)
	}
}
