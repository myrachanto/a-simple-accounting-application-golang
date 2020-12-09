package repository

import (
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
)
//Nortificationrepo ...
var (
	Nortificationrepo nortificationrepo = nortificationrepo{}
)

///curtesy to gorm
type nortificationrepo struct{}

func (nortificationRepo nortificationrepo) Create(nortification *model.Nortification) (*model.Nortification, *httperors.HttpError) {
	if err := nortification.Validate(); err != nil {
		return nil, err
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&nortification)
	IndexRepo.DbClose(GormDB)
	return nortification, nil
}
func (nortificationRepo nortificationrepo) GetOne(id int) (*model.Nortification, *httperors.HttpError) {
	ok := nortificationRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("nortification with that id does not exists!")
	}
	nortification := model.Nortification{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&nortification).Where("id = ?", id).First(&nortification)
	IndexRepo.DbClose(GormDB)
	
	return &nortification, nil
}

func (nortificationRepo nortificationrepo) GetAll(search string, page,pagesize int) ([]model.Nortification, *httperors.HttpError) {
	results := []model.Nortification{}
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

func (nortificationRepo nortificationrepo) Update(id int, nortification *model.Nortification) (*model.Nortification, *httperors.HttpError) {
	ok := nortificationRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("nortification with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&nortification).Where("id = ?", id).Save(&nortification)
	
	IndexRepo.DbClose(GormDB)

	return nortification, nil
}
func (nortificationRepo nortificationrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := nortificationRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("Product with that id does not exists!")
	}
	nortification := model.Nortification{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&nortification).Where("id = ?", id).First(&nortification)
	GormDB.Delete(nortification)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (nortificationRepo nortificationrepo)ProductUserExistByid(id int) bool {
	nortification := model.Nortification{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&nortification, "id =?", id)
	if res.Error != nil{
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
