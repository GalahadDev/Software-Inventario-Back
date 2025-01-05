package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"kings-house-back/API/models"
)

type CrearPedidoRequest struct {
	UsuarioID     string   `json:"usuario_id" binding:"required"`
	Descripcion   string   `json:"descripcion" binding:"required"`
	Imagen        string   `json:"imagen"`
	Precio        *float64 `json:"precio"`
	Nombre        string   `json:"nombre" binding:"required"`
	Observaciones string   `json:"observaciones"`
	Forma_Pago    string   `json:"forma_pago"`
	Direccion     string   `json:"direccion"`
}

// CrearPedidoHandler maneja la creación de un pedido.
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
			Fletero:       nil,
			Monto:         nil,
			Estado:        "No Entregado",
			Nombre:        req.Nombre,
			Observaciones: req.Observaciones,
			Forma_Pago:    req.Forma_Pago,
			Direccion:     req.Direccion,
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
