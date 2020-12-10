package repository

import (
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
)
//Payrectrasanrepo ...
var (
	Payrectrasanrepo payrectrasanrepo = payrectrasanrepo{}
)

///curtesy to gorm
type payrectrasanrepo struct{}

func (payrectrasanRepo payrectrasanrepo) Create(payrectrasan *model.Payrectrasan) (*model.Payrectrasan, *httperors.HttpError) {

	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	if payrectrasan.Receipt.ID > 0 {
		id := payrectrasan.Receipt.ID 
		payrectrasan.ReceiptID = id
		GormDB.Create(&payrectrasan) 
		return payrectrasan, nil
	}
	id := payrectrasan.Payment.ID
	payrectrasan.PaymentID = id
	GormDB.Create(&payrectrasan) 
	IndexRepo.DbClose(GormDB)
	return payrectrasan, nil
}
func (payrectrasanRepo payrectrasanrepo) View() (*model.Roptions, *httperors.HttpError) {
	c,e := Customerrepo.GetOptions();if e != nil {
		return nil,e
	}
	m,me := Supplierrepo.GetOptions();if me != nil {
		return nil,me
	}
	options := model.Roptions{}
	options.Customer = c
	options.Supplier = m
	return &options, nil
}
func (payrectrasanRepo payrectrasanrepo) GetOne(id int) (*model.Payrectrasan, *httperors.HttpError) {
	ok := payrectrasanRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("payrectrasan with that id does not exists!")
	}
	payrectrasan := model.Payrectrasan{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&payrectrasan).Where("id = ?", id).First(&payrectrasan)
	IndexRepo.DbClose(GormDB)
	
	return &payrectrasan, nil
}


func (payrectrasanRepo payrectrasanrepo) Update(id int, payrectrasan *model.Payrectrasan) (*model.Payrectrasan, *httperors.HttpError) {
	ok := payrectrasanRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("payrectrasan with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	// payrectrasan := model.Payrectrasan{}
	apayrectrasan := model.Payrectrasan{}
	
	GormDB.Model(&payrectrasan).Where("id = ?", id).First(&apayrectrasan)
	if payrectrasan.Name  == "" {
		payrectrasan.Name = apayrectrasan.Name
	}
	if payrectrasan.Title  == "" {
		payrectrasan.Title = apayrectrasan.Title
	}
	if payrectrasan.Description  == "" {
		payrectrasan.Description = apayrectrasan.Description
	}
	GormDB.Save(&payrectrasan)
	
	IndexRepo.DbClose(GormDB)

	return payrectrasan, nil
}
func (payrectrasanRepo payrectrasanrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := payrectrasanRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("Product with that id does not exists!")
	}
	payrectrasan := model.Payrectrasan{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&payrectrasan).Where("id = ?", id).First(&payrectrasan)
	GormDB.Delete(payrectrasan)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (payrectrasanRepo payrectrasanrepo)ProductUserExistByid(id int) bool {
	payrectrasan := model.Payrectrasan{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&payrectrasan, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}