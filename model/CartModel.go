package model

import (
	"github.com/jinzhu/gorm"
)
//Cart ..
type Cart struct {
	Product  []*Product `gorm:"many2many:cart_products;"`
	Code string `json:"code" json:"code"`
	Name string `gorm:"not null" json:"name"`
	Customername string `gorm:"not null" json:"customername"`
	Quantity float64 `gorm:"not null" json:"quantity"`
	SPrice float64 `gorm:"not null" son:"sprice"`
	Subtotal float64 `json:"subtotal"`
	Discountpercent float64 `json:"discountpercent"`
	Tax float64 `json:"tax"` 
	Discount float64 `json:"discount"`
	Taxpercent float64 `json:"taxpercent"`
	Total float64 
	Cartstatus bool 
	Picture string 
	gorm.Model
}
//Scart ..
type Scart struct {
	Product  []*Product `gorm:"many2many:scart_products;"`
	Code string `json:"code" json:"code"`
	Name string `gorm:"not null" json:"name"`
	Suppliername string `gorm:"not null" json:"suppliername"`
	Quantity float64 `gorm:"not null" json:"quantity"`
	BPrice float64 `gorm:"not null" son:"bprice"`
	Subtotal float64 `json:"subtotal"`
	Discountpercent float64 `json:"discountpercent"`
	Tax float64 `json:"tax"` 
	Discount float64 `json:"discount"`
	Taxpercent float64 `json:"taxpercent"`
	Total float64 
	SCartstatus bool 
	Picture string 
	gorm.Model
}