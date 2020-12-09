package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
)
//Customerservice ..
var (
	Customerservice customerservice = customerservice{}

) 
type customerservice struct {
	
}

func (service customerservice) Create(customer *model.Customer) (string, *httperors.HttpError) {
	if err := customer.Validate(); err != nil {
		return "", err
	}	
	s, err1 := r.Customerrepo.Create(customer)
	if err1 != nil {
		return "", err1
	}
	 return s, nil
 
}
func (service customerservice) Login(acustomer *model.Logincustomer) (*model.CustomnerAuth, *httperors.HttpError) {
	
	customer, err1 := r.Customerrepo.Login(acustomer)
	if err1 != nil {
		return nil, err1
	} 
	return customer, nil
}
func (service customerservice) Forgot(email string) (string, *httperors.HttpError) {
	
	s, err1 := r.Customerrepo.Forgot(email)
	if err1 != nil {
		return "", err1
	} 
	return s, nil
}
func (service customerservice) Logout(token string) (*httperors.HttpError) {
	err1 := r.Customerrepo.Logout(token)
	if err1 != nil {
		return err1
	}
	return nil
}
func (service customerservice) GetOne(id int) (*model.Customerdetails, *httperors.HttpError) {
	customer, err1 := r.Customerrepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return customer, nil
}

func (service customerservice) ViewReport() (*model.CustomerView, *httperors.HttpError) {
	options, err1 := r.Customerrepo.ViewReport()
	if err1 != nil {
		return nil, err1
	}
	return options, nil
}
func (service customerservice) GetAll(search string, page,pagesize int) ([]model.Customer, *httperors.HttpError) {
	results, err := r.Customerrepo.GetAll(search, page,pagesize)
	return results, err
}
func (service customerservice) Update(id int, customer *model.Customer) (*model.Customer, *httperors.HttpError) {
	customer, err1 := r.Customerrepo.Update(id, customer)
	if err1 != nil {
		return nil, err1
	}
	
	return customer, nil
}
func (service customerservice) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Customerrepo.Delete(id)
		return success, failure
}
///////deleting a batch////////////////////

//db.Where("age = ?", 20).Delete(&customer{})