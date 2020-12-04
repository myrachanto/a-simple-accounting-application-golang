package model

import (
  "gorm.io/gorm"
)
//Liatran is a structure to capture all transactions involving liabilities
type Liatran struct {
	Name string `gorm:"not null" json:"name"`
	LiabilityID uint  `gorm:"not null" json:"liabilityid"`
	Title string `gorm:"not null" json:"title"`
	Description string `gorm:"not null" json:"description"`
	Payment []Payment `gorm:"many2many:liatran_payments"`
	Interest float64 `gorm:"not null" json:"interest"`
	Amountpaid float64  `gorm:"not null" json:"amountpaid"`
	Balance float64 `gorm:"not null" json:"balance"`
	Status bool `gorm:"not null" json:"status"`
	gorm.Model
}