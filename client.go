package main 

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second
	pongWait = 60 * time.Second
	pingPeriod = (pongWait  * 9) / 10
	maxMsgSize = 512
)

var (
	newLine = []byte{"\n"}
	space = []byte{" "}
)

var upgrader = websocket.upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

type Client struct{
	hub *Hub
	conn *websocket.Conn
	send chan []byte
}

func (c *Client) readPump() {
	defer func(){
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMsgSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler( func(string) error { c.conn.SetReadDeadline(time.Now().Add(PongWait)); return nil})

	for{
		_, msg, err := c.conn.ReadMessage()

		if err != nil{
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure){
				log.Printf("error: %v",err)
			}
			break
		}
		msg = bytes.TrimSpace(bytes.Replace(msg, newline, space, -1))
		c.hub.broadcast <- msg
	}
}

func (c *Client) writePump(){
	timer := time.NewTicker(pingPeriod)
	defer func() {
		timer.Stop()
		c.conn.Close()
	}()

	for {
		select{
		case message, ok := <- c.send:
			

		}
	}
}
