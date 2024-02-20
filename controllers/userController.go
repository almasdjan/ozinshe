package controllers

import (
	"net/http"
	"os"
	"project1/initializers"
	"project1/middleware"

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

func Signup(c *gin.Context) {
	//get user parameters

	var body struct {
		Email        string `json:"email"`
		Password     string `json:"password"`
		Password2    string `json:"password2"`
		Name         string `json:"name"`
		Phone_number string `json:"phone_number"`
		Birthday     string `json:"birthday"`
	}

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
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*5, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
	//Respond
	//c.JSON(http.StatusOK, gin.H{})

}

func Login(c *gin.Context) {
	//get user parameters

	var body struct {
		Email    string
		Password string
	}

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
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*5, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})

}

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

func UpdateUserInfo(c *gin.Context) {
	middleware.RequireAuth(c)

	userid, _ := c.Get("user")

	var user models.User
	initializers.DB.First(&user, userid)

	var body struct {
		Email        string `json:"email"`
		Name         string `json:"name"`
		Phone_number string `json:"phone_number"`
		Birthday     string `json:"birthday"`
	}

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

func ChangePassword(c *gin.Context) {
	middleware.RequireAuth(c)

	userid, _ := c.Get("user")

	var user models.User
	initializers.DB.First(&user, userid)

	var body struct {
		Password     string `json:"password"`
		Newpassword  string `json:"newpassword"`
		Newpassword2 string `json:"newpassword2"`
	}

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

func Logout(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "/", "localhost", false, true)
	c.JSON(200, gin.H{"success": "user logged out"})
}

func DeleteProfile(c *gin.Context) {
	middleware.RequireAuth(c)

	userid, _ := c.Get("user")
	var user models.User
	initializers.DB.Delete(&user, userid)
}
