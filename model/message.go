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
	Tousercode string `json:"to"`
	Fromusercode string `json:"from"`
	gorm.Model
}
//MessageUnread ..
type MessageUnread struct{
	Num int `json:"num"`
	Messages []Message `json:"messages"`
}
//Validate ..
func (message Message) Validate() *httperors.HttpError{ 
	if message.Title == "" && len(message.Title) < 3 {
		return httperors.NewNotFoundError("Invalid Title")
	}
	
	if message.Description == "" && len(message.Description) < 10 {
		return httperors.NewNotFoundError("Invalid description")
	}
	return nil
}