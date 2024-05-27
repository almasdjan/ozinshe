package controllers

import (
	"context"
	"fmt"
	"net/http"
	"project1/initializers"
	"project1/middleware"
	"project1/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Summary Create material
// @Security BearerAuth
// @Tags admin
// @Accept json
// @Produce json
// @Param title formData string true "title"
// @Param posterr formData file true "poster"
// @Param description formData string true "description"
// @Param publish_year formData string true "publish year"
// @Param director formData string true "director"
// @Param producer formData string true "producer"
// @Param categories formData []string false "categories"
// @Param age_categories formData []string false "ages"
// @Param genres formData []string false "genres"
// @Param duration formData string true "duration"
// @Param image_srcs[] formData []file true "images"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /admin/material [post]
func CreateMaterial(c *gin.Context) {
	middleware.RequireAuth(c)

	userid, _ := c.Get("user")
	var user models.User

	initializers.DB.First(&user, userid)

	if !user.Isadmin {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This account is not admin",
		})
		return
	}

	title := c.PostForm("title")
	posterr, err := c.FormFile("posterr")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read poster",
		})
		return
	}
	description := c.PostForm("description")
	publish_year := c.PostForm("publish_year")
	director := c.PostForm("director")
	producer := c.PostForm("producer")
	categories := c.PostFormArray("categories")
	age := c.PostFormArray("age_categories")
	genre := c.PostFormArray("genres")
	duration := c.PostForm("duration")

	var material models.Material
	exist := initializers.DB.Where("title=?", title).First(&material)

	if exist.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This material is already exists",
		})
		return
	}

	publish_yearr, err := strconv.Atoi(publish_year)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error while getting publish year",
		})
		return
	}

	//save poster image
	path := "files//posters//" + posterr.Filename
	c.SaveUploadedFile(posterr, path)

	material = models.Material{
		Title:        title,
		Poster:       posterr.Filename,
		Description:  description,
		Publish_year: publish_yearr,
		Director:     director,
		Producer:     producer,
		Duration:     duration}

	result := initializers.DB.Create(&material)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create material",
		})
		return

	}

	//adding material_id, genre_id in material_genres
	var genresList []string
	for _, v := range genre {
		genresList = strings.Split(v, ",")
	}

	for _, v := range genresList {
		genre_id, err := strconv.Atoi(v)

		if err != nil {
			fmt.Println(err)
			fmt.Println(v)
			fmt.Println(genre_id)
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "genre_id is not correct",
			})
			return
		}
		material_genre := models.Material_genre{
			Material_id: material.ID,
			Genre_id:    uint(genre_id)}

		var genre models.Genre
		exist := initializers.DB.Where("id=?", v).First(&genre)
		if exist.RowsAffected == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "no such genre",
			})
			return

		}

		result1 := initializers.DB.Create(&material_genre)

		if result1.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to create material_genres",
			})
			return

		}

	}

	//adding material_id, category_id in material_categories
	//cate := strings.Split(categories[0], ",")
	var categoriesList []string
	for _, v := range categories {
		categoriesList = strings.Split(v, ",")
	}

	for _, v := range categoriesList {
		category_id, err := strconv.Atoi(v)
		if err != nil {
			fmt.Println(err)
			fmt.Println(v)
			fmt.Println(category_id)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "category_id is not correct",
			})
			return
		}
		material_category := models.Material_category{
			Material_id: material.ID,
			Category_id: uint(category_id)}

		var category models.Category
		exist := initializers.DB.Where("id=?", v).First(&category)
		if exist.RowsAffected == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "no such category",
			})
			return

		}

		result1 := initializers.DB.Create(&material_category)

		if result1.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to create material_category",
			})
			return

		}

	}

	//adding age category and material
	var agesList []string
	for _, v := range age {
		agesList = strings.Split(v, ",")
	}

	for _, v := range agesList {
		age_id, err := strconv.Atoi(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "age id is not correct",
			})
			return
		}
		material_age := models.Material_age{
			Material_id: material.ID,
			Age_id:      uint(age_id)}

		var age models.Age
		exist := initializers.DB.Where("id=?", v).First(&age)
		if exist.RowsAffected == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "no such age category",
			})
			return

		}

		result1 := initializers.DB.Create(&material_age)

		if result1.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to create material_ages",
			})
			return

		}

	}

	//adding material_id, image_src in image_srcs
	image_srcs, _ := c.MultipartForm()
	files := image_srcs.File["image_srcs[]"]
	fmt.Println(files)
	for _, file := range files {
		fmt.Println(file.Filename)
		//upload images into directory
		path := "files//images//" + file.Filename
		c.SaveUploadedFile(file, path)

		//adding filename in database
		image_src := models.Image_src{
			Material_id: material.ID,
			Image_src:   file.Filename}

		result1 := initializers.DB.Create(&image_src)

		if result1.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to create image",
			})
			return

		}
	}

	c.JSON(http.StatusOK, gin.H{
		"Action": "The material was succfully created",
	})

}

// @Summary Material by id
// @Security BearerAuth
// @Tags main
// @Accept json
// @Produce json
// @Param material_id path string true "material_id"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /main/material/{material_id} [get]
func GetMaterialById(c *gin.Context) {
	middleware.RequireAuth(c)

	userid, _ := c.Get("user")
	var user models.User

	initializers.DB.First(&user, userid)

	db, error := initializers.ConnectDb()
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to connect database",
		})
		return
	}

	material_id := c.Param("material_id")

	movie := db.QueryRow(context.Background(), "select poster, title, publish_year,  duration, description, director, producer, viewed from materials where id = $1	", material_id)
	if movie == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get materials info",
		})
		return
	}
	var movieInfo models.Movie

	err := movie.Scan(&movieInfo.Poster, &movieInfo.Title, &movieInfo.Publish_year, &movieInfo.Duration, &movieInfo.Description, &movieInfo.Director, &movieInfo.Producer, &movieInfo.Viewed)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get mmovies",
		})
		return
	}

	categoriesrows, err := db.Query(context.Background(), "select c.id, c.category_name from material_categories m join categories c on c.id = m.category_id where m.material_id = $1", material_id)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get categoriess",
		})
		return
	}

	var categories []models.Category
	for categoriesrows.Next() {
		var category models.Category
		err := categoriesrows.Scan(&category.ID, &category.CategoryName)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to get categories",
			})
			return
		}
		categories = append(categories, category)

	}

	agesrows, err := db.Query(context.Background(), "select a.id, a.age from material_ages m join ages a on a.id = m.age_id where m.material_id = $1", material_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get ages",
		})
		return
	}

	var ages []models.Age
	for agesrows.Next() {
		var age models.Age
		err := agesrows.Scan(&age.ID, &age.Age)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to get ages",
			})
			return
		}
		ages = append(ages, age)

	}

	genrerows, err := db.Query(context.Background(), "select g.id, g.genre_name from material_genres m join genres g on g.id = m.genre_id where m.material_id = $1", material_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get ages",
		})
		return
	}

	var genres []models.Genre
	for genrerows.Next() {
		var genre models.Genre
		err := genrerows.Scan(&genre.ID, &genre.GenreName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to get genres",
			})
			return
		}
		genres = append(genres, genre)

	}

	imagesrows, err := db.Query(context.Background(), "select image_src from image_srcs where material_id = $1", material_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get images",
		})
		return
	}

	var images []string
	for imagesrows.Next() {
		var image string
		err := imagesrows.Scan(&image)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to get imagess",
			})
			return
		}
		images = append(images, image)

	}

	row := db.QueryRow(context.Background(), "select count(*) from materials m join material_categories mc on m.id = mc.material_id join categories c on c.id = mc.category_id  where c.category_name like '%сериал%' and m.id = $1", material_id)

	if row == nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get materials info",
		})
		return
	}

	var isSerial int

	err = row.Scan(&isSerial)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to define serial or not",
		})
		return
	}

	if isSerial == 1 { //если сериал

		series := db.QueryRow(context.Background(), "select count(distinct sezon), count(series) from videos where material_id = $1", material_id)
		if series == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to define sezons and series",
			})
			return
		}

		var SezonsAndSeries models.SezonsAndSeries
		err = series.Scan(&SezonsAndSeries.SezonCount, &SezonsAndSeries.SeriesCount)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to define serial or not",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"categories":  categories,
			"movieinfo":   movieInfo,
			"screenshots": images,
			"sezons":      SezonsAndSeries.SezonCount,
			"series":      SezonsAndSeries.SeriesCount,
			"genres":      genres,
			"ages":        ages,
		})

	} else if isSerial == 0 {
		c.JSON(http.StatusOK, gin.H{
			"categories":  categories,
			"movieinfo":   movieInfo,
			"screenshots": images,
			"sezons":      nil,
			"series":      nil,
			"genres":      genres,
			"ages":        ages,
		})

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to define serial or not",
		})
		return

	}

}

type MainResponse struct {
	Recommends []models.Material_recommend `json:"recommends"`
	History    []models.Material_history   `json:"history"`
	Trends     []models.Material_get       `json:"trends"`
	Randoms    []models.Material_get       `json:"randoms"`
	Genre      []models.Genre              `json:"genres"`
	Ages       []models.Age                `json:"ages"`
}

// @Summary main page
// @Security BearerAuth
// @Tags main
// @Accept json
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /main [get]
func GetMainList(c *gin.Context) {
	middleware.RequireAuth(c)

	userid, _ := c.Get("user")
	var user models.User

	initializers.DB.First(&user, userid)

	recommended, err := GetRecommendedData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get reocommneds",
		})
		return
	}
	history, err := GetMaterialHistoryMain(userid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get history",
		})
		return
	}

	trends, err := GetTrendsMain()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get trends",
		})
		return
	}
	randoms, err := GetRandomMovieMain()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get randoms",
		})
		return
	}

	db, error := initializers.ConnectDb()
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to connect database",
		})
		return
	}

	rows, err := db.Query(context.Background(), `select * from genres`)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{

			"error": "Failed to get genres",
		})
		return
	}

	var genres []models.Genre
	for rows.Next() {
		var genre models.Genre
		err := rows.Scan(&genre.ID, &genre.GenreName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to genres",
			})
			return
		}
		genres = append(genres, genre)
	}

	agerows, err := db.Query(context.Background(), `select * from ages`)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{

			"error": "Failed to get ages",
		})
		return
	}

	var ages []models.Age
	for agerows.Next() {
		var age models.Age
		err := agerows.Scan(&age.ID, &age.Age)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to get ages",
			})
			return
		}
		ages = append(ages, age)
	}

	response := MainResponse{
		Recommends: recommended,
		History:    history,
		Trends:     trends,
		Randoms:    randoms,
		Genre:      genres,
		Ages:       ages,
	}

	c.JSON(http.StatusOK, response)

}

// @Summary Material by id
// @Security BearerAuth
// @Tags admin
// @Accept json
// @Produce json
// @Param material_id path string true "material_id"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /admin/material/{material_id} [delete]
func DeleteMaterial(c *gin.Context) {
	middleware.RequireAuth(c)

	userid, _ := c.Get("user")
	var user models.User

	initializers.DB.First(&user, userid)

	if !user.Isadmin {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This account is not admin",
		})
		return
	}

	db, error := initializers.ConnectDb()
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to connect database",
		})
		return
	}

	id := c.Param("material_id")
	_, err := db.Exec(context.Background(), `delete from material_ages where material_id=$1`, id)
	_, err1 := db.Exec(context.Background(), `delete from material_categories where material_id=$1;`, id)
	_, err2 := db.Exec(context.Background(), `delete from material_genres where material_id=$1`, id)
	_, err3 := db.Exec(context.Background(), `delete from image_srcs where material_id=$1`, id)
	_, err4 := db.Exec(context.Background(), `delete from user_favourites where material_id=$1`, id)
	_, err5 := db.Exec(context.Background(), `delete from videos where material_id=$1`, id)
	_, err6 := db.Exec(context.Background(), `delete from recommends where material_id=$1`, id)
	_, err7 := db.Exec(context.Background(), `delete from user_history where material_id=$1`, id)
	_, err8 := db.Exec(context.Background(), `delete from user_materials where material_id=$1`, id)
	_, err9 := db.Exec(context.Background(), `delete from materials where id=$1`, id)

	if err != nil || err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil || err7 != nil || err8 != nil || err9 != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete the material",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "the material was deleted",
	})

}

// @Summary edit material
// @Security BearerAuth
// @Tags admin
// @Accept json
// @Produce json
// @Param material_id path string true "material_id"
// @Param title formData string false "title"
// @Param posterr formData file false "poster"
// @Param description formData string false "description"
// @Param publish_year formData string false "publish year"
// @Param director formData string false "director"
// @Param producer formData string false "producer"
// @Param duration formData string false "producer"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /admin/material/{material_id} [patch]
func UpdateMaterial(c *gin.Context) {
	middleware.RequireAuth(c)

	userid, _ := c.Get("user")
	var user models.User

	initializers.DB.First(&user, userid)

	if !user.Isadmin {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This account is not admin",
		})
		return
	}

	db, error := initializers.ConnectDb()
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to connect database",
		})
		return
	}

	id := c.Param("material_id")
	var path string
	posterr, err := c.FormFile("posterr")

	if err != nil {
		path = ""
	} else {
		path = "files//posters//" + posterr.Filename
		c.SaveUploadedFile(posterr, path)
	}
	fmt.Print(path)

	title := c.PostForm("title")
	description := c.PostForm("description")
	publish_year := c.PostForm("publish_year")
	director := c.PostForm("director")
	fmt.Println(director)
	producer := c.PostForm("producer")
	fmt.Println(producer)
	duration := c.PostForm("duration")
	fmt.Println(duration)
	if title != "" {
		_, err := db.Exec(context.Background(), "update materials set title = $1 WHERE id = $2", title, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to update the material",
			})
			return
		}
	}

	if description != "" {
		_, err := db.Exec(context.Background(), "update materials set description = $1 WHERE id = $2", description, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to update the material",
			})
			return
		}
	}

	if director != "" {
		_, err := db.Exec(context.Background(), "update materials set director = $1 WHERE id = $2", director, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to update the material",
			})
			return
		}
	}

	if producer != "" {
		_, err := db.Exec(context.Background(), "update materials set producer = $1 WHERE id = $2", producer, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to update the material",
			})
			return
		}
	}

	if duration != "" {
		_, err := db.Exec(context.Background(), "update materials set duration = $1 WHERE id = $2", duration, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to update the material",
			})
			return
		}
	}

	if publish_year != "" {
		_, err := db.Exec(context.Background(), "update materials set publish_year = $1 WHERE id = $2", publish_year, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to update the material",
			})
			return
		}
	}

	if path != "" {
		_, err := db.Exec(context.Background(), "update materials set poster = $1 WHERE id = $2", posterr.Filename, id)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to update the material",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "the material was updated",
	})

}

// @Summary image
// @Security BearerAuth
// @Tags admin
// @Accept json
// @Produce json
// @Param image query string true "image"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /admin/materialimage [delete]
func DeleteImage(c *gin.Context) {
	middleware.RequireAuth(c)

	userid, _ := c.Get("user")
	var user models.User

	initializers.DB.First(&user, userid)

	if !user.Isadmin {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This account is not admin",
		})
		return
	}

	db, error := initializers.ConnectDb()
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to connect database",
		})
		return
	}

	image := c.Query("image")

	_, err := db.Exec(context.Background(), `delete from image_srcs WHERE image_src = $1`, image)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete image",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "The image was deleted",
	})

}

// @Summary add image to material
// @Security BearerAuth
// @Tags admin
// @Accept json
// @Produce json
// @Param material_id path string true "material id"
// @Param image_srcs[] formData []file true "images"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /admin/materialimage/{material_id} [post]
func AddImage(c *gin.Context) {
	middleware.RequireAuth(c)

	userid, _ := c.Get("user")
	var user models.User

	initializers.DB.First(&user, userid)

	if !user.Isadmin {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This account is not admin",
		})
		return
	}

	db, error := initializers.ConnectDb()
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to connect database",
		})
		return
	}
	id := c.Param("material_id")

	image_srcs, _ := c.MultipartForm()
	files := image_srcs.File["image_srcs[]"]

	for _, file := range files {
		//upload images into directory
		path := "files//images//" + file.Filename
		c.SaveUploadedFile(file, path)

		//adding filename in database
		_, err := db.Exec(context.Background(), `insert into image_srcs values ($1,$2)`, id, file.Filename)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to add images",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "The image is added",
	})

}

// @Summary search
// @Security BearerAuth
// @Tags main
// @Accept json
// @Produce json
// @Param search query string true "search"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /main/search [get]
func Search(c *gin.Context) {
	middleware.RequireAuth(c)

	userid, _ := c.Get("user")
	var user models.User

	initializers.DB.First(&user, userid)
	db, error := initializers.ConnectDb()
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to connect database",
		})
		return
	}

	search := c.Query("search")
	rows, err := db.Query(context.Background(), `select id,poster, title, category_name, publish_year from (
		select distinct on (id) *   from (
			select  m.id,m.poster, m.title, c.category_name, m.publish_year from materials m
			join material_categories mc on m.id = mc.material_id
			join categories c on mc.category_id = c.id
			where m.title like '%' || $1 || '%' 
		)
	 
	)`, search)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{

			"error": "Failed to get materials info",
		})
		return
	}

	var materials []models.Material_search
	for rows.Next() {
		var material models.Material_search
		err := rows.Scan(&material.Material_id, &material.Poster, &material.Title, &material.Category, &material.Publish_year)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to define materilas",
			})
			return
		}
		materials = append(materials, material)
	}
	c.JSON(http.StatusOK, gin.H{
		"found": materials,
	})

}
