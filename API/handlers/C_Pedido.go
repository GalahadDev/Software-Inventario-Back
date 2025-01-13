package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"kings-house-back/API/database"
	"kings-house-back/API/models"
)

// CrearPedidoHandler maneja la creación de un pedido.
func CrearPedidoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Leer los campos de form-data
		usuarioID := c.PostForm("usuario_id")
		descripcion := c.PostForm("descripcion")
		nombre := c.PostForm("nombre")
		observaciones := c.PostForm("observaciones")
		formaPago := c.PostForm("forma_pago")
		direccion := c.PostForm("direccion")

		// Parsear el precio
		precioStr := c.PostForm("precio")
		var precioFloat *float64
		if precioStr != "" {
			if p, err := strconv.ParseFloat(precioStr, 64); err == nil {
				precioFloat = &p
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Precio inválido"})
				return
			}
		}

		fileHeader, err := c.FormFile("imagen")
		var publicURL string

		if err == nil {
			// El usuario envió un archivo
			bucketName := "imagenes-pedidos"
			// Ejemplo: pedidos/<timestamp>_<nombreArchivo>
			filePath := fmt.Sprintf("pedidos/%d_%s", time.Now().Unix(), filepath.Base(fileHeader.Filename))

			// 2. Subir el archivo a Supabase
			publicURL, err = database.SubirAStorageSupabase(fileHeader, bucketName, filePath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al subir a Supabase", "details": err.Error()})
				return
			}
		}

		// 4. Crear el objeto Pedido
		nuevoPedido := models.Pedido{
			UsuarioID:     usuarioID,
			Descripcion:   descripcion,
			Imagen:        publicURL,
			FechaCreacion: time.Now(),
			Precio:        precioFloat,
			Fletero:       nil,
			Monto:         nil,
			Estado:        "No Entregado",

			Nombre:        nombre,
			Observaciones: observaciones,
			Forma_Pago:    formaPago,
			Direccion:     direccion,
		}

		// 5. Guardar en la base de datos
		if err := db.Create(&nuevoPedido).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el pedido"})
			return
		}

		// 6. Respuesta exitosa
		c.JSON(http.StatusOK, gin.H{
			"mensaje":    "Pedido creado exitosamente",
			"pedido_id":  nuevoPedido.ID,
			"usuario_id": nuevoPedido.UsuarioID,
			"imagen":     nuevoPedido.Imagen,
		})
	}
}
