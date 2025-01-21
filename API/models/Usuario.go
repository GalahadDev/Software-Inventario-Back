package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Modelo del Usuario
type Usuario struct {
	ID         string         `gorm:"type:uuid;default:uuid_generate_v4();primarykey"`
	Nombre     string         `gorm:"size:100;not null"`
	Email      string         `gorm:"size:100;not null;unique"`
	Contrasena string         `json:"-" gorm:"not null"`
	Rol        string         `gorm:"size:50;not null"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

// BeforeCreate se ejecuta antes de que GORM inserte el registro en la BD.
func (u *Usuario) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString() // Genera el UUID
	return
}
