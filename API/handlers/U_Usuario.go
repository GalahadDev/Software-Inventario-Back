package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"kings-house-back/API/models"
)

type ActualizarUsuarioRequest struct {
	Nombre     string `json:"nombre"`
	Rol        string `json:"rol"`
	Contrasena string `json:"contrasena"`
	Email      string `json:"email"`
}

func ActualizarUsuarioHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var req ActualizarUsuarioRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
			return
		}

		var usuario models.Usuario
		if err := db.First(&usuario, "id = ?", id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar el usuario"})
			}
			return
		}

		// Actualizamos sólo los campos que vengan con valor
		if req.Nombre != "" {
			usuario.Nombre = req.Nombre
		}
		if req.Rol != "" {
			usuario.Rol = req.Rol
		}
		if req.Email != "" {
			usuario.Email = req.Email
		}
		if req.Contrasena != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Contrasena), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al encriptar la contraseña"})
				return
			}
			usuario.Contrasena = string(hashedPassword)
		}

		if err := db.Save(&usuario).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar el usuario"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"mensaje": "Usuario actualizado correctamente"})
	}
}
