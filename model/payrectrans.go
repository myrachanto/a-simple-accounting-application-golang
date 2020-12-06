package model

import (
  "gorm.io/gorm"
)
//Payrectrasan all transactions revolving receipts and payments
type Payrectrasan struct {
	Name string `json:"name"` //customername or suppliername
	Title string `json:"title"` // tilte
	Description string `json:"description"`// description
	CLientcode string `json:"clientcode"` // customercode or suppliercode
	Invoicecode string `json:"invoicecode"` // customerinvoicecode or supplierinvoicecode
	Amount float64  `json:"amount"` //-/+ in regards to the type of trasaction + receipts - payments
	Usercode string `json:"usercode"` //usercode to refer the transaction entry user
	Status string  `json:"status"` //status of of the payment partial or fully paid
	Paymentform string `json:"paymentform"` //status of of the payment partial or fully paid
	CustomerID uint `json:"customerid"`
	Customer Customer `json:"customer"`
	Receipt Receipt  `json:"receipts"`
	ReceiptID uint  `json:"receiptID"`
	SupplierID uint  `json:"supplierid"`
	Supplier Supplier `json:"supplier"`
	Payment Payment `gorm:"foreignKey:PaymentID"`
	PaymentID uint  `json:"paymentID"`
	gorm.Model
}