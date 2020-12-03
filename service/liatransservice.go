package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
	"github.com/myrachanto/accounting/support"
)

var (
	LiatranService liatranService = liatranService{}

) 
type liatranService struct {
	
}

func (service liatranService) Create(liatran *model.Liatran) (*model.Liatran, *httperors.HttpError) {
	
	liatran, err1 := r.Liatranrepo.Create(liatran)
	if err1 != nil {
		return nil, err1
	}
	 return liatran, nil

}
func (service liatranService) GetOne(id int) (*model.Liatran, *httperors.HttpError) {
	liatran, err1 := r.Liatranrepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return liatran, nil
}

func (service liatranService) GetAll(liatran []model.Liatran,search *support.Search) ([]model.Liatran, *httperors.HttpError) {
	liatran, err := r.Liatranrepo.GetAll(liatran,search)
	if err != nil {
		return nil, err
	}
	return liatran, nil
}

func (service liatranService) Update(id int, liatran *model.Liatran) (*model.Liatran, *httperors.HttpError) {
	liatran, err1 := r.Liatranrepo.Update(id, liatran)
	if err1 != nil {
		return nil, err1
	}
	
	return liatran, nil
}
func (service liatranService) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Liatranrepo.Delete(id)
		return success, failure
}
