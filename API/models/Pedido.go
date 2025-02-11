package models

import (
	"time"

	"gorm.io/gorm"
)

// Orden de venta creada por un vendedor.
type Pedido struct {
	ID            uint      `gorm:"primaryKey"`
	Descripcion   string    `gorm:"size:255;not null"`
	Imagen        string    `gorm:"size:255"`
	FechaCreacion time.Time `gorm:"autoCreateTime"`
	Fletero       *string   `gorm:"size:100;default:null"`
	Monto         *float64  `gorm:"default:null"`
	Estado        string    `gorm:"size:50;default:'No Entregado'"`
	Precio        *float64  `gorm:"default:null"`
	Nombre        string    `gorm:"size:55"`
	Observaciones string    `gorm:"size:255"`
	Forma_Pago    string    `gorm:"size:55"`
	Direccion     string    `gorm:"size:255"`
	Nro_Tlf       string    `gorm:"size:12"`
	Pagado        string    `gorm:"size:9;default:'No Pagado'"`
	Atendido      bool      `gorm:"default:false"`

	Nombre_Vendedor string
	UsuarioID       string
	Usuario         Usuario `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
