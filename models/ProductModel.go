package models

type Product struct {
	Id                 uint `json:"id" gorm:"primaryKey; unique"`
	CategoryId         int  `json:"categoryId" gorm:"foriegnKey"`
	ProductVariationId int  `json:"ProductVariationId" gorm:"foriegnKey"`
	Type               string
	Name               string
	Description        string
	Image              string
	Price              int
	Stocks             int
}

type ProductVariation struct {
	ProductVariationId uint `json:"ProductVariationId" gorm:"PrimaryKey;unique"`
	Name               string
}
