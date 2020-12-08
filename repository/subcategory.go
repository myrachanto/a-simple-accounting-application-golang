package repository

import (
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
)
//Subcategoryrepo ...
var (
	Subcategoryrepo subcategoryrepo = subcategoryrepo{}
)

///curtesy to gorm
type subcategoryrepo struct{}

func (subcategoryRepo subcategoryrepo) Create(subcategory * model.Subcategory) (* model.Subcategory, *httperors.HttpError) {
	if err := subcategory.Validate(); err != nil {
		return nil, err
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&subcategory)
	IndexRepo.DbClose(GormDB)
	return subcategory, nil
}
func (subcategoryRepo subcategoryrepo) GetOne(id int) (* model.Subcategory, *httperors.HttpError) {
	ok := subcategoryRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("subcategory with that id does not exists!")
	}
	subcategory :=  model.Subcategory{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&subcategory).Where("id = ?", id).First(&subcategory)
	IndexRepo.DbClose(GormDB)
	
	return &subcategory, nil
}

func (subcategoryRepo subcategoryrepo) GetAll(search string) ([] model.Subcategory, *httperors.HttpError) {
	
	results := []model.Subcategory{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	if search == ""{
		GormDB.Find(&results)
	}
	GormDB.Where("name LIKE ?", "%"+search+"%").Or("title LIKE ?", "%"+search+"%").Or("description LIKE ?", "%"+search+"%").Find(&results)

	IndexRepo.DbClose(GormDB)
	return results, nil
}

func (subcategoryRepo subcategoryrepo) Update(id int, subcategory * model.Subcategory) (* model.Subcategory, *httperors.HttpError) {
	ok := subcategoryRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("subcategory with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&subcategory).Where("id = ?", id).Save(&subcategory)
	
	IndexRepo.DbClose(GormDB)

	return subcategory, nil
}
func (subcategoryRepo subcategoryrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := subcategoryRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("Product with that id does not exists!")
	}
	subcategory :=  model.Subcategory{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&subcategory).Where("id = ?", id).First(&subcategory)
	GormDB.Delete(subcategory)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (subcategoryRepo subcategoryrepo)ProductUserExistByid(id int) bool {
	subcategory :=  model.Subcategory{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&subcategory, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
