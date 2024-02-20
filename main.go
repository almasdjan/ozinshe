package main

import (
	"project1/controllers"
	"project1/initializers"
	"project1/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVaribles()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)

	r.GET("/requireauth", middleware.RequireAuth)
	r.GET("/getuser", controllers.GetUserInfo)
	r.GET("/logout", controllers.Logout)
	r.PATCH("/editprofile", controllers.UpdateUserInfo)
	r.DELETE("/deleteprofile", controllers.DeleteProfile)
	r.PATCH("/changepassword", controllers.ChangePassword)

	r.Run()
}

//go get -u github.com/gin-gonic/gin
//go get github.com/githubnemo/CompileDaemon
//go get github.com/joho/godotenv
//go get -u gorm.io/gorm
//go get -u gorm.io/driver/postgres
//go get github.com/badoux/checkmail
