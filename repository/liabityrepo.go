package repository

import (
	"strconv"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
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
	code, x := Liabilityrepo.GeneCode()
	if x != nil {
		return nil, x
	}
	liability.LiaCode = code
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
func (liabilityRepo liabilityrepo) View() (*model.LiabiltyView, *httperors.HttpError) {
	users, err := Userrepo.All()
	if err != nil {
		return nil, err
	}
	code,err4 := Liabilityrepo.GeneCode()
	if err4 != nil {
		return nil, httperors.NewNotFoundError("You got an error fetching customers")
	}
	return &model.LiabiltyView{
		Code:code,
		Users:users,
	}, nil
}

func (liabilityRepo liabilityrepo) All() (t []model.Liability, r *httperors.HttpError) {

	lia := model.Liability{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&lia).Find(&t)
	return t, nil

}
func (liabilityRepo liabilityrepo) GetAll(search string) ([]model.Liability, *httperors.HttpError) {
	results := []model.Liability{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	if search == ""{
		GormDB.Find(&results)
	}
	GormDB.Where("name LIKE ?", "%"+search+"%").Or("creditor LIKE ?", "%"+search+"%").Or("approveby LIKE ?", "%"+search+"%").Or("description LIKE ?", "%"+search+"%").Find(&results)

	IndexRepo.DbClose(GormDB)
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
	GormDB.Model(&liability).Where("id = ?", id).Save(&liability)
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

func (liabilityRepo liabilityrepo)GeneCode() (string, *httperors.HttpError) {
	liability := model.Liability{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	err := GormDB.Last(&liability)
	if err.Error != nil {
		var c1 uint = 1
		code := "LiaCode"+strconv.FormatUint(uint64(c1), 10)
		return code, nil
	 }
	c1 := liability.ID + 1
	code := "LiaCode"+strconv.FormatUint(uint64(c1), 10)
	IndexRepo.DbClose(GormDB)
	return code, nil
	
}