package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
	"github.com/myrachanto/accounting/support"
)

var (
	Paymentservice paymentservice = paymentservice{}
)

type paymentservice struct {
}

func (service paymentservice) Create(payment *model.Payment) (*model.Payment, *httperors.HttpError) {
	payment, err1 := r.Paymentrepo.Create(payment)
	if err1 != nil {
		return nil, err1
	}
	return payment, nil

}
func (service paymentservice) GetOne(id int) (*model.Payment, *httperors.HttpError) {
	payment, err1 := r.Paymentrepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return payment, nil
}

func (service paymentservice) GetAll(payments []model.Payment, search *support.Search) ([]model.Payment, *httperors.HttpError) {
	payments, err := r.Paymentrepo.GetAll(payments, search)
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (service paymentservice) Update(id int, payment *model.Payment) (*model.Payment, *httperors.HttpError) {
	payment, err1 := r.Paymentrepo.Update(id, payment)
	if err1 != nil {
		return nil, err1
	}

	return payment, nil
}
func (service paymentservice) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {

	success, failure := r.Paymentrepo.Delete(id)
	return success, failure
}
