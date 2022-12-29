package models

import (
	"gorm.io/gorm"
)

type Coupen struct {
	gorm.Model
	CoupenCode string `json:"coupencode"`
	Discount   int    `json:"discount"`
	ExpiryDate int64
	Status     bool `json:"status"`
	MinValue   int  `json:"minvalue"`
}
