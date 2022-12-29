package initializers

import (
	"github.com/francischacko/ecommerce/models"
)

func SyncDatabase() {
	DB.AutoMigrate(
		&models.User{},
		&models.Admin{},
		&models.Category{},
		&models.Product{},
		&models.ShoppingCart{},
		&models.ShoppingCartItem{},
		&models.Address{},
		&models.ShopOrders{},
		&models.Charge{},
		&models.Coupen{},
		&models.TotalOrders{},
	)

}
