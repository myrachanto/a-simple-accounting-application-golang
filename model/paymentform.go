package model
 
import (
  "gorm.io/gorm"
	"github.com/myrachanto/accounting/httperors"
)
//Paymentform type of payment credit card, cash ...
type Paymentform struct {
	Name string `gorm:"not null" json:"name"` 
	Title string `gorm:"not null" json:"title"`
	Description string `gorm:"not null" json:"description"` 
	ReceiptID uint `gorm:"not null" json:"receiptid"`
	Receipt []Receipt `gorm:"many2many:paymentform_receipts"`
	Payment []Payment `gorm:"many2many:paymentform_payments"`
	Amount float64  `json:"amount"`
	gorm.Model
}
//Validate ...
func (paymentform Paymentform) Validate() *httperors.HttpError{ 
	if paymentform.Name == "" && len(paymentform.Name) < 3 {
		return httperors.NewNotFoundError("Invalid Name")
	}
	if paymentform.Title == "" && len(paymentform.Title) < 3 {
		return httperors.NewNotFoundError("Invalid Title")
	}
	
	if paymentform.Description == "" && len(paymentform.Description) < 10 {
		return httperors.NewNotFoundError("Invalid description")
	}
	return nil
}