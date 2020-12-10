package model

import (
  "gorm.io/gorm"
	"github.com/myrachanto/accounting/httperors"
)
//Expencetrasan struture to capture all transactions involving expences
type Expencetrasan struct {
	Name string `gorm:"not null" json:"name"`
	ExpenceID uint  `json:"expenceID"`
	Code string `json:"code"`//expencecode, assetcode, liacode
	Type string `json:"type"` 
	Mode string `json:"mode"` //active or canceled
	Direct string `json:"direct"` //inderect expence or direct expence
	Title string `gorm:"not null" json:"title"`
	Description string `gorm:"not null" json:"description"`
	Amount float64 `gorm:"not null" json:"amount"`
	Usercode string `json:"usercode"`
	Status string `gorm:"not null" json:"status"`//type expence or asset purchase or liability payment
	Paid string `gorm:"not null" json:"paid"`
	gorm.Model
}
//ExpencetransView ...
type ExpencetransView struct {
Expence []Expence `json:"expences"`
Liability []Liability `json:"liabilitys"`
Asset []Asset `json:"assets"`
}
//ExpencesView ...
type ExpencesView struct {
	Expences []Expencetrasan `json:"expences"`
	Totalexpences SalesModule `json:"total"`
	Directexpences SalesModule `json:"directexpences"`
	InDirectexpences SalesModule `json:"indirectexpences"`
	Other SalesModule `json:"others"`
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