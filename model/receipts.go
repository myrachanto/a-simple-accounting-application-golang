package model

import (
	"gorm.io/gorm"
	"time"
	"github.com/myrachanto/accounting/httperors"
)
//Receipt ..
type Receipt struct {
	CustomerName string `json:"customername"`
	Customercode string `json:"customercode"`
	Usercode string `json:"usercode"`
	Description string `json:"description"`  
	Code string `json:"code"`
	CustomerID uint `json:"customerid"`
	Customer Customer  `json:"customer"`
	Invoice []Invoice `gorm:"many2many:receipt_invoices"`
	Paymentform []Paymentform `gorm:"many2many:receipt_paymentforms"`
	Type string `json:"type"`
	ClearanceDate time.Time `json:"clearancedate"`
	Amount float64 `json:"amount"`
	Status string `json:"status"`
	Allocated string `json:"allocated"`
	gorm.Model 
}
//ReceiptReport ...
type ReceiptReport struct {
	All []Receipt `json:"all"`
	ClearedRecipts SalesModule `json:"cleared"`
	PendingRecipts SalesModule `json:"pending"`
	CanceledRecipts SalesModule `json:"canceled"`
}
//ReceiptView ..structure to gather dat for receipts view
type ReceiptView struct {
	Code string `json:"code"`
	Customers []Customer `json:"customers"`
	Paymentform []Paymentform `json:"paymentforms"`
}
//ReceiptAlloc ..structure to gather dat for receipts allocation
type ReceiptAlloc struct {
	Receipt *Receipt `json:"receipt"`
	Invoice []Invoice `json:"invoices"` //unpaid invoices
}
//ReceiptOptions receipts view analysis
type ReceiptOptions struct {
	AllRecipts []Receipt `json:"allreceipts"`
	ClearedRecipts []Receipt `json:"cleared"`
	PendingRecipts []Receipt `json:"pending"`
	CanceledRecipts []Receipt `json:"canceled"`
}
//Validate ..
func (receipts Receipt) Validate() *httperors.HttpError{ 
	if receipts.CustomerName == "" && len(receipts.CustomerName) > 3 {
		return httperors.NewNotFoundError("Invalid customer Name")
	}
	if receipts.Description == "" && len(receipts.Description) > 3 {
		return httperors.NewNotFoundError("Invalid description")
	}
	if receipts.CustomerID == 0 {
		return httperors.NewNotFoundError("Invalid customer")
	}
	return nil
}