package main

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

type Hub struct {
	broadcast  chan []byte
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	publish chan []byte
	redisClient *redis.Client
}

func createHub() *Hub {

	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		redisUrl = "redis://localhost:6379"
	}

	opt, err := redis.ParseURL(redisUrl)

	rdb := redis.NewClient(opt)

	if err != nil {
		panic(err)
	}

	return &Hub{
		broadcast:   make(chan []byte),
		clients:     make(map[*Client]bool),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		publish:     make(chan []byte),
		redisClient: rdb,
	}
}

func (h *Hub) run() {

	go h.subscribeToRedis()

	for {
		select {

		case client := <-h.register:
			h.clients[client] = true

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		case message := <-h.publish:
			err := h.redisClient.Publish(context.Background(), "chat_room", message).Err()
			if err != nil{
				fmt.Println("Error While doing Redis Publish: ",err)
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) subscribeToRedis() {
	pubsub := h.redisClient.Subscribe(context.Background(), "chat_room")

	defer pubsub.Close()

	ch := pubsub.Channel()
	for msg := range ch{
		h.broadcast <- []byte(msg.Payload)
	}


}
