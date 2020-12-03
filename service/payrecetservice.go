package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
	"github.com/myrachanto/accounting/support"
)

var (
	PayrectrasanService payrectrasanService = payrectrasanService{}

) 
type payrectrasanService struct {
	
}

func (service payrectrasanService) Create(payrectrasan *model.Payrectrasan) (*model.Payrectrasan, *httperors.HttpError) {
	
	payrectrasan, err1 := r.Payrectrasanrepo.Create(payrectrasan)
	if err1 != nil {
		return nil, err1
	}
	 return payrectrasan, nil

}
func (service payrectrasanService) GetOne(id int) (*model.Payrectrasan, *httperors.HttpError) {
	payrectrasan, err1 := r.Payrectrasanrepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return payrectrasan, nil
}
func (service payrectrasanService) View() (*model.Roptions, *httperors.HttpError) {
	code, err1 := r.Payrectrasanrepo.View()
	if err1 != nil {
		return nil, err1
	}
	return code, nil
}
func (service payrectrasanService) GetAll(payrectrasans []model.Payrectrasan,search *support.Search) ([]model.Payrectrasan, *httperors.HttpError) {
	payrectrasans, err := r.Payrectrasanrepo.GetAll(payrectrasans,search)
	if err != nil {
		return nil, err
	}
	return payrectrasans, nil
}

func (service payrectrasanService) Update(id int, payrectrasan *model.Payrectrasan) (*model.Payrectrasan, *httperors.HttpError) {
	payrectrasan, err1 := r.Payrectrasanrepo.Update(id, payrectrasan)
	if err1 != nil {
		return nil, err1
	}
	
	return payrectrasan, nil
}
func (service payrectrasanService) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Payrectrasanrepo.Delete(id)
		return success, failure
}
