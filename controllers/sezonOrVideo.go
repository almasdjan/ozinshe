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

func GetSezonsOrVideo(c *gin.Context) {
	middleware.RequireAuth(c)

	userid, _ := c.Get("user")

	var user models.User
	initializers.DB.First(&user, userid)

	material_id := c.Param("material_id")

	if material_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read materialid",
		})
		return
	}

	db, error := initializers.ConnectDb()
	if error != nil {
		fmt.Println(error)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to connect database",
		})
		return
	}

	row := db.QueryRow(context.Background(), "select count(*) from materials m join material_categories mc on m.id = mc.material_id join categories c on c.id = mc.category_id  where c.category_name like '%сериал%' and m.id = $1", material_id)

	if row == nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get materials info",
		})
		return
	}

	var isSerial int

	err := row.Scan(&isSerial)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get materials from recommended",
		})
		return
	}

	if isSerial == 1 { //если сериал
		rows, err := db.Query(context.Background(), "select distinct sezon from videos where material_id = $1 order by sezon", material_id)

		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to get sezons",
			})
			return
		}

		var sezons []int

		for rows.Next() {
			var sezon int
			err := rows.Scan(&sezon)
			if err != nil {
				fmt.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Failed to get sezons",
				})
				return
			}

			sezons = append(sezons, sezon)
		}

		sezonParam := c.Param("sezon")
		if sezonParam == "" {
			sezonParam = "1"
		}

		rowsSeries, err := db.Query(context.Background(), "select series, image_src, video_src  from videos where material_id = $1 and sezon = $2 order by series ", material_id, sezonParam)

		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to get series",
			})
			return
		}

		var series []models.Series
		for rowsSeries.Next() {
			var series1 models.Series
			err := rowsSeries.Scan(&series1.Series, &series1.Image_src, &series1.Video_src)
			if err != nil {
				fmt.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Failed to output series",
				})
				return
			}
			series = append(series, series1)
		}

		c.JSON(http.StatusOK, gin.H{
			"sezons": sezons,
			"series": series,
		})

	} else if isSerial == 0 { //если не сериал
		row := db.QueryRow(context.Background(), "select video_src from videos where material_id = $1", material_id)

		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to get materials info",
			})
			return
		}

		var src string
		err := row.Scan(&src)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to get materials from recommended",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"video": src,
		})

	}

	//c.JSON(200, gin.H{"status": "success", "data": isSerial, "msg": "get recommends successfully"})

}

func AddViewed(c *gin.Context) {
	middleware.RequireAuth(c)

	userid, _ := c.Get("user")

	var user models.User
	initializers.DB.First(&user, userid)

	material_id := c.Param("material_id")

	if material_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read video info",
		})
		return
	}

	db, error := initializers.ConnectDb()
	if error != nil {
		fmt.Println(error)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to connect database",
		})
		return
	}

	_, err := db.Exec(context.Background(), "update materials set viewed = viewed+1 where id = $1", material_id)
	if err != nil {
		fmt.Println(error)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update viewed",
		})
		return
	}
	c.JSON(200, gin.H{"status": "success"})

}
