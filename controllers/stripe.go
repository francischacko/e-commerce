package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/francischacko/ecommerce/config"
	"github.com/francischacko/ecommerce/initializers"
	"github.com/francischacko/ecommerce/models"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
)

func StripePayment(c *gin.Context) {

	var payment models.Charge
	c.BindJSON(&payment)
	var grandtotal int64
	initializers.DB.Raw("select grand_total from total_orders where cid=?", ToInt).Scan(&grandtotal)
	var pname []string
	initializers.DB.Raw("select product_name from shopping_cart_items where cid=?", ToInt).Scan(&pname)
	fmt.Println(ToInt)
	descript := strings.Join(pname, ",")
	apiKey := config.EnvConf.StripeKey
	fmt.Println(apiKey + "asads")
	stripe.Key = apiKey
	_, err := charge.New(&stripe.ChargeParams{
		Amount:       &grandtotal,
		Currency:     stripe.String(string(stripe.CurrencyUSD)),
		Description:  &descript,
		Source:       &stripe.SourceParams{Token: stripe.String("tok_visa")},
		ReceiptEmail: stripe.String(payment.ReceiptEmail),
	})

	if err != nil {
		c.String(http.StatusBadRequest, "Payment Unsuccessfull")
		return
	}

	err1 := SavePayment(&payment)
	if err1 != nil {
		c.String(http.StatusNotFound, "error occured")
	} else {
		c.JSON(200, gin.H{
			"message": "payment succesfull",
		})
	}

}

func SavePayment(charge *models.Charge) (err error) {
	if err = initializers.DB.Create(charge).Error; err != nil {
		return err
	}
	return nil

}
