package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
)
//Paymentformservice ...
var (
	Paymentformservice paymentformservice = paymentformservice{}

) 
type paymentformservice struct {
	
}

func (service paymentformservice) Create(paymentform *model.Paymentform) (*model.Paymentform, *httperors.HttpError) {
	if err := paymentform.Validate(); err != nil {
		return nil, err
	}	
	paymentform, err1 := r.Paymentformrepo.Create(paymentform)
	if err1 != nil {
		return nil, err1
	}
	 return paymentform, nil

}
func (service paymentformservice) GetOne(id int) (*model.Paymentform, *httperors.HttpError) {
	paymentform, err1 := r.Paymentformrepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return paymentform, nil
}

func (service paymentformservice) GetAll(search string, page,pagesize int) ([]model.Paymentform, *httperors.HttpError) {
	results, err := r.Paymentformrepo.GetAll(search, page,pagesize)
	return results, err
}
func (service paymentformservice) Update(id int, paymentform *model.Paymentform) (*model.Paymentform, *httperors.HttpError) {
	paymentform, err1 := r.Paymentformrepo.Update(id, paymentform)
	if err1 != nil {
		return nil, err1
	}
	
	return paymentform, nil
}
func (service paymentformservice) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Paymentformrepo.Delete(id)
		return success, failure
}
