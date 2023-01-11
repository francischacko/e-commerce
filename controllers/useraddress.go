package controllers

import (
	"net/http"
	"strconv"

	"github.com/francischacko/ecommerce/initializers"
	"github.com/francischacko/ecommerce/middlewares"
	"github.com/francischacko/ecommerce/models"
	"github.com/gin-gonic/gin"
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
	toInt := middlewares.User(c)

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
	if c.Bind(&body) != nil {
		c.JSON(400, gin.H{"error": "failed to bind while edit product"})
		return
	}
	params := c.Query("id")
	bro, _ := strconv.Atoi(params)
	var address models.Address
	initializers.DB.First(&address, bro)
	initializers.DB.Model(&address).Updates(models.Address{
		StreetName:     body.StreetName,
		AddressLine1:   body.AddressLine1,
		AddressLine2:   body.AddressLine2,
		City:           body.City,
		State:          body.State,
		DefaultAddress: body.DefaultAddress,
	})

}
func DeleteUserAddress(c *gin.Context) {
	params := c.Query("id")
	bro, _ := strconv.Atoi(params)
	initializers.DB.Raw("Delete id from addresses where id=?", bro)
}

// showuseraddress function shows all the existing addresses of that particular user[user who is logged in]
func ShowUserAddress(c *gin.Context) {
	toInt := middlewares.User(c)
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
