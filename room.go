package main

type room struct {
	// forwardは他のクライアントに転送するためのメッセージを保持するチャネルです。
	forward chan []byte
	// joinはチャットルームに参加しようとしているクライアントの為のチャネルです。
	join chan *client
	// leaveはチャットルームから体質しようとしているクライアントの為のチャネルです。
	leave chan *client
	// clientsには在室している全てのクライアントが保持されます。
	clients map[*client]bool
}

