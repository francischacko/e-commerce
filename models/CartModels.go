package models

import "gorm.io/gorm"

type ShoppingCart struct {
	Cid    int
	UserId int `gorm:"foreignKey"`
}

type ShoppingCartItem struct {
	gorm.Model
	Cid           int
	ProductItemId int `gorm:"foriegnKey"`
	ProductName   string
	Quantity      int
	Total         int
}
