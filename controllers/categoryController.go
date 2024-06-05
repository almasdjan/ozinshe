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

// @Summary Create category
// @Security BearerAuth
// @Tags admin
// @Description Create category
// @Accept json
// @Produce json
// @Param category body models.Categoryjson true "Category"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /admin/categories [post]
func CreateCategory(c *gin.Context) {
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

	var body models.Categoryjson

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var category models.Category
	exist := initializers.DB.Where("name=?", body.Category).First(&category)

	if exist.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This category is already exists",
		})
		return
	}

	category = models.Category{CategoryName: body.Category}

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

// @Summary delete category
// @Security BearerAuth
// @Tags admin
// @Accept json
// @Produce json
// @Param category_id path string true "category id"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /admin/categories/{category_id} [delete]
func DeleteCategory(c *gin.Context) {
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

	id := c.Param("category_id")
	_, err := initializers.ConnPool.Exec(context.Background(), "DELETE FROM categories WHERE id = $1", id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "the category was deleted",
	})

}

// @Summary get all categories
// @Security BearerAuth
// @Tags main
// @Accept json
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /main/categories [get]
func GetCategories(c *gin.Context) {
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
	rows, err := initializers.ConnPool.Query(context.Background(), `select * from categories`)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{

			"error": "Failed to get categories",
		})
		return
	}

	var categories = []models.Category{}
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

// @Summary edit category
// @Security BearerAuth
// @Tags admin
// @Accept json
// @Produce json
// @Param category_id path string true "category id"
// @Param category body models.Categoryjson true "Category"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /admin/categories/{category_id} [patch]
func UpdateCategories(c *gin.Context) {
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

	id := c.Param("category_id")
	var body models.Categoryjson
	getCategoryQuery := `select category_name from categories where id = $1`
	getCategory := initializers.ConnPool.QueryRow(c, getCategoryQuery, id)

	var categoryName models.Categoryjson
	err := getCategory.Scan(&categoryName.Category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no such category",
		})
		return
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	category := body.Category

	_, err = initializers.ConnPool.Exec(context.Background(), "update categories set category_name = $1 WHERE id = $2", category, id)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update the category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "the category category was updated",
	})

}

// @Summary delete category from material
// @Security BearerAuth
// @Tags admin
// @Accept json
// @Produce json
// @Param material_id path string true "material id"
// @Param category_id path string true "category id"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /admin/categorymaterial/{material_id}/{category_id} [delete]
func DeleteGenreCategoryMaterial(c *gin.Context) {
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

	category := c.Param("category_id")
	id := c.Param("material_id")

	_, err := initializers.ConnPool.Exec(context.Background(), `delete from material_categories WHERE category_id = $1 and material_id = $2`, category, id)
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

// @Summary add category to material
// @Security BearerAuth
// @Tags admin
// @Accept json
// @Produce json
// @Param material_id path string true "material id"
// @Param category_id path string true "category id"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /admin/categorymaterial/{material_id}/{category_id} [post]
func AddCategoryToMaterial(c *gin.Context) {
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

	category := c.Param("category_id")

	_, err := initializers.ConnPool.Exec(context.Background(), `insert into material_categories values ($1,$2)`, id, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "The category is added",
	})

}
