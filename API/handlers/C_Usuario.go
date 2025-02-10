package handlers

import (
	"net/http"

	"math/rand"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"kings-house-back/API/models"
)

// Body de la peticion
type CrearUsuarioRequest struct {
	Nombre     string `json:"nombre" binding:"required"`
	Email      string `json:"email"`
	Contrasena string `json:"contrasena"`
	Rol        string `json:"rol" binding:"required"`
}

const usernameChars = "abcdefghijklmnopqrstuvwxyz0123456789"

func init() {
	// Seedea el generador de números aleatorios al iniciar el paquete
	rand.Seed(time.Now().UnixNano())
}

// GenerarUsername crea un username con base en el nombre y una parte aleatoria.
// name: "Samuel Llach"  => "S" + random(6) => "S0l0bpt"
func GenerarUsername(name string) string {
	// 1. Tomar la primera letra del nombre (en mayúscula o minúscula)
	firstLetter := ""
	name = strings.TrimSpace(name)
	if len(name) > 0 {
		firstLetter = strings.ToUpper(name[:1]) // "S"
	}

	// 2. Generar una parte aleatoria de longitud N
	randomPart := randomString(6) // ejemplo: "0l0bpt"

	// 3. Concatenar
	return firstLetter + randomPart
}

// randomString genera n caracteres aleatorios a partir de usernameChars
func randomString(n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = usernameChars[rand.Intn(len(usernameChars))]
	}
	return string(b)
}

func CrearUsuarioHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CrearUsuarioRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
			return
		}

		// Verificar rol
		switch req.Rol {
		case "administrador", "gestor":
			// Validar que el email y la contraseña no vengan vacíos
			if req.Email == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Email requerido para roles admin/gestor"})
				return
			}
			if req.Contrasena == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Contraseña requerida para roles admin/gestor"})
				return
			}

			// Hashear la contraseña
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Contrasena), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al encriptar la contraseña"})
				return
			}

			nuevoUsuario := models.Usuario{
				Nombre:     req.Nombre,
				Email:      req.Email, // email real
				Contrasena: string(hashedPassword),
				Rol:        req.Rol,
			}

			if err := db.Create(&nuevoUsuario).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el usuario"})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"mensaje":   "Usuario creado exitosamente",
				"usuarioID": nuevoUsuario.ID,
			})

		case "vendedor":
			// No debe mandar un email ni contraseña; las rechazamos si las envía
			if req.Email != "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "No se permite email para rol vendedor"})
				return
			}
			if req.Contrasena != "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "No se permite contraseña para rol vendedor"})
				return
			}

			// Generar un username aleatorio a partir del nombre
			autoUsername := GenerarUsername(req.Nombre)

			nuevoUsuario := models.Usuario{
				Nombre:     req.Nombre,
				Email:      autoUsername, // se usa como username
				Contrasena: "",           // sin password
				Rol:        "vendedor",
			}

			if err := db.Create(&nuevoUsuario).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el usuario"})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"mensaje":   "Usuario vendedor creado exitosamente",
				"usuarioID": nuevoUsuario.ID,
				"username":  autoUsername,
			})

		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Rol inválido"})
		}
	}
}
