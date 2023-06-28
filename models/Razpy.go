package models

import "gorm.io/gorm"

type RazorPay struct {
	gorm.Model
	UserID          uint
	RazorPaymentId  string
	Signature       string
	RazorPayOrderID string
	AmountPaid      uint
}
