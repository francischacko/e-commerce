package controllers

import (
	"net/http"

	"github.com/francischacko/ecommerce/initializers"
	"github.com/francischacko/ecommerce/models"
	"github.com/gin-gonic/gin"
)

func AddBanner(c *gin.Context) {
	var body struct {
		Name        string
		OfferPeriod string
	}

	if c.Bind(&body) != nil {
		c.JSON(400, gin.H{
			"error": "failed to bind banner body",
		})
		return
	}

	Banner := models.Banner{Name: body.Name, OfferPeriod: body.OfferPeriod}
	result := initializers.DB.Create(&Banner)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"error": "failed to add banner",
		})
		return
	}
	//respond
	c.JSON(http.StatusOK, gin.H{
		"message": "Banner added",
	})
}
