package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
)
//NortificationService
var (
	NortificationService nortificationService = nortificationService{}

) 
type nortificationService struct {
	
}

func (service nortificationService) Create(nortification *model.Nortification) (*model.Nortification, *httperors.HttpError) {
	if err := nortification.Validate(); err != nil {
		return nil, err
	}	
	nortification, err1 := r.Nortificationrepo.Create(nortification)
	if err1 != nil {
		return nil, err1
	}
	 return nortification, nil

}

func (service nortificationService) GetAllUnread() (*model.NortUnread, *httperors.HttpError) {
	results, err := r.Nortificationrepo.GetAllUnread()
	return results, err
}
func (service nortificationService) GetOne(id int) (*model.Nortification, *httperors.HttpError) {
	nortification, err1 := r.Nortificationrepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return nortification, nil
}

func (service nortificationService) GetAll(search string, page,pagesize int) ([]model.Nortification, *httperors.HttpError) {
	results, err := r.Nortificationrepo.GetAll(search, page,pagesize)
	return results, err
}
func (service nortificationService) Update(id int) (*model.Nortification, *httperors.HttpError) {
	nortification, err1 := r.Nortificationrepo.Update(id)
	if err1 != nil {
		return nil, err1
	}
	
	return nortification, nil
}
func (service nortificationService) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Nortificationrepo.Delete(id)
		return success, failure
}
