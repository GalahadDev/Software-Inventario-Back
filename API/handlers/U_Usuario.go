package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"kings-house-back/API/models"
)

type ActualizarUsuarioRequest struct {
	Nombre string `json:"nombre"`
	Rol    string `json:"rol"`
}

func ActualizarUsuarioHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
			return
		}

		var req ActualizarUsuarioRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
			return
		}

		var usuario models.Usuario
		if err := db.First(&usuario, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar el usuario"})
			}
			return
		}

		if req.Nombre != "" {
			usuario.Nombre = req.Nombre
		}
		if req.Rol != "" {
			usuario.Rol = req.Rol
		}

		if err := db.Save(&usuario).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar el usuario"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"mensaje": "Usuario actualizado correctamente"})
	}
}
