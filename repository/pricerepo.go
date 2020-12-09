package repository

import (
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
)
//Pricerepo ...
var (
	Pricerepo pricerepo = pricerepo{}
)

///curtesy to gorm
type pricerepo struct{}

func (priceRepo pricerepo) Create(price *model.Price) (*model.Price, *httperors.HttpError) {
	if err := price.Validate(); err != nil {
		return nil, err
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&price)
	IndexRepo.DbClose(GormDB)
	return price, nil
}
func (priceRepo pricerepo) View() ([]model.Product, *httperors.HttpError) {
	p, e := Productrepo.All()
	if e != nil{
		return nil, e
	}
	return p, nil
}
func (priceRepo pricerepo) All() (t []model.Price, r *httperors.HttpError) {

	price := model.Price{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&price).Find(&t)

	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (priceRepo pricerepo) GetOne(id int) (*model.Price, *httperors.HttpError) {
	ok := priceRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("price with that id does not exists!")
	}
	price := model.Price{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&price).Where("id = ?", id).First(&price)
	IndexRepo.DbClose(GormDB)
	
	return &price, nil
}
func (priceRepo pricerepo) GetOption(id int)([]model.Price, *httperors.HttpError){
	ok := Productrepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("product with that id does not exists!")
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	prices := []model.Price{}
	GormDB.Where("id = ? AND buy = ? ", id, true).Find(&prices)
	return prices, nil
}
func (priceRepo pricerepo) GetOptionsell(id int)([]model.Price, *httperors.HttpError){
	ok := Productrepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("product with that id does not exists!")
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	prices := []model.Price{}
	GormDB.Where("id = ? AND buy = ? ", id, false).Find(&prices)
	return prices, nil
}
func (priceRepo pricerepo) GetAll(search string, page,pagesize int) ([]model.Price, *httperors.HttpError) {
	results := []model.Price{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	if search == ""{
		GormDB.Find(&results)
	}
	// db.Scopes(Paginate(r)).Find(&users)
	GormDB.Scopes(Paginate(page,pagesize)).Where("name LIKE ?", "%"+search+"%").Or("title LIKE ?", "%"+search+"%").Or("description LIKE ?", "%"+search+"%").Find(&results)

	IndexRepo.DbClose(GormDB)
	return results, nil
}

func (priceRepo pricerepo) Update(id int, price *model.Price) (*model.Price, *httperors.HttpError) {
	ok := priceRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("price with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&price).Where("id = ?", id).Save(&price)

	return price, nil
}
func (priceRepo pricerepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := priceRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("Product with that id does not exists!")
	}
	price := model.Price{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&price).Where("id = ?", id).First(&price)
	GormDB.Delete(price)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (priceRepo pricerepo)ProductUserExistByid(id int) bool {
	price := model.Price{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&price, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
