package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
	"github.com/myrachanto/accounting/support"
)

var (
	TaxService taxService = taxService{}

) 
type taxService struct {
	
}

func (service taxService) Create(tax *model.Tax) (*model.Tax, *httperors.HttpError) {
	if err := tax.Validate(); err != nil {
		return nil, err
	}	
	tax, err1 := r.Taxrepo.Create(tax)
	if err1 != nil {
		return nil, err1
	}
	 return tax, nil

}
func (service taxService) GetOne(id int) (*model.Tax, *httperors.HttpError) {
	tax, err1 := r.Taxrepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return tax, nil
}

func (service taxService) GetAll(taxs []model.Tax,search *support.Search) ([]model.Tax, *httperors.HttpError) {
	taxs, err := r.Taxrepo.GetAll(taxs,search)
	if err != nil {
		return nil, err
	}
	return taxs, nil
}

func (service taxService) Update(id int, tax *model.Tax) (*model.Tax, *httperors.HttpError) {
	tax, err1 := r.Taxrepo.Update(id, tax)
	if err1 != nil {
		return nil, err1
	}
	
	return tax, nil
}
func (service taxService) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Taxrepo.Delete(id)
		return success, failure
}
