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
	initializers.ConnectDb()
	initializers.SyncDatabase()
}

// @title Ozinshe
// @version 1.0
// @description API Server for Ozinshe app

// @host ozinshetestapi.mobydev.kz
//ozinshetestapi.mobydev.kz
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	r := gin.Default()

	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := r.Group("/auth")
	auth.POST("/signup", controllers.Signup)
	auth.POST("/login", controllers.Login)
	auth.GET("/userinfo", controllers.GetUserInfo)
	auth.PATCH("/userinfo", controllers.UpdateUserInfo)
	auth.PATCH("/password", controllers.ChangePassword)
	auth.DELETE("/profile", controllers.DeleteProfile)

	r.GET("/logout", controllers.Logout)
	r.GET("/requireauth", middleware.RequireAuth)

	admin := r.Group("/admin")

	admin.POST("/age", controllers.CreateAge)
	admin.DELETE("/ages/:age_id", controllers.DeleteAge)
	admin.PATCH("/ages/:age_id", controllers.UpdateAge)

	admin.POST("/categories", controllers.CreateCategory)
	admin.PATCH("/categories/:category_id", controllers.UpdateCategories)
	admin.DELETE("/categories/:category_id", controllers.DeleteCategory)

	admin.POST("/genres", controllers.CreateGenre)
	admin.PATCH("/genres/:genre_id", controllers.UpdateGenre)
	admin.DELETE("/genres/:genre_id", controllers.DeleteGenre)

	admin.DELETE("categorymaterial/:material_id/:category_id", controllers.DeleteGenreCategoryMaterial)
	admin.DELETE("agematerial/:material_id/:age_id", controllers.DeleteAgeFromMaterial)
	admin.DELETE("genrematerial/:material_id/:genre_id", controllers.DeleteGenreFromMaterial)

	admin.POST("/genrematerial/:material_id/:genre_id", controllers.AddGenreToMaterial)
	admin.POST("/agematerial/:material_id/:age_id", controllers.AddAgeToMaterial)
	admin.POST("/categorymaterial/:material_id/:category_id", controllers.AddCategoryToMaterial)

	admin.POST("/videosrc", controllers.CreateVideo)
	admin.DELETE("/videosrc/:video_id", controllers.DeleteVideo)

	admin.POST("/recommends/:material_id", controllers.AddRecommend)
	admin.DELETE("/recommends/:material_id", controllers.DeleteFromRecommended)

	admin.POST("/material", controllers.CreateMaterial)
	admin.DELETE("/material/:material_id", controllers.DeleteMaterial)
	admin.PATCH("/material/:material_id", controllers.UpdateMaterial)
	admin.POST("/materialimage/:material_id", controllers.AddImage)
	admin.DELETE("materialimage", controllers.DeleteImage)

	admin.GET("material", controllers.GetAllMovies)

	main := r.Group("/main")
	main.GET("/ages", controllers.GetAges)
	main.GET("/categories", controllers.GetCategories)
	main.GET("/genres", controllers.GetGenres)

	main.GET("/", controllers.GetMainList)
	main.GET("/material/:material_id", controllers.GetMaterialById)
	main.GET("/genres/:genre_id", controllers.GetMovieByGenre)
	main.GET("/ages/:age_id", controllers.GetMovieByAge)
	main.GET("/history", controllers.GetMaterialHistory)
	main.GET("/trends", controllers.GetTrends)

	main.GET("/series/:material_id/:sezon", controllers.GetSezonsOrVideo)
	main.GET("/series/:material_id", controllers.GetSezonsOrVideo)
	main.GET("/series/serial/:material_id/:video_id", controllers.GetSerialSeries)

	main.GET("/foryou", controllers.GetRandomMovie)
	main.GET("/recommends", controllers.GetRecommended)

	main.GET("/search", controllers.Search)
	main.POST("/favourites/:material_id", controllers.AddFavouriteMovie)
	main.GET("/favourites", controllers.GetFavouriteMovies)
	main.DELETE("/favourites/:material_id", controllers.DeleteFromFavourites)

	r.PATCH("/updaterecommends/:queue/:material_id", controllers.UpdateRecommended)
	r.POST("/addhistory/:material_id", controllers.AddHistory)

	//r.GET("/getsezonsorvideo/:material_id/:sezon", controllers.GetSezonsOrVideo)
	//r.GET("/getsezonsorvideo/:material_id", controllers.GetSezonsOrVideo)

	r.PATCH("/updateviewed/:material_id", controllers.AddViewed)

	r.Run()
}
