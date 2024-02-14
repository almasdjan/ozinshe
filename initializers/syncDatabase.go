package initializers

import "project1/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
