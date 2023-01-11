package models

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model
	UserId        int
	WalletBalance int
}
