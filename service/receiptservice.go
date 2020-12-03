package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
	"github.com/myrachanto/accounting/support"
)

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
func (service receiptservice) GetOne(id int) (*model.Receipt, *httperors.HttpError) {
	receipt, err1 := r.Receiptrepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return receipt, nil
}

func (service receiptservice) GetAll(receipts []model.Receipt, search *support.Search) ([]model.Receipt, *httperors.HttpError) {
	receipts, err := r.Receiptrepo.GetAll(receipts, search)
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
