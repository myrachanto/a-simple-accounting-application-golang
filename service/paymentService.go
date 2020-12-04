package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
)
//Paymentservice
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
func (service paymentservice) Updatepayments( code,status string) (string, *httperors.HttpError) {
	cart, err1 := r.Paymentrepo.Updatepayments(code,status)
	if err1 != nil {
		return "", err1
	}
	 
	return cart, nil
}

func (service paymentservice) ViewReport() (*model.PaymentReport, *httperors.HttpError) {
	options, err1 := r.Paymentrepo.ViewReport()
	if err1 != nil {
		return nil, err1
	}
	return options, nil
}
func (service paymentservice) View() (*model.PaymentView, *httperors.HttpError) {
	code, err1 := r.Paymentrepo.View()
	if err1 != nil {
		return nil, err1
	}
	return code, nil
}
func (service paymentservice) GetOne(id int) (*model.Payment, *httperors.HttpError) {
	payment, err1 := r.Paymentrepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return payment, nil
}

func (service paymentservice) GetAll() (*model.PaymentOptions, *httperors.HttpError) {
	payments, err := r.Paymentrepo.GetAll()
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
