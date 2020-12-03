package model

import (
  "gorm.io/gorm"
	"github.com/myrachanto/accounting/httperors"
)
//Payment ..capturing all type of payments
type Payment struct {
	Name string  `json:"name"`
	Description string `json:"description"`
	Company string ` json:"company"`
	CustomerID uint `gorm:"not null" json:"customerid"`
	LiatranID uint `gorm:"not null" json:"liatransid"`
	Customer Customer `gorm:"not null" json:"customer"`
	InvoiceID  uint `gorm:"not null" json:"invoiceid"`
	PaymentMethod string `json:"paymentmethod"`
	Status string `json:"status"`
	Asstrans []Asstrans `gorm:"not null" json:"asstrans"`
	SInvoice []SInvoice `gorm:"many2many:payment_sinvoices"`
	gorm.Model
}
//Validate ...
func (payment Payment) Validate() *httperors.HttpError{ 
	if payment.Name == "" && len(payment.Name) > 3 {
		return httperors.NewNotFoundError("Invalid Name")
	}
	if payment.Description == "" && len(payment.Description) > 3 {
		return httperors.NewNotFoundError("Invalid description")
	}
	if payment.CustomerID == 0 {
		return httperors.NewNotFoundError("Invalid customer")
	}
	if payment.PaymentMethod == "" {
		return httperors.NewNotFoundError("Invalid Payment method")
	}
	return nil
}