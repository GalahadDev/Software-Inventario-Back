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

	UsuarioID uint
	Usuario   Usuario `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
