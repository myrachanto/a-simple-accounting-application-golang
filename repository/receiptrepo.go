package repository

import (
	"fmt"
	"strings"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/support"
)
//Receiptrepo ...
var (
	Receiptrepo receiptrepo = receiptrepo{}
)

///curtesy to gorm
type receiptrepo struct{} 

func (receiptRepo receiptrepo) Create(receipt *model.Receipt) (*model.Receipt, *httperors.HttpError) {
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&receipt)
	IndexRepo.DbClose(GormDB)
	return receipt, nil
}
func (receiptRepo receiptrepo) GetOne(id int) (*model.Receipt, *httperors.HttpError) {
	ok := receiptRepo.receiptUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("receipt with that id does not exists!")
	}
	receipt := model.Receipt{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&receipt).Where("id = ?", id).First(&receipt)
	IndexRepo.DbClose(GormDB)
	
	return &receipt, nil
}

func (receiptRepo receiptrepo) GetAll(receipts []model.Receipt,search *support.Search) ([]model.Receipt, *httperors.HttpError) {
	
	results, err1 := receiptRepo.Search(search, receipts)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}

func (receiptRepo receiptrepo) Update(id int, receipt *model.Receipt) (*model.Receipt, *httperors.HttpError) {
	ok := receiptRepo.receiptUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("receipt with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	areceipt := model.Receipt{}
	
	GormDB.Model(&areceipt).Where("id = ?", id).First(&areceipt)
	// if receipt.receipt  == "" {
	// 	receipt.receipt = areceipt.receipt
	// }
	// if receipt.Description  == "" {
	// 	receipt.Description = areceipt.Description
	// }
	// if receipt.Subtotal  == 0 {
	// 	receipt.Subtotal = areceipt.Subtotal
	// }
	// if receipt.Discount  == 0 {
	// 	receipt.Discount = areceipt.Discount
	// }	
	// if receipt.AmountPaid  == 0 {
	// 	receipt.AmountPaid = areceipt.AmountPaid
	// }
	GormDB.Save(&receipt)
	
	IndexRepo.DbClose(GormDB)

	return receipt, nil
}
func (receiptRepo receiptrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := receiptRepo.receiptUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("receipt with that id does not exists!")
	}
	receipt := model.Receipt{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&receipt).Where("id = ?", id).First(&receipt)
	GormDB.Delete(receipt)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (receiptRepo receiptrepo)receiptUserExistByid(id int) bool {
	receipt := model.Receipt{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&receipt, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}

func (receiptRepo receiptrepo) Search(Ser *support.Search, receipts []model.Receipt)([]model.Receipt, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	receipt := model.Receipt{}
	switch(Ser.Search_operator){
	case "all":
		GormDB.Model(&receipt).Order(Ser.Column+" "+Ser.Direction).Find(&receipts)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		
	break;
	case "equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&receipts);
		
	break;
	case "not_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&receipts);	
		
	break;
	case "less_than" :
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&receipts);	
		
	break;
	case "greater_than":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&receipts);	
		
	break;
	case "less_than_or_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&receipts);	
		
	break;
	case "greater_than_ro_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&receipts);	
		
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&receipts);
		
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		GormDB.Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&receipts);
		
	// break;
	case "like":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&receipts);
		
	break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&receipts);
		
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return receipts, nil
}