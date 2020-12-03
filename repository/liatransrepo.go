package repository

import (
	"fmt"
	"strings"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/support"
)
//Liatranrepo ...
var (
	Liatranrepo liatranrepo = liatranrepo{}
)

///curtesy to gorm
type liatranrepo struct{}

func (liatranRepo liatranrepo) Create(liatran *model.Liatran) (*model.Liatran, *httperors.HttpError) {

	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	// id := liatran.Liability.ID
	// liatran.LiabilityID = id
	GormDB.Create(&liatran)
	IndexRepo.DbClose(GormDB)
	return liatran, nil
}
func (liatranRepo liatranrepo) GetOne(id int) (*model.Liatran, *httperors.HttpError) {
	ok := liatranRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("liatran with that id does not exists!")
	}
	liatran := model.Liatran{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&liatran).Where("id = ?", id).First(&liatran)
	IndexRepo.DbClose(GormDB)
	
	return &liatran, nil
}

func (liatranRepo liatranrepo) GetAll(liatrans []model.Liatran,search *support.Search) ([]model.Liatran, *httperors.HttpError) {
	results, err1 := liatranRepo.Search(search, liatrans)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}

func (liatranRepo liatranrepo) Update(id int, liatran *model.Liatran) (*model.Liatran, *httperors.HttpError) {
	ok := liatranRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("liatran with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	// liatran := model.Liatran{}
	aliatran := model.Liatran{}
	
	GormDB.Model(&liatran).Where("id = ?", id).First(&aliatran)
	if liatran.Name  == "" {
		liatran.Name = aliatran.Name
	}
	if liatran.Title  == "" {
		liatran.Title = aliatran.Title
	}
	if liatran.Description  == "" {
		liatran.Description = aliatran.Description
	}
	if liatran.Interest < 0 {
		liatran.Interest = aliatran.Interest
	}
	GormDB.Save(&liatran)
	
	IndexRepo.DbClose(GormDB)

	return liatran, nil
}
func (liatranRepo liatranrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := liatranRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("Product with that id does not exists!")
	}
	liatran := model.Liatran{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&liatran).Where("id = ?", id).First(&liatran)
	GormDB.Delete(liatran)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (liatranRepo liatranrepo)ProductUserExistByid(id int) bool {
	liatran := model.Liatran{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&liatran, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}

func (liatranRepo liatranrepo) Search(Ser *support.Search, liatrans []model.Liatran)([]model.Liatran, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	liatran := model.Liatran{}
	switch(Ser.Search_operator){
	case "all":
		GormDB.Model(&liatran).Order(Ser.Column+" "+Ser.Direction).Find(&liatrans)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
	
	break;
	case "equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&liatrans);
	
	break;
	case "not_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&liatrans);	
	
	break;
	case "less_than" :
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&liatrans);	
	
	break;
	case "greater_than":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&liatrans);	
	
	break;
	case "less_than_or_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&liatrans);	
	
	break;
	case "greater_than_ro_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&liatrans);	
	
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&liatrans);
	
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		GormDB.Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&liatrans);
	
	// break;
	case "like":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&liatrans);
	
	break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&liatrans);
	
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return liatrans, nil
}