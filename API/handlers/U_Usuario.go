package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"kings-house-back/API/models"

	"github.com/golang-jwt/jwt/v5"
)

type ActualizarUsuarioRequest struct {
	Nombre        string `json:"nombre"`
	Rol           string `json:"rol"`
	Contrasena    string `json:"contrasena"`
	Email         string `json:"email"`
	Cedula        string `json:"cedula"`
	Numero_Cuenta string `json:"numero_cuenta"`
	Tipo_Cuenta   string `json:"tipo_cuenta"`
	Nombre_Banco  string `json:"nombre_banco"`
}

type ActualizarDatosBancariosRequest struct {
	Cedula        string `json:"cedula"`
	Numero_Cuenta string `json:"numero_cuenta"`
	Tipo_Cuenta   string `json:"tipo_cuenta"`
	Nombre_Banco  string `json:"nombre_banco"`
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
		if req.Cedula != "" {
			usuario.Cedula = req.Cedula
		}
		if req.Numero_Cuenta != "" {
			usuario.Numero_Cuenta = req.Numero_Cuenta
		}
		if req.Tipo_Cuenta != "" {
			usuario.Tipo_Cuenta = req.Tipo_Cuenta
		}
		if req.Nombre_Banco != "" {
			usuario.Nombre_Banco = req.Nombre_Banco
		}

		if err := db.Save(&usuario).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar el usuario"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"mensaje": "Usuario actualizado correctamente"})
	}
}

func ActualizarDatosBancariosHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Extraer claims del JWT para saber rol y userID
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

		// 2. Verificar que rol sea "vendedor"
		if rol != "vendedor" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Sólo los vendedores pueden editar sus datos bancarios"})
			return
		}

		// 3. Parsear el body JSON con los campos bancarios
		var req ActualizarDatosBancariosRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
			return
		}

		// 4. Buscar al usuario en la BD por su propio ID
		var usuario models.Usuario
		if err := db.First(&usuario, "id = ?", userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
			return
		}

		// 5. Actualizar sólo campos bancarios si vienen
		if req.Cedula != "" {
			usuario.Cedula = req.Cedula
		}
		if req.Numero_Cuenta != "" {
			usuario.Numero_Cuenta = req.Numero_Cuenta
		}
		if req.Tipo_Cuenta != "" {
			usuario.Tipo_Cuenta = req.Tipo_Cuenta
		}
		if req.Nombre_Banco != "" {
			usuario.Nombre_Banco = req.Nombre_Banco
		}

		// 6. Guardar
		if err := db.Save(&usuario).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar datos bancarios"})
			return
		}

		// 7. Respuesta exitosa
		c.JSON(http.StatusOK, gin.H{
			"mensaje": "Datos bancarios actualizados correctamente",
		})
	}
}
