package main

import (
	"log"

	"github.com/francischacko/ecommerce/config"
	_ "github.com/francischacko/ecommerce/docs"
	"github.com/francischacko/ecommerce/initializers"
	"github.com/francischacko/ecommerce/routes"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	config.InitEnvConfigs()
	initializers.ConnectToDb()
	initializers.SyncDatabase()

}

// @title E-Commerce API
// @version 2.0
// @description This is a REST API based on e-commerce logic, done as a personal project
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
// @schemes http

func main() {
	Route := gin.Default()
	Route.LoadHTMLGlob("template/*.html")
	routes.UserInfo(Route)
	routes.AdminInfo(Route)
	url := ginSwagger.URL("http://localhost:3000/swagger/doc.json") // The url pointing to API definition
	Route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	prt := config.EnvConf.LocalServerPort
	if err := Route.Run(prt); err != nil {
		log.Fatal("Error while starting server")
	}

}
