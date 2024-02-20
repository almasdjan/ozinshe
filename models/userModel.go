package models

import (
	//"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string `gorm:"unique"`
	Password     string
	Name         string
	Phone_number string `gorm:"unique"`
	Birthday     string
	Isadmin      bool
}
