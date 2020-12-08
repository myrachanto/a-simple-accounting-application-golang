package repository

import (
	"fmt"
	"strings"
	"strconv"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/support"
)
//Paymentformrepo ...
var (
	Paymentformrepo paymentformrepo = paymentformrepo{}
)

///curtesy to gorm
type paymentformrepo struct{}

func (paymentformRepo paymentformrepo) Create(paymentform *model.Paymentform) (*model.Paymentform, *httperors.HttpError) {
	if err := paymentform.Validate(); err != nil {
		return nil, err
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}

	code, x := Paymentrepo.GeneCode()
	if x != nil {
		return nil, x
	}
	paymentform.Paymentcode = code
	GormDB.Create(&paymentform)
	IndexRepo.DbClose(GormDB)
	return paymentform, nil
}
func (paymentformRepo paymentformrepo)GetPaymantbyCode(code string) *model.Invoice {
	invoice := model.Invoice{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil
	}
	GormDB.Where("code = ? ", code).First(&invoice)
	if invoice.ID == 0 {
	   return nil
	}
	IndexRepo.DbClose(GormDB)
	return &invoice
	
}

func (paymentformRepo paymentformrepo)GetPaymantformbyname(name string) *model.Paymentform {
	paymentform := model.Paymentform{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil
	}
	GormDB.Where("name = ? ", name).First(&paymentform)
	if paymentform.ID == 0 {
	   return nil
	}
	IndexRepo.DbClose(GormDB)
	return &paymentform
	
}
func (paymentformRepo paymentformrepo)GeneCode() (string, *httperors.HttpError) {
	p := model.Paymentform{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	err := GormDB.Last(&p)
	if err.Error != nil {
		var c1 uint = 1
		code := "PaymentformCode"+strconv.FormatUint(uint64(c1), 10)
		return code, nil
	 }
	c1 := p.ID + 1
	code := "PaymentformCode"+strconv.FormatUint(uint64(c1), 10)
	IndexRepo.DbClose(GormDB)
	return code, nil
	
}
func (paymentformRepo paymentformrepo) GetOne(id int) (*model.Paymentform, *httperors.HttpError) {
	ok := paymentformRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("paymentform with that id does not exists!")
	}
	paymentform := model.Paymentform{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&paymentform).Where("id = ?", id).First(&paymentform)
	IndexRepo.DbClose(GormDB)
	
	return &paymentform, nil
}

func (paymentformRepo paymentformrepo) GetAll(search string) ([]model.Paymentform, *httperors.HttpError) {
	results := []model.Paymentform{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	if search == ""{
		GormDB.Find(&results)
	}
	GormDB.Where("name LIKE ?", "%"+search+"%").Or("title LIKE ?", "%"+search+"%").Or("description LIKE ?", "%"+search+"%").Find(&results)

	IndexRepo.DbClose(GormDB)
	return results, nil
}
func (paymentformRepo paymentformrepo) All() (t []model.Paymentform, r *httperors.HttpError) {

	paymentform := model.Paymentform{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&paymentform).Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (paymentformRepo paymentformrepo) Update(id int, paymentform *model.Paymentform) (*model.Paymentform, *httperors.HttpError) {
	ok := paymentformRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("paymentform with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&paymentform).Where("id = ?", id).Save(&paymentform)
	
	IndexRepo.DbClose(GormDB)

	return paymentform, nil
}
func (paymentformRepo paymentformrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := paymentformRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("Product with that id does not exists!")
	}
	paymentform := model.Paymentform{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&paymentform).Where("id = ?", id).First(&paymentform)
	GormDB.Delete(paymentform)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (paymentformRepo paymentformrepo)ProductUserExistByid(id int) bool {
	paymentform := model.Paymentform{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&paymentform, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}

func (paymentformRepo paymentformrepo) Search(Ser *support.Search, paymentforms []model.Paymentform)([]model.Paymentform, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	paymentform := model.Paymentform{}
	switch(Ser.Search_operator){
	case "all":
		GormDB.Model(&paymentform).Order(Ser.Column+" "+Ser.Direction).Find(&paymentforms)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		
	break;
	case "equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&paymentforms);
		
	break;
	case "not_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&paymentforms);	
		
	break;
	case "less_than" :
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&paymentforms);	
		
	break;
	case "greater_than":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&paymentforms);	
		
	break;
	case "less_than_or_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&paymentforms);	
		
	break;
	case "greater_than_ro_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&paymentforms);	
		
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&paymentforms);
		
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		GormDB.Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&paymentforms);
		
	// break;
	case "like":
		// fmt.Println(Ser.Search_query_1)
		if Ser.Search_query_1 == "all" {
				//db.Order("name DESC")
		GormDB.Order(Ser.Column+" "+Ser.Direction).Find(&paymentforms)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		

		}else {

			GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&paymentforms);
		
		}
	break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&paymentforms);
		
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return paymentforms, nil
}