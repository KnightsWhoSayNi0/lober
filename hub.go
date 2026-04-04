package main

import (
	"encoding/json"
	"log"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	filter string
	send   chan []byte
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *Client
	mu         sync.Mutex
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *Client),
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Message is the filter string
		filter := string(message)
		c.filter = strings.TrimSpace(filter)
		log.Printf("Client filter updated: %q", c.filter)

		// When filter updates, send the current filtered dataset
		events := getEventsWithFilter(c.filter)
		eventsJSON, _ := json.Marshal(events)
		c.send <- eventsJSON
	}
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}
		}
	}
}

func (h *Hub) run() {
	for {
		select {
		case conn := <-h.register:
			client := &Client{hub: h, conn: conn, send: make(chan []byte, 256)}
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("New client registered")
			go client.writePump()
			go client.readPump()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("Client unregistered")
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			var event Event
			err := json.Unmarshal(message, &event)
			if err != nil {
				log.Printf("Error unmarshaling broadcast message: %v", err)
				continue
			}

			h.mu.Lock()
			for client := range h.clients {
				// Check if event matches client filter
				if eventMatchesFilter(event, client.filter) {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
			h.mu.Unlock()
		}
	}
}

func eventMatchesFilter(e Event, filter string) bool {
	if filter == "" {
		return true
	}
	f := strings.ToLower(filter)
	return strings.Contains(strings.ToLower(e.Command), f) ||
		strings.Contains(strings.ToLower(e.User), f) ||
		strings.Contains(strings.ToLower(e.C2), f) ||
		strings.Contains(strings.ToLower(e.Scope), f)
}

func (h *Hub) broadcastData(data []byte) {
	h.broadcast <- data
}
