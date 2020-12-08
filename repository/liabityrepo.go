package repository

import (
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
)
 //Liabilityrepo ...
var (
	Liabilityrepo liabilityrepo = liabilityrepo{}
)

///curtesy to gorm
type liabilityrepo struct{}

func (liabilityRepo liabilityrepo) Create(liability *model.Liability) (*model.Liability, *httperors.HttpError) {
	if err := liability.Validate(); err != nil {
		return nil, err
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&liability)
	IndexRepo.DbClose(GormDB)
	return liability, nil
}
func (liabilityRepo liabilityrepo) GetOne(id int) (*model.Liability, *httperors.HttpError) {
	ok := liabilityRepo.liabilityUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("liability with that id does not exists!")
	}
	liability := model.Liability{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Preload("Liatrans").Model(&liability).Where("id = ?", id).First(&liability)
	IndexRepo.DbClose(GormDB)
	
	return &liability, nil
}

func (liabilityRepo liabilityrepo) GetAll(search string) ([]model.Liability, *httperors.HttpError) {
	results := []model.Liability{}
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

func (liabilityRepo liabilityrepo) Update(id int, liability *model.Liability) (*model.Liability, *httperors.HttpError) {
	ok := liabilityRepo.liabilityUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("liability with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&liability).Where("id = ?", id).Save(&liability)
	IndexRepo.DbClose(GormDB)

	return liability, nil
}
func (liabilityRepo liabilityrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := liabilityRepo.liabilityUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("liability with that id does not exists!")
	}
	liability := model.Liability{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	} 
	GormDB.Model(&liability).Where("id = ?", id).First(&liability)
	GormDB.Delete(liability)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (liabilityRepo liabilityrepo)liabilityUserExistByid(id int) bool {
	liability := model.Liability{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&liability, "id =?", id)
	if res.Error != nil{
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}