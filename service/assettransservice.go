package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
	"github.com/myrachanto/accounting/support"
)

var (
	AsstransService asstransService = asstransService{}

)

type asstransService struct {
	
}

func (service asstransService) Create(asstrans *model.Asstrans) (*model.Asstrans, *httperors.HttpError) {	
	asstrans, err1 := r.Asstransrepo.Create(asstrans)
	if err1 != nil {
		return nil, err1
	}
	 return asstrans, nil

}
func (service asstransService) GetOne(id int) (*model.Asstrans, *httperors.HttpError) {
	asstrans, err1 := r.Asstransrepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return asstrans, nil
}

func (service asstransService) GetAll(asstranss []model.Asstrans,search *support.Search) ([]model.Asstrans, *httperors.HttpError) {
	asstranss, err := r.Asstransrepo.GetAll(asstranss,search)
	if err != nil {
		return nil, err
	}
	return asstranss, nil
}

func (service asstransService) Update(id int, asstrans *model.Asstrans) (*model.Asstrans, *httperors.HttpError) {
	asstrans, err1 := r.Asstransrepo.Update(id, asstrans)
	if err1 != nil {
		return nil, err1
	}
	
	return asstrans, nil
}
func (service asstransService) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Asstransrepo.Delete(id)
		return success, failure
}
