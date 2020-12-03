package repository

import (
	"fmt"
	"strings"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/support"
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
func (transactionRepo transactionrepo) GetAll(transactions []model.Transaction,search *support.Search) ([]model.Transaction, *httperors.HttpError) {
	results, err1 := transactionRepo.Search(search, transactions)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
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
	
	GormDB.Where("code = ? AND credit = ? AND status = ?", code, false, "invoice").Find(&t)
	IndexRepo.DbClose(GormDB)
	
	return t, nil
}
func (transactionRepo transactionrepo) GetTransactionscredit(code string) (t []model.Transaction, e *httperors.HttpError) {

	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	} 
	
	GormDB.Where("code = ? AND credit = ? AND status = ?", code, true, "credit").Find(&t)
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
func (transactionRepo transactionrepo) Search(Ser *support.Search, transactions []model.Transaction)([]model.Transaction, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	transaction := model.Transaction{}
	switch(Ser.Search_operator){
	case "all":
		GormDB.Model(&transaction).Order(Ser.Column+" "+Ser.Direction).Find(&transactions)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		
	break;
	case "equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&transactions);
		
	break;
	case "not_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&transactions);	
		
	break;
	case "less_than" :
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&transactions);	
		
	break;
	case "greater_than":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&transactions);	
		
	break;
	case "less_than_or_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&transactions);	
		
	break;
	case "greater_than_ro_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&transactions);	
		
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&transactions);
		
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		GormDB.Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&transactions);
		
	// break;
	case "like":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&transactions);
		
	break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&transactions);
		
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return transactions, nil
}