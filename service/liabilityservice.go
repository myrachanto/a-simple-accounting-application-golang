package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
	"github.com/myrachanto/accounting/support"
)

var (
	Liabilityservice liabilityservice = liabilityservice{}

) 
type liabilityservice struct {
	
}

func (service liabilityservice) Create(liability *model.Liability) (*model.Liability, *httperors.HttpError) {
	liability, err1 := r.Liabilityrepo.Create(liability)
	if err1 != nil {
		return nil, err1
	}
	 return liability, nil

}
func (service liabilityservice) GetOne(id int) (*model.Liability, *httperors.HttpError) {
	liability, err1 := r.Liabilityrepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return liability, nil
}

func (service liabilityservice) GetAll(liabilitys []model.Liability, search *support.Search) ([]model.Liability, *httperors.HttpError) {
	
	results, err := r.Liabilityrepo.GetAll(liabilitys, search)
	if err != nil {
		return nil, err
	}
	return results, nil 
}

func (service liabilityservice) Update(id int, liability *model.Liability) (*model.Liability, *httperors.HttpError) {
	liability, err1 := r.Liabilityrepo.Update(id, liability)
	if err1 != nil {
		return nil, err1
	}
	
	return liability, nil
}
func (service liabilityservice) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Liabilityrepo.Delete(id)
		return success, failure
}
///////deleting a batch////////////////////

//db.Where("age = ?", 20).Delete(&User{})