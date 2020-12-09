package repository

import (
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
)
//Messagerepo ... 
var (
	Messagerepo messagerepo = messagerepo{}
)

///curtesy to gorm
type messagerepo struct{}

func (messageRepo messagerepo) Create(message *model.Message) (*model.Message, *httperors.HttpError) {
	if err := message.Validate(); err != nil {
		return nil, err
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&message)
	IndexRepo.DbClose(GormDB)
	return message, nil
}
func (messageRepo messagerepo) GetOne(id int) (*model.Message, *httperors.HttpError) {
	ok := messageRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("message with that id does not exists!")
	}
	message := model.Message{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&message).Where("id = ?", id).First(&message)
	IndexRepo.DbClose(GormDB)
	
	return &message, nil
}


func (messageRepo messagerepo) GetAll(search string, page,pagesize int) ([]model.Message, *httperors.HttpError) {
	results := []model.Message{}
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

func (messageRepo messagerepo) Update(id int, message *model.Message) (*model.Message, *httperors.HttpError) {
	ok := messageRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("message with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&message).Where("id = ?", id).Save(&message)
	
	IndexRepo.DbClose(GormDB)

	return message, nil
}
func (messageRepo messagerepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := messageRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("Product with that id does not exists!")
	}
	message := model.Message{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&message).Where("id = ?", id).First(&message)
	GormDB.Delete(message)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (messageRepo messagerepo)ProductUserExistByid(id int) bool {
	message := model.Message{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&message, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
