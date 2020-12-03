package repository

import (
	"fmt"
	"strings"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/support"
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

func (payrectrasanRepo payrectrasanrepo) GetAll(payrectrasans []model.Payrectrasan,search *support.Search) ([]model.Payrectrasan, *httperors.HttpError) {
	results, err1 := payrectrasanRepo.Search(search, payrectrasans)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
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

func (payrectrasanRepo payrectrasanrepo) Search(Ser *support.Search, payrectrasans []model.Payrectrasan)([]model.Payrectrasan, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	payrectrasan := model.Payrectrasan{}
	switch(Ser.Search_operator){
	case "all":
	GormDB.Model(&payrectrasan).Order(Ser.Column+" "+Ser.Direction).Find(&payrectrasans)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
	
	break;
	case "equal_to":
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&payrectrasans);
	
	break;
	case "not_equal_to":
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&payrectrasans);	
	
	break;
	case "less_than" :
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&payrectrasans);	
	
	break;
	case "greater_than":
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&payrectrasans);	
	
	break;
	case "less_than_or_equal_to":
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&payrectrasans);	
	
	break;
	case "greater_than_ro_equal_to":
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&payrectrasans);	
	
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&payrectrasans);
	
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
	GormDB.Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&payrectrasans);
	
	// break;
	case "like":
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&payrectrasans);
	
	break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&payrectrasans);
	
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return payrectrasans, nil
}