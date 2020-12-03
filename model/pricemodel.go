package model

import (
  "gorm.io/gorm"
	"github.com/myrachanto/accounting/httperors"
)
//Price ..
type Price struct {
	Name string `gorm:"not null" json:"name"` 
	Title string ` json:"title"`
	Description string `gorm:"not null" json:"description"`
	Product string ` json:"product"`
	gorm.Model
}
//Validate ..
func (price Price) Validate() *httperors.HttpError{ 
	if price.Name == "" && len(price.Name) < 3 {
		return httperors.NewNotFoundError("Invalid Name")
	}
	if price.Title == "" && len(price.Title) < 3 {
		return httperors.NewNotFoundError("Invalid Title")
	}
	
	if price.Product == "" {
		return httperors.NewNotFoundError("Invalid Product")
	}
	
	if price.Description == "" && len(price.Description) < 10 {
		return httperors.NewNotFoundError("Invalid description")
	}
	return nil
}