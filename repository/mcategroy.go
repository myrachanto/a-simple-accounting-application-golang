package repository

import (
	"fmt"
	"strings"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/support"
)
//Majorcategoryrepo ...
var (
	Majorcategoryrepo majorcategoryrepo = majorcategoryrepo{}
)

///curtesy to gorm
type majorcategoryrepo struct{}

func (majorcategoryRepo majorcategoryrepo) Create(majorcategory *model.Majorcategory) (*model.Majorcategory, *httperors.HttpError) {
	if err := majorcategory.Validate(); err != nil {
		return nil, err
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&majorcategory)
	IndexRepo.DbClose(GormDB)
	return majorcategory, nil
}
func (majorcategoryRepo majorcategoryrepo) GetOne(id int) (*model.Majorcategory, *httperors.HttpError) {
	ok := majorcategoryRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("majorcategory with that id does not exists!")
	}
	majorcategory := model.Majorcategory{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&majorcategory).Where("id = ?", id).First(&majorcategory)
	IndexRepo.DbClose(GormDB)
	
	return &majorcategory, nil
}

func (majorcategoryRepo majorcategoryrepo) GetAll(majorcategorys []model.Majorcategory,search *support.Search) ([]model.Majorcategory, *httperors.HttpError) {

	results, err1 := majorcategoryRepo.Search(search, majorcategorys)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}
func (majorcategoryRepo majorcategoryrepo) All() (t []model.Majorcategory, r *httperors.HttpError) {

	majorcategory := model.Majorcategory{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&majorcategory).Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (majorcategoryRepo majorcategoryrepo) Update(id int, majorcategory *model.Majorcategory) (*model.Majorcategory, *httperors.HttpError) {
	ok := majorcategoryRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("majorcategory with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	amajorcategory := model.Majorcategory{}
	
	GormDB.Model(&amajorcategory).Where("id = ?", id).First(&amajorcategory)
	if majorcategory.Name  == "" {
		majorcategory.Name = amajorcategory.Name
	}
	if majorcategory.Title  == "" {
		majorcategory.Title = amajorcategory.Title
	}
	if majorcategory.Description  == "" {
		majorcategory.Description = amajorcategory.Description
	}
	fmt.Println(majorcategory)
	GormDB.Save(&majorcategory)
	
	IndexRepo.DbClose(GormDB)

	return majorcategory, nil
}
func (majorcategoryRepo majorcategoryrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := majorcategoryRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("Product with that id does not exists!")
	}
	majorcategory := model.Majorcategory{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&majorcategory).Where("id = ?", id).First(&majorcategory)
	GormDB.Delete(majorcategory)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (majorcategoryRepo majorcategoryrepo)ProductUserExistByid(id int) bool {
	majorcategory := model.Majorcategory{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&majorcategory, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}

func (majorcategoryRepo majorcategoryrepo) Search(Ser *support.Search, majorcategorys []model.Majorcategory)([]model.Majorcategory, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	majorcategory := model.Majorcategory{}
	switch(Ser.Search_operator){
	case "all":
		GormDB.Model(&majorcategory).Order(Ser.Column+" "+Ser.Direction).Find(&majorcategorys)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
	
	break;
	case "equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&majorcategorys);
	
	break;
	case "not_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&majorcategorys);	
	
	break;
	case "less_than" :
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&majorcategorys);	
	
	break;
	case "greater_than":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&majorcategorys);	
	
	break;
	case "less_than_or_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&majorcategorys);	
	
	break;
	case "greater_than_ro_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&majorcategorys);	
	
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&majorcategorys);
	
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		GormDB.Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&majorcategorys);
	
	// break;
	case "like":
		// fmt.Println(Ser.Search_query_1)
		if Ser.Search_query_1 == "all" {
				//db.Order("name DESC")
		GormDB.Order(Ser.Column+" "+Ser.Direction).Find(&majorcategorys)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		
		}else {

			GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&majorcategorys);
		
		}
	break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&majorcategorys);
	
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return majorcategorys, nil
}