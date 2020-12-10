package repository

import (
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
)
//Expencetrasanrepo ...
var (
	Expencetrasanrepo expencetrasanrepo = expencetrasanrepo{}
)

///curtesy to gorm
type expencetrasanrepo struct{}

func (expencetrasanRepo expencetrasanrepo) Create(expencetrasan *model.Expencetrasan) (*model.Expencetrasan, *httperors.HttpError) {
	if err := expencetrasan.Validate(); err != nil {
		return nil, err
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	expencetrasan.Mode = "invoice"
	// id := expencetrasan.Expence.ID
	// expencetrasan.ExpenceID = id 
	expencetrasan.Type = "direct"
	GormDB.Create(&expencetrasan)
	IndexRepo.DbClose(GormDB)
	return expencetrasan, nil
}
func (expencetrasanRepo expencetrasanrepo) CreateExp(expencetrasan *model.Expencetrasan) (*model.Expencetrasan, *httperors.HttpError) {
	if err := expencetrasan.Validate(); err != nil {
		return nil, err
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	expencetrasan.Mode = "other"
	// id := expencetrasan.Expence.ID
	// expencetrasan.ExpenceID = id
	GormDB.Create(&expencetrasan)
	IndexRepo.DbClose(GormDB)
	return expencetrasan, nil
}
func (expencetrasanRepo expencetrasanrepo) GetExpencesByCode(code string) (t []model.Expencetrasan, r *httperors.HttpError) {

	exps := model.Expencetrasan{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&exps).Where("code = ?", code).Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (expencetrasanRepo expencetrasanrepo) GetAll(search string, page,pagesize int) ([]model.Expencetrasan, *httperors.HttpError) {
	results := []model.Expencetrasan{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	if search == ""{
		GormDB.Find(&results)
	}
	// db.Scopes(Paginate(r)).Find(&users)
	GormDB.Scopes(Paginate(page,pagesize)).Where("name LIKE ?", "%"+search+"%").Or("title LIKE ?", "%"+search+"%").Or("description LIKE ?", "%"+search+"%").Find(&results)

	IndexRepo.DbClose(GormDB)
	return results, nil
}
func (expencetrasanRepo expencetrasanrepo) View(code string) ([]model.Expencetrasan, *httperors.HttpError) {
	mc, e := expencetrasanRepo.Getexpencestrans(code)
	if e != nil{
		return nil, e 
	}
	return mc, nil
}
func (expencetrasanRepo expencetrasanrepo) ViewExp() (*model.ExpencetransView, *httperors.HttpError) {

	expences, err := Expencerepo.All()
	if err != nil {
		return nil, err
	}
	lias, err := Liabilityrepo.All()
	if err != nil {
		return nil, err
	}
	
	assests, err := Assetrepo.All()
	if err != nil {
		return nil, err
	}
	return &model.ExpencetransView{
		Expence:expences,
		Liability:lias,
		Asset:assests,
	}, nil
}
func (expencetrasanRepo expencetrasanrepo) ViewReport() (*model.ExpencesView, *httperors.HttpError) {
	expence := model.Expencetrasan{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	expences := []model.Expencetrasan{}
	GormDB.Model(&expence).Find(&expences)
	var tes float64 = 0
	for _,te := range expences{
		tes += te.Amount
	}
	dexpences := []model.Expencetrasan{}
	GormDB.Model(&expence).Where("type = ?", "direct").Find(&dexpences)
	var dtes float64 = 0
	for _,dte := range dexpences{
		dtes += dte.Amount
	}
	idexpences := []model.Expencetrasan{}
	GormDB.Model(&expence).Where("type = ?", "indirect").Find(&idexpences)
	var idtes float64 = 0
	for _,idte := range idexpences{
		idtes += idte.Amount
	}
	other := []model.Expencetrasan{}
	GormDB.Model(&expence).Where("type = ?", "other").Find(&other)
	var o float64 = 0
	for _,ot := range other{
		o += ot.Amount
	}
	z := model.ExpencesView{}
	z.Expences = expences
	z.Totalexpences.Name = "Total expences"
	z.Totalexpences.Total = tes
	z.Totalexpences.Description = "Total expences incurred"
	//////////////////////////////////////////////////////////////
	z.Directexpences.Name = "Direct expences"
	z.Directexpences.Total = dtes
	z.Directexpences.Description = "Total Direct expences incurred"
	///////////////////////////////////////////////////////////////
	z.InDirectexpences.Name = "InDirect expences"
	z.InDirectexpences.Total = idtes
	z.InDirectexpences.Description = "Total InDirect expences incurred"
	////////////////////////////////////////////////////////////
	z.Other.Name = "Other expences"
	z.Other.Total = o
	z.Other.Description = "Total Other expences incurred"
	
	IndexRepo.DbClose(GormDB)
	return &z, nil
}
func (expencetrasanRepo expencetrasanrepo) Getexpencestransactions(code string) (t []model.Expencetrasan, e *httperors.HttpError) {
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Where("code = ?", code).Find(&t)
	IndexRepo.DbClose(GormDB)
	
	return t, nil
}
func (expencetrasanRepo expencetrasanrepo) All() (t []model.Expencetrasan, r *httperors.HttpError) {

	exptrans := model.Expencetrasan{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&exptrans).Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (expencetrasanRepo expencetrasanrepo) Getexpencestrans(code string) (t []model.Expencetrasan, e *httperors.HttpError) {

	exptrans := model.Expencetrasan{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&exptrans).Where("code = ?", code).Find(&t)
	IndexRepo.DbClose(GormDB)
	
	return t, nil
} 

func (expencetrasanRepo expencetrasanrepo) UpdateTrans(name,code string) (string, *httperors.HttpError) {
	ok := expencetrasanRepo.expenceExistByCode(code)
	if ok == false {
		return "", httperors.NewNotFoundError("Something went wrong with the Expence crediting!")
	}

	if name == "" {
		return "", httperors.NewNotFoundError("something went wrong")
	}
	exptrans := model.Expencetrasan{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	GormDB.Model(&exptrans).Where("name = ? AND code = ?", name, code).Update("mode","credit")
	
	IndexRepo.DbClose(GormDB)
	return "Credited succesifully", nil
}
func (expencetrasanRepo expencetrasanrepo) expenceExistByCode(code string) bool {
	exptr := model.Expencetrasan{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	GormDB.Where("code = ?", code).First(&exptr)
	if exptr.ID == 0 {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true

}
func (expencetrasanRepo expencetrasanrepo) GetOne(id int) (*model.Expencetrasan, *httperors.HttpError) {
	ok := expencetrasanRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("expencetrasan with that id does not exists!")
	}
	expencetrasan := model.Expencetrasan{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&expencetrasan).Where("id = ?", id).First(&expencetrasan)
	IndexRepo.DbClose(GormDB)
	
	return &expencetrasan, nil
}
func (expencetrasanRepo expencetrasanrepo) Update(id int, expencetrasan *model.Expencetrasan) (*model.Expencetrasan, *httperors.HttpError) {
	ok := expencetrasanRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("expencetrasan with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	aexpencetrasan := model.Expencetrasan{}
	
	GormDB.Model(&aexpencetrasan).Where("id = ?", id).First(&aexpencetrasan)
	if expencetrasan.Name  == "" {
		expencetrasan.Name = aexpencetrasan.Name
	}
	if expencetrasan.Amount  == 0 {
		expencetrasan.Amount = aexpencetrasan.Amount
	}
	GormDB.Save(&expencetrasan)
	IndexRepo.DbClose(GormDB)

	return expencetrasan, nil
}
func (expencetrasanRepo expencetrasanrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := expencetrasanRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("Product with that id does not exists!")
	}
	expencetrasan := model.Expencetrasan{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&expencetrasan).Where("id = ?", id).First(&expencetrasan)
	GormDB.Delete(expencetrasan)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (expencetrasanRepo expencetrasanrepo)ProductUserExistByid(id int) bool {
	expencetrasan := model.Expencetrasan{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&expencetrasan, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}