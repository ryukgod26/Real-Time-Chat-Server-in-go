package main



type Hub struct{
	broadcast chan []byte
	clients map[*Client] bool
	register chan *Client
	unregister chan *Client
}

func createHub() *Hub{
	return &Hub{
		broadcast: make(chan []byte),
		clients: make(map[*Client]bool),
		register: make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run(){
	for{
		select{

		case client := <-h.register:
			h.clients[client] = true
		
		case client := <-h.unregister:
			if _,ok := h.clients[client]; ok{
				delete(h.clients, client)
				close(client.send)
			}
		
		case message := <-h.broadcast:

		}
	}
}