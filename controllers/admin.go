package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/francischacko/ecommerce/initializers"
	"github.com/francischacko/ecommerce/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func SignupAdmin(C *gin.Context) {
	// get the email and password required
	var body struct {
		Email    string
		Password string
	}

	if C.Bind(&body) != nil {
		C.JSON(400, gin.H{
			"error": "failed to load admin details",
		})
		return
	}
	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		C.JSON(400, gin.H{
			"error": "failed to hash password of admin",
		})
		return
	}
	// create admin
	admin := models.Admin{Email: body.Email, Password: string(hash)}

	result := initializers.DB.Create(&admin)

	if result.Error != nil {
		C.JSON(400, gin.H{
			"error": "failed to create admin",
		})
		return
	}
	//respond
	C.JSON(http.StatusOK, gin.H{
		"message": "admin registered",
	})
}

func LoginAdmin(C *gin.Context) {
	// get email and password required off the body
	var body struct {
		Email    string
		Password string
	}

	if C.Bind(&body) != nil {
		C.JSON(400, gin.H{
			"error": "failed to load admin while logging in",
		})
		return
	}

	// lookup requested user
	var admin models.Admin
	initializers.DB.First(&admin, "email = ?", body.Email)
	if admin.ID == 0 {
		C.JSON(400, gin.H{
			"error": "invalid admin lookup",
		})
		return
	}

	// compare sent in hash  password with saved user hash password
	err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(body.Password))

	if err != nil {
		C.JSON(400, gin.H{
			"error": "invalid admin credentials",
		})
		return
	}
	// generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": admin.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		C.JSON(400, gin.H{
			"error": "failed to create token for admin",
		})
		return
	}
	// send it back
	C.SetSameSite(http.SameSiteLaxMode)
	C.SetCookie("AdminAuthorization", tokenString, 3600*24*30, "", "", false, true)

	C.JSON(http.StatusOK, gin.H{
		"message": admin,
		"token":   tokenString,
	})

	// return
}

func ValidateAdmin(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "admin is logged in.",
	})

}

func ListOrders(c *gin.Context) {
	var orders []models.ShopOrders
	initializers.DB.Raw("select * from shop_orders").Scan(&orders)
	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
	})
}

func CancelOrder(c *gin.Context) {
	id := c.Query("id")
	initializers.DB.Raw("delete from shop_orders where id=?", id)
	var cancelledorderquantity int
	initializers.DB.Raw("select quantity from shop_orders where id=?", id).Scan(&cancelledorderquantity)
	var productquantity int
	initializers.DB.Raw("select quantity from products where id=?", id).Scan(&productquantity)
	quantitytobeadded := productquantity + cancelledorderquantity
	var prod models.Product
	initializers.DB.Raw("update products set quantity=?", quantitytobeadded).Scan(&prod)
	c.JSON(200, gin.H{
		"message": "Order has been cancelled and updated the inventory",
	})
	change := "Order cancelled"
	var ord models.ShopOrders
	initializers.DB.Raw("update shop_orders set status=? where id=?", change, id).Scan(ord)
}
