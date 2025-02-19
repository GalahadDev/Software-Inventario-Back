package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"kings-house-back/API/models"

	"github.com/golang-jwt/jwt/v5"
)

// Devuelve todos los usuarios
func ListarUsuariosHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var usuarios []models.Usuario
		if err := db.Find(&usuarios).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener la lista de usuarios"})
			return
		}
		c.JSON(http.StatusOK, usuarios)
	}
}

// Devuelve un usuario por su ID
func ObtenerUsuarioHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var usuario models.Usuario
		if err := db.First(&usuario, "id = ?", id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar el usuario"})
			}
			return
		}

		c.JSON(http.StatusOK, usuario)
	}
}

// Devuelve todos los usuarios que son vendedores
func ListarVendedoresHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var vendedores []models.Usuario
		if err := db.Where("rol = ?", "vendedor").Find(&vendedores).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener la lista de vendedores"})
			return
		}
		c.JSON(http.StatusOK, vendedores)
	}
}

func ObtenerDatosVendedorHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Extraer claims del token (rol y usuario_id)
		claimsVal, existe := c.Get("claims")
		if !existe {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No se encontraron claims en el token"})
			return
		}
		claims, ok := claimsVal.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Claims inválidos"})
			return
		}

		rol, ok := claims["rol"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No se encontró rol en el token"})
			return
		}
		userID, ok := claims["usuario_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No se encontró usuario_id en el token"})
			return
		}

		// 2. Verificar que sea rol vendedor
		if rol != "vendedor" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Sólo un vendedor puede acceder a sus datos personales"})
			return
		}

		// 3. Buscar al usuario por su propio ID
		var usuario models.Usuario
		if err := db.First(&usuario, "id = ?", userID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar datos del usuario"})
			}
			return
		}

		// 4. Retornar la información del vendedor
		// Si deseas ocultar o filtrar algunos campos sensibles, crea un struct de respuesta
		c.JSON(http.StatusOK, usuario)
	}
}
