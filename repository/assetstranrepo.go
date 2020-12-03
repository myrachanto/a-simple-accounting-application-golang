package repository

import (
	"fmt"
	"strings"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/support"
)
//Asstransrepo ...
var (
	Asstransrepo asstransrepo = asstransrepo{}
)

///curtesy to gorm
type asstransrepo struct{}

func (asstransRepo asstransrepo) Create(asstrans *model.Asstrans) (*model.Asstrans, *httperors.HttpError) {

	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	// id := asstrans.Asset.ID
	// asstrans.AssetID = id
	// GormDB.Model(model.Invoice{}).Association("Transactions").Append(transact)
	GormDB.Create(&asstrans) 
	IndexRepo.DbClose(GormDB)
	return asstrans, nil
}
func (asstransRepo asstransrepo) GetOne(id int) (*model.Asstrans, *httperors.HttpError) {
	ok := asstransRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("asstrans with that id does not exists!")
	}
	asstrans := model.Asstrans{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&asstrans).Where("id = ?", id).First(&asstrans)
	IndexRepo.DbClose(GormDB)
	
	return &asstrans, nil
}

func (asstransRepo asstransrepo) GetAll(asstranss []model.Asstrans,search *support.Search) ([]model.Asstrans, *httperors.HttpError) {
	results, err1 := asstransRepo.Search(search, asstranss)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}

func (asstransRepo asstransrepo) Update(id int, asstrans *model.Asstrans) (*model.Asstrans, *httperors.HttpError) {
	ok := asstransRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("asstrans with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	// asstrans := model.Asstrans{}
	aasstrans := model.Asstrans{}
	
	GormDB.Model(&asstrans).Where("id = ?", id).First(&aasstrans)
	if asstrans.Name  == "" {
		asstrans.Name = aasstrans.Name
	}
	if asstrans.Title  == "" {
		asstrans.Title = aasstrans.Title
	}
	if asstrans.Description  == "" {
		asstrans.Description = aasstrans.Description
	}
	if asstrans.Depreciation < 0 {
		asstrans.Depreciation = aasstrans.Depreciation
	}
	GormDB.Save(&asstrans)
	
	IndexRepo.DbClose(GormDB)

	return asstrans, nil
}
func (asstransRepo asstransrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := asstransRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("Product with that id does not exists!")
	}
	asstrans := model.Asstrans{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&asstrans).Where("id = ?", id).First(&asstrans)
	GormDB.Delete(asstrans)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (asstransRepo asstransrepo)ProductUserExistByid(id int) bool {
	asstrans := model.Asstrans{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
 res := GormDB.First(&asstrans, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}

func (asstransRepo asstransrepo) Search(Ser *support.Search, asstranss []model.Asstrans)([]model.Asstrans, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	asstrans := model.Asstrans{}
	switch(Ser.Search_operator){
	case "all":
		GormDB.Model(&asstrans).Order(Ser.Column+" "+Ser.Direction).Find(&asstranss)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
	
	break;
	case "equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&asstranss);
	
	break;
	case "not_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&asstranss);	
	
	break;
	case "less_than" :
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&asstranss);	
	
	break;
	case "greater_than":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&asstranss);	
	
	break;
	case "less_than_or_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&asstranss);	
	
	break;
	case "greater_than_ro_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&asstranss);	
	
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&asstranss);
	
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		GormDB.Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&asstranss);
	
	// break;
	case "like":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&asstranss);
	
	break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&asstranss);
	
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return asstranss, nil
}