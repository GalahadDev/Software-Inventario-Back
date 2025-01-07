package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

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

		// 3. Manejar el archivo (imagen) si existe
		file, err := c.FormFile("imagen")
		var imagenRuta string
		if err == nil {
			// Se subió un archivo con la clave "imagen"
			// Ejemplo: guardar localmente en la carpeta "uploads"
			// (En producción, lo normal es subir a S3 u otro servicio)
			rutaArchivo := "./uploads/" + file.Filename
			if err := c.SaveUploadedFile(file, rutaArchivo); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar la imagen"})
				return
			}
			imagenRuta = rutaArchivo
		} else {
			// Si no se mandó archivo o hay error, puedes ignorar o manejarlo
			log.Printf("No se recibió archivo o error al recibir imagen: %v", err)
		}

		// 4. Crear el objeto Pedido
		nuevoPedido := models.Pedido{
			UsuarioID:     usuarioID,
			Descripcion:   descripcion,
			Imagen:        imagenRuta, // Guardamos la ruta donde está el archivo
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
