package controllers

import (
	"context"
	"fmt"
	"net/http"
	"project1/initializers"
	"project1/middleware"
	"project1/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Add video
// @Security BearerAuth
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "material id"
// @Param videos body []models.Videos false "videos"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /admin/videosrc/{id} [post]
func CreateVideo(c *gin.Context) {
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

	material_idd := c.Param("id")
	var series []models.Videos
	if err := c.ShouldBindJSON(&series); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	material_id, err1 := strconv.Atoi(material_idd)

	if err1 != nil {
		fmt.Print(err1.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read material id",
		})
		return
	}

	for _, v := range series {
		if v.Video_src == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "please, enter video source",
			})
			return
		}
	}

	/*
		if c.Bind(&body) != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to read body",
			})
			return
		}
	*/
	var material models.Material
	exist_material := initializers.DB.Where("id=?", material_id).First(&material)
	if exist_material.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no such material",
		})
		return
	}

	for _, v := range series {
		video := models.Video{Material_id: uint(material_id), Sezon: uint(v.Sezon), Series: uint(v.Series), Video_src: v.Video_src}

		result := initializers.DB.Create(&video)

		if result.Error != nil {
			fmt.Print(result.Error.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to create video",
			})
			return

		}
	}

	c.JSON(http.StatusOK, gin.H{
		"Action": "The video was succfully created",
	})

}

// @Summary delete video
// @Security BearerAuth
// @Tags admin
// @Accept json
// @Produce json
// @Param video_id path string true "video id"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /admin/videosrc/{video_id} [delete]
func DeleteVideo(c *gin.Context) {
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
	/*
		db, error := initializers.ConnectDb()
		if error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to connect database",
			})
			return
		}*/

	video_id := c.Param("video_id")

	_, err := initializers.ConnPool.Exec(context.Background(), `delete from videos WHERE id = $1 `, video_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete video",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"succes": "The video is deleted",
	})

}

// @Summary edit videos
// @Security BearerAuth
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "material id"
// @Param videos body []models.Videos false "videos"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /admin/videosrc/{id} [patch]
func EditVideos(c *gin.Context) {
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

	material_id := c.Param("id")

	row := initializers.ConnPool.QueryRow(context.Background(), "select m_type from materials where id = $1", material_id)

	if row == nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get materials info",
		})
		return
	}

	var isSerial string

	err := row.Scan(&isSerial)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to define serial or not",
		})
		return
	}

	var videos = []models.Videos{}

	if isSerial == "Сериалы" { //если сериал

		video, err := initializers.ConnPool.Query(context.Background(), "select sezon, series, video_src from videos where material_id = $1 order by sezon, series", material_id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed get videos",
			})
		}
		//var videos = []models.Video{}
		fmt.Print(videos)
		if video != nil {
			var videoss models.Videos
			for video.Next() {
				err = video.Scan(&videoss.Sezon, &videoss.Series, &videoss.Video_src)
				if err != nil {
					videos = []models.Videos{}
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

	} else if isSerial == "Фильмы" {
		video := initializers.ConnPool.QueryRow(context.Background(), "select sezon, series, video_src from videos where material_id = $1 AND SEZON = 0 AND SERIES = 0", material_id)

		fmt.Print(videos)
		if video != nil {
			var videoss models.Videos
			err = video.Scan(&videoss.Sezon, &videoss.Series, &videoss.Video_src)
			if err != nil {
				videos = []models.Videos{}
			}
			videos = append(videos, videoss)
		}

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to define serial or not",
		})
		return

	}

	var series []models.Videos
	if err := c.ShouldBindJSON(&series); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, v := range series {
		if v.Video_src == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "please, enter video source",
			})
			return
		}
	}

	videoMap := make(map[uint]map[uint]int)

	// Заполняем карту индексами элементов из второго массива
	for i, video := range videos {
		if _, exists := videoMap[video.Sezon]; !exists {
			videoMap[video.Sezon] = make(map[uint]int)
		}
		videoMap[video.Sezon][video.Series] = i
	}

	// Обрабатываем первый массив
	for _, video1 := range series {
		if seriesMap, exists := videoMap[video1.Sezon]; exists {
			if index, exists := seriesMap[video1.Series]; exists {
				// Элемент существует во втором массиве, проверяем Video_src
				if videos[index].Video_src != video1.Video_src {
					// Обновляем Video_src во втором массиве
					//videos[index].Video_src = video1.Video_src
					_, err := initializers.ConnPool.Exec(context.Background(), "Update videos set video_src = $4 where material_id =$1 and sezon =$2 and series = $3", material_id, video1.Sezon, video1.Series, video1.Video_src)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"error": "failed to update video src",
						})
						return
					}
				}
			} else {
				// Элемент с таким Sezon и Series не существует, добавляем его во второй массив
				videos = append(videos, video1)
				// Обновляем карту
				//videoMap[video1.Sezon][video1.Series] = len(videos) - 1
				_, err := initializers.ConnPool.Exec(context.Background(), "insert into videos (material_id, sezon, series,video_src)  VALUES ($1,$2,$3,$4);", material_id, video1.Sezon, video1.Series, video1.Video_src)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "failed to add new video",
					})
					return
				}
			}
		} else {
			// Элемент с таким Sezon не существует, добавляем его во второй массив
			videos = append(videos, video1)
			// Обновляем карту
			//videoMap[video1.Sezon] = make(map[uint]int)
			//videoMap[video1.Sezon][video1.Series] = len(videos) - 1
			_, err := initializers.ConnPool.Exec(context.Background(), "insert into videos (material_id, sezon, series,video_src)  VALUES ($1,$2,$3,$4);", material_id, video1.Sezon, video1.Series, video1.Video_src)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "failed to add new video",
				})
				return
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})

}
