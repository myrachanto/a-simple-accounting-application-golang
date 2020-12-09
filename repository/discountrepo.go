package repository

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
)
//Discountrepo ...
var (
	Discountrepo discountrepo = discountrepo{}
)

///curtesy to gorm
type discountrepo struct{}

func (discountRepo discountrepo) Create(discount *model.Discount) (*model.Discount, *httperors.HttpError) {
	if err := discount.Validate(); err != nil {
		return nil, err
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&discount)
	IndexRepo.DbClose(GormDB)
	return discount, nil
}
func (discountRepo discountrepo) All() (t []model.Discount, r *httperors.HttpError) {

	discount := model.Discount{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&discount).Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (discountRepo discountrepo) GetOne(id int) (*model.Discount, *httperors.HttpError) {
	ok := discountRepo.DiscountExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("discount with that id does not exists!")
	}
	discount := model.Discount{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&discount).Where("id = ?", id).First(&discount)
	IndexRepo.DbClose(GormDB)
	
	return &discount, nil
}
func (discountRepo discountrepo) GetOption(id int)([]model.Discount, *httperors.HttpError){
	ok := Productrepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("product with that id does not exists!")
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	discounts := []model.Discount{}
	GormDB.Where("id = ? AND buy = ? ", id, true).Find(&discounts)
	return discounts, nil
}
func (discountRepo discountrepo) GetOptionsell(id int)([]model.Discount, *httperors.HttpError){
	ok := Productrepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("product with that id does not exists!")
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	discounts := []model.Discount{}
	GormDB.Where("id = ? AND buy = ? ", id, false).Find(&discounts)
	return discounts, nil
}
func (discountRepo discountrepo) GetAll(search string, page,pagesize int) ([]model.Discount, *httperors.HttpError) {
	results := []model.Discount{}
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

func (discountRepo discountrepo) Update(id int, discount *model.Discount) (*model.Discount, *httperors.HttpError) {
	ok := discountRepo.DiscountExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("discount with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&discount).Where("id = ?", id).Save(&discount)
	return discount, nil
}
func (discountRepo discountrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := discountRepo.DiscountExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("Product with that id does not exists!")
	}
	discount := model.Discount{} 
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&discount).Where("id = ?", id).First(&discount)
	GormDB.Delete(discount)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (discountRepo discountrepo)DiscountExistByid(id int) bool {
	discount := model.Discount{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	
	GormDB.Where("id = ? ", id).First(&discount)
	if discount.ID == 0 {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}