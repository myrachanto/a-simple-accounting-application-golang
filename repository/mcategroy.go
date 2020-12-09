package repository

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
)
//Majorcategoryrepo ...
var (
	Majorcategoryrepo majorcategoryrepo = majorcategoryrepo{}
)

///curtesy to gorm
type majorcategoryrepo struct{}

func (majorcategoryRepo majorcategoryrepo) Create(majorcategory *model.Majorcategory) (*model.Majorcategory, *httperors.HttpError) {
	if err := majorcategory.Validate(); err != nil {
		return nil, err
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&majorcategory)
	IndexRepo.DbClose(GormDB)
	return majorcategory, nil
}
func (majorcategoryRepo majorcategoryrepo) GetOne(id int) (*model.Majorcategory, *httperors.HttpError) {
	ok := majorcategoryRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("majorcategory with that id does not exists!")
	}
	majorcategory := model.Majorcategory{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&majorcategory).Where("id = ?", id).First(&majorcategory)
	IndexRepo.DbClose(GormDB)
	
	return &majorcategory, nil
}

func (majorcategoryRepo majorcategoryrepo) GetAll(search string, page,pagesize int) ([]model.Majorcategory, *httperors.HttpError) {
	results := []model.Majorcategory{}
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
func (majorcategoryRepo majorcategoryrepo) All() (t []model.Majorcategory, r *httperors.HttpError) {

	majorcategory := model.Majorcategory{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&majorcategory).Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (majorcategoryRepo majorcategoryrepo) Update(id int, majorcategory *model.Majorcategory) (*model.Majorcategory, *httperors.HttpError) {
	ok := majorcategoryRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("majorcategory with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&majorcategory).Where("id = ?", id).Save(&majorcategory)
	
	IndexRepo.DbClose(GormDB)

	return majorcategory, nil
}
func (majorcategoryRepo majorcategoryrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := majorcategoryRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("Product with that id does not exists!")
	}
	majorcategory := model.Majorcategory{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&majorcategory).Where("id = ?", id).First(&majorcategory)
	GormDB.Delete(majorcategory)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (majorcategoryRepo majorcategoryrepo)ProductUserExistByid(id int) bool {
	majorcategory := model.Majorcategory{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&majorcategory, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
