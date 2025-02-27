package handlers

import (
	"encoding/json" // Importante para serializar a JSON
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"kings-house-back/API/database"
	"kings-house-back/API/models"
	"kings-house-back/API/ws"
)

// Estructura para notificar un nuevo pedido
type NotificacionPedido struct {
	Tipo    string        `json:"tipo"`
	Pedido  models.Pedido `json:"pedido"` // los campos del pedido
	Creador string        `json:"creador,omitempty"`
}

func CrearPedidoHandler(db *gorm.DB, hub *ws.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. Leer campos de form-data
		usuarioID := c.PostForm("usuario_id")
		descripcion := c.PostForm("descripcion")
		nombre := c.PostForm("nombre")
		observaciones := c.PostForm("observaciones")
		formaPago := c.PostForm("forma_pago")
		direccion := c.PostForm("direccion")
		numerotlf := c.PostForm("nro_tlf")
		tela := c.PostForm("tela")
		color := c.PostForm("color")
		subVendedor := c.PostForm("sub_vendedor")
		fechaEntrega := c.PostForm("fecha_entrega")

		// Parsear precio
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

		// Parsear comision
		comisionSugerida := c.PostForm("comision_sugerida")
		var comisionFloat *float64
		if comisionSugerida != "" {
			if p, err := strconv.ParseFloat(comisionSugerida, 64); err == nil {
				comisionFloat = &p
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Precio inválido"})
				return
			}
		}

		// Manejo de imagen si existe
		fileHeader, err := c.FormFile("imagen")
		var publicURL string
		if err == nil {
			bucketName := "imagenes-pedidos"
			rawFileName := filepath.Base(fileHeader.Filename)
			cleanFileName := strings.ReplaceAll(rawFileName, " ", "_")
			filePath := fmt.Sprintf("pedidos/%d_%s", time.Now().Unix(), cleanFileName)

			publicURL, err = database.SubirAStorageSupabase(fileHeader, bucketName, filePath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al subir a Supabase", "details": err.Error()})
				return
			}
		}

		// 4. Construir objeto Pedido
		nuevoPedido := models.Pedido{
			UsuarioID:         usuarioID,
			Descripcion:       descripcion,
			Imagen:            publicURL,
			FechaCreacion:     time.Now(),
			Precio:            precioFloat,
			Fletero:           nil,
			Monto:             nil,
			Estado:            "Pendiente",
			Nombre:            nombre,
			Observaciones:     observaciones,
			Forma_Pago:        formaPago,
			Direccion:         direccion,
			Nro_Tlf:           numerotlf,
			Pagado:            "No Pagado",
			Tela:              tela,
			Color:             color,
			Comision_Sugerida: comisionFloat,
			Sub_Vendedor:      subVendedor,
			Fecha_Entrega:     fechaEntrega,
		}

		// 5. Obtener el usuario para saber su nombre
		var usuarioCreador models.Usuario
		if err := db.First(&usuarioCreador, "id = ?", usuarioID).Error; err != nil {
			fmt.Printf("Usuario con ID %s no encontrado o error: %v\n", usuarioID, err)
		}

		creador := usuarioCreador.Nombre
		if creador == "" {
			creador = "Vendedor Desconocido"
		}

		// 6. Asignar Nombre_Vendedor con el nombre del usuario (vendedor)
		nuevoPedido.Nombre_Vendedor = creador

		// 7. Guardar en la BD
		if err := db.Create(&nuevoPedido).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el pedido"})
			return
		}

		// 8. Armar notificación como JSON
		notif := NotificacionPedido{
			Tipo:    "NUEVO_PEDIDO",
			Pedido:  nuevoPedido,
			Creador: creador,
		}

		notifBytes, _ := json.Marshal(notif)

		// 9. Enviar a admin/gestor
		hub.BroadcastMessage(string(notifBytes), "administrador", "gestor")

		// 10. Respuesta HTTP exitosa
		c.JSON(http.StatusOK, gin.H{
			"mensaje":    "Pedido creado exitosamente",
			"pedido_id":  nuevoPedido.ID,
			"usuario_id": nuevoPedido.UsuarioID,
			"imagen":     nuevoPedido.Imagen,
		})
	}
}
