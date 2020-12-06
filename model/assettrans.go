package model

import (
  "gorm.io/gorm"
)
//Asstrans ...
type Asstrans struct {
	AssetID uint `gorm:"not null" json:"assetid"`
	Name string `gorm:"not null" json:"name"`
	Title string `gorm:"not null" json:"title"`
	Description string `gorm:"not null" json:"description"`
	PaymentID uint `gorm:"not null" json:"paymentid"`
	Depreciation float64 `gorm:"not null" json:"depreciation"`
	Amount float64 `gorm:"not null" json:"amount"`
	Status bool `gorm:"not null" json:"status"`
	Usercode string `json:"usercode"`
	gorm.Model
}