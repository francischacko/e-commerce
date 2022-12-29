package controllers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/francischacko/ecommerce/initializers"
	"github.com/francischacko/ecommerce/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/segmentio/ksuid"
)

func PlaceOrder(c *gin.Context) {

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
	ToInt = int(GetId.(float64))
	var body struct {
		UserId        int
		OrderId       string
		PayMethod     string
		AddressId     int
		ProductItemId int
		Total         int
		Quantity      int
		Status        string
	}
	if c.Bind(&body) != nil {
		c.JSON(400, gin.H{
			"error": "failed to bind order body",
		})
		return
	}
	var proid []int
	initializers.DB.Raw("select id from shopping_cart_items where cid=?", ToInt).Scan(&proid)
	for _, i := range proid {
		orderid := (ksuid.New().String())
		fmt.Println(orderid)
		var cartitem models.ShoppingCartItem
		initializers.DB.Raw("select product_item_id, quantity, total from shopping_cart_items where id=?", i).Scan(&cartitem)
		var adid int
		initializers.DB.Raw("select id from addresses where user_id=? and default_address=true", ToInt).Scan(&adid)
		ordertable := models.ShopOrders{
			UserId:        ToInt,
			OrderId:       orderid,
			PayMethod:     body.PayMethod,
			AddressId:     adid,
			ProductItemId: cartitem.ProductItemId,
			Total:         cartitem.Total,
			Quantity:      cartitem.Quantity,
			Status:        body.Status,
		}

		result := initializers.DB.Create(&ordertable)

		if result.Error != nil {
			c.JSON(400, gin.H{
				"error": "failed to create order table",
			})
			return
		}
		var gt int
		initializers.DB.Raw("select sum(total) from shopping_cart_items where cid=?", ToInt).Scan(&gt)

		totalorders := models.TotalOrders{
			Cid:        ToInt,
			OrderId:    orderid,
			GrandTotal: gt,
		}
		resulta := initializers.DB.Create(&totalorders)

		if resulta.Error != nil {
			c.JSON(400, gin.H{
				"error": "failed to create order table",
			})
			return
		}

		var inventory int
		initializers.DB.Raw("select stocks from products where id=?", cartitem.ProductItemId).Scan(&inventory)
		fmt.Println(inventory)
		var orderedquantity int
		initializers.DB.Raw("select quantity from shopping_cart_items where id=?", cartitem.ProductItemId).Scan(&orderedquantity)
		fmt.Println(orderedquantity)
		UpdatedQuantity := inventory - orderedquantity
		var prod models.Product
		initializers.DB.Raw("update products  set stocks=?", UpdatedQuantity).Scan(&prod)
		fmt.Println(UpdatedQuantity)
	}
	// initializers.DB.Raw("delete id from shopping_cart_items where cid=?",ToInt)
	c.JSON(http.StatusOK, gin.H{"message": "order is placed and inventory have been updated"})

}
