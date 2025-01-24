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
	"kings-house-back/API/ws"
)

// CrearPedidoHandler maneja la creación de un pedido.
func CrearPedidoHandler(db *gorm.DB, hub *ws.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. Leer los campos de form-data
		usuarioID := c.PostForm("usuario_id")
		descripcion := c.PostForm("descripcion")
		nombre := c.PostForm("nombre")
		observaciones := c.PostForm("observaciones")
		formaPago := c.PostForm("forma_pago")
		direccion := c.PostForm("direccion")

		// 2. Parsear el precio (si viene)
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

		// 3. Manejo de la imagen (si viene)
		fileHeader, err := c.FormFile("imagen")
		var publicURL string
		if err == nil {
			// El usuario envió un archivo
			bucketName := "imagenes-pedidos"
			filePath := fmt.Sprintf("pedidos/%d_%s", time.Now().Unix(), filepath.Base(fileHeader.Filename))

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

		// 6. Buscar el usuario (vendedor) en la BD para obtener su nombre
		var usuarioCreador models.Usuario
		if err := db.First(&usuarioCreador, "id = ?", usuarioID).Error; err != nil {

			fmt.Printf("Usuario con ID %s no encontrado o error: %v\n", usuarioID, err)
		}

		// 7. Construir el mensaje con el nombre del vendedor (si existe)
		mensaje := "Se ha creado un nuevo pedido!"
		if usuarioCreador.ID != "" {
			mensaje = fmt.Sprintf("%s (%s) ha creado un nuevo pedido!", usuarioCreador.Nombre, usuarioCreador.ID)
		}

		// 8. Enviar notificación a Admin/Gestor
		hub.BroadcastMessage(mensaje, "administrador", "gestor")

		// 9. Respuesta exitosa
		c.JSON(http.StatusOK, gin.H{
			"mensaje":    "Pedido creado exitosamente",
			"pedido_id":  nuevoPedido.ID,
			"usuario_id": nuevoPedido.UsuarioID,
			"imagen":     nuevoPedido.Imagen,
		})
	}
}
