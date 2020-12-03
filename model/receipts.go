package model

import (
	"gorm.io/gorm"
	"time"
	"github.com/myrachanto/accounting/httperors"
)
//Receipt ..
type Receipt struct {
	CustomerName string `json:"customername"`
	Description string `json:"description"` 
	Code string `json:"code"`
	CustomerID uint `json:"customerid"`
	Customer Customer  `json:"customer"`
	Invoice []Invoice `gorm:"many2many:receipt_invoices"`
	Paymentform Paymentform `gorm:"not null"`
	Type string `json:"type"`
	ClearanceDate time.Time `json:"clearancedate"`
	Amount float64 `json:"amount"`
	Status string `json:"status"`
	gorm.Model
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