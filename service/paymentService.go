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

func (service paymentservice) ViewCleared() ([]model.Payment, *httperors.HttpError) {
	options, err1 := r.Paymentrepo.ViewCleared()
	if err1 != nil {
		return nil, err1
	}
	return options, nil
}
func (service paymentservice) ViewClearedExpence() ([]model.Payment, *httperors.HttpError) {
	options, err1 := r.Paymentrepo.ViewClearedExpence()
	if err1 != nil {
		return nil, err1
	}
	return options, nil
}
func (service paymentservice) AddReceiptTrans(clientcode,invoicecode,usercode,receiptcode string ,amount float64) (string, *httperors.HttpError) {
	options, err1 := r.Paymentrepo.AddReceiptTrans(clientcode,invoicecode,usercode,receiptcode,amount)
	return options, err1
}
func (service paymentservice) ViewInvoices(customercode string) (*model.PaymentAlloc, *httperors.HttpError) {
	invoices, err1 := r.Paymentrepo.ViewInvoices(customercode)
	if err1 != nil {
		return nil, err1
	}
	return invoices, nil
}
func (service paymentservice) ViewReport(dated,searchq2,searchq3 string) (*model.PaymentReport, *httperors.HttpError) {
	options, err1 := r.Paymentrepo.ViewReport(dated,searchq2,searchq3)
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
func (service paymentservice) ViewExpence() (*model.PaymentExpence, *httperors.HttpError) {
	code, err1 := r.Paymentrepo.ViewExpence()
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

func (service paymentservice) GetAll(dated,searchq2,searchq3 string) (*model.PaymentOptions, *httperors.HttpError) {
	payments, err := r.Paymentrepo.GetAll(dated,searchq2,searchq3)
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
