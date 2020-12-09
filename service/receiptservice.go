package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
)
//Receiptservice ...
var (
	Receiptservice receiptservice = receiptservice{}
)

type receiptservice struct {
}

func (service receiptservice) Create(receipt *model.Receipt) (*model.Receipt, *httperors.HttpError) {
	receipt, err1 := r.Receiptrepo.Create(receipt)
	if err1 != nil {
		return nil, err1
	}
	return receipt, nil

}
func (service receiptservice) UpdateReceipts( code,status string) (string, *httperors.HttpError) {
	cart, err1 := r.Receiptrepo.UpdateReceipts(code,status)
	if err1 != nil {
		return "", err1
	}
	 
	return cart, nil
}
func (service receiptservice) View() (*model.ReceiptView, *httperors.HttpError) {
	code, err1 := r.Receiptrepo.View()
	if err1 != nil {
		return nil, err1
	}
	return code, nil
}


func (service receiptservice) ViewReport() (*model.ReceiptReport, *httperors.HttpError) {
	options, err1 := r.Receiptrepo.ViewReport()
	if err1 != nil {
		return nil, err1
	}
	return options, nil
}

func (service receiptservice) ViewCleared() ([]model.Receipt, *httperors.HttpError) {
	options, err1 := r.Receiptrepo.ViewCleared()
	if err1 != nil {
		return nil, err1
	}
	return options, nil
}
func (service receiptservice) AddReceiptTrans(clientcode,invoicecode,usercode,receiptcode string ,amount float64) (string, *httperors.HttpError) {
	options, err1 := r.Receiptrepo.AddReceiptTrans(clientcode,invoicecode,usercode,receiptcode,amount)
	return options, err1
}
func (service receiptservice) ViewInvoices(customercode string) (*model.ReceiptAlloc, *httperors.HttpError) {
	invoices, err1 := r.Receiptrepo.ViewInvoices(customercode)
	if err1 != nil {
		return nil, err1
	}
	return invoices, nil
}
func (service receiptservice) GetOne(id int) (*model.Receipt, *httperors.HttpError) {
	receipt, err1 := r.Receiptrepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return receipt, nil
}

func (service receiptservice) GetAll() (*model.ReceiptOptions, *httperors.HttpError) {
	receipts, err := r.Receiptrepo.GetAll()
	if err != nil {
		return nil, err
	}
	return receipts, nil
}

func (service receiptservice) Update(id int, receipt *model.Receipt) (*model.Receipt, *httperors.HttpError) {
	receipt, err1 := r.Receiptrepo.Update(id, receipt)
	if err1 != nil {
		return nil, err1
	}

	return receipt, nil
}
func (service receiptservice) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {

	success, failure := r.Receiptrepo.Delete(id)
	return success, failure
}
