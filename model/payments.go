package model

import (
	"gorm.io/gorm"
	"time"
	"github.com/myrachanto/accounting/httperors"
)
//Payment ..
type Payment struct {
	SupplierName string `json:"suppliername"`
	Suppliercode string `json:"suppliercode"`
	Usercode string `json:"usercode"`
	Description string `json:"description"` 
	Code string `json:"code"`
	ChequeNo string `json:"chequeno"`
	Expirydate time.Time `json:"expirydate"`
	SupplierID uint `json:"supplierrid"`
	Supplier Supplier  `json:"supplier"`
	SInvoice []SInvoice `gorm:"many2many:payment_sinvoices"`
	Paymentform []Paymentform `gorm:"many2many:payment_paymentforms"`
	Liability []Liability `gorm:"many2many:payment_liabiltys"`
	Payment []Payment `gorm:"many2many:payment_liatrans"`
	Type string `json:"type"`
	ClearanceDate time.Time `json:"clearancedate"`
	Amount float64 `json:"amount"`
	Status string `json:"status"`
	Allocated string `json:"allocated"`
	gorm.Model 
}
//PaymentReport ...
type PaymentReport struct {
	All []Payment `json:"all"`
	ClearedPayments SalesModule `json:"cleared"`
	PendingPayments SalesModule `json:"pending"`
	CanceledPayments SalesModule `json:"canceled"`
}
//PaymentView ..structure to gather dat for payments view
type PaymentView struct {
	Code string `json:"code"`
	Suppliers []Supplier `json:"suppliers"`
	Paymentform []Paymentform `json:"paymentforms"`
}
//PaymentOptions payments view analysis
type PaymentOptions struct {
	AllPayments []Payment `json:"allpayments"`
	ClearedPayments []Payment `json:"cleared"`
	PendingPayments []Payment `json:"pending"`
	CanceledPayments []Payment `json:"canceled"`
}
//PaymentAlloc ..structure to gather dat for payments allocation
type PaymentAlloc struct {
	Payment *Payment `json:"payment"`
	SInvoice []SInvoice `json:"sinvoices"` //unpaid invoices
}
//Validate ..
func (payments Payment) Validate() *httperors.HttpError{ 
	if payments.SupplierName == "" && len(payments.SupplierName) > 3 {
		return httperors.NewNotFoundError("Invalid supplier Name")
	}
	if payments.Description == "" && len(payments.Description) > 3 {
		return httperors.NewNotFoundError("Invalid description")
	}
	return nil
}