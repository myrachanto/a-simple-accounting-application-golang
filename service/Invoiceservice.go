package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
)
//Invoiceservice service
var (
	Invoiceservice invoiceservice = invoiceservice{}

) 
type invoiceservice struct {
	
}

func (service invoiceservice) Create(invoice *model.Invoice) (string, *httperors.HttpError) {
	invoic, err1 := r.Invoicerepo.Create(invoice)
	if err1 != nil {
		return "", err1
	}
	 return invoic, nil

}
func (service invoiceservice) View() (*model.Cinvoiceoptions, *httperors.HttpError) {
	code, err1 := r.Invoicerepo.View()
	if err1 != nil {
		return nil, err1
	}
	return code, nil
}
func (service invoiceservice) GetOne(code string) (*model.InvoiceView, *httperors.HttpError) {
	invoice, err1 := r.Invoicerepo.GetOne(code)
	if err1 != nil {
		return nil, err1
	}
	return invoice, nil
}

func (service invoiceservice) GetAll(search,dated,searchq2,searchq3 string) ([]model.Invoice, *httperors.HttpError) {
	invoices, err := r.Invoicerepo.GetAll(search,dated,searchq2,searchq3)
	if err != nil {
		return nil, err
	}
	return invoices, nil
}
func (service invoiceservice) GetCredit(search,dated,searchq2,searchq3 string) ([]model.Invoice, *httperors.HttpError) {
	invoices, err := r.Invoicerepo.GetCredit(search,dated,searchq2,searchq3)
	if err != nil {
		return nil, err
	}
	return invoices, nil
}
func (service invoiceservice) Update(code string) (string, *httperors.HttpError) {
	invoice, err1 := r.Invoicerepo.Update(code)
	if err1 != nil {
		return "", err1
	}
	
	return invoice, nil
}
func (service invoiceservice) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Invoicerepo.Delete(id)
		return success, failure
}
