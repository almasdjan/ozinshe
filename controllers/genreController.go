package controllers

import (
	"context"
	"fmt"
	"net/http"
	"project1/initializers"
	"project1/middleware"
	"project1/models"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Summary Create genre
// @Security BearerAuth
// @Tags admin
// @Accept json
// @Produce json
// @Param genre formData string true "Genre"
// @Param image formData file true "Image"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /admin/genres [post]
func CreateGenre(c *gin.Context) {
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

	genreName := c.PostForm("genre")
	image, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read the image",
		})
		return
	}

	path := "files//genres//" + image.Filename
	c.SaveUploadedFile(image, path)

	body := models.Genre{
		GenreName: genreName,
		Image:     path,
	}

	var genre models.Genre
	exist := initializers.DB.Where("name=?", body.GenreName).First(&genre)

	if exist.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This genre is already exists",
		})
		return
	}

	genre = models.Genre{GenreName: body.GenreName, Image: path}

	result := initializers.DB.Create(&genre)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create genre",
		})
		return

	}

	c.JSON(http.StatusOK, gin.H{
		"Action": "The genre was succfully created",
	})

}

// @Summary get movies by genre
// @Security BearerAuth
// @Tags main
// @Accept json
// @Produce json
// @Param genre_id path string true "genre id"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /main/genres/{genre_id} [get]
func GetMovieByGenre(c *gin.Context) {
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

	genre_id := c.Param("genre_id")

	rows, err := initializers.ConnPool.Query(context.Background(), `select  m.id,m.title,m.poster,  c.genre_name from materials m
	join material_genres mc on m.id = mc.material_id
	join genres c on mc.genre_id = c.id
	where c.id = $1`, genre_id)

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
		err := rows.Scan(&material.Material_id, &material.Title, &material.Poster, &material.Category)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to define materilas",
			})
			return
		}
		materials = append(materials, material)
	}
	c.JSON(http.StatusOK, gin.H{
		"movies by genre": materials,
	})

}

// @Summary delete genre
// @Security BearerAuth
// @Tags admin
// @Accept json
// @Produce json
// @Param genre_id path string true "genre id"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /admin/genres/{genre_id} [delete]
func DeleteGenre(c *gin.Context) {
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

	id := c.Param("genre_id")
	_, err := initializers.ConnPool.Exec(context.Background(), "DELETE FROM genres WHERE id = $1", id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete genre",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "the genre was deleted",
	})

}

// @Summary get all genres
// @Security BearerAuth
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} []models.Genre
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /admin/genres [get]
func GetGenres(c *gin.Context) {
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
	rows, err := initializers.ConnPool.Query(context.Background(), `select * from genres`)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{

			"error": "Failed to get genre",
		})
		return
	}

	var genres = []models.Genre{}
	for rows.Next() {
		var genre models.Genre
		err := rows.Scan(&genre.ID, &genre.GenreName, &genre.Image)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to define genres",
			})
			return
		}
		genres = append(genres, genre)
	}
	c.JSON(http.StatusOK, genres)

}

// @Summary edit genre
// @Security BearerAuth
// @Tags admin
// @Accept json
// @Produce json
// @Param genre_id path string true "genre id"
// @Param genre formData string false "Genre"
// @Param image formData file false "Image"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /admin/genres/{genre_id} [patch]
func UpdateGenre(c *gin.Context) {
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
	id := c.Param("genre_id")
	getGenreQuery := `select id,genre_name, image from genres where id = $1`
	getGenre := initializers.ConnPool.QueryRow(c, getGenreQuery, id)

	var genre models.Genre
	err := getGenre.Scan(&genre.ID, &genre.GenreName, &genre.Image)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no such genre",
		})
		return
	}
	var emtyGenre models.Genre
	if genre == emtyGenre {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no such genre.",
		})
		return
	}
	genreName := c.PostForm("genre")
	image, err := c.FormFile("image")
	var path string
	if err != nil {
		path = genre.Image
	} else {
		path = "files//genres//" + image.Filename
		c.SaveUploadedFile(image, path)

	}

	input := models.Genre{
		GenreName: genreName,
		Image:     path}

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.GenreName != "" {
		setValues = append(setValues, fmt.Sprintf("genre_name = $%d", argId))
		args = append(args, input.GenreName)
		argId++
	}

	if input.Image != "" {
		setValues = append(setValues, fmt.Sprintf("image = $%d", argId))
		args = append(args, input.Image)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE genres SET %s WHERE id =$%d ",
		setQuery, argId)

	args = append(args, id)

	_, err = initializers.ConnPool.Exec(context.Background(), query, args...)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update the genre",
			//"err": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "the genre was updated",
	})

}

// @Summary delete genre from material
// @Security BearerAuth
// @Tags admin
// @Accept json
// @Produce json
// @Param material_id path string true "material id"
// @Param genre_id path string true "genre id"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /admin/genrematerial/{material_id}/{genre_id} [delete]
func DeleteGenreFromMaterial(c *gin.Context) {
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

	genre := c.Param("genre_id")
	id := c.Param("material_id")

	_, err := initializers.ConnPool.Exec(context.Background(), `delete from material_genres WHERE genre_id = $1 and material_id = $2`, genre, id)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete genre",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"succes": "The genre is deleted",
	})

}

// @Summary add genre to material
// @Security BearerAuth
// @Tags admin
// @Accept json
// @Produce json
// @Param material_id path string true "material id"
// @Param genre_id path string true "genre id"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /admin/genrematerial/{material_id}/{genre_id} [post]
func AddGenreToMaterial(c *gin.Context) {
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

	genre := c.Param("genre_id")

	_, err := initializers.ConnPool.Exec(context.Background(), `insert into material_genres values ($1,$2)`, id, genre)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to add genre",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "The genre is added",
	})

}
