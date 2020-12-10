package model

import (
  "gorm.io/gorm"
	"github.com/myrachanto/accounting/httperors"
)
//Asset ...
type Asset struct {
	Name string `gorm:"not null" json:"name"`
	Assetcode string `json:"assetcode"`
	Pincode string `json:"pin"`
	Description string ` json:"description"`
	Ownership string `gorm:"not null" json:"ownership"`
	Asstrans []Asstrans `gorm:"not null" json:"asstrans"`
	Price float64 `gorm:"not null" json:"price"`
	Depreciationtype string `gorm:"not null" json:"depreciationtype"`
	Depreciationrate float64 `gorm:"not null" json:"depreciationrate"`
	ExpectedUsage float64 `gorm:"not null" json:"expected usage"`
	Liscence string `gorm:"not null" json:"liscence"`
	Usercode string `json:"usercode"`
	Picture string `json:"picture"`
	gorm.Model
}
//Validate ...
func (asset Asset) Validate() *httperors.HttpError{ 
	if asset.Name == "" && len(asset.Name) < 3 {
		return httperors.NewNotFoundError("Invalid Name")
	}
	if asset.Description == "" && len(asset.Description) < 3 {
		return httperors.NewNotFoundError("Invalid description")
	}
	if asset.Liscence == "" {
		return httperors.NewNotFoundError("Invalid liscence")
	}
	if asset.Depreciationtype == "" {
		return httperors.NewNotFoundError("Invalid depreciation type")
	}
	if asset.Depreciationrate < 0 {
		return httperors.NewNotFoundError("Invalid depreciation rate")
	}
	if asset.Price < 0 {
		return httperors.NewNotFoundError("Invalid Price")
	}
	return nil
}