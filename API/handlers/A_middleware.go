package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Buscar token en encabezado Authorization: "Bearer <token>"
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No se proporcionó token"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido"})
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return secret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			return
		}

		// Extraer los claims y pasarlos al contexto
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("claims", claims)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			return
		}

		c.Next()
	}
}

// AuthMiddlewareQuery es la NUEVA función que lee el token de la URL "?token=..."
func AuthMiddlewareQuery(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Obtener token de la query param
		tokenString := c.Query("token")
		if tokenString == "" {
			// Si no hay token en la query, devolvemos 401
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No se proporcionó token en la URL"})
			return
		}

		// 2. Parsear el token
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido o firma incorrecta"})
			return
		}

		// 3. Extraer claims y guardarlos en el contexto
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("claims", claims)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Claims inválidos"})
			return
		}

		c.Next()
	}
}
