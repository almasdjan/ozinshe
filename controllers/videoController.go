package controllers

import (
	"context"
	"net/http"
	"project1/initializers"
	"project1/middleware"
	"project1/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Create video
// @Security BearerAuth
// @Tags admin
// @Accept json
// @Produce json
// @Param material_id formData string true "material id"
// @Param sezon formData string false "sezon"
// @Param series formData string false "series"
// @Param video_src formData string true "video_src"
// @Param image_src formData file true "image"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /admin/videosrc [post]
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

	material_idd := c.PostForm("material_id")
	sezonn := c.PostForm("sezon")
	seriess := c.PostForm("series")
	video_src := c.PostForm("video_src")
	image_src, err := c.FormFile("image_src")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read image_src",
		})
		return
	}
	path := "files//videoposters//" + image_src.Filename
	c.SaveUploadedFile(image_src, path)

	if seriess == "" {
		seriess = "1"
	}
	if sezonn == "" {
		sezonn = "1"
	}
	material_id, err1 := strconv.Atoi(material_idd)
	sezon, err2 := strconv.Atoi(sezonn)
	series, err3 := strconv.Atoi(seriess)
	if err1 != nil || err2 != nil || err3 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read sezon or series",
		})
		return
	}

	if video_src == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "please, enter video source",
		})
		return
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

	var video models.Video
	exist := initializers.DB.Where("material_id=?", material_id).Where("sezon = ?", sezon).Where("series = ?", series).First(&video)

	if exist.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This series is already exists",
		})
		return
	}

	video = models.Video{Material_id: uint(material_id), Sezon: uint(sezon), Series: uint(series), Image_src: image_src.Filename, Video_src: video_src}

	result := initializers.DB.Create(&video)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create video",
		})
		return

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

	db, error := initializers.ConnectDb()
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to connect database",
		})
		return
	}

	video_id := c.Param("video_id")

	_, err := db.Exec(context.Background(), `delete from videos WHERE id = $1 `, video_id)
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
