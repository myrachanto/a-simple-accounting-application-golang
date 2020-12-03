package repository

import (
	"fmt"
	"strings"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/support"
)
//Expencerepo ..
var (
	Expencerepo expencerepo = expencerepo{}
)

///curtesy to gorm
type expencerepo struct{}

func (expenceRepo expencerepo) Create(expence *model.Expence) (*model.Expence, *httperors.HttpError) {
	if err := expence.Validate(); err != nil {
		return nil, err
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&expence)
	IndexRepo.DbClose(GormDB)
	return expence, nil
}
func (expenceRepo expencerepo) GetOne(id int) (*model.Expence, *httperors.HttpError) {
	ok := expenceRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("expence with that id does not exists!")
	}
	expence := model.Expence{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Preload("Expencetrasans").Model(&expence).Where("id = ?", id).First(&expence)
	IndexRepo.DbClose(GormDB)
	
	return &expence, nil
}
func (expenceRepo expencerepo) All() (t []model.Expence, r *httperors.HttpError) {

	expence := model.Expence{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&expence).Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (expenceRepo expencerepo) GetAll(expences []model.Expence,search *support.Search) ([]model.Expence, *httperors.HttpError) {
	results, err1 := expenceRepo.Search(search, expences)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}

func (expenceRepo expencerepo) Update(id int, expence *model.Expence) (*model.Expence, *httperors.HttpError) {
	ok := expenceRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("expence with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	aexpence := model.Expence{}
	
	GormDB.Model(&aexpence).Where("id = ?", id).First(&aexpence)
	if expence.Name  == "" {
		expence.Name = aexpence.Name
	}
	GormDB.Save(&expence)
	
	IndexRepo.DbClose(GormDB)

	return expence, nil
}
func (expenceRepo expencerepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := expenceRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("Product with that id does not exists!")
	}
	expence := model.Expence{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&expence).Where("id = ?", id).First(&expence)
	GormDB.Delete(expence)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (expenceRepo expencerepo)ProductUserExistByid(id int) bool {
	expence := model.Expence{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&expence, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}

func (expenceRepo expencerepo) Search(Ser *support.Search, expences []model.Expence)([]model.Expence, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	expence := model.Expence{}
	switch(Ser.Search_operator){
	case "all":
		GormDB.Preload("Expencetrasans").Model(&expence).Order(Ser.Column+" "+Ser.Direction).Find(&expences)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		
	break;
	case "equal_to":
		GormDB.Preload("Expencetrasans").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&expences);
		
	break;
	case "not_equal_to":
		GormDB.Preload("Expencetrasans").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&expences);	
		
	break;
	case "less_than" :
		GormDB.Preload("Expencetrasans").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&expences);	
		
	break;
	case "greater_than":
		GormDB.Preload("Expencetrasans").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&expences);	
		
	break;
	case "less_than_or_equal_to":
		GormDB.Preload("Expencetrasans").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&expences);	
		
	break;
	case "greater_than_ro_equal_to":
		GormDB.Preload("Expencetrasans").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&expences);	
		
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		GormDB.Preload("Expencetrasans").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&expences);
		
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		GormDB.Preload("Expencetrasans").Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&expences);
		
	// break;
case "like":
	// fmt.Println(Ser.Search_query_1)
	if Ser.Search_query_1 == "all" {
			//db.Order("name DESC")
	GormDB.Order(Ser.Column+" "+Ser.Direction).Find(&expences)
	///////////////////////////////////////////////////////////////////////////////////////////////////////
	///////////////find some other paginator more effective one///////////////////////////////////////////
	
	}else {

		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&expences);
	
	}
break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		GormDB.Preload("Expencetrasans").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&expences);
		
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return expences, nil
}