package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
	"github.com/myrachanto/accounting/support"
)

var (
	ExpenceService expenceService = expenceService{}

) 
type expenceService struct {
	
}

func (service expenceService) Create(expence *model.Expence) (*model.Expence, *httperors.HttpError) {
	if err := expence.Validate(); err != nil {
		return nil, err
	}	
	expence, err1 := r.Expencerepo.Create(expence)
	if err1 != nil {
		return nil, err1
	}
	 return expence, nil

}
func (service expenceService) GetOne(id int) (*model.Expence, *httperors.HttpError) {
	expence, err1 := r.Expencerepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return expence, nil
}

func (service expenceService) GetAll(expences []model.Expence,search *support.Search) ([]model.Expence, *httperors.HttpError) {
	expences, err := r.Expencerepo.GetAll(expences,search)
	if err != nil {
		return nil, err
	}
	return expences, nil
}

func (service expenceService) Update(id int, expence *model.Expence) (*model.Expence, *httperors.HttpError) {
	expence, err1 := r.Expencerepo.Update(id, expence)
	if err1 != nil {
		return nil, err1
	}
	
	return expence, nil
}
func (service expenceService) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Expencerepo.Delete(id)
		return success, failure
}
