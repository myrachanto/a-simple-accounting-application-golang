package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
	"github.com/myrachanto/accounting/support"
)

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
func (service nortificationService) GetOne(id int) (*model.Nortification, *httperors.HttpError) {
	nortification, err1 := r.Nortificationrepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return nortification, nil
}

func (service nortificationService) GetAll(nortifications []model.Nortification,search *support.Search) ([]model.Nortification, *httperors.HttpError) {
	nortifications, err := r.Nortificationrepo.GetAll(nortifications,search)
	if err != nil {
		return nil, err
	}
	return nortifications, nil
}

func (service nortificationService) Update(id int, nortification *model.Nortification) (*model.Nortification, *httperors.HttpError) {
	nortification, err1 := r.Nortificationrepo.Update(id, nortification)
	if err1 != nil {
		return nil, err1
	}
	
	return nortification, nil
}
func (service nortificationService) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Nortificationrepo.Delete(id)
		return success, failure
}
