package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/francischacko/ecommerce/initializers"
	"github.com/francischacko/ecommerce/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AddUserAddress(c *gin.Context) {
	var body struct {
		UserId         float64
		StreetName     string
		AddressLine1   string
		AddressLine2   string
		City           string
		State          string
		DefaultAddress bool
	}

	if c.Bind(&body) != nil {
		c.JSON(400, gin.H{
			"error": "failed binding body of address entry body",
		})

	}
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})
	claims := token.Claims.(jwt.MapClaims)
	GetId := claims["sub"]
	toInt := GetId.(float64)

	useraddress := models.Address{
		UserId:         toInt,
		StreetName:     body.StreetName,
		AddressLine1:   body.AddressLine1,
		AddressLine2:   body.AddressLine2,
		City:           body.City,
		State:          body.State,
		DefaultAddress: body.DefaultAddress,
	}
	initializers.DB.Create(&useraddress)
	c.JSON(http.StatusOK, gin.H{
		"message": "user address is added",
	})
}

func EditUserAddress(c *gin.Context) {
	var body struct {
		StreetName     string
		AddressLine1   string
		AddressLine2   string
		City           string
		State          string
		DefaultAddress bool
	}
	params := c.Query("id")
	bro, _ := strconv.Atoi(params)
	var address models.Address
	if c.Bind(&body) != nil {
		c.JSON(400, gin.H{"error": "failed to bind while edit product"})
		return
	}
	if body.StreetName != "" {
		initializers.DB.Raw("update addresses SET street_name=? WHERE id=?", body.StreetName, bro).Scan(&address)
	}
	if body.StreetName != "" {
		initializers.DB.Raw("update addresses SET address_line1=? WHERE id=?", body.AddressLine1, bro).Scan(&address)
	}
	if body.StreetName != "" {
		initializers.DB.Raw("update addresses SET address_line2=? WHERE id=?", body.AddressLine2, bro).Scan(&address)
	}
	if body.StreetName != "" {
		initializers.DB.Raw("update addresses SET city=? WHERE id=?", body.City, bro).Scan(&address)
	}
	if body.StreetName != "" {
		initializers.DB.Raw("update addresses SET state=? WHERE id=?", body.State, bro).Scan(&address)
	}
	if !body.DefaultAddress {
		initializers.DB.Raw("update addresses SET default_address=? WHERE id=?", body.StreetName, bro).Scan(&address)
	}
}
func DeleteUserAddress(c *gin.Context) {
	params := c.Query("id")
	bro, _ := strconv.Atoi(params)
	initializers.DB.Raw("Delete id from addresses where id=?", bro)
}

// showuseraddress function shows all the existing addresses of that particular user[user who is logged in]
func ShowUserAddress(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})
	claims := token.Claims.(jwt.MapClaims)
	GetId := claims["sub"]
	toInt := GetId.(float64)
	var address []models.Address
	initializers.DB.Raw("select * from addresses where user_id=?", toInt).Scan(&address)
	c.JSON(http.StatusOK, address)
}

func ChooseAddress(c *gin.Context) {
	params := c.Query("id")
	inp, _ := strconv.Atoi(params)
	var choose models.Address
	initializers.DB.Raw("select * from addresses where id=?", inp).Scan(&choose)
	c.JSON(200, gin.H{
		"address choosed": choose,
	})
}
