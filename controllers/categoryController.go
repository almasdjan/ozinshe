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

func CreateCategory(c *gin.Context) {
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
		CategoryName string `json:"category_name"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var category models.Category
	exist := initializers.DB.Where("name=?", body.CategoryName).First(&category)

	if exist.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This category is already exists",
		})
		return
	}

	category = models.Category{CategoryName: body.CategoryName}

	result := initializers.DB.Create(&category)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create category",
		})
		return

	}

	c.JSON(http.StatusOK, gin.H{
		"Action": "The category was succfully created",
	})

}

func DeleteCategory(c *gin.Context) {
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

	id := c.Param("category_id")
	_, err := db.Exec(context.Background(), "DELETE FROM categories WHERE id = $1", id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "the category was deleted",
	})

}

func GetCategories(c *gin.Context) {
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
	rows, err := db.Query(context.Background(), `select * from categories`)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{

			"error": "Failed to get categories",
		})
		return
	}

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.CategoryName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to define categories",
			})
			return
		}
		categories = append(categories, category)
	}
	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}

func UpdateCategories(c *gin.Context) {
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

	id := c.Param("category_id")
	category := c.PostForm("category")
	_, err := db.Exec(context.Background(), "update categories set category_name = $1 WHERE id = $2", category, id)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update the category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "the category category was updated",
	})

}

func DeleteGenreCategoryMaterial(c *gin.Context) {
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

	category := c.Param("category_id")
	id := c.Param("material_id")

	_, err := db.Exec(context.Background(), `delete from material_categories WHERE category_id = $1 and material_id = $2`, category, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete category",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"succes": "The category is deleted",
	})

}

func AddCategoryToMaterial(c *gin.Context) {
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

	category := c.Param("category_id")

	_, err := db.Exec(context.Background(), `insert into material_categories values ($1,$2)`, id, category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to add category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "The category is added",
	})

}
