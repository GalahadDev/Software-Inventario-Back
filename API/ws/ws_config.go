package ws

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Rol  string
	// Podrías almacenar userID si lo deseas
}

type Hub struct {
	// Lista de clientes conectados
	Clients map[*Client]bool

	// Mutex para proteger el acceso concurrente a Clients
	mu sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Clients: make(map[*Client]bool),
	}
}

// RegisterClient agrega un nuevo cliente al hub
func (h *Hub) RegisterClient(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.Clients[c] = true
	fmt.Println("Nuevo cliente registrado con rol:", c.Rol)
}

// UnregisterClient elimina el cliente del hub
func (h *Hub) UnregisterClient(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.Clients, c)
	c.Conn.Close()
	fmt.Println("Cliente desconectado con rol:", c.Rol)
}

// BroadcastMessage envía un mensaje a todos los clientes que cumplan cierta condición (ej. rol)
func (h *Hub) BroadcastMessage(msg string, rolesPermitidos ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for client := range h.Clients {
		if inArray(client.Rol, rolesPermitidos) {
			err := client.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				fmt.Println("Error al enviar mensaje:", err)
				client.Conn.Close()
				delete(h.Clients, client)
			}
		}
	}
}

func inArray(val string, arr []string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}
