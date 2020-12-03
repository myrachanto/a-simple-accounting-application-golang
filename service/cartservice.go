package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
)
//Cartservice ..
var (
	Cartservice cartservice = cartservice{}

) 
type cartservice struct {
	
}

func (service cartservice) Create(cart *model.Cart) (string, *httperors.HttpError) {
	car, err1 := r.Cartrepo.Create(cart)
	if err1 != nil {
		return "", err1
	}
	 return car, nil

}
func (service cartservice) View(code string) ([]model.Cart, *httperors.HttpError) {
	options, err1 := r.Cartrepo.View(code)
	if err1 != nil {
		return nil, err1
	}
	return options, nil
}

func (service cartservice) Getcredits(code string) ([]model.Transaction, *httperors.HttpError) {
	options, err1 := r.Cartrepo.Getcredits(code)
	if err1 != nil {
		return nil, err1
	}
	return options, nil
}


func (service cartservice) GetcreditsList(code string) ([]model.Transaction, *httperors.HttpError) {
	options, err1 := r.Cartrepo.GetcreditsList(code)
	if err1 != nil {
		return nil, err1
	}
	return options, nil 
}
func (service cartservice) GetOne(id int) (*model.Cart, *httperors.HttpError) {
	cart, err1 := r.Cartrepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return cart, nil
}

func (service cartservice) GetAll(carts []model.Cart) ([]model.Cart, *httperors.HttpError) {
	carts, err := r.Cartrepo.GetAll(carts)
	if err != nil {
		return nil, err
	}
	return carts, nil
}

func (service cartservice) Update(qty float64, name, code string) (string, *httperors.HttpError) {
	cart, err1 := r.Cartrepo.Update(qty, name,code)
	if err1 != nil {
		return "", err1
	}
	 
	return cart, nil
}
func (service cartservice) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Cartrepo.Delete(id)
		return success, failure
}
func (service cartservice) DeleteALL(code string) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Cartrepo.DeleteAll(code)
		return success, failure
}
///////deleting a batch////////////////////

//db.Where("age = ?", 20).Delete(&User{})