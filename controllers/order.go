package controllers

import (
	"net/http"

	"github.com/francischacko/ecommerce/initializers"
	"github.com/francischacko/ecommerce/middlewares"
	"github.com/francischacko/ecommerce/models"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

func PlaceOrder(c *gin.Context) {

	id, err := middlewares.User(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	ToInt := int(id)

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
		var cartitem models.ShoppingCartItem
		initializers.DB.Raw("select product_item_id, quantity, total from shopping_cart_items where id=?", i).Scan(&cartitem)
		var adid int
		initializers.DB.Raw("select id from addresses where user_id=? and default_address=?", ToInt, true).Scan(&adid)
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
		var orderedquantity int
		initializers.DB.Raw("select quantity from shopping_cart_items where id=?", cartitem.ProductItemId).Scan(&orderedquantity)
		UpdatedQuantity := inventory - orderedquantity
		var prod models.Product
		initializers.DB.Raw("update products  set stocks=? where id=?", UpdatedQuantity, cartitem.ProductItemId).Scan(&prod)
	}
	// initializers.DB.Raw("delete id from shopping_cart_items where cid=?",ToInt)
	c.JSON(http.StatusOK, gin.H{"message": "order is placed and inventory have been updated"})

}

func ReturnOrder(c *gin.Context) {
	toInt, err := middlewares.User(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	proid := c.Query("productid")
	var Pro int
	initializers.DB.Raw("select total from shop_orders where product_item_id=? and user_id=?", proid, toInt).Scan(&Pro)

	walletTable := models.Wallet{
		UserId:        int(toInt),
		WalletBalance: Pro,
	}
	result := initializers.DB.Create(&walletTable)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"error": "failed to create wallet table",
		})
		return
	}
	var sto int
	initializers.DB.Raw("select quantity from shop_orders where product_item_id=? and user_id=?", proid, toInt).Scan(&sto)
	var prevq int
	initializers.DB.Raw("select stocks from products where id=?", proid).Scan(&prevq)
	new := sto + prevq
	var proq models.Product
	initializers.DB.Raw("update products set stocks=? where id=?", new, proid).Scan(&proq)
	c.JSON(200, gin.H{
		"msg": "product returned and amount updated to wallet",
	})
}
