package model

import (
	"github.com/jinzhu/gorm"
)
//Cart ..
type Cart struct {
	Product  []*Product `gorm:"many2many:cart_products;"`
	Code string `json:"code" json:"code"`
	Productcode string `json:"productcode"`
	Name string `gorm:"not null" json:"name"`
	Customername string `gorm:"not null" json:"customername"`
	Customercode string `json:"customercode"`
	Quantity float64 `gorm:"not null" json:"quantity"`
	SPrice float64 `gorm:"not null" son:"sprice"`
	Subtotal float64 `json:"subtotal"`
	Discountpercent float64 `json:"discountpercent"`
	Tax float64 `json:"tax"` 
	Discount float64 `json:"discount"`
	Taxpercent float64 `json:"taxpercent"`
	Usercode string `json:"usercode"`
	Total float64 
	CostPrice float64 `json:"costprice"`
	Cost float64 `json:"cost"`
	Cartstatus bool 
	Picture string 
	gorm.Model
}
//Scart ..
type Scart struct {
	Product  []*Product `gorm:"many2many:scart_products;"`
	Code string `json:"code" json:"code"`
	Name string `gorm:"not null" json:"name"`
	Productcode string `json:"productcode"`
	Suppliername string `gorm:"not null" json:"suppliername"`
	Suppliercode string `json:"suppliercode"`
	Usercode string `json:"usercode"`
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