package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/francischacko/ecommerce/config"
	"github.com/francischacko/ecommerce/initializers"
	"github.com/francischacko/ecommerce/middlewares"
	"github.com/francischacko/ecommerce/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// @Summary User signup.
// @Description .Can provide user details to create an accout , later used to login
// @Tags User signup and login
// @Accept application/json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [post]
func Signup(C *gin.Context) {
	// get the email and password required
	var body struct {
		Email    string
		Password string
		Name     string
		Phone    string
	}

	if C.Bind(&body) != nil {
		C.JSON(400, gin.H{
			"error": "failed to load",
		})
		return
	}
	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		C.JSON(400, gin.H{
			"error": "failed to hash password",
		})
		return
	}
	// create user
	user := models.User{Email: body.Email, Name: body.Name, Phone: body.Phone, Password: string(hash)}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		C.JSON(400, gin.H{
			"error": "failed to create user",
		})
		return
	}
	//respond
	C.JSON(http.StatusOK, gin.H{
		"message": "user registered",
	})
}

func Login(C *gin.Context) {
	// get email and password required off the body
	var body struct {
		Email    string
		Password string
	}

	if C.Bind(&body) != nil {
		C.JSON(400, gin.H{
			"error": "failed to load user",
		})
		return
	}

	// lookup requested user
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)
	if user.Id == 0 {
		C.JSON(400, gin.H{
			"error": "failed to lookup user",
		})
		return
	}

	// compare sent in hash  password with saved user hash password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		C.JSON(400, gin.H{
			"error": "invalid user or password",
		})
		return
	}
	if user.Status {
		C.JSON(400, gin.H{
			"error": "you are blocked",
		})
		return
	}
	// generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	secret := config.EnvConf.JWT
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		C.JSON(400, gin.H{
			"error": "failed to create token for user ",
		})
		return
	}
	// send it back
	C.SetSameSite(http.SameSiteLaxMode)
	C.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	C.JSON(http.StatusOK, gin.H{
		"message": user,
		"token":   tokenString,
	})

	// return
}

func GetAllUsers(c *gin.Context) {
	var user []models.User
	initializers.DB.Find(&user)
	for _, i := range user {
		c.JSON(http.StatusOK, gin.H{
			"user id ":    i.Id,
			"user email":  i.Email,
			"user status": i.Status,
		})

	}
}

var page int

func BlockUser(c *gin.Context) {

	params := c.Param("id")

	page, _ = strconv.Atoi(params)
	var user models.User
	var stat string
	initializers.DB.Raw("UPDATE users SET status=true WHERE id=?", page).Scan(&user)
	initializers.DB.Raw("SELECT status FROM users WHERE id=?", page).Scan(&stat)

	c.JSON(http.StatusOK, gin.H{

		"Block status": stat,
	})

}

func UnblockUser(c *gin.Context) {

	params := c.Param("id")
	page, _ = strconv.Atoi(params)
	var user models.User
	var stat string
	initializers.DB.Raw("UPDATE users SET status=false WHERE id=?", page).Scan(&user)
	initializers.DB.Raw("SELECT status FROM users WHERE id=?", page).Scan(&stat)
	c.JSON(http.StatusOK, gin.H{

		"Block status": stat,
	})

}

func ChangePassword(c *gin.Context) {
	id, err := middlewares.User(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	toInt := int(id)
	var body struct {
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(400, gin.H{
			"error": "failed to load",
		})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	var reset models.User
	initializers.DB.Raw("update users set password=? where id=?", string(hash), toInt).Scan(&reset)
}
