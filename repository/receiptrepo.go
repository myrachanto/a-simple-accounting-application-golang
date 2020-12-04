package repository

import (
	"fmt"
	"strings"
	"strconv"
	"gorm.io/gorm"
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
	paymentform := model.Paymentform{}
	p := model.Paymentform{}
	if (receipt.Status == "cleared"){
		////////////begin transaction/////////////////////
	GormDB.Transaction(func(tx *gorm.DB) error {

		fmt.Println("level 1")
		tx.Create(&receipt)
		tx.Transaction(func(tx2 *gorm.DB) error {
			fmt.Println("level 2")
			tx2.Model(&paymentform).Where("name = ?", receipt.Paymentform).First(&p)
			updatedamount := p.Amount + receipt.Amount 
			tx2.Model(&paymentform).Where("name = ?", receipt.Paymentform).Update("amount", updatedamount)
			return nil
		})

		return nil
	})
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
func (receiptRepo receiptrepo) UpdateReceipts(code,status string) (string, *httperors.HttpError) {
	ok := Receiptrepo.ReceiptExistByCode(code)
	if ok == false {
		return "", httperors.NewNotFoundError("That receipt does not exist!")
	}
	r := model.Receipt{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	paymentform := model.Paymentform{}
	p := model.Paymentform{}
	if (r.Status == "cleared"){
		////////////begin transaction/////////////////////
	GormDB.Transaction(func(tx *gorm.DB) error {

		fmt.Println("level 1")
		tx.Model(&r).Where("code = ?", code).Update("status",status)

		tx.Transaction(func(tx2 *gorm.DB) error {
			fmt.Println("level 2")
			tx2.Model(&paymentform).Where("name = ?", r.Paymentform).First(&p)
			updatedamount := p.Amount + r.Amount 
			tx2.Model(&paymentform).Where("name = ?", r.Paymentform).Update("amount", updatedamount)
			return nil
		})

		return nil
	})
	}
	GormDB.Model(&r).Where("code = ?", code).Update("status",status)
	
	IndexRepo.DbClose(GormDB)
	return "Receipt updated succesifully", nil
}
func (receiptRepo receiptrepo) ViewReport() (*model.ReceiptReport, *httperors.HttpError) {
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	receipts := model.Receipt{}
	all := []model.Receipt{}
	cleared := []model.Receipt{}
	pending := []model.Receipt{}
	canceled := []model.Receipt{}
	GormDB.Model(&receipts).Find(&all)
	GormDB.Model(&receipts).Where("status = ?", "cancel").Find(&canceled)
	GormDB.Model(&receipts).Where("status = ?", "pending").Find(&pending)
	GormDB.Model(&receipts).Where("status = ?", "cleared").Find(&cleared)
	var clear float64 = 0
	for _,cl := range cleared {
		clear += cl.Amount
	}
	var pend float64 = 0
	for _,pen := range pending {
		pend += pen.Amount
	}
	var can float64 = 0
	for _,cn := range canceled {
		can += cn.Amount
	}
	
	z := model.ReceiptReport{}
	z.All = all
	z.ClearedRecipts.Name = "Cleared Receipts"
	z.ClearedRecipts.Total = clear
	z.ClearedRecipts.Description = "Total Amount Receipts cleared"
	//////////////////////////////////////////////////////////////
	z.PendingRecipts.Name = "Pending receipts"
	z.PendingRecipts.Total = pend
	z.PendingRecipts.Description = "Total Amount Receipts pending"
	///////////////////////////////////////////////////////////////
	z.CanceledRecipts.Name = "Cancelled receipts"
	z.CanceledRecipts.Total = can
	z.CanceledRecipts.Description = "Total Amount Receipts Cancelled"
	
	IndexRepo.DbClose(GormDB)
	return &z, nil
}
func (receiptRepo receiptrepo) ReceiptExistByCode(code string) bool {
	r := model.Receipt{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}

	res := GormDB.First(&r, "code =?", code)
	if res.Error != nil {
		return false
	}
	IndexRepo.DbClose(GormDB)
	return true

}
func (receiptRepo receiptrepo) View() (*model.ReceiptView, *httperors.HttpError) {
	r := &model.ReceiptView{}

	customers,err1 := Customerrepo.All()
	if err1 != nil {
		return nil, httperors.NewNotFoundError("You got an error fetching customers")
	}
	paymentforms,err7 := Paymentformrepo.All()
	if err7 != nil {
		return nil, httperors.NewNotFoundError("You got an error fetching customers")
	}
	code,err4 := Receiptrepo.GeneCode()
	if err4 != nil {
		return nil, httperors.NewNotFoundError("You got an error fetching customers")
	}
	r.Code = code
	r.Customers = customers
	r.Paymentform = paymentforms
	return r, nil
} 
func (receiptRepo receiptrepo) GetAll() (*model.ReceiptOptions, *httperors.HttpError) {
	
	receipts := model.Receipt{}
	all := []model.Receipt{}
	cleared := []model.Receipt{}
	pending := []model.Receipt{}
	canceled := []model.Receipt{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&receipts).Find(&all)
	GormDB.Model(&receipts).Where("status = ?", "cancel").Find(&canceled)
	GormDB.Model(&receipts).Where("status = ?", "pending").Find(&pending)
	GormDB.Model(&receipts).Where("status = ?", "cleared").Find(&cleared)
	if err1 != nil {
			return nil, err1
		}
		
	IndexRepo.DbClose(GormDB)
	return &model.ReceiptOptions{
		AllRecipts: all,
		ClearedRecipts: cleared,
		PendingRecipts: pending,
		CanceledRecipts: canceled,
	}, nil
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

func (receiptRepo receiptrepo)GeneCode() (string, *httperors.HttpError) {
	r := model.Receipt{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	err := GormDB.Last(&r)
	if err.Error != nil {
		var c1 uint = 1
		code := "ReceiptNo"+strconv.FormatUint(uint64(c1), 10)
		return code, nil
	 }
	c1 := r.ID + 1
	code := "ReceiptNo"+strconv.FormatUint(uint64(c1), 10)
	IndexRepo.DbClose(GormDB)
	return code, nil
	
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
		// fmt.Println(Ser.Search_query_1)
		if Ser.Search_query_1 == "all" {
			//db.Order("name DESC")
			GormDB.Order(Ser.Column + " " + Ser.Direction).Find(&receipts)

		} else {

			GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column + " " + Ser.Direction).Find(&receipts)
		}
		break
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