package main

import (
	"github.com/gorilla/websocket"
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


