package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
)
//Supplierservice ...
var (
	Supplierservice supplierservice = supplierservice{}

) 
type supplierservice struct {
	
}

func (service supplierservice) Create(supplier *model.Supplier) (string, *httperors.HttpError) {
	if err := supplier.Validate(); err != nil {
		return "", err
	}	
	s, err1 := r.Supplierrepo.Create(supplier)
	 return s, err1
 
}
func (service supplierservice) Login(asupplier *model.Loginsupplier) (*model.SupplierAuth, *httperors.HttpError) {
	supplier, err1 := r.Supplierrepo.Login(asupplier)
	return supplier, err1
}
func (service supplierservice) Forgot(email string) (string, *httperors.HttpError) {
	s, err1 := r.Supplierrepo.Forgot(email)
	return s, err1
}
func (service supplierservice) Logout(token string) (*httperors.HttpError) {
	err1 := r.Supplierrepo.Logout(token)
	return err1
}
func (service supplierservice) GetOne(id int) (*model.Supplierdetails, *httperors.HttpError) {
	supplier, err1 := r.Supplierrepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return supplier, nil
}

func (service supplierservice) ViewReport() (*model.SupplierView, *httperors.HttpError) {
	options, err1 := r.Supplierrepo.ViewReport()
	if err1 != nil {
		return nil, err1
	}
	return options, nil
}
func (service supplierservice) GetAll(search string) ([]model.Supplier, *httperors.HttpError) {
	
	results, err := r.Supplierrepo.GetAll(search)
	if err != nil { 
		return nil, err
	}
	return results, nil 
}

func (service supplierservice) Update(id int, supplier *model.Supplier) (*model.Supplier, *httperors.HttpError) {
	supplier, err1 := r.Supplierrepo.Update(id, supplier)
	if err1 != nil {
		return nil, err1
	}
	
	return supplier, nil
}
func (service supplierservice) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Supplierrepo.Delete(id)
		return success, failure
}
///////deleting a batch////////////////////

//db.Where("age = ?", 20).Delete(&supplier{})