package main

import (
	"fmt"
	"net/http"

	"github.com/durotimicodes/penguine-chatapp/pkg/websocket"
)


func serveWebsocket(pool *websocket.Pool, w http.ResponseWriter, r *http.Request ){
	fmt.Println("Websocket endpoint reached")

	conn, err := websocket.Upgrade(w,r)

	if err != nil {
		fmt.Fprintf(w, "%+V/n", err)
	}
	
	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
}

func setupRoutes() {
	//create a new websocket
	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWebsocket(pool, w, r)
	})

}

func main() {
	fmt.Println("Start Penguine project")
	setupRoutes()
	http.ListenAndServe(":9000", nil)

}
