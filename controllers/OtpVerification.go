package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/francischacko/ecommerce/config"
	"github.com/francischacko/ecommerce/initializers"
	"github.com/francischacko/ecommerce/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/verify/v2"
)

var (
	uname      string
	pword      string
	serviceSid string

	client *twilio.RestClient
)

var Phone string

func Verify(c *gin.Context) {

	uname = config.EnvConf.Twillio1
	pword = config.EnvConf.Twillio2
	serviceSid = config.EnvConf.Twillio3
	client = twilio.NewRestClientWithParams(twilio.ClientParams{

		Username: uname,
		Password: pword,
	})
	params := c.Query("phone")
	var user models.User
	initializers.DB.First(&user, "phone= ?", params)
	if user.Phone == "" {
		c.JSON(400, gin.H{
			"error": "Not Registered",
		})
	}
	c.JSON(http.StatusOK, gin.H{

		"message": "found",
	})
	num := "+91" + params
	Phone = num

	sendOtp(num)
}

func sendOtp(no string) {

	params := &openapi.CreateVerificationParams{}

	params.SetTo(no)
	params.SetChannel("sms")
	serviceSid = config.EnvConf.Twillio3
	resp, err := client.VerifyV2.CreateVerification(serviceSid, params)

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("error")
	} else {
		fmt.Printf("Sent verification '%s'\n", *resp.Sid)

	}
}

func CheckOtp(c *gin.Context) {
	to := c.Query("code")
	var user models.User
	params := &openapi.CreateVerificationCheckParams{}
	params.SetTo(Phone)
	params.SetCode(to)
	serviceSid = config.EnvConf.Twillio3
	resp, err := client.VerifyV2.CreateVerificationCheck(serviceSid, params)

	if err != nil {
		fmt.Println(err.Error())
	} else if *resp.Status == "approved" {
		fmt.Println("Code is Correct!")
		c.JSON(http.StatusOK, gin.H{
			"message": "code has been verified",
		})
		// generate JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.Id,
			"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
		})
		// Sign and get the complete encoded token as a string using the secret
		secret := config.EnvConf.JWT
		tokenString, err := token.SignedString([]byte(secret))
		if err != nil {
			c.JSON(400, gin.H{
				"error": "failed to create token for user ",
			})
			return
		}
		// send it back
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

		c.JSON(http.StatusOK, gin.H{
			"message": user,
			"token":   tokenString,
		})

	} else {
		fmt.Println("Incorrect!")
	}

}
