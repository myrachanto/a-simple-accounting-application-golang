package repository

import (
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
)
//Categoryrepo ...
var (
	Categoryrepo categoryrepo = categoryrepo{}
)


///curtesy to gorm 
type categoryrepo struct{}

func (categoryRepo categoryrepo) Create(category *model.Category) (*model.Category, *httperors.HttpError) {
	if err := category.Validate(); err != nil {
		return nil, err
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&category)
	IndexRepo.DbClose(GormDB)
	return category, nil
}
func (categoryRepo categoryrepo) View() ([]model.Majorcategory, *httperors.HttpError) {
	mc, e := Majorcategoryrepo.All()
	if e != nil{
		return nil, e
	}
	return mc, nil
}
func (categoryRepo categoryrepo) GetMajorcat(name string) (*model.Category, *httperors.HttpError) {
	category := model.Category{}	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&category).Where("name = ?", name).First(&category)
	IndexRepo.DbClose(GormDB)
	
	return &category, nil
}
func (categoryRepo categoryrepo) GetOne(id int) (*model.Category, *httperors.HttpError) {
	ok := categoryRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("category with that id does not exists!")
	}
	category := model.Category{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&category).Where("id = ?", id).First(&category)
	IndexRepo.DbClose(GormDB)
	
	return &category, nil
}

func (categoryRepo categoryrepo) GetAll(search string, page,pagesize int) ([]model.Category, *httperors.HttpError) {
	results := []model.Category{}
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
func (categoryRepo categoryrepo) All() (t []model.Category, r *httperors.HttpError) {

	category := model.Category{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&category).Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (categoryRepo categoryrepo) Update(id int, category *model.Category) (*model.Category, *httperors.HttpError) {
	ok := categoryRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("category with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&category).Where("id = ?", id).Save(&category)
	
	IndexRepo.DbClose(GormDB)

	return category, nil
}
func (categoryRepo categoryrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := categoryRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("Product with that id does not exists!")
	}
	category := model.Category{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&category).Where("id = ?", id).First(&category)
	GormDB.Delete(category)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (categoryRepo categoryrepo)ProductUserExistByid(id int) bool {
	category := model.Category{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&category, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
