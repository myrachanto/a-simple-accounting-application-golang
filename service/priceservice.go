package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
)
//PriceService ...
var (
	PriceService priceService = priceService{}

) 
type priceService struct {
	
}

func (service priceService) Create(price *model.Price) (*model.Price, *httperors.HttpError) {
	if err := price.Validate(); err != nil {
		return nil, err
	}	
	price, err1 := r.Pricerepo.Create(price)
	if err1 != nil {
		return nil, err1
	}
	 return price, nil

}
func (service priceService) View() ([]model.Product, *httperors.HttpError) {
	options, err1 := r.Pricerepo.View()
	if err1 != nil {
		return nil, err1
	}
	return options, nil
}
func (service priceService) GetOne(id int) (*model.Price, *httperors.HttpError) {
	price, err1 := r.Pricerepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return price, nil
}

func (service priceService) GetAll(search string, page,pagesize int) ([]model.Price, *httperors.HttpError) {
	results, err := r.Pricerepo.GetAll(search, page,pagesize)
	return results, err
}

func (service priceService) Update(id int, price *model.Price) (*model.Price, *httperors.HttpError) {
	price, err1 := r.Pricerepo.Update(id, price)
	if err1 != nil {
		return nil, err1
	}
	
	return price, nil
}
func (service priceService) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Pricerepo.Delete(id)
		return success, failure
}
