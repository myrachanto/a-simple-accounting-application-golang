package model

import (
	"time"
  "gorm.io/gorm"
)
//SInvoice ... supplier invoice structure
type SInvoice struct {
	SupplierID uint `gorm:"not null" json:"supplierid"`
	Suppliername string `gorm:"not null" json:"name"`
	Usercode string `json:"usercode"`
	Suppliercode string `json:"suppliercode"`
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
	STransactionID uint `json:"stransactionid"`
	Cash string `json:"cash"`
	Terms string `json:"terms"`
	Instructions string `json:"instructions"`
	Payment []Payment `gorm:"many2many:sinvoice_payments"`
	Supplier []Supplier `gorm:"many2many:sinvoice_suppliers"`
	gorm.Model
}
//SInvoiceItems ..structure...
type SInvoiceItems struct {
	Description string `json:"description"`
	Qty int `json:"qty"`
	Price float64 `json:"price"`
}
//Sinvoiceoptions structure of the view invoice route
type Sinvoiceoptions struct {
	Code string `json:"code"`
	Suppliers []Supplier `json:"suppliers"`
	Products []Product `json:"products"`
	Taxs []Tax `json:"taxs"`
	Prices []Price `json:"prices"`
	Discounts []Discount `json:"discounts"`
	Paymentform []Paymentform `json:"paymentforms"`
	Expences []Expence `json:"expences"`
}
//SInvoiceView what to gather for viewing invoice
type SInvoiceView struct {
	SInvoice SInvoice `json:"sinvoices"`
	Supplier *Supplier `json:"supplier"`
	STransactions []STransaction `json:"stransactions"`
	ExpencesTrans []Expencetrasan `json:"expences"`
	Grn []STransaction `json:"grns"`
}