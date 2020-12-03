package repository

import (
	"fmt"
	"strings"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/support"
)
//Subcategoryrepo ...
var (
	Subcategoryrepo subcategoryrepo = subcategoryrepo{}
)

///curtesy to gorm
type subcategoryrepo struct{}

func (subcategoryRepo subcategoryrepo) Create(subcategory * model.Subcategory) (* model.Subcategory, *httperors.HttpError) {
	if err := subcategory.Validate(); err != nil {
		return nil, err
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&subcategory)
	IndexRepo.DbClose(GormDB)
	return subcategory, nil
}
func (subcategoryRepo subcategoryrepo) GetOne(id int) (* model.Subcategory, *httperors.HttpError) {
	ok := subcategoryRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("subcategory with that id does not exists!")
	}
	subcategory :=  model.Subcategory{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&subcategory).Where("id = ?", id).First(&subcategory)
	IndexRepo.DbClose(GormDB)
	
	return &subcategory, nil
}

func (subcategoryRepo subcategoryrepo) GetAll(subcategorys [] model.Subcategory,search *support.Search) ([] model.Subcategory, *httperors.HttpError) {
	
	results, err1 := subcategoryRepo.Search(search, subcategorys)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}

func (subcategoryRepo subcategoryrepo) Update(id int, subcategory * model.Subcategory) (* model.Subcategory, *httperors.HttpError) {
	ok := subcategoryRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("subcategory with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	asubcategory :=  model.Subcategory{}
	
	GormDB.Model(&asubcategory).Where("id = ?", id).First(&asubcategory)
	if subcategory.Name  == "" {
		subcategory.Name = asubcategory.Name
	}
	if subcategory.Title  == "" {
		subcategory.Title = asubcategory.Title
	}
	if subcategory.Description  == "" {
		subcategory.Description = asubcategory.Description
	}
	GormDB.Save(&subcategory)
	
	IndexRepo.DbClose(GormDB)

	return subcategory, nil
}
func (subcategoryRepo subcategoryrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := subcategoryRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("Product with that id does not exists!")
	}
	subcategory :=  model.Subcategory{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&subcategory).Where("id = ?", id).First(&subcategory)
	GormDB.Delete(subcategory)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (subcategoryRepo subcategoryrepo)ProductUserExistByid(id int) bool {
	subcategory :=  model.Subcategory{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&subcategory, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}

func (subcategoryRepo subcategoryrepo) Search(Ser *support.Search, subcategorys [] model.Subcategory)([] model.Subcategory, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	subcategory :=  model.Subcategory{}
	switch(Ser.Search_operator){
	case "all":
	GormDB.Model(&subcategory).Order(Ser.Column+" "+Ser.Direction).Find(&subcategorys)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
	
	break;
	case "equal_to":
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&subcategorys);
	
	break;
	case "not_equal_to":
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&subcategorys);	
	
	break;
	case "less_than" :
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&subcategorys);	
	
	break;
	case "greater_than":
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&subcategorys);	
	
	break;
	case "less_than_or_equal_to":
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&subcategorys);	
	
	break;
	case "greater_than_ro_equal_to":
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&subcategorys);	
	
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&subcategorys);
	
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
	GormDB.Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&subcategorys);
	
	// break;
	case "like":
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&subcategorys);
	
	break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&subcategorys);
	
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return subcategorys, nil
}