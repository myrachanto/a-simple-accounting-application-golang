package model

import (
  "gorm.io/gorm"
)
//STransaction ...list the contents of an invoice
type STransaction struct {
	Product []*Product `gorm:"many2many:stransaction_products;"`
	Productname string `json:"name"`
	SInvoiceID uint  `json:"sinvoiceid"`
	Credit bool  `json:"Good return note"`
	Code string `gorm:"not null" json:"code"`
	Title string `gorm:"not null" json:"title"`
	Quantity float64 `gorm:"not null" json:"quantity"`
	Price float64  `gorm:"not null" json:"price"`
	Tax float64  `gorm:"not null" json:"tax"`
	Taxpercent float64  `json:"taxpercent"`
	Subtotal float64  `gorm:"not null" json:"subtotal"`
	Discount float64  `gorm:"not null" json:"discount"`
	Discountpercent float64  `json:"discountpercent"`
	Total float64  `gorm:"not null" json:"total"`
	AmountPaid float64  `gorm:"not null" json:"amountpaid"`
	Balance float64 `gorm:"not null" json:"balance"`
	Status string  `gorm:"not null" json:"status"`
	gorm.Model
} 
//CreditTransaction recordeing customer flow
type CreditTransaction struct {
	Code string `gorm:"not null" json:"code"`
	Description string `gorm:"not null" json:"description"`
	Suppliername string `gorm:"not null" json:"suppliername"`
	Supplier []*Supplier `gorm:"many2many:creditTransaction_suppliers;"`
	Amount float64 `gorm:"not null" json:"amount"`
	gorm.Model
}