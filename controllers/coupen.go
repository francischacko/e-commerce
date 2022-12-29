package controllers

import (
	"fmt"
	"net/http"
	"os"

	"time"

	"github.com/francischacko/ecommerce/initializers"
	"github.com/francischacko/ecommerce/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AddCoupen(c *gin.Context) {
	var body struct {
		CoupenCode string
		Discount   int
		ExpiryDate int64
		Status     bool
		MinValue   int
	}

	if c.Bind(&body) != nil {
		c.JSON(400, gin.H{
			"error": "failed to bind coupen body",
		})
		return
	}
	Coupen := models.Coupen{CoupenCode: body.CoupenCode, Discount: body.Discount, ExpiryDate: time.Now().Add(time.Hour * 24 * 30).Unix(), Status: body.Status, MinValue: body.MinValue}

	result := initializers.DB.Create(&Coupen)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"error": "failed to add coupen",
		})
		return
	}
	//respond
	c.JSON(http.StatusOK, gin.H{
		"message": "Coupen added",
	})
}
func RedeemCoupen(c *gin.Context) {
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
	Coupen := c.Query("code")
	// Coupe,_ := strconv.Atoi(Coupen)
	var coup models.Coupen
	// var coupe []string
	initializers.DB.Raw("select *from coupens where coupen_code=?", Coupen).Scan(&coup)
	if coup.ID == 0 {
		c.JSON(400, gin.H{
			"msg": "coupen does not exist",
		})

	}
	fmt.Println(coup)
	fmt.Println(toInt)

	var totl int
	var current int
	initializers.DB.Raw("select sum(total) from shopping_cart_items where cid=?", toInt).Scan(&totl)
	fmt.Println(totl)
	fmt.Println(time.Now().Unix())
	fmt.Println(coup.ExpiryDate)
	if !coup.Status && (time.Now().Unix()) < coup.ExpiryDate && totl >= coup.MinValue {
		fmt.Println("conditions working")
		current = totl - coup.Discount
		var upd models.TotalOrders
		initializers.DB.Raw("update total_orders set grand_total=? where cid=?", current, toInt).Scan(&upd)

	} else {
		c.JSON(400, gin.H{
			"msg": "coupon not valid",
		})
	}
	c.JSON(200, gin.H{
		"total":    totl,
		"discount": current,
	})
}
