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
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/verify/v2"
)

var client *twilio.RestClient = twilio.NewRestClientWithParams(twilio.ClientParams{
	Username: os.Getenv("TWILIO_ACCOUNT_SID"),
	Password: os.Getenv("TWILIO_AUTH_TOKEN"),
})

var Phone string

func Verify(c *gin.Context) {

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
	fmt.Println(params)

	resp, err := client.VerifyV2.CreateVerification(os.Getenv("VERIFY_SERVICE_SID"), params)

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("error")
	} else {
		fmt.Printf("Sent verification '%s'\n", *resp.Sid)

	}
}

// func CodeVerify(c *gin.Context) {
// 	params := c.Param("code")
// 	checkOtp(params)
// }

func CheckOtp(c *gin.Context) {
	to := c.Query("code")
	var user models.User
	params := &openapi.CreateVerificationCheckParams{}
	params.SetTo(Phone)
	params.SetCode(to)

	resp, err := client.VerifyV2.CreateVerificationCheck(os.Getenv("VERIFY_SERVICE_SID"), params)

	if err != nil {
		fmt.Println(err.Error())
	} else if *resp.Status == "approved" {
		fmt.Println("Correct!")
		// generate JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.Id,
			"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
		})
		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
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
