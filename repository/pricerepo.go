package repository

import (
	"fmt"
	"strings"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/support"
)
//Pricerepo ...
var (
	Pricerepo pricerepo = pricerepo{}
)

///curtesy to gorm
type pricerepo struct{}

func (priceRepo pricerepo) Create(price *model.Price) (*model.Price, *httperors.HttpError) {
	if err := price.Validate(); err != nil {
		return nil, err
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&price)
	IndexRepo.DbClose(GormDB)
	return price, nil
}
func (priceRepo pricerepo) View() ([]model.Product, *httperors.HttpError) {
	p, e := Productrepo.All()
	if e != nil{
		return nil, e
	}
	return p, nil
}
func (priceRepo pricerepo) All() (t []model.Price, r *httperors.HttpError) {

	price := model.Price{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&price).Find(&t)

	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (priceRepo pricerepo) GetOne(id int) (*model.Price, *httperors.HttpError) {
	ok := priceRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("price with that id does not exists!")
	}
	price := model.Price{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&price).Where("id = ?", id).First(&price)
	IndexRepo.DbClose(GormDB)
	
	return &price, nil
}
func (priceRepo pricerepo) GetOption(id int)([]model.Price, *httperors.HttpError){
	ok := Productrepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("product with that id does not exists!")
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	prices := []model.Price{}
	GormDB.Where("id = ? AND buy = ? ", id, true).Find(&prices)
	return prices, nil
}
func (priceRepo pricerepo) GetOptionsell(id int)([]model.Price, *httperors.HttpError){
	ok := Productrepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("product with that id does not exists!")
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	prices := []model.Price{}
	GormDB.Where("id = ? AND buy = ? ", id, false).Find(&prices)
	return prices, nil
}
func (priceRepo pricerepo) GetAll(prices []model.Price,search *support.Search) ([]model.Price, *httperors.HttpError) {
	results, err1 := priceRepo.Search(search, prices)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}

func (priceRepo pricerepo) Update(id int, price *model.Price) (*model.Price, *httperors.HttpError) {
	ok := priceRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("price with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	aprice := model.Price{}
	
	GormDB.Model(&price).Where("id = ?", id).First(&aprice)
	if price.Name  == "" {
		price.Name = aprice.Name
	}
	if price.Title  == "" {
		price.Title = aprice.Title
	}
	if price.Description  == "" {
		price.Description = aprice.Description
	}
	if price.Product  == "" {
		price.Product = aprice.Product
	}
	GormDB.Save(&price)

	return price, nil
}
func (priceRepo pricerepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := priceRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("Product with that id does not exists!")
	}
	price := model.Price{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&price).Where("id = ?", id).First(&price)
	GormDB.Delete(price)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (priceRepo pricerepo)ProductUserExistByid(id int) bool {
	price := model.Price{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&price, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}

func (priceRepo pricerepo) Search(Ser *support.Search, prices []model.Price)([]model.Price, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	price := model.Price{}
	switch(Ser.Search_operator){
	case "all":
		GormDB.Model(&price).Order(Ser.Column+" "+Ser.Direction).Find(&prices)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		
	break;
	case "equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&prices);
		
	break;
	case "not_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&prices);	
		
	break;
	case "less_than" :
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&prices);	
		
	break;
	case "greater_than":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&prices);	
		
	break;
	case "less_than_or_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&prices);	
		
	break;
	case "greater_than_ro_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&prices);	
		
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&prices);
		
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		GormDB.Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&prices);
		
		break
	case "like":
		// fmt.Println(Ser.Search_query_1)
		if Ser.Search_query_1 == "all" {
				//db.Order("name DESC")
		GormDB.Order(Ser.Column+" "+Ser.Direction).Find(&prices)
		
	
		}else {
	
			GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&prices);
			
		}
	break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&prices);
		
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return prices, nil
}