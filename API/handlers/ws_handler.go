package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"kings-house-back/API/ws"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Ajusta según tu política de orígenes
		return true
	},
}

// Ajustes de heartbeat
const (
	pongWait   = 60 * time.Second    // Cuánto esperamos sin recibir un pong
	pingPeriod = (pongWait * 9) / 10 // Cada cuánto enviamos ping
)

// WSHandler maneja la conexión WebSocket y configura ping/pong
func WSHandler(hub *ws.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Verificar claims en el contexto (AuthMiddleware)
		claimsVal, existe := c.Get("claims")
		if !existe {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No se encontraron claims"})
			return
		}
		claims, ok := claimsVal.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Claims inválidos"})
			return
		}
		rol, ok := claims["rol"].(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Rol no presente en token"})
			return
		}

		// 2. Hacer upgrade a WebSocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println("Error al convertir a WebSocket:", err)
			return
		}

		// 3. Configurar deadlines y pong handler
		conn.SetReadDeadline(time.Now().Add(pongWait))
		conn.SetPongHandler(func(appData string) error {
			// Al recibir PONG, renovamos el ReadDeadline
			conn.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})

		// 4. Iniciar goroutine que envía pings periódicos
		go func() {
			ticker := time.NewTicker(pingPeriod)
			defer ticker.Stop()

			for {
				<-ticker.C
				// Enviar ping
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Println("Error al enviar ping:", err)
					conn.Close()
					return
				}
			}
		}()

		// 5. Registrar el cliente en tu hub
		client := &ws.Client{
			Conn: conn,
			Rol:  rol,
		}
		hub.RegisterClient(client)

		// 6. Bucle para leer mensajes desde el cliente
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				// Cuando ocurre un error (desconexión, timeout, etc.), desregistramos
				hub.UnregisterClient(client)
				break
			}
			// Si deseas, puedes renovar Deadline aquí también:
			conn.SetReadDeadline(time.Now().Add(pongWait))

			// Manejar mensajes entrantes si es necesario...
		}
	}
}
