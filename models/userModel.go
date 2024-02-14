package models

import (
	//"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string
	/*
		Name         string    `json:"name"`
		Phone_number string    `json:"phone_number"`
		Birthday     time.Time `json:"birthday"`
		Isadmin      bool      `json:"isadmin"`
	*/
}
