package controllers

import (
	"context"
	"fmt"
	"net/http"
	"project1/initializers"
	"project1/middleware"
	"project1/models"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Summary Create material
// @Security BearerAuth
// @Tags admin
// @Param title formData string true "title"
// @Param description formData string true "description"
// @Param publish_year formData string true "publish year"
// @Param director formData string true "director"
// @Param producer formData string true "producer"
// @Param categories formData []string false "categories"
// @Param age_categories formData []string false "ages"
// @Param genres formData []string false "genres"
// @Param duration formData string true "duration"
// @Param keywords formData string true "keywords"
// @Param type formData string true "type" Enums(Фильмы, Сериалы)
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /admin/material [post]
func CreateMaterial(c *gin.Context) {
	middleware.RequireAuth(c)
	if c.IsAborted() {
		return
	}
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

	description := c.PostForm("description")
	publish_year := c.PostForm("publish_year")
	if publish_year == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "publish year is required",
		})
		return
	}
	director := c.PostForm("director")
	producer := c.PostForm("producer")
	categories := c.PostFormArray("categories")
	age := c.PostFormArray("age_categories")
	genre := c.PostFormArray("genres")
	duration := c.PostForm("duration")
	keywords := c.PostForm("keywords")
	m_type := c.PostForm("type")

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
		fmt.Print(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error while getting publish year",
		})
		return
	}

	//save poster image

	material = models.Material{
		Title:        title,
		Description:  description,
		Publish_year: publish_yearr,
		Director:     director,
		Producer:     producer,
		Duration:     duration,
		Keywords:     keywords,
		M_type:       m_type}

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

	c.JSON(http.StatusOK, gin.H{
		"Action": "The material was succfully created",
	})

}

// @Summary add poster and screenshots
// @Security BearerAuth
// @Tags admin
// @Param id path string true "material id"
// @Param posterr formData file true "poster"
// @Param image_srcs[] formData []file true "images"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /admin/material/screens/{id} [post]
func AddScreens(c *gin.Context) {
	middleware.RequireAuth(c)
	if c.IsAborted() {
		return
	}
	userid, _ := c.Get("user")
	var user models.User

	initializers.DB.First(&user, userid)

	if !user.Isadmin {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This account is not admin",
		})
		return
	}

	material_id := c.Param("id")

	id, err := strconv.Atoi(material_id)
	if err != nil {
		fmt.Print(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read id",
		})
		return
	}

	posterr, err := c.FormFile("posterr")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read poster",
		})
		return
	}

	//save poster image
	path := "files//posters//" + posterr.Filename
	c.SaveUploadedFile(posterr, path)

	_, err = initializers.ConnPool.Exec(context.Background(), "update materials set poster = $1 WHERE id = $2", path, material_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to add poster",
		})
		return
	}

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
			Material_id: uint(id),
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
		"Action": "The images were succfully added",
	})

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
// @Router /admin/material/{material_id} [get]
func GetMaterialById(c *gin.Context) {
	middleware.RequireAuth(c)
	if c.IsAborted() {
		return
	}
	userid, _ := c.Get("user")
	var user models.User

	initializers.DB.First(&user, userid)
	/*
		db, error := initializers.ConnectDb()
		if error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to connect database",
			})
			return
		}*/

	material_id := c.Param("material_id")

	movie := initializers.ConnPool.QueryRow(context.Background(), "select id, poster, title, publish_year,  duration, description, director, producer from materials where id = $1	", material_id)
	if movie == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get materials info",
		})
		return
	}
	var movieInfo models.Movie

	err := movie.Scan(&movieInfo.Id, &movieInfo.Poster, &movieInfo.Title, &movieInfo.Publish_year, &movieInfo.Duration, &movieInfo.Description, &movieInfo.Director, &movieInfo.Producer)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get mmovies",
		})
		return
	}

	categoriesrows, err := initializers.ConnPool.Query(context.Background(), "select c.id, c.category_name from material_categories m join categories c on c.id = m.category_id where m.material_id = $1", material_id)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get categoriess",
		})
		return
	}

	var categories = []models.Category{}
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

	agesrows, err := initializers.ConnPool.Query(context.Background(), "select a.id, a.age from material_ages m join ages a on a.id = m.age_id where m.material_id = $1", material_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get ages",
		})
		return
	}

	var ages = []models.Age{}
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

	genrerows, err := initializers.ConnPool.Query(context.Background(), "select g.id, g.genre_name from material_genres m join genres g on g.id = m.genre_id where m.material_id = $1", material_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get ages",
		})
		return
	}

	var genres = []models.Genre{}
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

	imagesrows, err := initializers.ConnPool.Query(context.Background(), "select id, material_id, image_src from image_srcs where material_id = $1", material_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get images",
		})
		return
	}

	var images = []models.Image_src{}
	for imagesrows.Next() {
		var image models.Image_src
		err := imagesrows.Scan(&image.Id, &image.Material_id, &image.Image_src)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to get imagess",
			})
			return
		}
		images = append(images, image)

	}

	row := initializers.ConnPool.QueryRow(context.Background(), "select m_type from materials where id = $1", material_id)

	if row == nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get materials info",
		})
		return
	}

	var isSerial string

	err = row.Scan(&isSerial)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to define serial or not",
		})
		return
	}

	if isSerial == "Сериалы" { //если сериал

		video, err := initializers.ConnPool.Query(context.Background(), "select id, material_id, sezon, series, video_src, viewed  from videos where material_id = $1 order by sezon, series", material_id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed get videos",
			})
		}
		var videos = []models.Video{}
		fmt.Print(videos)
		if video != nil {
			var videoss models.Video
			for video.Next() {
				err = video.Scan(&videoss.Id, &videoss.Material_id, &videoss.Sezon, &videoss.Series, &videoss.Video_src, &videoss.Viewed)
				if err != nil {
					videos = []models.Video{}
				}
				videos = append(videos, videoss)
			}

		}

		series := initializers.ConnPool.QueryRow(context.Background(), "select count(*) from videos where material_id = $1 AND sezon = 1", material_id)

		var serie int
		err = series.Scan(&serie)
		if err != nil {
			serie = 0
		}

		sezons := initializers.ConnPool.QueryRow(context.Background(), "select count(distinct(sezon)) from videos where material_id = $1", material_id)

		if sezons == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to define sezons and series",
			})
			return
		}

		var sezon int
		err = sezons.Scan(&sezon)
		if err != nil {
			sezon = 0
		}
		c.JSON(http.StatusOK, gin.H{
			"categories":  categories,
			"movieinfo":   movieInfo,
			"screenshots": images,
			"videos":      videos,
			"sezons":      sezon,
			"series":      serie,
			"genres":      genres,
			"ages":        ages,
		})

	} else if isSerial == "Фильмы" {
		video := initializers.ConnPool.QueryRow(context.Background(), "select id, material_id, sezon, series, video_src, viewed  from videos where material_id = $1 AND SEZON = 0 AND SERIES = 0", material_id)

		var videos = []models.Video{}
		fmt.Print(videos)
		if video != nil {
			var videoss models.Video
			err = video.Scan(&videoss.Id, &videoss.Material_id, &videoss.Sezon, &videoss.Series, &videoss.Video_src, &videoss.Viewed)
			if err != nil {
				videos = []models.Video{}
			}
			videos = append(videos, videoss)
		}

		c.JSON(http.StatusOK, gin.H{
			"categories":  categories,
			"movieinfo":   movieInfo,
			"screenshots": images,
			"sezons":      0,
			"series":      0,
			"videos":      videos,
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
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /main [get]
func GetMainList(c *gin.Context) {
	middleware.RequireAuth(c)
	if c.IsAborted() {
		return
	}
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
	/*
		db, error := initializers.ConnectDb()
		if error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to connect database",
			})
			return
		}*/

	rows, err := initializers.ConnPool.Query(context.Background(), `select * from genres`)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{

			"error": "Failed to get genres",
		})
		return
	}

	var genres = []models.Genre{}
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

	agerows, err := initializers.ConnPool.Query(context.Background(), `select * from ages`)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{

			"error": "Failed to get ages",
		})
		return
	}

	var ages = []models.Age{}
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
	if c.IsAborted() {
		return
	}
	userid, _ := c.Get("user")
	var user models.User

	initializers.DB.First(&user, userid)

	if !user.Isadmin {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This account is not admin",
		})
		return
	}
	/*
		db, error := initializers.ConnectDb()
		if error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to connect database",
			})
			return
		}*/

	id := c.Param("material_id")
	_, err := initializers.ConnPool.Exec(context.Background(), `delete from material_ages where material_id=$1`, id)
	_, err1 := initializers.ConnPool.Exec(context.Background(), `delete from material_categories where material_id=$1;`, id)
	_, err2 := initializers.ConnPool.Exec(context.Background(), `delete from material_genres where material_id=$1`, id)
	_, err3 := initializers.ConnPool.Exec(context.Background(), `delete from image_srcs where material_id=$1`, id)
	_, err4 := initializers.ConnPool.Exec(context.Background(), `delete from user_favourites where material_id=$1`, id)
	_, err5 := initializers.ConnPool.Exec(context.Background(), `delete from videos where material_id=$1`, id)
	_, err6 := initializers.ConnPool.Exec(context.Background(), `delete from recommends where material_id=$1`, id)
	_, err7 := initializers.ConnPool.Exec(context.Background(), `delete from user_history where material_id=$1`, id)
	_, err8 := initializers.ConnPool.Exec(context.Background(), `delete from user_materials where material_id=$1`, id)
	_, err9 := initializers.ConnPool.Exec(context.Background(), `delete from materials where id=$1`, id)

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
	if c.IsAborted() {
		return
	}
	userid, _ := c.Get("user")
	var user models.User

	initializers.DB.First(&user, userid)

	if !user.Isadmin {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This account is not admin",
		})
		return
	}
	/*
		db, error := initializers.ConnectDb()
		if error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to connect database",
			})
			return
		}
	*/
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
		_, err := initializers.ConnPool.Exec(context.Background(), "update materials set title = $1 WHERE id = $2", title, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to update the material",
			})
			return
		}
	}

	if description != "" {
		_, err := initializers.ConnPool.Exec(context.Background(), "update materials set description = $1 WHERE id = $2", description, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to update the material",
			})
			return
		}
	}

	if director != "" {
		_, err := initializers.ConnPool.Exec(context.Background(), "update materials set director = $1 WHERE id = $2", director, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to update the material",
			})
			return
		}
	}

	if producer != "" {
		_, err := initializers.ConnPool.Exec(context.Background(), "update materials set producer = $1 WHERE id = $2", producer, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to update the material",
			})
			return
		}
	}

	if duration != "" {
		_, err := initializers.ConnPool.Exec(context.Background(), "update materials set duration = $1 WHERE id = $2", duration, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to update the material",
			})
			return
		}
	}

	if publish_year != "" {
		_, err := initializers.ConnPool.Exec(context.Background(), "update materials set publish_year = $1 WHERE id = $2", publish_year, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to update the material",
			})
			return
		}
	}

	if path != "" {
		_, err := initializers.ConnPool.Exec(context.Background(), "update materials set poster = $1 WHERE id = $2", posterr.Filename, id)

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
	if c.IsAborted() {
		return
	}
	userid, _ := c.Get("user")
	var user models.User

	initializers.DB.First(&user, userid)

	if !user.Isadmin {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This account is not admin",
		})
		return
	}
	/*
		db, error := initializers.ConnectDb()
		if error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to connect database",
			})
			return
		}*/

	image := c.Query("image")

	_, err := initializers.ConnPool.Exec(context.Background(), `delete from image_srcs WHERE image_src = $1`, image)
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
	if c.IsAborted() {
		return
	}
	userid, _ := c.Get("user")
	var user models.User

	initializers.DB.First(&user, userid)

	if !user.Isadmin {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This account is not admin",
		})
		return
	}
	/*
		db, error := initializers.ConnectDb()
		if error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to connect database",
			})
			return
		}*/
	id := c.Param("material_id")

	image_srcs, _ := c.MultipartForm()
	files := image_srcs.File["image_srcs[]"]

	for _, file := range files {
		//upload images into directory
		path := "files//images//" + file.Filename
		c.SaveUploadedFile(file, path)

		//adding filename in database
		_, err := initializers.ConnPool.Exec(context.Background(), `insert into image_srcs values ($1,$2)`, id, file.Filename)
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
	if c.IsAborted() {
		return
	}
	userid, _ := c.Get("user")
	var user models.User

	initializers.DB.First(&user, userid)
	/*
		db, error := initializers.ConnectDb()
		if error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to connect database",
			})
			return
		}*/

	search := c.Query("search")
	rows, err := initializers.ConnPool.Query(context.Background(), `select id,poster, title, category_name, publish_year from (
		select distinct on (id) *   from (
			select  m.id,m.poster, m.title, c.category_name, m.publish_year from materials m
			join material_categories mc on m.id = mc.material_id
			join categories c on mc.category_id = c.id
			where m.title ILIKE '%' || $1 || '%' 
		)as foo
	 
	)as foor`, search)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{

			"error": "Failed to get materials info",
		})
		return
	}

	var materials = []models.Material_search{}
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

/*
func GetAllMovies(c *gin.Context) {
	middleware.RequireAuth(c)
	if c.IsAborted() {
		return
	}


	sort := c.DefaultQuery("sort", "Популярные")
	category := c.DefaultQuery("category", "Все категории")
	materialType := c.DefaultQuery("type", "Фильмы и сериалы")
	year := c.DefaultQuery("year", "Выберите год")

	// Construct the base query
	query := `
	select id, title, poster,category_name, viewed from (
		select distinct on (id) *   from (
			select  m.id,m.poster, m.title, c.category_name, m.viewed from materials m
			join material_categories mc on m.id = mc.material_id
			join categories c on mc.category_id = c.id
		)as foo
	)as foor
	WHERE 1=1
	`

	countQuery := `
	select count(*) from (
		select distinct on (id) *   from (
			select  m.id,m.poster, m.title, c.category_name, m.viewed from materials m
			join material_categories mc on m.id = mc.material_id
			join categories c on mc.category_id = c.id
		)as foo
	)as foor
	WHERE 1=1
	`

	// Add filters to the query
	var params []interface{}
	paramIndex := 1

	if category != "Все категории" {
		query += fmt.Sprintf(" AND category_name = $%d", paramIndex)
		countQuery += fmt.Sprintf(" AND category_name = $%d", paramIndex)
		params = append(params, category)
		paramIndex++
	}
	if materialType != "Фильмы и сериалы" {
		query += fmt.Sprintf(" AND m_type = $%d", paramIndex)
		countQuery += fmt.Sprintf(" AND m_type = $%d", paramIndex)
		params = append(params, materialType)
		paramIndex++
	}
	if year != "Выберите год" {
		query += fmt.Sprintf(" AND publish_year = $%d", paramIndex)
		countQuery += fmt.Sprintf(" AND m_type = $%d", paramIndex)
		params = append(params, year)
		paramIndex++
	}
	if sort != "Популярные" {
		if sort == "По дате регистрации" {
			query += " ORDER BY created_at DESC"
		} else if sort == "По имени" {
			query += " ORDER BY title"
		}
	} else {
		query += " ORDER BY viewed DESC"
	}
	fmt.Print(query)

	// Add sorting to the query (you might need to adjust this depending on your sort logic)

	rows, err := initializers.ConnPool.Query(context.Background(), query, params...)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{

			"error": "Failed to get materials info",
		})
		return
	}

	var materials = []models.Material_get{}
	for rows.Next() {
		var material models.Material_get
		err := rows.Scan(&material.Material_id, &material.Title, &material.Poster, &material.Category, &material.Viewed)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to define materilas",
			})
			return
		}
		materials = append(materials, material)
	}

	countRows := initializers.ConnPool.QueryRow(context.Background(), countQuery, params...)

	var count int
	err = countRows.Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to count movies",
		})
		return
	}

	yearsQuery := `select distinct publish_year from materials order by publish_year desc`

	yearRows, err := initializers.ConnPool.Query(context.Background(), yearsQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get years",
		})
		return
	}

	var years = []int{}
	for yearRows.Next() {
		var year int
		err := yearRows.Scan(&year)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to define years",
			})
			return
		}
		years = append(years, year)
	}

	c.JSON(http.StatusOK, gin.H{
		"movies": materials,
		"count":  count,
		"years":  years,
	})

}
*/

// @Summary Get All movies
// @Security BearerAuth
// @Tags admin
// @Accept json
// @Produce json
// @Param sort query string false "Sort order" default(Популярные) Enums(По дате регистрации, По дате обновления)
// @Param category query string false "Category" default(Все категории)
// @Param type query string false "Type" default(Фильмы и сериалы) Enums(Фильмы, Сериалы)
// @Param year query string false "Year" default(Выберите год)
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /admin/material [get]
func GetAll(c *gin.Context) {
	middleware.RequireAuth(c)
	if c.IsAborted() {
		return
	}
	/*
		userid, _ := c.Get("user")

		var user models.User
		initializers.DB.First(&user, userid)
	*/

	sortt := c.DefaultQuery("sort", "Популярные")
	category := c.DefaultQuery("category", "Все категории")
	materialType := c.DefaultQuery("type", "Фильмы и сериалы")
	year := c.DefaultQuery("year", "Выберите год")

	// Construct the base query
	query := `
	select m.id,m.poster, m.title, m.created_at,m.updated_at,c.category_name from materials m 
	join material_categories mc on m.id=mc.material_id
	join categories c on c.id = mc.category_id 
	WHERE 1=1
	`

	countQuery := `
	select count(*) from (
		select distinct on (m.id) * from materials m 
			join material_categories mc on m.id=mc.material_id
			join categories c on c.id = mc.category_id 
			WHERE 1=1
	`

	// Add filters to the query
	var params []interface{}
	paramIndex := 1

	if category != "Все категории" {
		query += fmt.Sprintf(" AND c.category_name = $%d", paramIndex)
		countQuery += fmt.Sprintf(" AND c.category_name = $%d", paramIndex)
		params = append(params, category)
		paramIndex++
	}
	if materialType != "Фильмы и сериалы" {
		query += fmt.Sprintf(" AND m.m_type = $%d", paramIndex)
		countQuery += fmt.Sprintf(" AND m.m_type = $%d", paramIndex)
		params = append(params, materialType)
		paramIndex++
	}
	if year != "Выберите год" {
		query += fmt.Sprintf(" AND m.publish_year = $%d", paramIndex)
		countQuery += fmt.Sprintf(" AND m.publish_year = $%d", paramIndex)
		params = append(params, year)
		paramIndex++
	}

	countQuery += " )as foo"
	fmt.Print(query)

	// Add sorting to the query (you might need to adjust this depending on your sort logic)

	rows, err := initializers.ConnPool.Query(context.Background(), query, params...)
	fmt.Print(params...)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{

			"error": "Failed to get materials info",
		})
		return
	}

	var materials = []models.Material_get{}
	for rows.Next() {
		var material models.Material_get
		err := rows.Scan(&material.Material_id, &material.Poster, &material.Title, &material.CreatedAt, &material.UpdatedAt, &material.Category)
		if err != nil {
			fmt.Print(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{

				"error": "Failed to define materilas",
			})
			return
		}
		materials = append(materials, material)
	}

	moviesMap := make(map[uint]*models.Materials_get)
	for _, movie := range materials {
		if _, found := moviesMap[movie.Material_id]; !found {

			moviesMap[movie.Material_id] = &models.Materials_get{
				Material_id: movie.Material_id,
				Title:       movie.Title,
				CreatedAt:   movie.CreatedAt,
				UpdatedAt:   movie.UpdatedAt,
				Poster:      movie.Poster,
			}

		}

		category := movie.Category
		moviesMap[movie.Material_id].Category = append(moviesMap[movie.Material_id].Category, category)

	}
	type series struct {
		Material_id uint
		Viewed      int
		Sezon       int
		Series      int
	}
	var seriesNumber []series
	result := make([]models.Materials_get, 0, len(moviesMap))
	if sortt == "Популярные" {
		fmt.Print("Популярные")
		query := `
		WITH MaxViewed AS (
			SELECT 
				material_id, MAX(viewed) AS max_viewed 
			FROM videos 
			GROUP BY material_id
		)
		SELECT distinct on (m.id)
			m.id, mv.max_viewed AS viewed, v.sezon, v.series 
		FROM materials m
		JOIN MaxViewed mv ON m.id = mv.material_id
		JOIN videos v ON m.id = v.material_id AND mv.max_viewed = v.viewed
		;	
		`

		rows, err := initializers.ConnPool.Query(context.Background(), query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get series",
			})
			return
		}
		for rows.Next() {
			var seria series
			err = rows.Scan(&seria.Material_id, &seria.Viewed, &seria.Sezon, &seria.Series)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to get series",
				})
				return
			}
			seriesNumber = append(seriesNumber, seria)

		}
		fmt.Print(seriesNumber)
		for _, movie := range seriesNumber {
			if _, found := moviesMap[movie.Material_id]; found {
				/*
					moviesMap[movie.Material_id] = &models.Materials_get{
						Sezon:  movie.Sezon,
						Series: movie.Series,
					}*/

				moviesMap[movie.Material_id].Sezon = movie.Sezon
				moviesMap[movie.Material_id].Series = movie.Series
				moviesMap[movie.Material_id].Viewed = movie.Viewed

			}
		}
		for _, v := range moviesMap {
			result = append(result, *v)
		}
		sort.Slice(result, func(i, j int) bool {
			return result[i].Viewed > result[j].Viewed
		})

	}

	if sortt == "По дате регистрации" || sortt == "По дате обновления" || sortt == "По имени" {
		fmt.Print("Не популярные")
		query := `
		WITH MaxSeries AS(
			WITH MaxSezon AS (
						SELECT 
							material_id, MAX(sezon) AS sezon
						FROM videos 
						GROUP BY material_id
					)
					SELECT --distinct on (m.id)
						m.id, v.sezon, max(v.series) series
					FROM materials m
					JOIN MaxSezon ms ON m.id = ms.material_id
					JOIN videos v ON m.id = v.material_id AND ms.sezon = v.sezon
					group by m.id, v.sezon
			)
			Select v.material_id, v.sezon, v.series, v.viewed
			from MaxSeries mss join videos v on mss.id = v.material_id where mss.sezon = v.sezon and mss.series = v.series		
		;	
		`

		rows, err := initializers.ConnPool.Query(context.Background(), query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get series",
			})
			return
		}
		for rows.Next() {
			var seria series
			err = rows.Scan(&seria.Material_id, &seria.Sezon, &seria.Series, &seria.Viewed)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to get series",
				})
				return
			}
			seriesNumber = append(seriesNumber, seria)

		}
		fmt.Print(seriesNumber)
		for _, movie := range seriesNumber {
			if _, found := moviesMap[movie.Material_id]; found {
				/*
					moviesMap[movie.Material_id] = &models.Materials_get{
						Sezon:  movie.Sezon,
						Series: movie.Series,
					}*/

				moviesMap[movie.Material_id].Sezon = movie.Sezon
				moviesMap[movie.Material_id].Series = movie.Series
				moviesMap[movie.Material_id].Viewed = movie.Viewed

			}
		}

		for _, v := range moviesMap {
			result = append(result, *v)
		}

	}

	if sortt == "По дате регистрации" {
		fmt.Print("По дате регистрации")
		sort.Slice(result, func(i, j int) bool {
			return result[i].CreatedAt.After(result[j].CreatedAt)
		})
	}

	if sortt == "По дате обновления" {
		fmt.Print("По дате обновления")
		sort.Slice(result, func(i, j int) bool {
			return result[i].UpdatedAt.After(result[j].UpdatedAt)
		})
	}

	if sortt == "По имени" {
		fmt.Print("По имени")
		sort.Slice(result, func(i, j int) bool {
			return result[i].Title < result[j].Title
		})
	}

	countRows := initializers.ConnPool.QueryRow(context.Background(), countQuery, params...)

	var count int
	err = countRows.Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to count movies",
		})
		return
	}

	yearsQuery := `select distinct publish_year from materials order by publish_year desc`

	yearRows, err := initializers.ConnPool.Query(context.Background(), yearsQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get years",
		})
		return
	}

	var years = []int{}
	for yearRows.Next() {
		var year int
		err := yearRows.Scan(&year)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to define years",
			})
			return
		}
		years = append(years, year)
	}

	c.JSON(http.StatusOK, gin.H{
		"movies": result,
		"count":  count,
		"years":  years,
	})

}
