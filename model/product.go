package model

import (
  "gorm.io/gorm"
	"github.com/myrachanto/accounting/httperors"
)

//Product structure
type Product struct {
	Name string `gorm:"not null" json:"name"`
	Title string `gorm:"not null" json:"title"`
	Description string `gorm:"not null" json:"description"`
	CategoryID uint
	Category string ` json:"category"`
	Majorcategory string ` json:"majorcategory"`
	Picture string `json:"picture"`
	BPrice float64 `json:"bprice"`
	SPrice float64 `json:"sprice"`
	Usercode string `json:"usercode"`
	Productcode string `json:"productcode"`
	Quantity float64 `json:"quantity"`
	Cart  []*Cart `gorm:"many2many:product_carts;"`
	Scart  []*Scart `gorm:"many2many:product_carts;"`
	Product  []*Product `gorm:"many2many:cart_products;"`
	Transaction []*Transaction `gorm:"many2many:product_transactions;"`
	STransaction []*STransaction `gorm:"many2many:product_transactions;"`
	gorm.Model
}
//Info well reltion view
type Info struct {
	Category []Category `gorm:"foreignKey:ategories; not null"`
	Majorcategory []Majorcategory `gorm:"foreignKey:Majorcategoryies; not null"`
	Subcategory []Subcategory `gorm:"foreignKey:Subcategories; not null"`
}
//ProductView ...
 type ProductView struct {
	 Product Product `json:"product"`
	 Sold []Transaction `json:"sold"`
	 Bought []STransaction `json:"bought"`
 }
 //ProductReport ..
 type ProductReport struct {
	Products []Product `json:"products"`
	Product SalesModule `json:"product"`
	Sold SalesModule `json:"sold"`
	Bought SalesModule `json:"bought"`
}
//Options stuff required to create a product...
type Options struct {
	Price []Price `gorm:"foreignKey:prices; not null"`
	Tax []Tax `gorm:"foreignKey:taxs; not null"`
	Discount []Discount `gorm:"foreignKey:Discounts; not null"`
}
//Validate ...
func (product Product) Validate() *httperors.HttpError{ 
	if product.Name == "" && len(product.Name) > 3 {
		return httperors.NewNotFoundError("Invalid Name")
	}
	if product.Title == "" && len(product.Title) > 3 {
		return httperors.NewNotFoundError("Invalid Title")
	}
	
	if product.Description == "" && len(product.Description) > 10 {
		return httperors.NewNotFoundError("Invalid description")
	}
	return nil
}