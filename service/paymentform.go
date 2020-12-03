package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/support"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
)

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

func (service paymentformservice) GetAll(paymentforms []model.Paymentform,search *support.Search) ([]model.Paymentform, *httperors.HttpError) {
	paymentforms, err := r.Paymentformrepo.GetAll(paymentforms,search)
	if err != nil {
		return nil, err
	}
	return paymentforms, nil
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
