package repository

import (
	"fmt"
	"strings"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/support"
)
 //Liabilityrepo ...
var (
	Liabilityrepo liabilityrepo = liabilityrepo{}
)

///curtesy to gorm
type liabilityrepo struct{}

func (liabilityRepo liabilityrepo) Create(liability *model.Liability) (*model.Liability, *httperors.HttpError) {
	if err := liability.Validate(); err != nil {
		return nil, err
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&liability)
	IndexRepo.DbClose(GormDB)
	return liability, nil
}
func (liabilityRepo liabilityrepo) GetOne(id int) (*model.Liability, *httperors.HttpError) {
	ok := liabilityRepo.liabilityUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("liability with that id does not exists!")
	}
	liability := model.Liability{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Preload("Liatrans").Model(&liability).Where("id = ?", id).First(&liability)
	IndexRepo.DbClose(GormDB)
	
	return &liability, nil
}

func (liabilityRepo liabilityrepo) GetAll(liabilitys []model.Liability,search *support.Search) ([]model.Liability, *httperors.HttpError) {

	results, err1 := liabilityRepo.Search(search, liabilitys)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}

func (liabilityRepo liabilityrepo) Update(id int, liability *model.Liability) (*model.Liability, *httperors.HttpError) {
	ok := liabilityRepo.liabilityUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("liability with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	aliability := model.Liability{}
	
	GormDB.Model(&aliability).Where("id = ?", id).First(&aliability)
	if liability.Name  == "" {
		liability.Name = aliability.Name
	}
	if liability.Approvedby  == "" {
		liability.Approvedby = aliability.Approvedby
	}
	if liability.Description  == "" {
		liability.Description = aliability.Description
	}
	if liability.Creditor  == "" {
		liability.Creditor = aliability.Creditor
	}
	if liability.Interestrate  < 0 {
		liability.Interestrate = aliability.Interestrate
	}
	if liability.Monthlypayment  < 0 {
		liability.Monthlypayment = aliability.Monthlypayment
	}
	if liability.Interestrate  < 0 {
		liability.Interestrate = aliability.Interestrate
	}
	GormDB.Save(&liability)
	IndexRepo.DbClose(GormDB)

	return liability, nil
}
func (liabilityRepo liabilityrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := liabilityRepo.liabilityUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("liability with that id does not exists!")
	}
	liability := model.Liability{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	} 
	GormDB.Model(&liability).Where("id = ?", id).First(&liability)
	GormDB.Delete(liability)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (liabilityRepo liabilityrepo)liabilityUserExistByid(id int) bool {
	liability := model.Liability{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&liability, "id =?", id)
	if res.Error != nil{
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}

func (liabilityRepo liabilityrepo) Search(Ser *support.Search, liabilitys []model.Liability)([]model.Liability, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	liability := model.Liability{}
	switch(Ser.Search_operator){
	case "all":
		//db.Order("name DESC")
		GormDB.Preload("Liatrans").Model(&liability).Order(Ser.Column+" "+Ser.Direction).Find(&liabilitys)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
	
	break;
	case "equal_to":
		GormDB.Preload("Liatrans").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&liabilitys);
	
	break;
	case "not_equal_to":
		GormDB.Preload("Liatrans").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&liabilitys);	
	
	break;
	case "less_than" :
		GormDB.Preload("Liatrans").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&liabilitys);	
	
	break;
	case "greater_than":
		GormDB.Preload("Liatrans").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&liabilitys);	
	
	break;
	case "less_than_or_equal_to":
		GormDB.Preload("Liatrans").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&liabilitys);	
	
	break;
	case "greater_than_ro_equal_to":
		GormDB.Preload("Liatrans").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&liabilitys);	
	
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		GormDB.Preload("Liatrans").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&liabilitys);
	
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		GormDB.Preload("Liatrans").Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&liabilitys);
	
	// break;
	case "like":
		GormDB.Preload("Liatrans").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&liabilitys);
	
	break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		GormDB.Preload("Liatrans").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&liabilitys);
	
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return liabilitys, nil
}
////////////subject to futher scrutiny/////////////////////////////////
// func (liabilityRepo liabilityrepo)paginator(q *gorm.DB, Ser *support.Search, liabilitys []model.Liability) ([]model.Liability, *httperors.HttpError) {
// 	p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
// 	p.SetPage(Ser.Page)
// 	// fmt.Println(Ser.Per_page)
// 	err3 := p.Results(&liabilitys)
// 	if err3 != nil {
// 		return nil, httperors.NewNotFoundError("something went wrong paginating!")
// 	}
// 	return liabilitys, nil
	
// }