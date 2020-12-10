package repository

import (
	"strconv"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
)
//Expencerepo ..
var (
	Expencerepo expencerepo = expencerepo{}
)

///curtesy to gorm
type expencerepo struct{}

func (expenceRepo expencerepo) Create(expence *model.Expence) (*model.Expence, *httperors.HttpError) {
	if err := expence.Validate(); err != nil {
		return nil, err
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	code, x := expenceRepo.GeneCode()
	if x != nil {
		return nil, x
	}
	expence.ExpenceCode = code
	GormDB.Create(&expence)
	IndexRepo.DbClose(GormDB)
	return expence, nil
}
func (expenceRepo expencerepo) GetOne(id int) (*model.Expence, *httperors.HttpError) {
	ok := expenceRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("expence with that id does not exists!")
	}
	expence := model.Expence{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Preload("Expencetrasans").Model(&expence).Where("id = ?", id).First(&expence)
	IndexRepo.DbClose(GormDB)
	
	return &expence, nil
}
func (expenceRepo expencerepo) All() (t []model.Expence, r *httperors.HttpError) {

	expence := model.Expence{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&expence).Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (expenceRepo expencerepo) GetAll(search string, page,pagesize int) ([]model.Expence, *httperors.HttpError) {
	results := []model.Expence{}
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

func (expenceRepo expencerepo) Update(id int, expence *model.Expence) (*model.Expence, *httperors.HttpError) {
	ok := expenceRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("expence with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&expence).Where("id = ?", id).Save(&expence)
	
	IndexRepo.DbClose(GormDB)

	return expence, nil
}
func (expenceRepo expencerepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := expenceRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("Product with that id does not exists!")
	}
	expence := model.Expence{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&expence).Where("id = ?", id).First(&expence)
	GormDB.Delete(expence)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (expenceRepo expencerepo)ProductUserExistByid(id int) bool {
	expence := model.Expence{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false 
	}
	res := GormDB.First(&expence, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (expenceRepo expencerepo)GeneCode() (string, *httperors.HttpError) {
	ex := model.Expence{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	err := GormDB.Last(&ex)
	if err.Error != nil {
		var c1 uint = 1
		code := "ExpenceCode"+strconv.FormatUint(uint64(c1), 10)
		return code, nil
	 }
	c1 := ex.ID + 1
	code := "ExpenceCode"+strconv.FormatUint(uint64(c1), 10)
	IndexRepo.DbClose(GormDB)
	return code, nil
	
}
