package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/francischacko/ecommerce/config"
	"github.com/francischacko/ecommerce/initializers"
	"github.com/francischacko/ecommerce/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func UserAuth(c *gin.Context) {

	// Get the cookie off req
	tokenString, err := c.Cookie("Authorization")

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	// Decode/Validate it

	// Parse takes the token string and a function for looking up the key. The latter is especially

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		secret := config.EnvConf.JWT
		return []byte(secret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		// Check the expm
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// Find the user with token sub
		var user models.User
		initializers.DB.First(&user, claims["sub"])

		if user.Id == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// Attach to request
		c.Set("user", user)

		// Continue

		c.Next()

		fmt.Println(claims["foo"], claims["nbf"])
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}

func AdminAuth(c *gin.Context) {

	// Get the cookie off req
	tokenString, err := c.Cookie("AdminAuthorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	// Decode/Validate it

	// Pzarse takes the token string and a function for looking up the key. The latter is especially

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		secret := config.EnvConf.JWT
		return []byte(secret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		// Check the expm
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		// Find the user with token sub
		var admin models.Admin
		initializers.DB.First(&admin, claims["sub"])
		if admin.ID == 0 {
			fmt.Println("check")
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// Attach to request
		c.Set("admin", admin)

		// Continue

		c.Next()

		// fmt.Println(claims["foo"], claims["nooo"])
	} else {

		c.AbortWithStatus(http.StatusUnauthorized)
	}

}
