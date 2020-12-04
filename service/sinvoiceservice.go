package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/support"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
)
//SInvoiceservice service
var (
	SInvoiceservice sInvoiceservice = sInvoiceservice{}

) 
type sInvoiceservice struct {
	
}

func (service sInvoiceservice) Create(sInvoice *model.SInvoice) (string, *httperors.HttpError) {
	sInvoic, err1 := r.SInvoicerepo.Create(sInvoice)
	if err1 != nil {
		return "", err1
	}
	 return sInvoic, nil

}
func (service sInvoiceservice) View() (*model.Sinvoiceoptions, *httperors.HttpError) {
	code, err1 := r.SInvoicerepo.View()
	if err1 != nil {
		return nil, err1
	}
	return code, nil
}
func (service sInvoiceservice) GetOne(code string) (*model.SInvoiceView, *httperors.HttpError) {
	sInvoice, err1 := r.SInvoicerepo.GetOne(code)
	if err1 != nil {
		return nil, err1
	}
	return sInvoice, nil
}

func (service sInvoiceservice) GetAll(sInvoices []model.SInvoice,search *support.Search) ([]model.SInvoice, *httperors.HttpError) {
	sInvoices, err := r.SInvoicerepo.GetAll(sInvoices,search)
	if err != nil {
		return nil, err
	}
	return sInvoices, nil
}
func (service sInvoiceservice) GetCredit(sInvoices []model.SInvoice,search *support.Search) ([]model.SInvoice, *httperors.HttpError) {
	sInvoices, err := r.SInvoicerepo.GetCredit(sInvoices,search)
	if err != nil {
		return nil, err
	}
	return sInvoices, nil
}
func (service sInvoiceservice) Update(code string) (string, *httperors.HttpError) {
	sInvoice, err1 := r.SInvoicerepo.Update(code)
	if err1 != nil {
		return "", err1
	}
	
	return sInvoice, nil
}
func (service sInvoiceservice) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.SInvoicerepo.Delete(id)
		return success, failure
}
