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
	Phone_number string
	Birthday     string
	Isadmin      bool
	Favourites   []*Material `gorm:"many2many:user_materials;"`
}

type Userjson struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	Password2    string `json:"password2"`
	Name         string `json:"name"`
	Phone_number string `json:"phone_number"`
	Birthday     string `json:"birthday"`
}

type User_favourites struct {
	User_id     uint
	Material_id uint
}
