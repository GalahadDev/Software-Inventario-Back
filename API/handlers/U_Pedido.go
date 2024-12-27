package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"kings-house-back/API/models"
)

type ActualizarPedidoRequest struct {
	Descripcion string   `json:"descripcion"`
	Imagen      string   `json:"imagen"`
	Fletero     *string  `json:"fletero"`
	Monto       *float64 `json:"monto"`
	Estado      string   `json:"estado"`
	Precio      *float64 `json:"precio"`
}

func ActualizarPedidoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
			return
		}

		var req ActualizarPedidoRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
			return
		}

		var pedido models.Pedido
		if err := db.First(&pedido, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Pedido no encontrado"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar el pedido"})
			}
			return
		}

		// Actualizar los campos que vengan en el request
		if req.Descripcion != "" {
			pedido.Descripcion = req.Descripcion
		}
		if req.Imagen != "" {
			pedido.Imagen = req.Imagen
		}
		if req.Fletero != nil {
			pedido.Fletero = req.Fletero
		}
		if req.Monto != nil {
			pedido.Monto = req.Monto
		}
		if req.Estado != "" {
			pedido.Estado = req.Estado
		}
		if req.Precio != nil {
			pedido.Precio = req.Precio
		}

		if err := db.Save(&pedido).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar el pedido"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"mensaje": "Pedido actualizado correctamente"})
	}
}
