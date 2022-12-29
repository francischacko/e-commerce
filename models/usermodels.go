package models

type User struct {
	Id       uint   `json:"id" gorm:"primaryKey;unique"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique; not null"`
	Phone    string `json:"phone" validate:"required,min=4 ,max=10" gorm:"unique"`
	Password string `json:"password" gorm:"not null"`
	Status   bool   `json:"status"`
}

type Address struct {
	Id             uint `json:"id"`
	UserId         float64
	StreetName     string `json:"streetname,omitempty"`
	AddressLine1   string `json:"addressLine1" gorm:"not null"`
	AddressLine2   string `json:"addressLine2,omitempty"`
	DefaultAddress bool   `json:"default,omitempty"`
	City           string `json:"city" gorm:"not null"`
	State          string `json:"state" gorm:"not null"`
}
