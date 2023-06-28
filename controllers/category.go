package controllers

import (
	"net/http"
	"strconv"

	"github.com/francischacko/ecommerce/initializers"
	"github.com/francischacko/ecommerce/models"
	"github.com/gin-gonic/gin"
)

func AddCategory(c *gin.Context) {
	var body struct {
		CategoryName string
	}

	if c.Bind(&body) != nil {
		c.JSON(400, gin.H{
			"error": "failed to bind category body",
		})
		return
	}
	Category := models.Category{CategoryName: body.CategoryName}

	result := initializers.DB.Create(&Category)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"error": "failed to add category",
		})
		return
	}
	//respond
	c.JSON(http.StatusOK, gin.H{
		"message": "Category added",
	})
}

func ListCategories(c *gin.Context) {
	var cat []models.Category
	initializers.DB.Find(&cat)
	for _, i := range cat {
		c.JSON(http.StatusOK, gin.H{
			"Category id ":  i.CategoryId,
			"Category name": i.CategoryName,
		})

	}
}

func DeleteCategory(c *gin.Context) {
	var body struct {
		Id uint
	}
	if c.Bind(&body) != nil {
		c.JSON(400, gin.H{
			"error": "failed to bind category body",
		})
		return
	}
	var category models.Category
	initializers.DB.Raw(" DELETE FROM categories WHERE category_id=?", body.Id).Scan(&category)

	c.JSON(http.StatusOK, gin.H{
		"message": "selected category is deleted",
	})

}

func EditCategory(g *gin.Context) {
	params := g.Param("id")

	page, _ := strconv.Atoi(params)

	var body struct {
		CategoryName string
	}
	if g.Bind(&body) != nil {

		g.JSON(400, gin.H{"error": "failed to bind for category editing"})
		return
	}
	var EditCategory models.Category
	if body.CategoryName != "" {
		initializers.DB.Raw("UPDATE categories SET category_name=? WHERE category_id=?", body.CategoryName, page).Scan(&EditCategory)
		g.JSON(http.StatusOK, gin.H{
			"message": "category name is edited",
		})
	}
}
