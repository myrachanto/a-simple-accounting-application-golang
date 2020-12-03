package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
	"github.com/myrachanto/accounting/support"
)

var (
	 AssetService assetService = assetService{}

) 
type redirectCategroy interface{
	Create(asset *model.Asset) (*model.Asset, *httperors.HttpError)
	GetOne(id int) (*model.Asset, *httperors.HttpError)
	GetAll(assets []model.Asset,search *support.Search) ([]model.Asset, *httperors.HttpError)
	Update(id int, asset *model.Asset) (*model.Asset, *httperors.HttpError)
	Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError)
}


type assetService struct {
	
}

func (service assetService) Create(asset *model.Asset) (*model.Asset, *httperors.HttpError) {
	if err := asset.Validate(); err != nil {
		return nil, err
	}	
	asset, err1 := r.Assetrepo.Create(asset)
	if err1 != nil {
		return nil, err1
	}
	 return asset, nil

}
func (service assetService) GetOne(id int) (*model.Asset, *httperors.HttpError) {
	asset, err1 := r.Assetrepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return asset, nil
}

func (service assetService) GetAll(assets []model.Asset,search *support.Search) ([]model.Asset, *httperors.HttpError) {
	assets, err := r.Assetrepo.GetAll(assets,search)
	if err != nil {
		return nil, err
	}
	return assets, nil
}

func (service assetService) Update(id int, asset *model.Asset) (*model.Asset, *httperors.HttpError) {
	asset, err1 := r.Assetrepo.Update(id, asset)
	if err1 != nil {
		return nil, err1
	}
	
	return asset, nil
}
func (service assetService) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Assetrepo.Delete(id)
		return success, failure
}
