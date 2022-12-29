package main

import (
	"github.com/francischacko/ecommerce/initializers"
	"github.com/francischacko/ecommerce/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.Loadvariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	Route := gin.Default()
	routes.UserInfo(Route)
	routes.AdminInfo(Route)
	Route.Run()
}
