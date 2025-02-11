package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"kings-house-back/API/models"
	"kings-house-back/API/ws"
)

func ActualizarPedidoHandler(db *gorm.DB, hub *ws.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
			return
		}

		// Leer campos de form-data
		descripcion := c.PostForm("descripcion")
		fletero := c.PostForm("fletero")
		montoStr := c.PostForm("monto")
		estado := c.PostForm("estado")
		precioStr := c.PostForm("precio")
		nombre := c.PostForm("nombre")
		observaciones := c.PostForm("observaciones")
		formaPago := c.PostForm("forma_pago")
		direccion := c.PostForm("direccion")
		numerotlf := c.PostForm("nro_tlf")
		pagado := c.PostForm("pagado")
		atendidoStr := c.PostForm("atendido")

		// Parsear monto y precio
		var montoFloat *float64
		if montoStr != "" {
			if m, err := strconv.ParseFloat(montoStr, 64); err == nil {
				montoFloat = &m
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Monto inválido"})
				return
			}
		}

		var precioFloat *float64
		if precioStr != "" {
			if p, err := strconv.ParseFloat(precioStr, 64); err == nil {
				precioFloat = &p
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Precio inválido"})
				return
			}
		}

		// Manejo imagen (opcional)
		file, errFile := c.FormFile("imagen")
		var imagenRuta string
		if errFile == nil {
			rutaArchivo := "./uploads/" + file.Filename
			if err := c.SaveUploadedFile(file, rutaArchivo); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar la imagen"})
				return
			}
			imagenRuta = rutaArchivo
		} else {
			log.Printf("No se recibió nueva imagen o error al recibir imagen: %v", errFile)
		}

		// Buscar el pedido existente
		var pedido models.Pedido
		if err := db.First(&pedido, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Pedido no encontrado"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar el pedido"})
			}
			return
		}

		// Actualizar campos
		if descripcion != "" {
			pedido.Descripcion = descripcion
		}
		if imagenRuta != "" {
			pedido.Imagen = imagenRuta
		}
		if fletero != "" {
			pedido.Fletero = &fletero
		}
		if montoFloat != nil {
			pedido.Monto = montoFloat
		}
		if estado != "" {
			pedido.Estado = estado
		}
		if precioFloat != nil {
			pedido.Precio = precioFloat
		}
		if nombre != "" {
			pedido.Nombre = nombre
		}
		if observaciones != "" {
			pedido.Observaciones = observaciones
		}
		if formaPago != "" {
			pedido.Forma_Pago = formaPago
		}
		if direccion != "" {
			pedido.Direccion = direccion
		}
		if numerotlf != "" {
			pedido.Nro_Tlf = numerotlf
		}
		if pagado != "" {
			pedido.Pagado = pagado
		}
		if atendidoStr != "" {
			pedido.Atendido = true
		}

		// Guardar en la BD
		if err := db.Save(&pedido).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar el pedido"})
			return
		}

		// ---- ENVIAR NOTIFICACIÓN WS: "PEDIDO_ACTUALIZADO" ----
		notif := NotificacionPedido{
			Tipo:   "PEDIDO_ACTUALIZADO",
			Pedido: pedido,
		}
		notifBytes, _ := json.Marshal(notif)
		// Avisamos a admin/gestor (o a quien quieras)
		hub.BroadcastMessage(string(notifBytes), "administrador", "gestor")

		c.JSON(http.StatusOK, gin.H{"mensaje": "Pedido actualizado correctamente"})
	}
}
