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

	var body struct {
		Age string `json:"age"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var age models.Age
	exist := initializers.DB.Where("name=?", body.Age).First(&age)

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

func GetMovieByAge(c *gin.Context) {
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

	age_id := c.Param("age_id")

	rows, err := db.Query(context.Background(), `select  m.id,m.poster, m.title, c.age from materials m
	join material_ages mc on m.id = mc.material_id
	join ages c on mc.age_id = c.id
	where c.id = $1`, age_id)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{

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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to connect database",
		})
		return
	}

	id := c.Param("age_id")
	_, err := db.Exec(context.Background(), "DELETE FROM ages WHERE id = $1", id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete age category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "age category was deleted",
	})

}

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
	age := c.PostForm("age")
	_, err := db.Exec(context.Background(), "update ages set age = $1 WHERE id = $2", age, id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update age category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "age category was updated",
	})

}

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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to add age",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "The image is added",
	})

}
