package model

import (
  "gorm.io/gorm"
	"github.com/myrachanto/accounting/httperors"
)
//Subcategory ..
type Subcategory struct {
	Name string `gorm:"not null"`
	Title string `gorm:"not null"`
	Description string `gorm:"not null"`
	CategoryID uint 
	gorm.Model
}
//Validate ..
func (subcategory Subcategory) Validate() *httperors.HttpError{ 
	if subcategory.Name == "" && len(subcategory.Name) < 3 {
		return httperors.NewNotFoundError("Invalid Name")
	}
	if subcategory.Title == "" && len(subcategory.Title) < 3 {
		return httperors.NewNotFoundError("Invalid Title")
	}
	
	if subcategory.Description == "" && len(subcategory.Description) < 10 {
		return httperors.NewNotFoundError("Invalid description")
	}
	return nil
}