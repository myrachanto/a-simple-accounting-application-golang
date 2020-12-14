package repository

import (
	"fmt"
	"strconv"
	"time"
	"gorm.io/gorm"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
)
//Paymentrepo ...
var (
	Paymentrepo paymentrepo = paymentrepo{}
)

///curtesy to gorm
type paymentrepo struct{} 

func (paymentRepo paymentrepo) Create(payment *model.Payment) (*model.Payment, *httperors.HttpError) {
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	} 
	ok := Supplierrepo.SupplierExistByname(payment.ItemName)
	itemcode := ""
	if ok {
		sup := Supplierrepo.Getsupplier(payment.ItemName)
		itemcode = sup.Suppliercode
		payment.Itemcode = itemcode
		payment.Allocated = "notallocated"
		payment.Mode = "invoice"
		paymentform := model.Paymentform{}
		p := model.Paymentform{}
		if (payment.Status == "cleared"){
			////////////begin transaction/////////////////////
		GormDB.Transaction(func(tx *gorm.DB) error {

			fmt.Println("level 1")
			tx.Create(&payment)
			tx.Transaction(func(tx2 *gorm.DB) error {
				fmt.Println("level 2")
				tx2.Model(&paymentform).Where("name = ?", payment.Paymentform).First(&p)
				updatedamount := p.Amount - payment.Amount 
				tx2.Model(&paymentform).Where("name = ?", payment.Paymentform).Update("amount", updatedamount)
				return nil
			})

		return nil
	})
	}
	GormDB.Create(&payment)
	}
	ok = Expencetrasanrepo.ExpenceExistByCode(payment.ItemName)
	if ok {
		sup := Expencetrasanrepo.ExpenceExistByNameGet(payment.ItemName)
		itemcode = sup.Code
		payment.Itemcode = itemcode
		payment.Allocated = "notallocated"
		payment.Mode = "other"
		paymentform := model.Paymentform{}
		p := model.Paymentform{}
		if (payment.Status == "cleared"){
			////////////begin transaction/////////////////////
		GormDB.Transaction(func(tx *gorm.DB) error {
	
			fmt.Println("level 1")
			tx.Create(&payment)
			tx.Transaction(func(tx2 *gorm.DB) error {
				fmt.Println("level 2")
				tx2.Model(&paymentform).Where("name = ?", payment.Paymentform).First(&p)
				updatedamount := p.Amount - payment.Amount 
				tx2.Model(&paymentform).Where("name = ?", payment.Paymentform).Update("amount", updatedamount)
				return nil
			})
	
			return nil
		})
		}
		Nortificationrepo.Create(&model.Nortification{
			Title:"you have a pending Payment",
			Description: "You have a peding payment worth" + strconv.FormatUint(uint64(payment.Amount), 10),
		})
		GormDB.Create(&payment)
	}
	IndexRepo.DbClose(GormDB)
	return payment, nil
}
func (paymentRepo paymentrepo) GetOne(id int) (*model.Payment, *httperors.HttpError) {
	ok := paymentRepo.paymentUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("payment with that id does not exists!")
	}
	payment := model.Payment{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&payment).Where("id = ?", id).First(&payment)
	IndexRepo.DbClose(GormDB)
	
	return &payment, nil
}

func (paymentRepo paymentrepo) AddReceiptTrans(clientcode,invoicecode,usercode,receiptcode string ,amount float64) (string, *httperors.HttpError) {
	ok := Supplierrepo.SupplierExistbycode(clientcode)
	if !ok {
		return "", httperors.NewNotFoundError("That customer does not exist")
	}
	supplier := Supplierrepo.Getsupplierwithcode(clientcode)
	ok = SInvoicerepo.SInvoiceExistByCode(invoicecode)
	if !ok {
		return "", httperors.NewNotFoundError("That invoice does not exist")
	}
	
	invo := SInvoicerepo.GetInvoicebyCode(invoicecode)
	stats := ""
	if invo.Total == amount {
		stats = "fullypaid"
	}
	stats = "partialpaid"
	ok = Paymentrepo.ReceiptExistByCode(receiptcode)
	if !ok {
		return "", httperors.NewNotFoundError("That receipt does not exist")
	}
	payment := Paymentrepo.GetpaymentwithCode(receiptcode)
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	transact := model.Payrectrasan{}
	transact.Name = supplier.Name
	transact.Title = "Payments"
	transact.Description = "Payments to  Supplier"
	transact.CLientcode = clientcode
	transact.Invoicecode = invoicecode
	transact.Amount = amount
	transact.Usercode = usercode
	transact.Paymentform = payment.Type
	transact.Status = stats
	paymentf := Paymentformrepo.GetPaymantformbyname(payment.Type)

	invoic := model.SInvoice{}
	paymentform := model.Paymentform{}
	reciep := model.Payment{}
	////////////begin transaction/////////////////////
	GormDB.Transaction(func(tx *gorm.DB) error {
		
		fmt.Println("level 1")
		tx.Create(&transact)

		
		tx.Transaction(func(tx2 *gorm.DB) error { 
		
			fmt.Println("level 2")
			bal := invo.Total - amount
			tx2.Model(&invoic).Where("code = ?", invoicecode).Updates(model.SInvoice{Paidstatus: stats, AllPaidstatus: stats, AmountPaid: amount, Balance:bal})
			return nil
		})
			remaining := paymentf.Amount - amount
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
func (paymentRepo paymentrepo) ReceiptExistByCode(code string) bool {
	r := model.Payment{}
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
func (paymentRepo paymentrepo)GetpaymentwithCode(code string) *model.Payment {
	p := model.Payment{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil
	}
	GormDB.Where("code = ? ", code).First(&p)
	if p.ID == 0 {
	   return nil
	}
	IndexRepo.DbClose(GormDB)
	return &p
	
}
func (paymentRepo paymentrepo) Updatepayments(code,status string) (string, *httperors.HttpError) {
	ok := Paymentrepo.paymentExistByCode(code)
	if ok == false {
		return "", httperors.NewNotFoundError("That payment does not exist!")
	}
	r := model.Payment{}
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
			updatedamount := p.Amount - r.Amount 
			tx2.Model(&paymentform).Where("name = ?", r.Paymentform).Update("amount", updatedamount)
			return nil
		})

		return nil
	})
	}
	GormDB.Model(&r).Where("code = ?", code).Update("status",status)
	
	IndexRepo.DbClose(GormDB)
	return "payment updated succesifully", nil
}
func (paymentRepo paymentrepo) ViewCleared() ([]model.Payment, *httperors.HttpError) {
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1  
	}
	payment := model.Payment{}
	cleared := []model.Payment{}
	GormDB.Model(&payment).Where("status = ? AND allocated = ? AND mode = ?", "cleared", "notallocated", "invoice").Find(&cleared)
	IndexRepo.DbClose(GormDB)
	return cleared, nil

}
func (paymentRepo paymentrepo) ViewClearedExpence() ([]model.Payment, *httperors.HttpError) {
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1  
	}
	payment := model.Payment{}
	cleared := []model.Payment{}
	GormDB.Model(&payment).Where("allocated = ? AND mode = ?", "notallocated", "other").Find(&cleared)
	IndexRepo.DbClose(GormDB)
	return cleared, nil

}
func (paymentRepo paymentrepo) ViewInvoices(code string) (*model.PaymentAlloc, *httperors.HttpError) {
	fmt.Println(code)
	rec := Paymentrepo.GetreceiptwithCode(code)
	if rec == nil {
		return nil, httperors.NewNotFoundError("That Payment does not exist")
	}
	invoices, err := SInvoicerepo.InvoiceByCustomercodenotpaid(rec.Itemcode)
	if err != nil {
		return nil,err
	}
	
	return &model.PaymentAlloc{
		Payment: rec,
		SInvoice:invoices,
	}, nil
}
func (paymentRepo paymentrepo)GetreceiptwithCode(code string) *model.Payment {
	payment := model.Payment{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil
	}
	GormDB.Where("code = ? ", code).First(&payment)
	if payment.ID == 0 {
	   return nil
	}
	IndexRepo.DbClose(GormDB)
	return &payment
	
}
func (paymentRepo paymentrepo) paymentExistByCode(code string) bool {
	r := model.Payment{}
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
func (paymentRepo paymentrepo) View() (*model.PaymentView, *httperors.HttpError) {
	r := &model.PaymentView{}

	suppliers,err1 := Supplierrepo.All()
	if err1 != nil {
		return nil, httperors.NewNotFoundError("You got an error fetching customers")
	}
	paymentforms,err7 := Paymentformrepo.All()
	if err7 != nil {
		return nil, httperors.NewNotFoundError("You got an error fetching customers")
	}
	code,err4 := Paymentrepo.GeneCode()
	if err4 != nil {
		return nil, httperors.NewNotFoundError("You got an error fetching customers")
	}
	r.Code = code
	r.Suppliers = suppliers
	r.Paymentform = paymentforms
	return r, nil
} 
func (paymentRepo paymentrepo) ViewExpence() (*model.PaymentExpence, *httperors.HttpError) {
	r := &model.PaymentExpence{}

	expences,err1 := Expencetrasanrepo.All()
	if err1 != nil {
		return nil, httperors.NewNotFoundError("You got an error fetching customers")
	}
	paymentforms,err7 := Paymentformrepo.All()
	if err7 != nil {
		return nil, httperors.NewNotFoundError("You got an error fetching customers")
	}
	code,err4 := Paymentrepo.GeneCode()
	if err4 != nil {
		return nil, httperors.NewNotFoundError("You got an error fetching customers")
	}
	r.Code = code
	r.Expences = expences
	r.Paymentform = paymentforms
	return r, nil
} 
func (paymentRepo paymentrepo) GetAll(dated,searchq2,searchq3 string) (*model.PaymentOptions, *httperors.HttpError) {
	
	all, err1 := Paymentrepo.AllSearch(dated,searchq2,searchq3)
	if err1 != nil {
		return nil, err1
	}
	cleared, err2 := Paymentrepo.AllCleared(dated,searchq2,searchq3)
	if err2 != nil {
		return nil, err2
	}
	
	pending, err3 := Paymentrepo.AllPending(dated,searchq2,searchq3)
	if err3 != nil {
		return nil, err3
	}
	canceled, err4 := Paymentrepo.AllCanceled(dated,searchq2,searchq3)
	if err4 != nil {
		return nil, err4
	}
	return &model.PaymentOptions{
		AllPayments: all,
		ClearedPayments: cleared,
		PendingPayments: pending,
		CanceledPayments: canceled,
	}, nil
}
func (paymentRepo paymentrepo) AllSearch(dated,searchq2,searchq3 string) (results []model.Payment, r *httperors.HttpError) {

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
func (paymentRepo paymentrepo) AllCleared(dated,searchq2,searchq3 string) (results []model.Payment, r *httperors.HttpError) {

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
func (paymentRepo paymentrepo) AllPending(dated,searchq2,searchq3 string) (results []model.Payment, r *httperors.HttpError) {

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
func (paymentRepo paymentrepo) AllCanceled(dated,searchq2,searchq3 string) (results []model.Payment, r *httperors.HttpError) {

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
func (paymentRepo paymentrepo) ViewReport(dated,searchq2,searchq3 string) (*model.PaymentReport, *httperors.HttpError) {
	all, err1 := Paymentrepo.AllSearch(dated,searchq2,searchq3)
	if err1 != nil {
		return nil, err1
	}
	cleared, err2 := Paymentrepo.AllCleared(dated,searchq2,searchq3)
	if err2 != nil {
		return nil, err2
	}
	
	pending, err3 := Paymentrepo.AllPending(dated,searchq2,searchq3)
	if err3 != nil {
		return nil, err3
	}
	canceled, err4 := Paymentrepo.AllCanceled(dated,searchq2,searchq3)
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
	
	z := model.PaymentReport{}
	z.All = all
	z.ClearedPayments.Name = "Cleared Payments"
	z.ClearedPayments.Total = clear
	z.ClearedPayments.Description = "Total Amount Payments cleared"
	//////////////////////////////////////////////////////////////
	z.PendingPayments.Name = "Pending payments"
	z.PendingPayments.Total = pend
	z.PendingPayments.Description = "Total Amount Payments pending"
	///////////////////////////////////////////////////////////////
	z.CanceledPayments.Name = "Cancelled Payments"
	z.CanceledPayments.Total = can
	z.CanceledPayments.Description = "Total Amount Payments Cancelled"
	
	return &z, nil
}
func (paymentRepo paymentrepo) All() (t []model.Payment, r *httperors.HttpError) {

	rec := model.Payment{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&rec).Where("status = ?", "cleared").Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (paymentRepo paymentrepo) Update(id int, payment *model.Payment) (*model.Payment, *httperors.HttpError) {
	ok := paymentRepo.paymentUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("payment with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	apayment := model.Payment{}
	
	GormDB.Model(&apayment).Where("id = ?", id).First(&apayment)
	// if payment.payment  == "" {
	// 	payment.payment = apayment.payment
	// }
	// if payment.Description  == "" {
	// 	payment.Description = apayment.Description
	// }
	// if payment.Subtotal  == 0 {
	// 	payment.Subtotal = apayment.Subtotal
	// }
	// if payment.Discount  == 0 {
	// 	payment.Discount = apayment.Discount
	// }	
	// if payment.AmountPaid  == 0 {
	// 	payment.AmountPaid = apayment.AmountPaid
	// }
	GormDB.Save(&payment)
	
	IndexRepo.DbClose(GormDB)

	return payment, nil
}
func (paymentRepo paymentrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := paymentRepo.paymentUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("payment with that id does not exists!")
	}
	payment := model.Payment{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&payment).Where("id = ?", id).First(&payment)
	GormDB.Delete(payment)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (paymentRepo paymentrepo)paymentUserExistByid(id int) bool {
	payment := model.Payment{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&payment, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}

func (paymentRepo paymentrepo)GeneCode() (string, *httperors.HttpError) {
	r := model.Payment{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	err := GormDB.Last(&r)
	if err.Error != nil {
		var c1 uint = 1
		code := "PaymentNo"+strconv.FormatUint(uint64(c1), 10)
		return code, nil
	 }
	c1 := r.ID + 1
	code := "PaymentNo"+strconv.FormatUint(uint64(c1), 10)
	IndexRepo.DbClose(GormDB)
	return code, nil
	
}