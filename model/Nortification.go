package model

import (
  "gorm.io/gorm"
	"github.com/myrachanto/accounting/httperors"
)
//Nortification ..
type Nortification struct {
	Name string `gorm:"not null" json:"name"`
	Title string  `gorm:"not null" json:"title"`
	Users []*User `gorm:"many2many:nortification_users;"`
	Description string  `gorm:"not null" json:"descriptio"`
	Read bool  `gorm:"not null" json:"read"`
	Usercode string `json:"usercode"`
	gorm.Model
}
//Validate ..
func (nortification Nortification) Validate() *httperors.HttpError{ 
	if nortification.Name == "" && len(nortification.Name) < 3 {
		return httperors.NewNotFoundError("Invalid Name")
	}
	if nortification.Title == "" && len(nortification.Title) < 3 {
		return httperors.NewNotFoundError("Invalid Title")
	}
	
	if nortification.Description == "" && len(nortification.Description) < 10 {
		return httperors.NewNotFoundError("Invalid description")
	}
	return nil
}