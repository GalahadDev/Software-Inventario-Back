package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// RoleMiddleware se encarga de verificar que el rol del usuario
func RoleMiddleware(rolesPermitidos ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Obtener los claims del token (establecidos en AuthMiddleware).
		claimsVal, existe := c.Get("claims")
		if !existe {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "No se encontraron claims en el contexto",
			})
			return
		}

		// 2. Verificar que sean del tipo correcto
		claims, ok := claimsVal.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Claims inválidos",
			})
			return
		}

		// 3. Extraer el rol de los claims
		rol, ok := claims["rol"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "No se encontró rol en el token",
			})
			return
		}

		// 4. Verificar si ese rol está en la lista de rolesPermitidos
		for _, rolPermitido := range rolesPermitidos {
			if rol == rolPermitido {
				// Rol coincide, continuar con la ruta
				c.Next()
				return
			}
		}

		// Si no coincide con ninguno, rechazar
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "No tienes permisos para acceder a esta ruta",
		})
	}
}
