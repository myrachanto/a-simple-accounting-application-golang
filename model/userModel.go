package model

import (
  "gorm.io/gorm"
	"time"
	"golang.org/x/crypto/bcrypt"
	jwt "github.com/dgrijalva/jwt-go"
	"regexp"
	"github.com/myrachanto/accounting/httperors"
)
//ExpiresAt ..
var ExpiresAt = time.Now().Add(time.Minute * 100000).Unix()
//User ..
type User struct {
	FName string `json:"fname"`
	LName string `json:"lname"`
	UName string `gorm:"not null" json:"uname"`
	Usercode string `json:"usercode"`
	Phone string  `json:"phone"`
	Address string  `json:"address"`
	Dob *time.Time   
	Picture string  `json:"picture"`
	Email string `gorm:"not null;unique" json:"email"`
	Password string `gorm:"not null" json:"password"`
	Admin bool `json:"admin"`
	Employee bool `json:"employee"`
	Supervisor bool `json:"supervisor"`
	Message []*Message `gorm:"many2many:user_messages;"`
	Nortification []*Nortification `gorm:"many2many:user_nortifications;"`
	gorm.Model
}
//Auth ..
type Auth struct {
	//User User `gorm:"foreignKey:UserID; not null"`
	UserID uint `json:"userid"`
	UName string `json:"uname"`
	Usercode string `json:"usercode"`
	Picture string `json:"picture"`
	Token string `gorm:"size:500;not null"`
	Admin bool `json:"admin"`
	Employee bool `json:"employee"`
	Supervisor bool `json:"supervisor"`
	gorm.Model
}
//LoginUser ..
type LoginUser struct {
	Email string `gorm:"not null"`
	Password string `gorm:"not null"`
}
//Token struct declaration
type Token struct {
	UserID uint
	UName string `json:"uname"`
	Email  string `json:"email"`
	Usercode string `json:"usercode"`
	Admin bool `json:"admin"`
	Employee bool `json:"employee"`
	Supervisor bool `json:"supervisor"`
	*jwt.StandardClaims
}
//UserProfile user profile and messages
type UserProfile struct{
 User User `json:"user"`
 Inbox []Message `json:"inbox"`
 Sent []Message `json:"sent"`
 Users []User `json:"users"`
 Nortification []Nortification `json:"nortifications"`
}
//ValidateEmail ..
func (user User)ValidateEmail(email string) (matchedString bool) {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&amp;'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	matchedString = re.MatchString(email)
	return matchedString
}
//ValidatePassword ...
func (user User)ValidatePassword(password string) (bool, *httperors.HttpError) {
	if len(password) < 5 {
		return false, httperors.NewBadRequestError("your password need more characters!")
	} else if len(password) > 32 {
		return false, httperors.NewBadRequestError("your password is way too long!")
	}
	return true, nil
}
//HashPassword ..
func (user User)HashPassword(password string)(string, *httperors.HttpError){
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			return "", httperors.NewNotFoundError("type a stronger password!")
		}
		return string(pass),nil 
		
	}
	//Compare ..
func (user User) Compare(p1,p2 string) bool{
	err := bcrypt.CompareHashAndPassword([]byte(p2), []byte(p1))
	if err != nil {
		return false
	}
	return true
}
//Validate ..
func (loginuser LoginUser) Validate() *httperors.HttpError{ 
	if loginuser.Email == "" {
		return httperors.NewNotFoundError("Invalid Email")
	}
	if loginuser.Password == "" {
		return httperors.NewNotFoundError("Invalid password")
	}
	return nil
}
//Validate ..
func (user User) Validate() *httperors.HttpError{
	if user.FName == "" {
		return httperors.NewNotFoundError("Invalid first Name")
	}
	if user.LName == "" {
		return httperors.NewNotFoundError("Invalid last name")
	}
	if user.UName == "" {
		return httperors.NewNotFoundError("Invalid username")
	}
	if user.Phone == "" {
		return httperors.NewNotFoundError("Invalid phone number")
	}
	if user.Email == "" {
		return httperors.NewNotFoundError("Invalid Email")
	}
	if user.Address == "" {
		return httperors.NewNotFoundError("Invalid Address")
	}
	if user.Password == "" {
		return httperors.NewNotFoundError("Invalid password")
	}
	// if user.Picture == "" {
	// 	return httperors.NewNotFoundError("Invalid picture")
	// }
	if user.Email == "" {
		return httperors.NewNotFoundError("Invalid picture")
	}
	return nil
}