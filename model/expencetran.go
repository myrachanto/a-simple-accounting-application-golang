package model

import (
  "gorm.io/gorm"
	"github.com/myrachanto/accounting/httperors"
)
//Expencetrasan struture to capture all transactions involving expences
type Expencetrasan struct {
	Name string `gorm:"not null" json:"name"`
	ExpenceID uint  `json:"expenceID"`
	Code string `json:"code"`
	Usercode string `json:"usercode"`
	Mode string `json:"mode"`
	Title string `gorm:"not null" json:"title"`
	Description string `gorm:"not null" json:"description"`
	Amount float64 `gorm:"not null" json:"amount"`
	Status bool `gorm:"not null" json:"status"`
	Paid string `gorm:"not null" json:"paid"`
	Type string `json:"type"`
	gorm.Model
}
//ExpencesView ...
type ExpencesView struct {
	Expences []Expencetrasan `json:"expences"`
	Totalexpences SalesModule `json:"total"`
	Directexpences SalesModule `json:"directexpences"`
	InDirectexpences SalesModule `json:"indirectexpences"`
}
//Validate ...
func (expence Expencetrasan) Validate() *httperors.HttpError{ 
	if expence.Name == "" && len(expence.Name) < 3 {
		return httperors.NewNotFoundError("Invalid Name")
	}
	if expence.Description == "" && len(expence.Description) < 3 {
		return httperors.NewNotFoundError("Invalid description")
	}
	if expence.Amount == 0 {
		return httperors.NewNotFoundError("Invalid customer")
	}
	return nil
}