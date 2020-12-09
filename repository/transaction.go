package repository

import (
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
)
//Transactionrepo...
var (
	Transactionrepo transactionrepo = transactionrepo{}
)

///curtesy to gorm
type transactionrepo struct{}

func (transactionRepo transactionrepo) Create(transaction *model.Transaction) (*model.Transaction, *httperors.HttpError) {
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&transaction)
	IndexRepo.DbClose(GormDB)
	return transaction, nil
}
func (transactionRepo transactionrepo) GetOne(id int) (*model.Transaction, *httperors.HttpError) {
	ok := transactionRepo.transactionUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("transaction with that id does not exists!")
	}
	transaction := model.Transaction{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	} 
	
	GormDB.Model(&transaction).Where("id = ?", id).First(&transaction)
	IndexRepo.DbClose(GormDB)
	
	return &transaction, nil
}
func (transactionRepo transactionrepo) All() (t []model.Transaction, r *httperors.HttpError) {

	transaction := model.Transaction{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&transaction).Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}

func (transactionRepo transactionrepo) Update(id int, transaction *model.Transaction) (*model.Transaction, *httperors.HttpError) {
	ok := transactionRepo.transactionUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("transaction with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	atransaction := model.Transaction{}
	
	GormDB.Model(&atransaction).Where("id = ?", id).First(&atransaction)
	// if transaction.Name  == "" {
	// 	transaction.Name = atransaction.Name
	// }
	// if transaction.Qty  == 0 {
	// 	transaction.Qty = atransaction.Qty
	// }
	// if transaction.Price  == 0 {
	// 	transaction.Price = atransaction.Price
	// }
	
	// if transaction.Discount  == 0 {
	// 	transaction.Discount = atransaction.Discount
	// }
	// if transaction.Tax  == 0 {
	// 	transaction.Tax = atransaction.Tax
	// }
	GormDB.Save(&transaction)
	
	IndexRepo.DbClose(GormDB)

	return transaction, nil
}
func (transactionRepo transactionrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := transactionRepo.transactionUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("transaction with that id does not exists!")
	}
	transaction := model.Transaction{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&transaction).Where("id = ?", id).First(&transaction)
	GormDB.Delete(transaction)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (transactionRepo transactionrepo)transactionUserExistByid(id int) bool {
	transaction := model.Transaction{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&transaction, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (transactionRepo transactionrepo) GetTransactionsinvoice(code string) (t []model.Transaction, e *httperors.HttpError) {
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	 
	GormDB.Where("code = ? AND status = ?", code, "invoice").Find(&t)
	IndexRepo.DbClose(GormDB)
	
	return t, nil
}
func (transactionRepo transactionrepo) GetTransactionscredit(code string) (t []model.Transaction, e *httperors.HttpError) {

	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	} 
	
	GormDB.Where("code = ? AND status = ?", code, "credit").Find(&t)
	IndexRepo.DbClose(GormDB)
	
	return t, nil
}
func (transactionRepo transactionrepo) GetTransactionspedingcredit(code string) (t []model.Transaction, e *httperors.HttpError) {
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Where("code = ? AND credit = ? AND status = ?", code, false, "pending").Find(&t)
	IndexRepo.DbClose(GormDB)
	
	return t, nil
}