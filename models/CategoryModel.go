package models

type Category struct {
	CategoryId   uint   `json:"categoryid" gorm:"PrimaryKey; unique"`
	CategoryName string `json:"categoryname"`
}
