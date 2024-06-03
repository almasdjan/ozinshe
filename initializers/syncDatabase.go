package initializers

import "project1/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Material{})
	DB.AutoMigrate(&models.Image_src{})

}




