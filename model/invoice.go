package model

import (
	"time"
  "gorm.io/gorm"
)
//Invoice ...
type Invoice struct {
	CustomerID uint `gorm:"not null" json:"customerid"`
	Customername string `gorm:"not null" json:"name"`
	Code string `json:"code"`
	Title  string ` json:"title"`
	Description string `json:"description"`
	Dated time.Time `gorm:"not null" json:"dated"`
	Duedate time.Time `gorm:"not null" json:"duedate"`
	Subtotal float64  `gorm:"not null" json:"subtotal"`
	Discount float64  `gorm:"not null" json:"discount"`
	Tax float64  `gorm:"not null" json:"tax"`
	Total float64  `gorm:"not null" json:"total"`
	// PaidStatus bool  `gorm:"not null" json:"paidstatus"`
	AmountPaid float64  `json:"amountpaid"`
	Expences float64  `json:"expence"`
	Balance float64 `gorm:"not null" json:"balance"`
	Status string `gorm:"not null" json:"status"`
	Cn bool `gorm:"not null" json:"cn"`
	TransactionID uint `json:"transactionid"`
	Cash string `json:"cash"`
	Terms string `json:"terms"`
	Instructions string `json:"instructions"`
	Receipt []Receipt `gorm:"many2many:invoice_receipts"`
	gorm.Model
}
//InvoiceItems ..structure...
type InvoiceItems struct {
	Description string `json:"description"`
	Qty int `json:"qty"`
	Price float64 `json:"price"`
}
//Cinvoiceoptions structure of the view invoice route
type Cinvoiceoptions struct {
	Code string `json:"code"`
	Customers []Customer `json:"customers"`
	Products []Product `json:"products"`
	Taxs []Tax `json:"taxs"`
	Prices []Price `json:"prices"`
	Discounts []Discount `json:"discounts"`
	Paymentform []Paymentform `json:"paymentforms"`
	Expences []Expence `json:"expences"`
}
//Roptions ..
type Roptions struct {
	Customer []Customer
	Supplier []Supplier
}
//InvoiceView what to gather for viewing invoice
type InvoiceView struct {
	Invoice Invoice `json:"invoices"`
	Customer *Customer `json:"customer"`
	Transactions []Transaction `json:"transactions"`
	ExpencesTrans []Expencetrasan `json:"expences"`
	Credits []Transaction `json:"credits"`
}