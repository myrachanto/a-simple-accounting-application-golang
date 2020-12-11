package repository

import (
	"fmt"
	"strings"
	"time"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/support"
)
//STransactionrepo... supplier transaction repository
var (
	STransactionrepo sTransactionrepo = sTransactionrepo{}
)

///curtesy to gorm
type sTransactionrepo struct{}

func (sTransactionRepo sTransactionrepo) Create(sTransaction *model.STransaction) (*model.STransaction, *httperors.HttpError) {
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&sTransaction)
	IndexRepo.DbClose(GormDB)
	return sTransaction, nil
}
func (sTransactionRepo sTransactionrepo) GetOne(id int) (*model.STransaction, *httperors.HttpError) {
	ok := sTransactionRepo.sTransactionUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("sTransaction with that id does not exists!")
	}
	sTransaction := model.STransaction{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	} 
	
	GormDB.Model(&sTransaction).Where("id = ?", id).First(&sTransaction)
	IndexRepo.DbClose(GormDB)
	
	return &sTransaction, nil
}
func (sTransactionRepo sTransactionrepo) All() (t []model.STransaction, r *httperors.HttpError) {

	sTransaction := model.STransaction{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&sTransaction).Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (sTransactionRepo sTransactionrepo) GetAll(sTransactions []model.STransaction,search *support.Search) ([]model.STransaction, *httperors.HttpError) {
	results, err1 := sTransactionRepo.Search(search, sTransactions)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}

func (sTransactionRepo sTransactionrepo) Update(id int, sTransaction *model.STransaction) (*model.STransaction, *httperors.HttpError) {
	ok := sTransactionRepo.sTransactionUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("sTransaction with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	asTransaction := model.STransaction{}
	
	GormDB.Model(&asTransaction).Where("id = ?", id).First(&asTransaction)
	// if sTransaction.Name  == "" {
	// 	sTransaction.Name = asTransaction.Name
	// }
	// if sTransaction.Qty  == 0 {
	// 	sTransaction.Qty = asTransaction.Qty
	// }
	// if sTransaction.Price  == 0 {
	// 	sTransaction.Price = asTransaction.Price
	// }
	
	// if sTransaction.Discount  == 0 {
	// 	sTransaction.Discount = asTransaction.Discount
	// }
	// if sTransaction.Tax  == 0 {
	// 	sTransaction.Tax = asTransaction.Tax
	// }
	GormDB.Save(&sTransaction)
	
	IndexRepo.DbClose(GormDB)

	return sTransaction, nil
}
func (sTransactionRepo sTransactionrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := sTransactionRepo.sTransactionUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("sTransaction with that id does not exists!")
	}
	sTransaction := model.STransaction{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&sTransaction).Where("id = ?", id).First(&sTransaction)
	GormDB.Delete(sTransaction)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (sTransactionRepo sTransactionrepo)sTransactionUserExistByid(id int) bool {
	sTransaction := model.STransaction{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&sTransaction, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (sTransactionRepo sTransactionrepo) ProductsBought(code,dated,searchq2,searchq3 string) (results []model.STransaction, r *httperors.HttpError) {
	now := time.Now()
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}

	if dated != "custom"{
		if dated == "In the last 24hrs"{
			d := now.AddDate(0, 0, -1)
			GormDB.Where("productcode = ? AND updated_at > ? AND credit = ?",code, d,false).Find(&results)
		}
		if dated == "In the last 7days"{
			d := now.AddDate(0, 0, -7)
			GormDB.Where("productcode = ? AND updated_at > ? AND credit = ?",code, d,false).Find(&results)
		}
		if dated == "In the last 15day"{
			d := now.AddDate(0, 0, -15)
			GormDB.Where("productcode = ? AND updated_at > ? AND credit = ?",code, d,false).Find(&results)
		}
		if dated == "In the last 30days"{
			d := now.AddDate(0, 0, -30)
			GormDB.Where("productcode = ? AND updated_at > ? AND credit = ?",code, d,false).Find(&results)
		}
	}
	if dated == "custom"{
		start,err := time.Parse(Layout,searchq2)
		if err != nil {
			return nil, httperors.NewNotFoundError("Something went wrong parsing date1!")
		}
		end,err1 := time.Parse(Layout,searchq3)
		if err1 != nil {
			return nil, httperors.NewNotFoundError("Something went wrong parsing date1!")
		}
		GormDB.Where("productcode = ? AND credit = ? AND updated_at BETWEEN ? AND ?",code,false, start, end).Find(&results)
	}
	IndexRepo.DbClose(GormDB)
	return results, nil

}
func (sTransactionRepo sTransactionrepo) GetsTransactionsinvoice(code string) (t []model.STransaction, e *httperors.HttpError) {
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Where("code = ? AND credit = ? AND status = ?", code, false, "invoice").Find(&t)
	IndexRepo.DbClose(GormDB)
	
	return t, nil
}
func (sTransactionRepo sTransactionrepo) GetsTransactionscredit(code string) (t []model.STransaction, e *httperors.HttpError) {

	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	} 
	
	GormDB.Where("code = ? AND credit = ? AND status = ?", code, true, "credit").Find(&t)
	IndexRepo.DbClose(GormDB)
	
	return t, nil 
}
func (sTransactionRepo sTransactionrepo) GetsTransactionspedingcredit(code string) (t []model.STransaction, e *httperors.HttpError) {
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Where("code = ? AND credit = ? AND status = ?", code, false, "pending").Find(&t)
	IndexRepo.DbClose(GormDB)
	
	return t, nil
}
func (sTransactionRepo sTransactionrepo) Search(Ser *support.Search, sTransactions []model.STransaction)([]model.STransaction, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	sTransaction := model.STransaction{}
	switch(Ser.Search_operator){
	case "all":
		GormDB.Model(&sTransaction).Order(Ser.Column+" "+Ser.Direction).Find(&sTransactions)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		
	break;
	case "equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&sTransactions);
		
	break;
	case "not_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&sTransactions);	
		
	break;
	case "less_than" :
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&sTransactions);	
		
	break;
	case "greater_than":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&sTransactions);	
		
	break;
	case "less_than_or_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&sTransactions);	
		
	break;
	case "greater_than_ro_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&sTransactions);	
		
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&sTransactions);
		
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		GormDB.Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&sTransactions);
		
	// break;
	case "like":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&sTransactions);
		
	break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&sTransactions);
		
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return sTransactions, nil
}