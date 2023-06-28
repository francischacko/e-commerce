package controllers

import (
	"net/http"

	"github.com/francischacko/ecommerce/config"
	"github.com/francischacko/ecommerce/initializers"
	"github.com/francischacko/ecommerce/middlewares"
	"github.com/francischacko/ecommerce/models"
	"github.com/gin-gonic/gin"
	razorpay "github.com/razorpay/razorpay-go"
)

func RazorPay(c *gin.Context) {
	id, err := middlewares.User(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	ToInt = int(id)
	var grandtotal int64
	initializers.DB.Raw("select grand_total from total_orders where cid=?", ToInt).Scan(&grandtotal)
	client := razorpay.NewClient(config.EnvConf.RzpKey, config.EnvConf.RzpSecret)
	data := map[string]interface{}{
		"amount":   grandtotal,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}
	body, err := client.Order.Create(data, nil)

	orderid := body["id"]
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "Error creating orderId",
		})
	}
	var useremail string
	var phone string
	initializers.DB.Raw("SELECT email FROM users WHERE id=?", ToInt).Scan(&useremail)
	initializers.DB.Raw("SELECT phone FROM users WHERE id=?", ToInt).Scan(&phone)

	c.HTML(200, "app.html", gin.H{
		"UserID":       ToInt,
		"total":        grandtotal,
		"orderid":      orderid,
		"amount":       grandtotal,
		"Email":        useremail,
		"Phone_Number": phone,
	})
	if err != nil {
		c.JSON(200, gin.H{
			"msg": orderid,
		})
	}
}

func RazorpaySuccess(c *gin.Context) {

	id, err := middlewares.User(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	ToInt = int(id)

	userID := ToInt
	orderID := c.Query("order_id")
	signature := c.Query("signature")

	var grandtotal int64
	initializers.DB.Raw("select grand_total from total_orders where cid=?", ToInt).Scan(&grandtotal)

	var oid string
	initializers.DB.Raw("SELECT order_id FROM total_orders WHERRE cid=?", userID).First(&oid)

	rpay := models.RazorPay{
		UserID:          uint(userID),
		RazorPaymentId:  orderID,
		Signature:       signature,
		RazorPayOrderID: oid,
		AmountPaid:      uint(grandtotal),
	}
	if err := initializers.DB.Create(&rpay).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}

func Success(c *gin.Context) {

	c.HTML(200, "accepted.html", nil)

}
