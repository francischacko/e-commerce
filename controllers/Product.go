package controllers

import (
	"net/http"
	"strconv"

	"github.com/francischacko/ecommerce/initializers"
	"github.com/francischacko/ecommerce/models"
	"github.com/gin-gonic/gin"
)

func AddProduct(c *gin.Context) {
	var body struct {
		CategoryId         int
		ProductVariationId int
		Type               string
		Name               string
		Description        string
		Image              string
		Price              int
		Stocks             int
	}

	if c.Bind(&body) != nil {
		c.JSON(400, gin.H{
			"error": "failed to bind Product body",
		})
		return
	}
	Product := models.Product{CategoryId: body.CategoryId, ProductVariationId: body.ProductVariationId, Type: body.Type, Name: body.Name, Description: body.Description, Image: body.Image, Price: body.Price, Stocks: body.Stocks}

	result := initializers.DB.Create(&Product)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"error": "failed to add Product",
		})
		return
	}
	//respond
	c.JSON(http.StatusOK, gin.H{
		"message": "Product added",
	})
}

func DeleteProduct(c *gin.Context) {
	var body struct {
		Id uint
	}
	if c.Bind(&body) != nil {
		c.JSON(400, gin.H{
			"error": "failed to bind Product body",
		})
		return
	}
	var product models.Product
	initializers.DB.Raw(" DELETE FROM products WHERE id=?", body.Id).Scan(&product)

	c.JSON(http.StatusOK, gin.H{
		"message": "Product is deleted",
	})

}

func EditProduct(g *gin.Context) {
	var body struct {
		CategoryId         int
		ProductVariationId int
		Type               string
		Name               string
		Description        string
		Image              string
		Price              int
		Stocks             int
	}

	params := g.Query("id")
	page, _ := strconv.Atoi(params)

	var product models.Product
	if g.Bind(&body) != nil {
		g.JSON(400, gin.H{"error": "failed to bind while edit product"})
		return
	}
	initializers.DB.First(&product, page)
	initializers.DB.Model(&product).Updates(models.Product{
		CategoryId:         body.CategoryId,
		ProductVariationId: body.ProductVariationId,
		Type:               body.Type,
		Name:               body.Name,
		Stocks:             body.Stocks,
		Image:              body.Image,
		Price:              body.Price,
		Description:        body.Description,
	})

}

// product listing with pagination , we can input page number to list products in different pages
func ListProducts(c *gin.Context) {
	pagestring := c.Query("page")
	page, _ := strconv.Atoi(pagestring)
	offset := (page - 1) * 3
	var product []models.Product
	initializers.DB.Limit(3).Offset(offset).Find(&product)
	c.JSON(http.StatusOK, product)
}
