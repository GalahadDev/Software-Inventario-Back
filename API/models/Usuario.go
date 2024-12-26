package models

import (
	"time"

	"gorm.io/gorm"
)

// Modelo del Usuario
type Usuario struct {
	ID         uint   `gorm:"primaryKey"`
	Nombre     string `gorm:"size:100;not null"`
	Email      string `gorm:"size:100;not null;unique"`
	Contrasena string `gorm:"not null"`
	Rol        string `gorm:"size:50;not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
