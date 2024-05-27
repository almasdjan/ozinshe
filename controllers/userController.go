package controllers

import (
	"net/http"
	"os"
	"project1/initializers"
	"project1/middleware"
	"strconv"

	"project1/models"
	"time"

	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"golang.org/x/crypto/bcrypt"
)

func isEmailValid(email string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(email)
}

func isPasswordValid(password string) bool {
	isnumeric := regexp.MustCompile(`^[0-9]+$`).MatchString(password)
	isLetter := regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(password)

	if len(password) < 6 || isLetter || isnumeric {
		return false
	}
	return true

}

// @Summary SignUp
// @Tags auth
// @Description Create account
// @Accept json
// @Produce json
// @Param user body models.Userjson true "User information"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /auth/signup [post]
func Signup(c *gin.Context) {
	//get user parameters

	var body models.Userjson

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	//check email exists
	var user models.User
	exist := initializers.DB.Where("email=?", body.Email).First(&user)

	if exist.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This email is already exists",
		})
		return
	}

	//check email format
	if !isEmailValid(body.Email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect email format",
		})
		return
	}

	//check password format
	if !isPasswordValid(body.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect password format",
		})
		return

	}

	//check matching two password
	if body.Password != body.Password2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password do not match",
		})
		return

	}

	//hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return

	}

	//create user
	user = models.User{Email: body.Email, Password: string(hash), Name: body.Name, Phone_number: body.Phone_number, Birthday: body.Birthday}
	//user := models.User{Email: body.Email, Password: string(hash)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return

	}

	initializers.DB.First(&user, "email = ?", body.Email)

	//generate jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 5).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create new token",
		})
		return

	}

	//c.Header("Token", tokenString)
	//c.SetSameSite(http.SameSiteLaxMode)
	//c.SetCookie("Authorization", tokenString, 3600*5, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
	//Respond
	//c.JSON(http.StatusOK, gin.H{})

}

// @Summary Login
// @Tags auth
// @Description SignIn
// @Accept json
// @Produce json
// @Param user body models.Login true "User information"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /auth/login [post]
func Login(c *gin.Context) {
	//get user parameters

	var body models.Login

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	//look up for user
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email",
		})
		return

	}

	//compare password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid password",
		})
		return

	}

	//generate jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		//"exp": time.Now().Add(time.Hour * 5).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create new token",
		})
		return

	}

	//c.Header("Token", tokenString)
	//c.SetSameSite(http.SameSiteLaxMode)
	//c.SetCookie("Authorization", tokenString, 3600*5, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})

}

func Logout(c *gin.Context) {
	middleware.RequireAuth(c)

	//tokenString := c.GetHeader("Authorization")
	//tokenString = tokenString[7:]

}

// @Summary User Info
// @Tags auth
// @Description See User Info
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /auth/userinfo [get]
func GetUserInfo(c *gin.Context) {
	middleware.RequireAuth(c)

	userid, _ := c.Get("user")

	var user models.User
	initializers.DB.First(&user, userid)

	c.JSON(http.StatusOK, gin.H{
		"name":         user.Name,
		"email":        user.Email,
		"phone_number": user.Phone_number,
		"birthday":     user.Birthday,
	})

}

// @Summary Update User Info
// @Tags auth
// @Description Update user info
// @Accept json
// @Produce json
// @Param userinfo body models.Userupdate true "User information"
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /auth/userinfo [patch]
func UpdateUserInfo(c *gin.Context) {
	middleware.RequireAuth(c)

	userid, _ := c.Get("user")

	var user models.User
	initializers.DB.First(&user, userid)

	var body models.Userupdate

	c.Bind(&body)

	exist := initializers.DB.Where("email=?", body.Email).First(&user)

	if exist.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This email is already exists",
		})
		return
	}

	if !isEmailValid(body.Email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect email format",
		})
		return
	}

	user.Name = body.Name
	user.Email = body.Email
	user.Phone_number = body.Phone_number
	user.Birthday = body.Birthday

	initializers.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"name":         user.Name,
		"email":        user.Email,
		"phone_number": user.Phone_number,
		"birthday":     user.Birthday,
	})

}

// @Summary Change password
// @Tags auth
// @Description Change password
// @Accept json
// @Produce json
// @Param password body models.Changepasswd true "New password"
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /auth/password [patch]
func ChangePassword(c *gin.Context) {
	middleware.RequireAuth(c)

	userid, _ := c.Get("user")

	var user models.User
	initializers.DB.First(&user, userid)

	var body models.Changepasswd

	c.Bind(&body)

	//check correctness of password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid password",
		})
		return

	}

	//check new password format
	if !isPasswordValid(body.Newpassword) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect password format",
		})
		return

	}

	//check matching two password
	if body.Newpassword != body.Newpassword2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password do not match",
		})
		return

	}

	//get hash of new password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Newpassword), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return

	}
	user.Password = string(hash)
	initializers.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"result": "Password was succesfully changed",
	})

}

// @Summary Delete account
// @Tags auth
// @Description Delete profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /auth/profile [delete]
func DeleteProfile(c *gin.Context) {
	middleware.RequireAuth(c)

	userid, _ := c.Get("user")
	var user models.User
	initializers.DB.Delete(&user, userid)
}

// @Summary Add to favourite
// @Description Add the movie to favourite list
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Param material_id path string true "Material ID of the movie"
// @Router /favourites/{material_id} [post]
func AddFavouriteMovie(c *gin.Context) {
	middleware.RequireAuth(c)

	userid, _ := c.Get("user")

	var user models.User
	initializers.DB.First(&user, userid)

	material_idd := c.Param("material_id")

	material_id, err := strconv.Atoi(material_idd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get material_id",
		})
		return
	}

	user_favourites := models.User_favourites{
		User_id:     user.ID,
		Material_id: uint(material_id)}

	result := initializers.DB.Create(&user_favourites)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create material",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Action": "The material was succfully added to favorites",
	})

}
