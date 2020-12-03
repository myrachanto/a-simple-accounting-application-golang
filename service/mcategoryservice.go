package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/support"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
)

var (
	Mcategoryservice mcategoryservice = mcategoryservice{}

) 
type mcategoryservice struct {
	
}

func (service mcategoryservice) Create(majorcategory *model.Majorcategory) (*model.Majorcategory, *httperors.HttpError) {
	if err := majorcategory.Validate(); err != nil {
		return nil, err
	}	
	majorcategory, err1 := r.Majorcategoryrepo.Create(majorcategory)
	if err1 != nil {
		return nil, err1
	}
	 return majorcategory, nil

}
func (service mcategoryservice) GetOne(id int) (*model.Majorcategory, *httperors.HttpError) {
	majorcategory, err1 := r.Majorcategoryrepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return majorcategory, nil
}

func (service mcategoryservice) GetAll(majorcategorys []model.Majorcategory,search *support.Search) ([]model.Majorcategory, *httperors.HttpError) {
	majorcategorys, err := r.Majorcategoryrepo.GetAll(majorcategorys,search)
	if err != nil {
		return nil, err
	}
	return majorcategorys, nil
}

func (service mcategoryservice) Update(id int, majorcategory *model.Majorcategory) (*model.Majorcategory, *httperors.HttpError) {
	majorcategory, err1 := r.Majorcategoryrepo.Update(id, majorcategory)
	if err1 != nil {
		return nil, err1
	}
	
	return majorcategory, nil
}
func (service mcategoryservice) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Majorcategoryrepo.Delete(id)
		return success, failure
}
