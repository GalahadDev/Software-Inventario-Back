package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"kings-house-back/API/models"
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
