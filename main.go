package main

import (
	"project1/controllers"
	"project1/initializers"
	"project1/middleware"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "project1/docs"

	swaggerFiles "github.com/swaggo/files"
)

func init() {
	initializers.LoadEnvVaribles()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

// @title Ozinshe
// @version 1.0
// @description API Server for Ozinshe app

// @host localhost:8080
// @BasePath /

// securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	r := gin.Default()

	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := r.Group("/auth")
	auth.POST("/signup", controllers.Signup)

	r.POST("/login", controllers.Login)
	r.GET("/requireauth", middleware.RequireAuth)
	r.GET("/getuser", controllers.GetUserInfo)
	r.GET("/logout", controllers.Logout)
	r.PATCH("/editprofile", controllers.UpdateUserInfo)
	r.DELETE("/deleteprofile", controllers.DeleteProfile)
	r.PATCH("/changepassword", controllers.ChangePassword)

	r.POST("/addfavourite", controllers.AddFavouriteMovie)
	r.POST("/addrecommend/:material_id", controllers.AddRecommend)
	r.GET("/getrecommends", controllers.GetRecommended)
	r.DELETE("/deletefromrecommends/:queue", controllers.DeleteFromRecommended)
	r.PATCH("/updaterecommends/:queue/:material_id", controllers.UpdateRecommended)
	r.POST("/addhistory/:material_id", controllers.AddHistory)
	r.GET("/gethistory", controllers.GetMaterialHistory)
	r.GET("/gettrends", controllers.GetTrends)

	r.GET("/foryou", controllers.GetRandomMovie)

	r.GET("/main", controllers.GetMainList)
	r.GET("/main/genre/:genre_id", controllers.GetMovieByGenre)
	r.GET("/main/age/:age_id", controllers.GetMovieByAge)
	r.GET("/getsezonsorvideo/:material_id/:sezon", controllers.GetSezonsOrVideo)
	r.GET("/getsezonsorvideo/:material_id", controllers.GetSezonsOrVideo)
	r.PATCH("/updateviewed/:material_id", controllers.AddViewed)

	r.GET("/getmaterialbyid/:material_id", controllers.GetMaterialById)

	r.GET("/search", controllers.Search)

	r.POST("/addmaterial", controllers.CreateMaterial)
	r.DELETE("/delete/:material_id", controllers.DeleteMaterial)
	r.PATCH("/update/:material_id", controllers.UpdateMaterial)
	r.POST("/addimage/:material_id", controllers.AddImage)
	r.DELETE("deleteimage/:image", controllers.DeleteImage)
	r.POST("/addgenretomaterial/:material_id/:genre_id", controllers.AddGenreToMaterial)
	r.DELETE("deletegenrefrommaterial/:material_id/:genre_id", controllers.DeleteGenreFromMaterial)
	r.POST("/addagetomaterial/:material_id/:age_id", controllers.AddAgeToMaterial)
	r.DELETE("deleteagefrommaterial/:material_id/:age_id", controllers.DeleteAgeFromMaterial)
	r.POST("/addcategorytomaterial/:material_id/:category_id", controllers.AddCategoryToMaterial)
	r.DELETE("deletecategoryfrommaterial/:material_id/:category_id", controllers.DeleteGenreCategoryMaterial)

	r.POST("/createcategory", controllers.CreateCategory)
	r.GET("/getcategories", controllers.GetCategories)
	r.PATCH("/updatecategories/:category_id", controllers.UpdateCategories)
	r.DELETE("/deletecategories/:category_id", controllers.DeleteCategory)

	r.POST("/creategenre", controllers.CreateGenre)
	r.GET("/getgenres", controllers.GetGenres)
	r.PATCH("/updategenre/:genre_id", controllers.UpdateGenre)
	r.DELETE("/deletegenre/:genre_id", controllers.DeleteGenre)

	r.POST("/createagecategory", controllers.CreateAge)
	r.GET("/getages", controllers.GetAges)
	r.PATCH("/updateages/:age_id", controllers.UpdateAge)
	r.DELETE("/deleteage/:age_id", controllers.DeleteAge)

	r.POST("/addvideosrc", controllers.CreateVideo)
	r.DELETE("/deletevideosrc/:video_id", controllers.DeleteVideo)

	r.POST("/adddirector", controllers.CreateCategory)

	r.Run()
}
