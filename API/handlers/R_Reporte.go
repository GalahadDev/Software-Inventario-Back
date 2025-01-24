package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"kings-house-back/API/models"
)

func SumarMontosPorVendedor(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Extraer el ID del vendedor desde la ruta (o param).
		vendedorID := c.Param("vendedor_id")
		if vendedorID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No se especificó el ID del vendedor"})
			return
		}

		// 2. Leer las fechas de query params: ?start_date=2025-01-01&end_date=2025-01-31
		startDateStr := c.Query("start_date")
		endDateStr := c.Query("end_date")

		if startDateStr == "" || endDateStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Debe proporcionar start_date y end_date"})
			return
		}

		// 3. Parsear las fechas
		layout := "2006-01-02" // Ajusta según el formato que envíes
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

		// 4. Consulta para sumar montos en ese rango de fechas
		var total float64

		// OJO: Como Monto es *float64, conviene usar COALESCE o algo similar,
		// pero GORM permite hacer un select sum(...) sin problemas.
		// Ejemplo con GORM:
		//  - Filtrar por usuario y fecha_creacion
		//  - Hacer un .Select("COALESCE(SUM(monto), 0)").Scan(&total)
		//  - Asegúrate de comparar con "fecha_creacion BETWEEN ? AND ?"

		// Filtramos también que 'monto' no sea nulo (o si deseas contar null como 0, no hace falta).
		err = db.Model(&models.Pedido{}).
			Where("usuario_id = ?", vendedorID).
			Where("fecha_creacion BETWEEN ? AND ?", startDate, endDate).
			Select("COALESCE(SUM(monto), 0)").Scan(&total).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al calcular la suma de montos"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"vendedor_id": vendedorID,
			"start_date":  startDate.Format(layout),
			"end_date":    endDate.Format(layout),
			"total_monto": total,
		})
	}
}
