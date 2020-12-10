package repository

import (
	"strconv"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
)
//Assetrepo...
var ( 
	Assetrepo assetrepo = assetrepo{}
)
///curtesy to gorm
type assetrepo struct{}

func (assetRepo assetrepo) Create(asset *model.Asset) (*model.Asset, *httperors.HttpError) {
	if err := asset.Validate(); err != nil {
		return nil, err
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	code, x := assetRepo.GeneCode()
	if x != nil {
		return nil, x
	}
	asset.Assetcode = code
	GormDB.Create(&asset)
	IndexRepo.DbClose(GormDB)
	return asset, nil
}
func (assetRepo assetrepo) All() (t []model.Asset, r *httperors.HttpError) {

	asset := model.Asset{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&asset).Find(&t)
	return t, nil

}
func (assetRepo assetrepo) View() (string, *httperors.HttpError) {
	
	code,err4 := Assetrepo.GeneCode()
	if err4 != nil {
		return "", httperors.NewNotFoundError("You got an error fetching customers")
	}
	return code, nil
}
func (assetRepo assetrepo) GetOne(id int) (*model.Asset, *httperors.HttpError) {
	ok := assetRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("asset with that id does not exists!")
	}
	asset := model.Asset{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Preload("Asstranss").Model(&asset).Where("id = ?", id).First(&asset)
	IndexRepo.DbClose(GormDB)
	
	return &asset, nil
}

func (assetRepo assetrepo) GetAll(search string, page,pagesize int) ([]model.Asset, *httperors.HttpError) {
	results := []model.Asset{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	if search == ""{
		GormDB.Find(&results)
	}
	// db.Scopes(Paginate(r)).Find(&users)
	GormDB.Scopes(Paginate(page,pagesize)).Where("name LIKE ?", "%"+search+"%").Or("description LIKE ?", "%"+search+"%").Or("ownership LIKE ?", "%"+search+"%").Find(&results)

	IndexRepo.DbClose(GormDB)
	return results, nil
}

func (assetRepo assetrepo) Update(id int, asset *model.Asset) (*model.Asset, *httperors.HttpError) {
	ok := assetRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("asset with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	// asset := model.Asset{}
	aasset := model.Asset{}
	
	GormDB.Model(&asset).Where("id = ?", id).First(&aasset)
	if asset.Name  == "" {
		asset.Name = aasset.Name
	}
	if asset.Liscence  == "" {
		asset.Liscence = aasset.Liscence
	}
	if asset.Ownership  == "" {
		asset.Ownership = aasset.Ownership
	}
	if asset.Depreciationtype  == "" {
		asset.Depreciationtype = aasset.Depreciationtype
	}
	if asset.Depreciationrate  < 0 {
		asset.Depreciationrate = aasset.Depreciationrate
	}
	if asset.Price  < 0 {
		asset.Price = aasset.Price
	}
	if asset.ExpectedUsage  < 0 {
		asset.ExpectedUsage = aasset.ExpectedUsage
	}
	if asset.Description  == "" {
		asset.Description = aasset.Description
	}
	GormDB.Save(&asset)
	IndexRepo.DbClose(GormDB)

	return asset, nil
}
func (assetRepo assetrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := assetRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("Product with that id does not exists!")
	}
	asset := model.Asset{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&asset).Where("id = ?", id).First(&asset)
	GormDB.Delete(asset)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (assetRepo assetrepo)ProductUserExistByid(id int) bool {
	asset := model.Asset{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&asset, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (assetRepo assetrepo)GeneCode() (string, *httperors.HttpError) {
	asset := model.Asset{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	err := GormDB.Last(&asset)
	if err.Error != nil {
		var c1 uint = 1
		code := "AssetCode"+strconv.FormatUint(uint64(c1), 10)
		return code, nil
	 }
	c1 := asset.ID + 1
	code := "AssetCode"+strconv.FormatUint(uint64(c1), 10)
	IndexRepo.DbClose(GormDB)
	return code, nil
	
}