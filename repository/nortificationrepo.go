package repository

import (
	"fmt"
	"strings"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/support"
)
//Nortificationrepo ...
var (
	Nortificationrepo nortificationrepo = nortificationrepo{}
)

///curtesy to gorm
type nortificationrepo struct{}

func (nortificationRepo nortificationrepo) Create(nortification *model.Nortification) (*model.Nortification, *httperors.HttpError) {
	if err := nortification.Validate(); err != nil {
		return nil, err
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&nortification)
	IndexRepo.DbClose(GormDB)
	return nortification, nil
}
func (nortificationRepo nortificationrepo) GetOne(id int) (*model.Nortification, *httperors.HttpError) {
	ok := nortificationRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("nortification with that id does not exists!")
	}
	nortification := model.Nortification{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&nortification).Where("id = ?", id).First(&nortification)
	IndexRepo.DbClose(GormDB)
	
	return &nortification, nil
}

func (nortificationRepo nortificationrepo) GetAll(nortifications []model.Nortification,search *support.Search) ([]model.Nortification, *httperors.HttpError) {
	results, err1 := nortificationRepo.Search(search, nortifications)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}

func (nortificationRepo nortificationrepo) Update(id int, nortification *model.Nortification) (*model.Nortification, *httperors.HttpError) {
	ok := nortificationRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("nortification with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	anortification := model.Nortification{}
	
	GormDB.Model(&nortification).Where("id = ?", id).First(&anortification)
	if nortification.Name  == "" {
		nortification.Name = anortification.Name
	}
	if nortification.Title  == "" {
		nortification.Title = anortification.Title
	}
	if nortification.Description  == "" {
		nortification.Description = anortification.Description
	}
	GormDB.Save(&nortification)
	
	IndexRepo.DbClose(GormDB)

	return nortification, nil
}
func (nortificationRepo nortificationrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := nortificationRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("Product with that id does not exists!")
	}
	nortification := model.Nortification{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&nortification).Where("id = ?", id).First(&nortification)
	GormDB.Delete(nortification)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (nortificationRepo nortificationrepo)ProductUserExistByid(id int) bool {
	nortification := model.Nortification{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&nortification, "id =?", id)
	if res.Error != nil{
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}

func (nortificationRepo nortificationrepo) Search(Ser *support.Search, nortifications []model.Nortification)([]model.Nortification, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	nortification := model.Nortification{}
	switch(Ser.Search_operator){
	case "all":
		GormDB.Model(&nortification).Order(Ser.Column+" "+Ser.Direction).Find(&nortifications)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		
	break;
	case "equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&nortifications);
		
	break;
	case "not_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&nortifications);	
		
	break;
	case "less_than" :
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&nortifications);	
		
	break;
	case "greater_than":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&nortifications);	
		
	break;
	case "less_than_or_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&nortifications);	
		
	break;
	case "greater_than_ro_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&nortifications);	
		
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&nortifications);
		
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		GormDB.Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&nortifications);
		
	// break;
	case "like":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&nortifications);
		
	break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&nortifications);
		
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return nortifications, nil
}