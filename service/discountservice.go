package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
)
//DiscountService ...
var (
	DiscountService discountService = discountService{}

) 
type discountService struct {
	
}

func (service discountService) Create(discount *model.Discount) (*model.Discount, *httperors.HttpError) {
	if err := discount.Validate(); err != nil {
		return nil, err
	}	
	discount, err1 := r.Discountrepo.Create(discount)
	if err1 != nil {
		return nil, err1
	}
	 return discount, nil

}
func (service discountService) GetOne(id int) (*model.Discount, *httperors.HttpError) {
	discount, err1 := r.Discountrepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return discount, nil
}


func (service discountService) GetAll(search string, page,pagesize int) ([]model.Discount, *httperors.HttpError) {
	results, err := r.Discountrepo.GetAll(search, page,pagesize)
	return results, err
}
func (service discountService) Update(id int, discount *model.Discount) (*model.Discount, *httperors.HttpError) {
	discount, err1 := r.Discountrepo.Update(id, discount)
	if err1 != nil {
		return nil, err1
	}
	
	return discount, nil
}
func (service discountService) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Discountrepo.Delete(id)
		return success, failure
}
