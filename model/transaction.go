package model

import (
  "gorm.io/gorm"
)
//Transaction ...list the contents of an invoice
type Transaction struct {
	Product []*Product `gorm:"many2many:transaction_products;"`
	Productname string `json:"name"`
	InvoiceID uint  `json:"invoiceid"`
	Credit bool  `json:"credit note"`
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
	Usercode string `json:"usercode"`
	Customercode string `json:"customercode"`
	gorm.Model
} 
//DebtTransaction recordeing customer flow
type DebtTransaction struct {
	Code string `gorm:"not null" json:"code"`
	Description string `gorm:"not null" json:"description"`
	Customername string `gorm:"not null" json:"customername"`
	Customers []*Customer `gorm:"many2many:debtTransaction_customers;"`
	Amount float64 `gorm:"not null" json:"amount"`
	gorm.Model
}