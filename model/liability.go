package model

import (
  "gorm.io/gorm"
	"github.com/myrachanto/accounting/httperors"
)
//Liability ..s
type Liability struct {
	Name string  `gorm:"not null" json:"name"`
	Description string  `gorm:"not null" json:"description"`
	Creditor string  `gorm:"not null" json:"creditor"`
	Approvedby string  `gorm:"not null" json:"approvedby"`
	Amount float64  `gorm:"not null" json:"amount"`
	Interestrate float64 `gorm:"not null" json:"interestrate"`
	Paymentperiod float64 `gorm:"not null" json:"paymentperiod"`
	Amoutinterest float64 `gorm:"not null" json:"amountinterest"`
	Monthlypayment float64 `gorm:"not null" json:"monthlypayment"`
	Liatran []Liatran `gorm:"not null" json:"liatrans"`
	Payment []Payment `gorm:"many2many:liabilty_payments"`
	gorm.Model
}
//Validate ..
func (liability Liability) Validate() *httperors.HttpError{ 
	if liability.Name == "" && len(liability.Name) < 3 {
		return httperors.NewNotFoundError("Invalid Name")
	}
	if liability.Description == "" && len(liability.Description) < 3 {
		return httperors.NewNotFoundError("Invalid description")
	}
	if liability.Creditor == "" {
		return httperors.NewNotFoundError("Invalid Creditor name")
	}
	if liability.Approvedby == "" {
		return httperors.NewNotFoundError("Invalid Approved name")
	}
	if liability.Paymentperiod < 0 {
		return httperors.NewNotFoundError("Invalid payment period")
	}
	if liability.Interestrate < 0 {
		return httperors.NewNotFoundError("Invalid interst rate")
	}
	if liability.Amount < 0 {
		return httperors.NewNotFoundError("Invalid amount")
	}
	return nil
}