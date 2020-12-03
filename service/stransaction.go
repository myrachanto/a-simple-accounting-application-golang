package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/support"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
)
//STransactionservice service
var (
	STransactionservice sTransactionservice = sTransactionservice{}

) 
type sTransactionservice struct {
	
}

func (service sTransactionservice) Create(sTransaction *model.STransaction) (*model.STransaction, *httperors.HttpError) {
	sTransaction, err1 := r.STransactionrepo.Create(sTransaction)
	if err1 != nil {
		return nil, err1
	}
	 return sTransaction, nil

}
func (service sTransactionservice) GetOne(id int) (*model.STransaction, *httperors.HttpError) {
	sTransaction, err1 := r.STransactionrepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return sTransaction, nil
}

func (service sTransactionservice) GetAll(sTransactions []model.STransaction,search *support.Search) ([]model.STransaction, *httperors.HttpError) {
	sTransactions, err := r.STransactionrepo.GetAll(sTransactions,search)
	if err != nil {
		return nil, err
	}
	return sTransactions, nil
}

func (service sTransactionservice) Update(id int, sTransaction *model.STransaction) (*model.STransaction, *httperors.HttpError) {
	sTransaction, err1 := r.STransactionrepo.Update(id, sTransaction)
	if err1 != nil {
		return nil, err1
	}
	
	return sTransaction, nil
}
func (service sTransactionservice) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.STransactionrepo.Delete(id)
		return success, failure
}
