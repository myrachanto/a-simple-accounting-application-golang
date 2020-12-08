package repository

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
)
//Taxrepo ...
var (
	Taxrepo taxrepo = taxrepo{}
)

///curtesy to gorm
type taxrepo struct{}

func (taxRepo taxrepo) Create(tax *model.Tax) (*model.Tax, *httperors.HttpError) {
	if err := tax.Validate(); err != nil {
		return nil, err
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&tax)
	IndexRepo.DbClose(GormDB)
	return tax, nil
}
func (taxRepo taxrepo) GetOne(id int) (*model.Tax, *httperors.HttpError) {
	ok := taxRepo.TaxExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("tax with that id does not exists!")
	}
	tax := model.Tax{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&tax).Where("id = ?", id).First(&tax)
	IndexRepo.DbClose(GormDB)
	
	return &tax, nil
}
func (taxRepo taxrepo) All() (t []model.Tax, r *httperors.HttpError) {

	tax := model.Tax{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&tax).Find(&t)

	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (taxRepo taxrepo) GetOption()([]model.Tax, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	tax := model.Tax{}
	taxs := []model.Tax{}
	GormDB.Model(&tax).Find(&taxs)
	return taxs, nil
}
func (taxRepo taxrepo) GetAll(search string) ([]model.Tax, *httperors.HttpError) {

	results := []model.Tax{}
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

func (taxRepo taxrepo) Update(id int, tax *model.Tax) (*model.Tax, *httperors.HttpError) {
	ok := taxRepo.TaxExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("tax with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&tax).Where("id = ?", id).Save(&tax)
	
	IndexRepo.DbClose(GormDB)

	return tax, nil
}
func (taxRepo taxrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := taxRepo.TaxExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("Product with that id does not exists!")
	}
	tax := model.Tax{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&tax).Where("id = ?", id).First(&tax)
	GormDB.Delete(tax)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (taxRepo taxrepo)TaxExistByid(id int) bool {
	tax := model.Tax{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	GormDB.Where("id = ? ", id).First(&tax)
	if tax.ID == 0 {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}