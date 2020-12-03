package service

import (
	"fmt"
	"github.com/myrachanto/accounting/support"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
)

var (
	UserService userService = userService{}

) 
type redirectUser interface{
	Create(customer *model.User) (*model.User, *httperors.HttpError)
	Login(auser *model.LoginUser) (*model.Auth, *httperors.HttpError)
	Logout(token string) (*httperors.HttpError)
	GetOne(id int) (*model.User, *httperors.HttpError)
	GetAll(users []model.User, search *support.Search) ([]model.User, *httperors.HttpError)
	Update(id int, user *model.User) (*model.User, *httperors.HttpError)
	Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError)
}


type userService struct {
}

func (service userService) Create(user *model.User) (string, *httperors.HttpError) {
	if err := user.Validate(); err != nil {
		return "", err
	}	
	s, err1 := r.Userrepo.Create(user)
	if err1 != nil {
		return "", err1
	}
	 return s, nil
 
}
func (service userService) Login(auser *model.LoginUser) (*model.Auth, *httperors.HttpError) {
	
	user, err1 := r.Userrepo.Login(auser)
	if err1 != nil {
		return nil, err1
	} 
	return user, nil
}
func (service userService) Logout(token string) (*httperors.HttpError) {
	err1 := r.Userrepo.Logout(token)
	if err1 != nil {
		return err1
	}
	return nil
}
func (service userService) GetOne(id int) (*model.User, *httperors.HttpError) {
	user, err1 := r.Userrepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return user, nil
}

func (service userService) GetAll(users []model.User, search *support.Search) ([]model.User, *httperors.HttpError) {
	users, err := r.Userrepo.GetAll(users, search)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (service userService) Update(id int, user *model.User) (*model.User, *httperors.HttpError) {
	
	fmt.Println("update1-controller")
	fmt.Println(id)
	user, err1 := r.Userrepo.Update(id, user)
	if err1 != nil {
		return nil, err1
	}
	
	return user, nil
}
func (service userService) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Userrepo.Delete(id)
		return success, failure
}
