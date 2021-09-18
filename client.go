package main

import (
	"github.com/gorilla/websocket"
	"log"
)

// Clientはチャットを行なっている1人のユーザを表します
type client struct {
	// socketはクライアントと通信を行うためのWebSocketを参照する。
	socket *websocket.Conn
	// sendはメッセージが送られるチャネルです。
	// 受信したメッセージが待ち行列のように蓄積され、WebSocketを通じてユーザのブラウザに送られるのを待機する。
	send chan []byte
	// roomはこのクライアントが参加しているチャットルームへの参照が保持される。
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
