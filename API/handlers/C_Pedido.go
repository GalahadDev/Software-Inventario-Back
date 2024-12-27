package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"kings-house-back/API/models"
)

// Estructura del pedido
type CrearPedidoRequest struct {
	UsuarioID   uint     `json:"usuario_id" binding:"required"`
	Descripcion string   `json:"descripcion" binding:"required"`
	Imagen      string   `json:"imagen"`
	Precio      *float64 `json:"precio"`
}

// Creación de un pedido.
func CrearPedidoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CrearPedidoRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
			return
		}

		nuevoPedido := models.Pedido{
			UsuarioID:     req.UsuarioID,
			Descripcion:   req.Descripcion,
			Imagen:        req.Imagen,
			FechaCreacion: time.Now(),
			Precio:        req.Precio,
			Fletero:       nil,            // Por defecto en nil
			Monto:         nil,            // Por defecto en nil
			Estado:        "No Entregado", // Estado inicial
		}

		if err := db.Create(&nuevoPedido).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el pedido"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"mensaje":    "Pedido creado exitosamente",
			"pedido_id":  nuevoPedido.ID,
			"usuario_id": nuevoPedido.UsuarioID,
		})
	}
}
