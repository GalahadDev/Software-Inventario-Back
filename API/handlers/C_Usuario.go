package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"kings-house-back/API/models"
)

// Body de la peticion
type CrearUsuarioRequest struct {
	Nombre     string `json:"nombre" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Contrasena string `json:"contrasena" binding:"required"`
	Rol        string `json:"rol" binding:"required"`
}

// Creación de un usuario
func CrearUsuarioHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CrearUsuarioRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
			return
		}

		// Hashear la contraseña antes de guardar
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Contrasena), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al encriptar la contraseña"})
			return
		}

		nuevoUsuario := models.Usuario{
			Nombre:     req.Nombre,
			Email:      req.Email,
			Contrasena: string(hashedPassword),
			Rol:        req.Rol,
		}

		// Intentar crear el usuario en la BD
		if err := db.Create(&nuevoUsuario).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el usuario"})
			return
		}

		// Respuesta exitosa
		c.JSON(http.StatusOK, gin.H{
			"mensaje":   "Usuario creado exitosamente",
			"usuarioID": nuevoUsuario.ID,
		})
	}
}
