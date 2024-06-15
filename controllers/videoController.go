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
