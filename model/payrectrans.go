package model

import (
  "gorm.io/gorm"
)
//Payrectrasan all transactions revolving receipts and payments
type Payrectrasan struct {
	Name string `json:"name"`
	Title string `json:"title"`
	Code string 
	CustomerID uint `json:"customerid"`
	Customer Customer `json:"customer"`
	Receipt Receipt  `json:"receipts"`
	ReceiptID uint  `json:"receiptID"`
	SupplierID uint  `json:"supplierid"`
	Supplier Supplier `json:"supplier"`
	Payment Payment `gorm:"foreignKey:PaymentID"`
	PaymentID uint  `json:"paymentID"`
	Description string `json:"description"`
	Amount float64  `json:"amount"`
	Status string  `json:"status"`
	gorm.Model
}