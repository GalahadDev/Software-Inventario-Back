package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"kings-house-back/API/models"
)

// SumarMontosPorVendedor maneja la suma de montos por vendedor en un rango de fechas
func SumarMontosPorVendedor(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Extraer el ID del vendedor desde la ruta
		vendedorID := c.Param("vendedor_id")
		if vendedorID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No se especificó el ID del vendedor"})
			return
		}

		// 2. Leer las fechas de query params: ?start_date=2025-01-15&end_date=2025-01-15
		startDateStr := c.Query("start_date")
		endDateStr := c.Query("end_date")

		if startDateStr == "" || endDateStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Debe proporcionar start_date y end_date"})
			return
		}

		// 3. Parsear las fechas
		layout := "2006-01-02" // Formato esperado
		startDate, err := time.Parse(layout, startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "start_date con formato inválido (YYYY-MM-DD)"})
			return
		}
		endDate, err := time.Parse(layout, endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "end_date con formato inválido (YYYY-MM-DD)"})
			return
		}

		// 4. Ajustar endDate para incluir todo el día
		endDate = endDate.AddDate(0, 0, 1) // Agregar un día

		// Log para verificar fechas
		log.Printf("DEBUG: VendedorID: %s, StartDate: %s, EndDate: %s", vendedorID, startDate, endDate)

		// 5. Consulta para sumar montos en ese rango de fechas y estado "Entregado"
		var total float64

		err = db.Model(&models.Pedido{}).
			Where("usuario_id = ?", vendedorID).
			Where("fecha_creacion >= ? AND fecha_creacion < ?", startDate, endDate).
			Where("estado = ?", "Entregado"). // Añadir esta línea para filtrar por estado
			Select("COALESCE(SUM(monto), 0)").Scan(&total).Error

		if err != nil {
			log.Printf("ERROR: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al calcular la suma de montos"})
			return
		}

		// Log para verificar el total calculado
		log.Printf("DEBUG: Total monto para VendedorID %s entre %s y %s: %f", vendedorID, startDate, endDate, total)

		c.JSON(http.StatusOK, gin.H{
			"vendedor_id": vendedorID,
			"start_date":  startDate.Format(layout),
			"end_date":    endDate.AddDate(0, 0, -1).Format(layout), // Restar un día para mostrar el end_date original
			"total_monto": total,
		})
	}
}
