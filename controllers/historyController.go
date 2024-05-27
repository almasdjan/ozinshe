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

func AddHistory(c *gin.Context) {
	middleware.RequireAuth(c)

	userid, _ := c.Get("user")

	var user models.User
	initializers.DB.First(&user, userid)

	material_id := c.Param("material_id")
	/*
		db, error := initializers.ConnectDb()
		if error != nil {
			fmt.Println(error)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to connect db",
			})
			return
		}*/

	_, error := initializers.ConnPool.Exec(context.Background(), "delete from user_history where material_id = $1 and user_id= $2", material_id, user.ID)
	if error != nil {
		fmt.Println(error)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": error,
			"err":   "failed to delete",
		})
		return
	}

	_, err := initializers.ConnPool.Exec(context.Background(), "INSERT INTO user_history (material_id, user_id) VALUES ($1,$2)", material_id, user.ID)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(200, gin.H{"status": "success"})

}

func AddHistoryFunc(userId interface{}, material_id string) error {
	/*
		db, error := initializers.ConnectDb()
		if error != nil {

			return error
		}*/

	_, error := initializers.ConnPool.Exec(context.Background(), "delete from user_history where material_id = $1 and user_id= $2", material_id, userId)
	if error != nil {
		fmt.Println(error)
		return error
	}

	_, err := initializers.ConnPool.Exec(context.Background(), "INSERT INTO user_history (material_id, user_id) VALUES ($1,$2)", material_id, userId)

	if err != nil {
		fmt.Println(err)
		return error
	}

	return nil

}

// @Summary watch history
// @Security BearerAuth
// @Tags main
// @Accept json
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /main/history [get]
func GetMaterialHistory(c *gin.Context) {
	middleware.RequireAuth(c)

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

	rows, err := initializers.ConnPool.Query(context.Background(), "select m.id, m.poster, m.title from user_history us join materials m on m.id = us.material_id where user_id=$1 order by  us.id desc", userid)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get material history",
		})
		return
	}

	var materials []models.Material_history
	for rows.Next() {
		var material models.Material_history
		err := rows.Scan(&material.Id, &material.Poster, &material.Title)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to output materials",
			})
			return
		}
		materials = append(materials, material)
	}

	c.JSON(http.StatusOK, gin.H{
		"history": materials,
	})
}

func GetMaterialHistoryMain(userid interface{}) ([]models.Material_history, error) {

	var user models.User
	initializers.DB.First(&user, userid)
	/*
		db, error := initializers.ConnectDb()
		if error != nil {
			fmt.Println(error)
			return nil, error
		}*/

	rows, err := initializers.ConnPool.Query(context.Background(), "select m.id, m.poster, m.title from user_history us join materials m on m.id = us.material_id where user_id=$1 order by  us.id desc LIMIT 5", userid)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var materials []models.Material_history
	for rows.Next() {
		var material models.Material_history
		err := rows.Scan(&material.Id, &material.Poster, &material.Title)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		materials = append(materials, material)
	}

	return materials, nil
}

const MaxInt = int(MaxUint >> 1)
const MaxUint = ^uint(0)

// @Summary trends
// @Security BearerAuth
// @Tags main
// @Accept json
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /main/trends [get]
func GetTrends(c *gin.Context) {
	middleware.RequireAuth(c)

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
	/*
		if li == 0 {
			limitt := c.Query("limit")
			li, _ := strconv.Atoi(limitt)
			fmt.Println(li)
			if li == 0 {
				li = MaxInt
			}

		}

		fmt.Println(li)*/

	rows, err := initializers.ConnPool.Query(context.Background(), `select id,poster, title, category_name from (
		select distinct on (id) *   from (
			select  m.id,m.poster, m.title, c.category_name, m.viewed from materials m
			join material_categories mc on m.id = mc.material_id
			join categories c on mc.category_id = c.id
		)
	 
	)order by viewed desc limit null `)

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
		"trends": materials,
	})

}

func GetTrendsMain() ([]models.Material_get, error) {
	/*
		db, error := initializers.ConnectDb()
		if error != nil {
			fmt.Println(error)

			return nil, error
		}*/
	/*
		if li == 0 {
			limitt := c.Query("limit")
			li, _ := strconv.Atoi(limitt)
			fmt.Println(li)
			if li == 0 {
				li = MaxInt
			}

		}

		fmt.Println(li)*/

	rows, err := initializers.ConnPool.Query(context.Background(), `select id,poster, title, category_name from (
		select distinct on (id) *   from (
			select  m.id,m.poster, m.title, c.category_name, m.viewed from materials m
			join material_categories mc on m.id = mc.material_id
			join categories c on mc.category_id = c.id
		)
	 
	)order by viewed desc limit 5`)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var materials []models.Material_get
	for rows.Next() {
		var material models.Material_get
		err := rows.Scan(&material.Material_id, &material.Title, &material.Poster, &material.Category)
		if err != nil {
			fmt.Print(err)
			return nil, err
		}
		materials = append(materials, material)
	}
	return materials, err

}
