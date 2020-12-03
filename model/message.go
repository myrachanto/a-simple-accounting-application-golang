package model

import (
  "gorm.io/gorm"
	"github.com/myrachanto/accounting/httperors"
)
//Message structure
type Message struct {
	Name string  `gorm:"not null" json:"name"`
	Title string  `gorm:"not null" json:"title"`	
	Sender []*User `gorm:"many2many:message_users;"`
	Description string  `gorm:"not null" json:"description"`
	Read bool  `gorm:"not null" json:"read"`
	gorm.Model
}
//Validate ..
func (message Message) Validate() *httperors.HttpError{ 
	if message.Name == "" && len(message.Name) < 3 {
		return httperors.NewNotFoundError("Invalid Name")
	}
	if message.Title == "" && len(message.Title) < 3 {
		return httperors.NewNotFoundError("Invalid Title")
	}
	
	if message.Description == "" && len(message.Description) < 10 {
		return httperors.NewNotFoundError("Invalid description")
	}
	return nil
}