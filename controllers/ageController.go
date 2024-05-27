package controllers

import (
	"context"
	"fmt"
	"net/http"
	"project1/initializers"
	"project1/middleware"
	"project1/models"

	"github.com/gin-gonic/gin"
)

// @Summary Create age category
// @Security BearerAuth
// @Tags admin
// @Description Create age category
// @Accept json
// @Produce json
// @Param age body models.Agejson true "Age category"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /admin/age [post]
func CreateAge(c *gin.Context) {
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

	var body models.Agejson

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var age models.Age
	exist := initializers.DB.Where("age=?", body.Age).First(&age)

	if exist.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This age category is already exists",
		})
		return
	}

	age = models.Age{Age: body.Age}

	result := initializers.DB.Create(&age)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create age category",
		})
		return

	}

	c.JSON(http.StatusOK, gin.H{
		"Action": "The age category was succfully created",
	})

}

// @Summary get movies by age category
// @Security BearerAuth
// @Tags main
// @Accept json
// @Produce json
// @Param age_id path string true "age id"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /main/ages/{age_id} [get]
func GetMovieByAge(c *gin.Context) {
	middleware.RequireAuth(c)

	userid, _ := c.Get("user")
	var user models.User

	initializers.DB.First(&user, userid)

	db, error := initializers.ConnectDb()
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to connect database",
		})
		return
	}

	age_id := c.Param("age_id")

	rows, err := db.Query(context.Background(), `select  m.id,m.poster, m.title, c.age from materials m
	join material_ages mc on m.id = mc.material_id
	join ages c on mc.age_id = c.id
	where c.id = $1`, age_id)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{

			"error": "Failed to get materials info",
		})
		return
	}

	var materials []models.Material_get
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
		"movies by age": materials,
	})
}

// @Summary delete age category
// @Security BearerAuth
// @Tags admin
// @Accept json
// @Produce json
// @Param age_id path string true "age id"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /admin/ages/{age_id} [delete]
func DeleteAge(c *gin.Context) {
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

	id := c.Param("age_id")
	_, err := db.Exec(context.Background(), "DELETE FROM ages WHERE id = $1", id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete age category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "age category was deleted",
	})

}

// @Summary get all age categories
// @Security BearerAuth
// @Tags main
// @Accept json
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /main/ages [get]
func GetAges(c *gin.Context) {
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
	rows, err := db.Query(context.Background(), `select * from ages`)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{

			"error": "Failed to get ages",
		})
		return
	}

	var ages []models.Age
	for rows.Next() {
		var age models.Age
		err := rows.Scan(&age.ID, &age.Age)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to define ages",
			})
			return
		}
		ages = append(ages, age)
	}
	c.JSON(http.StatusOK, gin.H{
		"ages": ages,
	})
}

// @Summary edit age category
// @Security BearerAuth
// @Tags admin
// @Accept json
// @Produce json
// @Param age_id path string true "age id"
// @Param age body models.Agejson true "Age category"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /admin/ages/{age_id} [patch]
func UpdateAge(c *gin.Context) {
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

	id := c.Param("age_id")

	var body models.Agejson

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	age := body.Age

	fmt.Println(age)
	fmt.Println(id)

	_, err := db.Exec(context.Background(), "update ages set age =$1 WHERE id = $2", age, id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			//"error": "Failed to update age category",
			"err": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "age category was updated",
	})

}

// @Summary delete age category from material
// @Security BearerAuth
// @Tags admin
// @Accept json
// @Produce json
// @Param material_id path string true "material id"
// @Param age_id path string true "age id"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /admin/agematerial/{material_id}/{age_id} [delete]
func DeleteAgeFromMaterial(c *gin.Context) {
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

	age := c.Param("age_id")
	id := c.Param("material_id")

	_, err := db.Exec(context.Background(), `delete from material_ages WHERE age_id = $1 and material_id = $2`, age, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete image",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"succes": "The age is deleted",
	})

}

// @Summary add age category to material
// @Security BearerAuth
// @Tags admin
// @Accept json
// @Produce json
// @Param material_id path string true "material id"
// @Param age_id path string true "age id"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /admin/agematerial/{material_id}/{age_id} [post]
func AddAgeToMaterial(c *gin.Context) {
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

	age := c.Param("age_id")

	_, err := db.Exec(context.Background(), `insert into material_ages values ($1,$2)`, id, age)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add age",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "The image is added",
	})

}
