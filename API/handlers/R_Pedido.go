package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"kings-house-back/API/models"
)

func ListarPedidosHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var pedidos []models.Pedido
		if err := db.Find(&pedidos).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener la lista de pedidos"})
			return
		}
		c.JSON(http.StatusOK, pedidos)
	}
}

func ObtenerPedidoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
			return
		}

		var pedido models.Pedido
		if err := db.Preload("Usuario").First(&pedido, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Pedido no encontrado"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar el pedido"})
			}
			return
		}

		c.JSON(http.StatusOK, pedido)
	}
}

func ListarPedidosPorUsuarioHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		usuarioID := c.Param("usuario_id")

		var pedidos []models.Pedido
		if err := db.Where("usuario_id = ?", usuarioID).Find(&pedidos).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener los pedidos"})
			return
		}
		c.JSON(http.StatusOK, pedidos)
	}
}
