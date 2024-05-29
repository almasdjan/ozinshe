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

// @Summary add to recommended list
// @Security BearerAuth
// @Tags admin
// @Accept json
// @Produce json
// @Param material_id path string true "material id"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /admin/recommends/{material_id} [post]
func AddRecommend(c *gin.Context) {
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
			fmt.Println(error)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to connect database",
			})
			return
		}*/

	material_id := c.Param("material_id")
	_, err := initializers.ConnPool.Exec(context.Background(), "insert into recommends (material_id) values ($1)", material_id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to add material into recommended",
		})
		return

	}

	c.JSON(http.StatusOK, gin.H{
		"Action": "The mareial is succfully added in recommended",
	})

}

// @Summary recommended list
// @Security BearerAuth
// @Tags main
// @Accept json
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /main/recommends [get]
func GetRecommended(c *gin.Context) {
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
			fmt.Println(error)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to connect database",
			})
			return
		}*/
	rows, err := initializers.ConnPool.Query(context.Background(), "select r.material_id, m.poster, m.title, m.description from materials m join recommends r on m.id = r.material_id order by r.queue desc")

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get materials from recommended",
		})
		return
	}
	var materials = []models.Material_recommend{}

	for rows.Next() {
		var materialRecommed models.Material_recommend
		err := rows.Scan(&materialRecommed.Material_id, &materialRecommed.Poster, &materialRecommed.Title, &materialRecommed.Description)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to get materials from recommended",
			})
			return
		}

		materials = append(materials, materialRecommed)

	}
	c.JSON(http.StatusOK, gin.H{
		"recommended": materials,
	})

}

func GetRecommendedData() ([]models.Material_recommend, error) {
	/*
		db, err := initializers.ConnectDb()
		if err != nil {
			return nil, err
		}*/
	rows, err := initializers.ConnPool.Query(context.Background(), "SELECT r.material_id, m.poster, m.title, m.description FROM materials m JOIN recommends r ON m.id = r.material_id ORDER BY r.queue DESC LIMIT 5")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var materials = []models.Material_recommend{}
	for rows.Next() {
		var material models.Material_recommend
		err := rows.Scan(&material.Material_id, &material.Poster, &material.Title, &material.Description)
		if err != nil {
			return nil, err
		}
		materials = append(materials, material)
	}
	return materials, nil
}

// @Summary Get random movies
// @Security BearerAuth
// @Tags main
// @Accept json
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /main/foryou [get]
func GetRandomMovie(c *gin.Context) {
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
			fmt.Println(error)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to connect database",
			})
			return
		}*/
	rows, err := initializers.ConnPool.Query(context.Background(), `select id, title,poster,
	
	
	
	
	category_name from (
		select distinct on (id) *   from (
			select  m.id,m.poster, m.title, c.category_name, m.viewed from materials m
			join material_categories mc on m.id = mc.material_id
			join categories c on mc.category_id = c.id
		)
	 
	)order by random()`)

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
		"material for you": materials,
	})

}

func GetRandomMovieMain() ([]models.Material_get, error) {
	//middleware.RequireAuth(c)

	//userid, _ := c.Get("user")

	//var user models.User
	//initializers.DB.First(&user, userid)
	/*
		db, error := initializers.ConnectDb()
		if error != nil {
			fmt.Println(error)

			return nil, error
		}*/
	rows, err := initializers.ConnPool.Query(context.Background(), `select id, title, poster,category_name from (
		select distinct on (id) *   from (
			select  m.id,m.poster, m.title, c.category_name, m.viewed from materials m
			join material_categories mc on m.id = mc.material_id
			join categories c on mc.category_id = c.id
		)
	 
	)order by random() LIMIT 5 `)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var materials = []models.Material_get{}
	for rows.Next() {
		var material models.Material_get
		err := rows.Scan(&material.Material_id, &material.Title, &material.Poster, &material.Category)
		if err != nil {
			return nil, err
		}
		materials = append(materials, material)
	}
	return materials, err

}

// @Summary form recommended list
// @Security BearerAuth
// @Tags admin
// @Accept json
// @Produce json
// @Param material_id path string true "material id"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /admin/recommends/{material_id} [delete]
func DeleteFromRecommended(c *gin.Context) {
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
			fmt.Println(error)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to connect database",
			})
			return
		}*/

	material := c.Param("material_id")
	_, err := initializers.ConnPool.Exec(context.Background(), `delete from recommends where material_id = $1`, material)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete the material from recommends",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "the material was deleted",
	})
}

func UpdateRecommended(c *gin.Context) {
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
			fmt.Println(error)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to connect database",
			})
			return
		}*/

	queue := c.Param("queue")
	material_id := c.Param("material_id")
	_, err := initializers.ConnPool.Exec(context.Background(), `update recommends set material_id = $1 where queue = $2 `, material_id, queue)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update the material in recommends",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "the material was edited",
	})
}
