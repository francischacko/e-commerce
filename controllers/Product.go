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

	params := g.Param("id")
	page, _ := strconv.Atoi(params)

	var product models.Product
	if g.Bind(&body) != nil {
		g.JSON(400, gin.H{"error": "failed to bind while edit product"})
		return
	}
	if body.CategoryId != 0 {
		initializers.DB.Raw("update products SET CategoryId=? WHERE id=?", body.CategoryId, page).Scan(&product)
	}
	if body.ProductVariationId != 0 {
		initializers.DB.Raw("update products SET ProductVariationId =? WHERE id=?", body.ProductVariationId, page).Scan(&product)
	}
	if body.Type != "" {
		initializers.DB.Raw("update products SET type=? WHERE id=?", body.Type, page).Scan(&product)
	}
	if body.Name != "" {
		initializers.DB.Raw("update products SET productname=? WHERE id=?", body.Name, page).Scan(&product)
	}
	if body.Description != "" {
		initializers.DB.Raw("update products SET description=? WHERE id=?", body.Description, page).Scan(&product)
	}
	if body.Image != "" {
		initializers.DB.Raw("update products SET immage=? WHERE id=?", body.Image, page).Scan(&product)
	}
	if body.Price != 0 {
		initializers.DB.Raw("update products SET price=? WHERE id=?", body.Image, page).Scan(&product)
	}
	if body.Stocks != 0 {
		initializers.DB.Raw("update products SET stocks=? WHERE id=?", body.Image, page).Scan(&product)
	}
}

func ListProducts(c *gin.Context) {
	var product []models.Product
	initializers.DB.Find(&product)
	c.JSON(http.StatusOK, product)
}
