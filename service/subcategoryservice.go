package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
)
//Subcategoryservice ...
var (
	Subcategoryservice subcategoryservice = subcategoryservice{}

) 
type subcategoryservice struct {
	
}

func (service subcategoryservice) Create(subcategory *model.Subcategory) (*model.Subcategory, *httperors.HttpError) {
	if err := subcategory.Validate(); err != nil {
		return nil, err
	}	
	subcategory, err1 := r.Subcategoryrepo.Create(subcategory)
	if err1 != nil {
		return nil, err1
	}
	 return subcategory, nil

}
func (service subcategoryservice) GetOne(id int) (*model.Subcategory, *httperors.HttpError) {
	subcategory, err1 := r.Subcategoryrepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return subcategory, nil
}
func (service subcategoryservice) GetAll(search string, page,pagesize int) ([]model.Subcategory, *httperors.HttpError) {
	results, err := r.Subcategoryrepo.GetAll(search, page,pagesize)
	return results, err
}


func (service subcategoryservice) Update(id int, subcategory *model.Subcategory) (*model.Subcategory, *httperors.HttpError) {
	subcategory, err1 := r.Subcategoryrepo.Update(id, subcategory)
	if err1 != nil {
		return nil, err1
	}
	
	return subcategory, nil
}
func (service subcategoryservice) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Subcategoryrepo.Delete(id)
		return success, failure
}
