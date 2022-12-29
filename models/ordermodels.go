package models

import "gorm.io/gorm"

type ShopOrders struct {
	gorm.Model
	UserId        int
	OrderId       string
	PayMethod     string `json:"paymethod"`
	AddressId     int
	ProductItemId int
	Total         int
	Quantity      int
	Status        string `json:"status"`
}
type TotalOrders struct {
	gorm.Model
	Cid        int
	OrderId    string
	GrandTotal int
}
