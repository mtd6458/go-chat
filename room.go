package main

import (
	"github.com/chat/trace"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type room struct {
	// forwardは他のクライアントに転送するためのメッセージを保持するチャネルです。
	forward chan []byte
	// joinはチャットルームに参加しようとしているクライアントの為のチャネルです。
	join chan *client
	// leaveはチャットルームから体質しようとしているクライアントの為のチャネルです。
	leave chan *client
	// clientsには在室している全てのクライアントが保持されます。
	clients map[*client]bool
	// tracerはチャットルーム上で行われた操作のログを受け取ります。
	tracer trace.Tracer
}

// newRoomはすぐに利用できるチャットルームを生成して返します。
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// 参加
			r.clients[client] = true
			r.tracer.Trace("新しいクライアントが参加しました")
		case client := <-r.leave:
			//退室
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("クライアントが退室しました")
		case msg := <-r.forward:
			r.tracer.Trace("メッセージを受信しました: ", string(msg))
			// 全てのクライアントにメッセージを転送
			for client := range r.clients {
				select {
				case client.send <- msg:
					//メッセージを送信
					r.tracer.Trace(" -- クライアントに送信されました")
				default:
					//送信に失敗
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace(" -- 送信に失敗しました。クライアントをクリーンアップします")
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(write http.ResponseWriter, request *http.Request) {
	socket, err := upgrader.Upgrade(write, request, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
	}

	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}

	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
