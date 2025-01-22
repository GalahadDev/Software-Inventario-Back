package handlers

import (
	"fmt"
	"net/http"

	"kings-house-back/API/ws"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Asegúrate de permitir orígenes adecuados o filtrar por CORS
		return true
	},
}

func WSHandler(hub *ws.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Revisar si hay claims en el contexto (AuthMiddleware)
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

		// "upgrade" la conexión HTTP a WebSocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println("Error al convertir a WebSocket:", err)
			return
		}

		client := &ws.Client{
			Conn: conn,
			Rol:  rol,
		}
		hub.RegisterClient(client)

		// Leer en un loop hasta que se cierre
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				// Cuando se produce un error, desconectamos
				hub.UnregisterClient(client)
				break
			}
			// Puedes manejar mensajes entrantes si es necesario
		}
	}
}
