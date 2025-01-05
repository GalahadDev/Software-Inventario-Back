package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"kings-house-back/API/models"

	"gorm.io/gorm"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type Credenciales struct {
	Email      string `json:"email" binding:"required"`
	Contrasena string `json:"contrasena" binding:"required"`
}

// LoginHandler procesa el login y genera un token JWT
func LoginHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var creds Credenciales
		if err := c.ShouldBindJSON(&creds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
			return
		}

		var usuario models.Usuario
		// Buscar usuario por Email
		if err := db.Where("email = ?", creds.Email).First(&usuario).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
			return
		}

		// Verificar contraseña hasheada
		if err := bcrypt.CompareHashAndPassword([]byte(usuario.Contrasena), []byte(creds.Contrasena)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
			return
		}

		//Creacion un token JWT
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"usuario_id": usuario.ID,
			"rol":        usuario.Rol,
			"exp":        time.Now().Add(time.Hour * 24).Unix(), // 1 día de validez
		})

		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo generar el token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"mensaje":    "Inicio de sesión exitoso",
			"token":      tokenString,
			"usuario_id": usuario.ID,
			"nombre":     usuario.Nombre,
		})
	}
}
