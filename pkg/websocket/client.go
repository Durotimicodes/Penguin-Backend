package websocket

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
	Mu   sync.Mutex
}

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	//read all the messages from the connection
	for {
		msgType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		message := Message{
			Type: msgType,
			Body: string(p),
		}
		c.Pool.Broadcast <- message
		fmt.Printf("Message received:%+v\n", message)
	}

}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("size of the connection pool:", len(pool.Clients))
			for client, _ := range pool.Clients {
				fmt.Println(client)
				client.Conn.WriteJSON(Message{
					Type: 1,
					Body: "New User has Joined....",
				})
			} 
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("size of the connection pool:", len(pool.Clients))
			for client, _ := range pool.Clients {
				client.Conn.WriteJSON(Message{
					Type: 1,
					Body: "User Disconnected...",
				})
			}
		case message := <-pool.Broadcast:
			fmt.Println("Sending message to all clients in the pool")
			for client, _ := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err!= nil {
					fmt.Println(err)
					return
				}
			}

		}
	}
}
