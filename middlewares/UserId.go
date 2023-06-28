package middlewares

import (
	"fmt"
	"net/http"

	"github.com/francischacko/ecommerce/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func User(c *gin.Context) (float64, error) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		secret := config.EnvConf.JWT
		return []byte(secret), nil
	})

	claims := token.Claims.(jwt.MapClaims)
	GetId := claims["sub"]
	toInt := GetId.(float64)
	fmt.Println(toInt)
	return toInt, nil
}
