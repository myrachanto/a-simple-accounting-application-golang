package model

import (
  "gorm.io/gorm"
	"github.com/myrachanto/accounting/httperors"
)
//Expence stucture to ...
type Expence struct {
	Name string `gorm:"not null" json:"name"`
	Title string  `gorm:"not null" json:"title"`
	Description string `gorm:"not null" json:"description"`
	Expencetrasan []Expencetrasan
	gorm.Model
}
//Validate ..
func (expence Expence) Validate() *httperors.HttpError{ 
	if expence.Name == "" && len(expence.Name) < 3 {
		return httperors.NewNotFoundError("Invalid Name")
	}
	if expence.Description == "" && len(expence.Description) < 3 {
		return httperors.NewNotFoundError("Invalid description")
	}
	if expence.Title == "" {
		return httperors.NewNotFoundError("Invalid title")
	}
	return nil
}