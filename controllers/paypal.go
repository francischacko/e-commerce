package controllers

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	paypalsdk "github.com/netlify/PayPal-Go-SDK"
	//  "github.com/plutov/paypal"
	// "github.com/plutov/paypal"
)

func CheckOut(C *gin.Context) {
	clientID := os.Getenv("CLIENT_ID")
	secretID := os.Getenv("PAYAL_SECRET")
	c, err := paypalsdk.NewClient(clientID, secretID, paypalsdk.APIBaseSandBox)
	if err != nil {
		C.JSON(400, gin.H{
			"error": "failed paypal attempt 1 ",
		})
		return
	}
	// accessToken, err :=
	c.GetAccessToken()
	// fmt.Println(accessToken)

	amount := paypalsdk.Amount{
		Total:    "7.00",
		Currency: "USD",
	}
	redirectURI := "http://youtube.com/redirect-uri"
	cancelURI := "http://spotify.com/cancel-uri"
	description := "Description for this payment"
	paymentResult, err := c.CreateDirectPaypalPayment(amount, redirectURI, cancelURI, description)
	if err != nil {
		C.JSON(400, gin.H{
			"error": "failed paypal attempt 1 ",
		})
		return
	}
	fmt.Println(paymentResult)
}
