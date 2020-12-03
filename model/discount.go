package model

import (
  "gorm.io/gorm"
	"github.com/myrachanto/accounting/httperors"
)
//Discount structure
type Discount struct {
	Name string `gorm:"not null" json:"name"` 
	Title string `gorm:"not null" json:"title"`
	Description string `gorm:"not null" json:"description"`
	gorm.Model
}
//Validate ..
func (discount Discount) Validate() *httperors.HttpError{ 
	if discount.Name == "" && len(discount.Name) < 3 {
		return httperors.NewNotFoundError("Invalid Name")
	}
	if discount.Title == "" && len(discount.Title) < 3 {
		return httperors.NewNotFoundError("Invalid Title")
	}
	
	if discount.Description == "" && len(discount.Description) < 10 {
		return httperors.NewNotFoundError("Invalid description")
	}
	return nil
}