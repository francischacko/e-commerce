package controllers

import (
	"fmt"
	"net/http"

	"github.com/francischacko/ecommerce/initializers"
	"github.com/francischacko/ecommerce/middlewares"
	"github.com/francischacko/ecommerce/models"
	"github.com/gin-gonic/gin"
)

var ToInt int
var bro int
var Price int
var Qty int

func AddToCart(c *gin.Context) {

	id, err := middlewares.User(c)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	ToInt := int(id)
	cart := models.ShoppingCart{
		UserId: ToInt,
		Cid:    ToInt,
	}
	result := initializers.DB.Create(&cart)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"error": "failed to create type ShoppingCart",
		})
		return
	}

	var body struct {
		Cid           int
		ProductItemId int
		ProductName   string
		Quantity      int
		Total         int
	}
	var price int
	initializers.DB.Raw("SELECT price FROM products WHERE id =?", body.ProductItemId).Scan(&price)
	initializers.DB.Raw(" SELECT cid FROM shopping_carts WHERE user_id=?", ToInt).Scan(&bro)

	if c.Bind(&body) != nil {
		c.JSON(400, gin.H{
			"error": "failed to bind ShoppingCartItem body",
		})
		return
	}
	var StockCheck int
	CartItems := models.ShoppingCartItem{Cid: bro, ProductItemId: body.ProductItemId, Quantity: body.Quantity, Total: body.Quantity * price}
	initializers.DB.Raw("SELECT stocks FROM products WHERE id=? ", body.ProductItemId).Scan(&StockCheck)

	if body.Quantity > StockCheck {

		fmt.Println(body.Quantity)
		c.JSON(400, gin.H{
			"error": "not enough stock",
		})
		return
	} else {
		resulta := initializers.DB.Create(&CartItems)

		if resulta.Error != nil {
			c.JSON(400, gin.H{
				"error": "failed to add to cart",
			})
			return
		}
	}
	//respond
	c.JSON(http.StatusOK, gin.H{
		"message": "Product is added to cart ",
	})
	var pname string
	initializers.DB.Raw("select name from products where id=?", body.ProductItemId).Scan(&pname)
	initializers.DB.Raw("SELECT price FROM products WHERE id=?", body.ProductItemId).Scan(&Price)
	initializers.DB.Raw("SELECT quantity FROM shopping_cart_items WHERE  product_item_id=?", body.ProductItemId).Scan(&Qty)
	var crtname models.ShoppingCartItem
	initializers.DB.Raw("update shopping_cart_items set product_name=? where product_item_id=?", pname, body.ProductItemId).Scan(&crtname)

	var update models.ShoppingCartItem
	resulta := Price * Qty
	initializers.DB.Raw("update shopping_cart_items set total=? where id in(select max(id) from shopping_cart_items)", resulta).Scan(&update)
}

func QuantityUpdation(c *gin.Context) {

	var body struct {
		ProductItemId int
		Quantity      int
	}
	if c.Bind(&body) != nil {
		c.JSON(400, gin.H{
			"error": "failed to bind ShoppingCartItem body for quantity updation",
		})
		return
	}

	var shoppingcartitem models.ShoppingCartItem
	if body.Quantity != 0 {
		initializers.DB.Raw("update shopping_cart_items SET quantity=? WHERE product_item_id=?", body.Quantity, body.ProductItemId).Scan(&shoppingcartitem)
	}

}

func ListCart(c *gin.Context) {
	toInt, err := middlewares.User(c)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	var cartitem []models.ShoppingCartItem
	initializers.DB.Raw("select * from shopping_cart_items where cid=?", toInt).Scan(&cartitem)
	c.JSON(http.StatusOK, gin.H{
		"items": cartitem,
	})

	var grandtotal int
	initializers.DB.Raw("select sum(total) from shopping_cart_items").Scan(&grandtotal)
	c.JSON(http.StatusOK, gin.H{
		"Grand Total": grandtotal,
	})

}

func CartItemDeletion(c *gin.Context) {
	var body struct {
		ProductItemId int
	}
	if c.Bind(&body) != nil {
		c.JSON(400, gin.H{
			"error": "failed to bind ShoppingCartItem body for cart item deletion",
		})
		return
	}
	var itemdeletion models.ShoppingCartItem
	initializers.DB.Raw("DELETE  FROM shopping_cart_items  WHERE product_item_id=?", body.ProductItemId).Scan(&itemdeletion)

}
