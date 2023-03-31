package main

import (
	"github.com/francischacko/ecommerce/controllers"
	docs "github.com/francischacko/ecommerce/docs"
	"github.com/francischacko/ecommerce/initializers"
	"github.com/francischacko/ecommerce/routes"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	initializers.Loadvariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	Route := gin.Default()
	docs.SwaggerInfo.Title = "E-Commerce API"
	docs.SwaggerInfo.Description = "E-commerce API which in written in GO."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8000"
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	routes.UserInfo(Route))
	routes.AdminInfo(Route)
	Route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	Route.GET("/home", controllers.Hello)
	Route.Run()

}
