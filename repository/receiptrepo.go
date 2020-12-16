package repository

import (
	"fmt"
	"strconv"
	"time"
	"gorm.io/gorm"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
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
	cust := Customerrepo.Getcustomer(receipt.CustomerName)
	receipt.Customercode = cust.Customercode
	receipt.Allocated = "notallocated"
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
	
	Nortificationrepo.Create(&model.Nortification{
		Title:"you have a pending Recipts",
		Description: "You have a peding receipts worth" + strconv.FormatUint(uint64(receipt.Amount), 10),
	})
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
func (receiptRepo receiptrepo) ViewInvoices(code string) (*model.ReceiptAlloc, *httperors.HttpError) {
	fmt.Println(code)
	rec := Receiptrepo.GetreceiptwithCode(code)
	if rec == nil {
		return nil, httperors.NewNotFoundError("That Receipt does not exist")
	}
	invoices, err := Invoicerepo.InvoiceByCustomercodenotpaid(rec.Customercode)
	if err != nil {
		return nil,err
	}
	
	return &model.ReceiptAlloc{
		Receipt: rec,
		Invoice:invoices,
	}, nil
}

func (receiptRepo receiptrepo) AddReceiptTrans(clientcode,invoicecode,usercode,receiptcode string ,amount float64) (string, *httperors.HttpError) {
	ok := Customerrepo.CustomerExistbycode(clientcode)
	if !ok {
		return "", httperors.NewNotFoundError("That customer does not exist")
	}
	customer := Customerrepo.GetcustomerwithCode(clientcode)
	ok = Invoicerepo.InvoiceExistByCode(invoicecode)
	if !ok {
		return "", httperors.NewNotFoundError("That invoice does not exist")
	}
	invo := Invoicerepo.GetInvoicebyCode(invoicecode)
	stats := ""
	if invo.Total == amount {
		stats = "fullypaid"
	}
	stats = "partialpaid"
fmt.Println(receiptcode)
	ok = Receiptrepo.ReceiptExistByCode(receiptcode)
	if !ok {
		return "", httperors.NewNotFoundError("That receipt does not exist")
	}
	receipt := Receiptrepo.GetreceiptwithCode(receiptcode)
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	transact := model.Payrectrasan{}
	transact.Name = customer.Name
	transact.Title = "Receipts"
	transact.Description = "Receipts from customer"
	transact.CLientcode = clientcode
	transact.Invoicecode = invoicecode
	transact.Amount = amount
	transact.Usercode = usercode
	transact.Paymentform = receipt.Type
	transact.Status = stats
	paymentf := Paymentformrepo.GetPaymantformbyname(receipt.Type)

	invoic := model.Invoice{}
	paymentform := model.Paymentform{}
	reciep := model.Receipt{}
	////////////begin transaction/////////////////////
	GormDB.Transaction(func(tx *gorm.DB) error {
		
		fmt.Println("level 1")
		tx.Create(&transact)

		
		tx.Transaction(func(tx2 *gorm.DB) error { 
		
			fmt.Println("level 2")
			bal := invo.Total - amount
			tx2.Model(&invoic).Where("code = ?", invoicecode).Updates(model.Invoice{Paidstatus: stats, AllPaidstatus: stats, AmountPaid: amount, Balance:bal})
			return nil
		})
			remaining := paymentf.Amount + amount
			tx.Transaction(func(tx4 *gorm.DB) error {
				fmt.Println("level 4")
				tx4.Model(&paymentform).Where("paymentcode = ?", paymentf.Paymentcode).Update("amount",remaining)
				return nil
			})
			tx.Transaction(func(tx4 *gorm.DB) error {
				fmt.Println("level 4")
				tx4.Model(&reciep).Where("code = ?", receiptcode).Update("allocated","allocated")
				return nil
			})
			return nil
		})
	
	IndexRepo.DbClose(GormDB)
	return "transaction completed succesifully", nil
}

func (receiptRepo receiptrepo)GetreceiptwithCode(code string) *model.Receipt {
	receipt := model.Receipt{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil
	}
	GormDB.Where("code = ? ", code).First(&receipt)
	if receipt.ID == 0 {
	   return nil
	}
	IndexRepo.DbClose(GormDB)
	return &receipt
	
}
func (receiptRepo receiptrepo)GetreceiptwithCodeCustomer(code, customercode string) *model.Receipt {
	receipt := model.Receipt{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil
	}
	GormDB.Where("code = ? AND customercode = ?", code, customercode).First(&receipt)
	if receipt.ID == 0 {
	   return nil
	}
	IndexRepo.DbClose(GormDB)
	return &receipt
	
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
func (receiptRepo receiptrepo) ViewReport(dated,searchq2,searchq3 string) (*model.ReceiptReport, *httperors.HttpError) {
	all, err1 := Receiptrepo.AllSearch(dated,searchq2,searchq3)
	if err1 != nil {
		return nil, err1
	}
	cleared, err2 := Receiptrepo.AllCleared(dated,searchq2,searchq3)
	if err2 != nil {
		return nil, err2
	}
	
	pending, err3 := Receiptrepo.AllPending(dated,searchq2,searchq3)
	if err3 != nil {
		return nil, err3
	}
	canceled, err4 := Receiptrepo.AllCanceled(dated,searchq2,searchq3)
	if err4 != nil {
		return nil, err4
	}
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
	
	return &z, nil
}

func (receiptRepo receiptrepo) All() (t []model.Receipt, r *httperors.HttpError) {

	receipt := model.Receipt{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&receipt).Where("status = ?", "cleared").Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

} 
func (receiptRepo receiptrepo) ViewCleared() ([]model.Receipt, *httperors.HttpError) {
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1  
	}
	receipts := model.Receipt{}
	cleared := []model.Receipt{}
	GormDB.Model(&receipts).Where("status = ? AND allocated = ?", "cleared", "notallocated").Find(&cleared)
	IndexRepo.DbClose(GormDB)
	return cleared, nil

}
func (receiptRepo receiptrepo) ReceiptExistByCode(code string) bool {
	r := model.Receipt{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}

	GormDB.Where("code = ? ", code).First(&r)
	if r.ID == 0 {
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
func (receiptRepo receiptrepo) GetAll(dated,searchq2,searchq3 string) (*model.ReceiptOptions, *httperors.HttpError) {
	
	all, err1 := Receiptrepo.AllSearch(dated,searchq2,searchq3)
	if err1 != nil {
		return nil, err1
	}
	cleared, err2 := Receiptrepo.AllCleared(dated,searchq2,searchq3)
	if err2 != nil {
		return nil, err2
	}
	
	pending, err3 := Receiptrepo.AllPending(dated,searchq2,searchq3)
	if err3 != nil {
		return nil, err3
	}
	canceled, err4 := Receiptrepo.AllCanceled(dated,searchq2,searchq3)
	if err4 != nil {
		return nil, err4
	}
	return &model.ReceiptOptions{
		AllRecipts: all,
		ClearedRecipts: cleared,
		PendingRecipts: pending,
		CanceledRecipts: canceled,
	}, nil
}
func (receiptRepo receiptrepo) AllSearch(dated,searchq2,searchq3 string) (results []model.Receipt, r *httperors.HttpError) {

	now := time.Now()
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	if dated != "custom"{ 
		if dated == "In the last 24hrs"{
			d := now.AddDate(0, 0, -1)
			GormDB.Where("updated_at > ?", d).Find(&results)
		}
		if dated == "In the last 7days"{
			d := now.AddDate(0, 0, -7)
			GormDB.Where("updated_at > ?",d).Find(&results)
		}
		if dated == "In the last 15day"{
			d := now.AddDate(0, 0, -15)
			GormDB.Where("updated_at > ?",d).Find(&results)
		}
		if dated == "In the last 30days"{
			d := now.AddDate(0, 0, -30)
			GormDB.Where("updated_at > ?",d).Find(&results)
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
		GormDB.Where("updated_at BETWEEN ? AND ?",start, end).Find(&results)
	}
	IndexRepo.DbClose(GormDB)
	return results, nil

}
func (receiptRepo receiptrepo) AllCleared(dated,searchq2,searchq3 string) (results []model.Receipt, r *httperors.HttpError) {

	now := time.Now()
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}

	if dated != "custom"{
		if dated == "In the last 24hrs"{
			d := now.AddDate(0, 0, -1)
			GormDB.Where("updated_at > ? AND status = ?", d,"cleared").Find(&results)
		}
		if dated == "In the last 7days"{
			d := now.AddDate(0, 0, -7)
			GormDB.Where("updated_at > ? AND status = ?", d,"cleared").Find(&results)
		}
		if dated == "In the last 15day"{
			d := now.AddDate(0, 0, -15)
			GormDB.Where("updated_at > ? AND status = ?", d,"cleared").Find(&results)
		}
		if dated == "In the last 30days"{
			d := now.AddDate(0, 0, -30)
			GormDB.Where("updated_at > ? AND status = ?", d,"cleared").Find(&results)
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
		GormDB.Where("status = ? AND updated_at BETWEEN ? AND ?","cleared", start, end).Find(&results)
	}
	IndexRepo.DbClose(GormDB)
	return results, nil

}
func (receiptRepo receiptrepo) AllPending(dated,searchq2,searchq3 string) (results []model.Receipt, r *httperors.HttpError) {

	now := time.Now()
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}

	if dated != "custom"{
		if dated == "In the last 24hrs"{
			d := now.AddDate(0, 0, -1)
			GormDB.Where("updated_at > ? AND status = ?", d,"pending").Find(&results)
		}
		if dated == "In the last 7days"{
			d := now.AddDate(0, 0, -7)
			GormDB.Where("updated_at > ? AND status = ?", d,"pending").Find(&results)
		}
		if dated == "In the last 15day"{
			d := now.AddDate(0, 0, -15)
			GormDB.Where("updated_at > ? AND status = ?", d,"pending").Find(&results)
		}
		if dated == "In the last 30days"{
			d := now.AddDate(0, 0, -30)
			GormDB.Where("updated_at > ? AND status = ?", d,"pending").Find(&results)
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
		GormDB.Where("status = ? AND updated_at BETWEEN ? AND ?","pending", start, end).Find(&results)
	}
	IndexRepo.DbClose(GormDB)
	return results, nil

}
func (receiptRepo receiptrepo) AllCanceled(dated,searchq2,searchq3 string) (results []model.Receipt, r *httperors.HttpError) {

	now := time.Now()
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}

	if dated != "custom"{
		if dated == "In the last 24hrs"{
			d := now.AddDate(0, 0, -1)
			GormDB.Where("updated_at > ? AND status = ?", d,"cancel").Find(&results)
		}
		if dated == "In the last 7days"{
			d := now.AddDate(0, 0, -7)
			GormDB.Where("updated_at > ? AND status = ?", d,"cancel").Find(&results)
		}
		if dated == "In the last 15day"{
			d := now.AddDate(0, 0, -15)
			GormDB.Where("updated_at > ? AND status = ?", d,"cancel").Find(&results)
		}
		if dated == "In the last 30days"{
			d := now.AddDate(0, 0, -30)
			GormDB.Where("updated_at > ? AND status = ?", d,"cancel").Find(&results)
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
		GormDB.Where("status = ? AND updated_at BETWEEN ? AND ?","cancel", start, end).Find(&results)
	}
	IndexRepo.DbClose(GormDB)
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
func (receiptRepo receiptrepo) AllSearchCustomer(code,dated,searchq2,searchq3 string) (results []model.Receipt, r *httperors.HttpError) {

	now := time.Now()
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	if dated != "custom"{ 
		if dated == "In the last 24hrs"{
			d := now.AddDate(0, 0, -1)
			GormDB.Where("updated_at > ? AND customercode = ?", d,code).Find(&results)
		}
		if dated == "In the last 7days"{
			d := now.AddDate(0, 0, -7)
			GormDB.Where("updated_at > ? AND customercode = ?",d,code).Find(&results)
		}
		if dated == "In the last 15day"{
			d := now.AddDate(0, 0, -15)
			GormDB.Where("updated_at > ? AND customercode = ?",d,code).Find(&results)
		}
		if dated == "In the last 30days"{
			d := now.AddDate(0, 0, -30)
			GormDB.Where("updated_at > ? AND customercode = ?",d,code).Find(&results)
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
		GormDB.Where("customercode = ? AND updated_at BETWEEN ? AND ?",code,start, end).Find(&results)
	}
	IndexRepo.DbClose(GormDB)
	return results, nil

}