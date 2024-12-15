package domain

import (
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	Email    string `gorm:"UNIQUE"`
	Name     string
	Password string
}
