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
func (paymentRepo paymentrepo) GetAll() (*model.PaymentOptions, *httperors.HttpError) {
	
	payments := model.Payment{}
	all := []model.Payment{}
	cleared := []model.Payment{}
	pending := []model.Payment{}
	canceled := []model.Payment{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&payments).Find(&all)
	GormDB.Model(&payments).Where("status = ?", "cancel").Find(&canceled)
	GormDB.Model(&payments).Where("status = ?", "pending").Find(&pending)
	GormDB.Model(&payments).Where("status = ?", "cleared").Find(&cleared)
	if err1 != nil {
			return nil, err1
		}
		
	IndexRepo.DbClose(GormDB)
	return &model.PaymentOptions{
		AllPayments: all,
		ClearedPayments: cleared,
		PendingPayments: pending,
		CanceledPayments: canceled,
	}, nil
}
func (paymentRepo paymentrepo) ViewReport() (*model.PaymentReport, *httperors.HttpError) {
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	payment := model.Payment{}
	all := []model.Payment{}
	cleared := []model.Payment{}
	pending := []model.Payment{}
	canceled := []model.Payment{}
	GormDB.Model(&payment).Find(&all)
	GormDB.Model(&payment).Where("status = ?", "cancel").Find(&canceled)
	GormDB.Model(&payment).Where("status = ?", "pending").Find(&pending)
	GormDB.Model(&payment).Where("status = ?", "cleared").Find(&cleared)
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
	
	IndexRepo.DbClose(GormDB)
	return &z, nil
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
func (paymentRepo paymentrepo) Search(Ser *support.Search, payments []model.Payment)([]model.Payment, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	payment := model.Payment{}
	switch(Ser.Search_operator){
	case "all":
		GormDB.Model(&payment).Order(Ser.Column+" "+Ser.Direction).Find(&payments)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		
	break;
	case "equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&payments);
		
	break;
	case "not_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&payments);	
		
	break;
	case "less_than" :
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&payments);	
		
	break;
	case "greater_than":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&payments);	
		
	break;
	case "less_than_or_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&payments);	
		
	break;
	case "greater_than_ro_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&payments);	
		
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&payments);
		
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		GormDB.Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&payments);
		
	// break;
	case "like":
		// fmt.Println(Ser.Search_query_1)
		if Ser.Search_query_1 == "all" {
			//db.Order("name DESC")
			GormDB.Order(Ser.Column + " " + Ser.Direction).Find(&payments)

		} else {

			GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column + " " + Ser.Direction).Find(&payments)
		}
		break
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&payments);
		
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return payments, nil
}