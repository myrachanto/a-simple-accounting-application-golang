package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
)
//ExpencetrasanService ...
var (
	ExpencetrasanService expencetrasanService = expencetrasanService{}

) 
type expencetrasanService struct {
	
}

func (service expencetrasanService) Create(expencetrasan *model.Expencetrasan) (*model.Expencetrasan, *httperors.HttpError) {
	if err := expencetrasan.Validate(); err != nil {
		return nil, err
	}	
	expencetrasan, err1 := r.Expencetrasanrepo.Create(expencetrasan)
	if err1 != nil {
		return nil, err1
	}
	 return expencetrasan, nil

}
func (service expencetrasanService) CreateExp(expencetrasan *model.Expencetrasan) (*model.Expencetrasan, *httperors.HttpError) {

	expencetrasan, err1 := r.Expencetrasanrepo.CreateExp(expencetrasan)
	 return expencetrasan, err1

}
func (service expencetrasanService) View(code string) ([]model.Expencetrasan, *httperors.HttpError) {
	options, err1 := r.Expencetrasanrepo.View(code)
	if err1 != nil {
		return nil, err1
	}
	return options, nil
}

func (service expencetrasanService) ViewExp() (*model.ExpencetransView, *httperors.HttpError) {
	options, err1 := r.Expencetrasanrepo.ViewExp()
	if err1 != nil {
		return nil, err1
	}
	return options, nil
}
func (service expencetrasanService) ViewReport(dated,searchq2,searchq3 string) (*model.ExpencesView, *httperors.HttpError) {
	options, err1 := r.Expencetrasanrepo.ViewReport(dated,searchq2,searchq3)
	if err1 != nil {
		return nil, err1
	}
	return options, nil
}
func (service expencetrasanService) GetOne(id int) (*model.Expencetrasan, *httperors.HttpError) {
	expencetrasan, err1 := r.Expencetrasanrepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return expencetrasan, nil
}

func (service expencetrasanService) UpdateTrans( name, code string) (string, *httperors.HttpError) {
	cart, err1 := r.Expencetrasanrepo.UpdateTrans(name,code)
	if err1 != nil {
		return "", err1
	}
	 
	return cart, nil
}

func (service expencetrasanService) GetAll(search string, page,pagesize int) ([]model.Expencetrasan, *httperors.HttpError) {
	results, err := r.Expencetrasanrepo.GetAll(search, page,pagesize)
	return results, err
}

func (service expencetrasanService) Update(id int, expencetrasan *model.Expencetrasan) (*model.Expencetrasan, *httperors.HttpError) {
	expencetrasan, err1 := r.Expencetrasanrepo.Update(id, expencetrasan)
	if err1 != nil {
		return nil, err1
	}
	
	return expencetrasan, nil
}
func (service expencetrasanService) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Expencetrasanrepo.Delete(id)
		return success, failure
}
