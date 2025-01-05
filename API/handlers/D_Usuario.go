package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"kings-house-back/API/models"
)

func EliminarUsuarioHandler(db *gorm.DB) gin.HandlerFunc {
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

		if err := db.Delete(&usuario).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar el usuario"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"mensaje": "Usuario eliminado correctamente"})
	}
}
