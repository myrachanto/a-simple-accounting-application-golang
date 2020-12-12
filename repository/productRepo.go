package repository

import (
	"strconv"
	"time"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
)
//Productrepo ...
var (
	Productrepo productrepo = productrepo{}
)

///curtesy to gorm
type productrepo struct{}

func (productRepo productrepo) Create(product *model.Product) (*model.Product, *httperors.HttpError) {
	if err := product.Validate(); err != nil {
		return nil, err
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	cat, err := Categoryrepo.GetMajorcat(product.Category)
	if err != nil {
		return nil, httperors.NewNotFoundError("category with that name does not exists!")
	}
	product.Majorcategory = cat.Majorcategory 

	code, x := Productrepo.GeneCode()
	if x != nil {
		return nil, x
	}
	product.Productcode = code
	GormDB.Create(&product)
	IndexRepo.DbClose(GormDB)
	return product, nil
}
func (productRepo productrepo) GetOne(id int,dated,searchq2,searchq3 string) (*model.ProductView, *httperors.HttpError) {
	ok := productRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("product with that id does not exists!")
	}
	product := model.Product{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&product).Where("id = ?", id).First(&product)
	IndexRepo.DbClose(GormDB)
	sold, e := Transactionrepo.ProductsSold(product.Productcode,dated,searchq2,searchq3)
	if e != nil {
		return nil, e
	}
	// credits, er := Invoicerepo.CustomerCredits(customer.Name)
	bought, er := STransactionrepo.ProductsBought(product.Productcode,dated,searchq2,searchq3)
	if er != nil {
		return nil, er
	}
	
	return &model.ProductView{
		Product:product,
		Sold:sold,
		Bought:bought,
	}, nil
}

func (productRepo productrepo) View() ([]model.Category, *httperors.HttpError) {
	mc, e := Categoryrepo.All()
	if e != nil{
		return nil, e
	}
	return mc, nil
}
func (productRepo productrepo) ViewReport(dated,searchq2,searchq3 string) (*model.ProductReport, *httperors.HttpError) {

	products, er := Productrepo.AllSearch(dated,searchq2,searchq3)
	if er != nil {
		return nil, er
	}

	sold, er1 := Transactionrepo.Allsearch(dated,searchq2,searchq3)
	if er1 != nil {
		return nil, er1
	}
	var s float64 = 0
	for _, so := range sold {
		s += so.Total
	}
	bought, er1 := STransactionrepo.Allsearch(dated,searchq2,searchq3)
	if er != nil {
		return nil, er
	}
	var b float64 = 0
	for _, bo := range bought {
		b += bo.Total
	}
	z := model.ProductReport{}
	z.Products = products
	z.Product.Name = "All Products"
	z.Product.Total = float64(len(products))
	z.Product.Description = "Total Products registered"
	//////////////////////////////////////////////////////////////
	z.Sold.Name = "Total sales"
	z.Sold.Total = s
	z.Sold.Description = "Total sales  in this search period"
	///////////////////////////////////////////////////////////////
	z.Bought.Name = "Total Purchases"
	z.Bought.Total = b
	z.Bought.Description = "Total purchase in this search period"
	
	return &z, nil
}

func (productRepo productrepo) AllSearch(dated,searchq2,searchq3 string) (results []model.Product, r *httperors.HttpError) {

	now := time.Now()
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	if dated != "custom"{ 
		if dated == "In the last 24hrs"{
			d := now.AddDate(0, 0, -1)
			GormDB.Where("updated_at > ?", d).Find(&results)
		}
		if dated == "In the last 7days"{
			d := now.AddDate(0, 0, -7)
			GormDB.Where("updated_at > ?",d).Find(&results)
		}
		if dated == "In the last 15day"{
			d := now.AddDate(0, 0, -15)
			GormDB.Where("updated_at > ?",d).Find(&results)
		}
		if dated == "In the last 30days"{
			d := now.AddDate(0, 0, -30)
			GormDB.Where("updated_at > ?",d).Find(&results)
		}
	}
	if dated == "custom"{
		start,err := time.Parse(Layout,searchq2)
		if err != nil {
			return nil, httperors.NewNotFoundError("Something went wrong parsing date1!")
		}
		end,err1 := time.Parse(Layout,searchq3)
		if err1 != nil {
			return nil, httperors.NewNotFoundError("Something went wrong parsing date1!")
		}
		GormDB.Where("updated_at BETWEEN ? AND ?",start, end).Find(&results)
	}
	IndexRepo.DbClose(GormDB)
	return results, nil

}
func (productRepo productrepo) All() (t []model.Product, r *httperors.HttpError) {

	product := model.Product{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&product).Find(&t)
	
	IndexRepo.DbClose(GormDB)
	return t, nil

}

func (productRepo productrepo)GeneCode() (string, *httperors.HttpError) {
	product := model.Product{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	err := GormDB.Last(&product)
	if err.Error != nil {
		var c1 uint = 1
		code := "ProductCode"+strconv.FormatUint(uint64(c1), 10)
		return code, nil
	 }
	c1 := product.ID + 1
	code := "ProductCode"+strconv.FormatUint(uint64(c1), 10)
	IndexRepo.DbClose(GormDB)
	return code, nil
	
}
func (productRepo productrepo) SearchProducts(search string) ( []model.Product, *httperors.HttpError) {
	products := []model.Product{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Where("name LIKE ?", "%"+ search +"%").Find(&products)
	if err1 != nil {
			return nil, err1
		}
	return products, nil
}
// func (productRepo productrepo) GetProducts(products []model.Product,search *support.Productsearch) ([]model.Product, *httperors.HttpError) {
// 	results, err1 := productRepo.SearchFront(search, products)
// 	if err1 != nil {
// 			return nil, err1
// 		}
// 	return results, nil
// }
func (productRepo productrepo) GetAll(search string, page,pagesize int) ([]model.Product, *httperors.HttpError) {
	results := []model.Product{}
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

func (productRepo productrepo) UpdateQty(id int,quantity float64, usercode string) (string, *httperors.HttpError) {
	product := model.Product{}
	ok := Productrepo.ProductUserExistByid(id)
	if !ok {
		return "", httperors.NewNotFoundError("Product with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	GormDB.Model(&product).Where("id = ?", id).Update("quantity",quantity)
	IndexRepo.DbClose(GormDB)

	return "user updated succesifully", nil
}
func (productRepo productrepo) Update(id int, product *model.Product) (*model.Product, *httperors.HttpError) {
	ok := productRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("product with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	cat, err := Categoryrepo.GetMajorcat(product.Category)
	if err != nil {
		return nil, httperors.NewNotFoundError("category with that name does not exists!")
	}
	product.Majorcategory = cat.Majorcategory 
	
	GormDB.Model(&product).Where("id = ?", id).Save(&product)
	
	
	IndexRepo.DbClose(GormDB)

	return product, nil
}
func (productRepo productrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := productRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("Product with that id does not exists!")
	}
	product := model.Product{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&product).Where("id = ?", id).First(&product)
	GormDB.Delete(product)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (productRepo productrepo)ProductUserExistByid(id int) bool {
	product := model.Product{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&product, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (productRepo productrepo)Productexist(name string) bool {
	product := model.Product{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	GormDB.Where("name = ?", name).First(&product)
	if product.Name == "" {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (productRepo productrepo)Productqty(name string) *model.Product {
	product := model.Product{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil
	}
	GormDB.Where("name = ? ", name).First(&product)
	if product.Name == "" {
	   return nil
	}
	IndexRepo.DbClose(GormDB)
	return &product
	
}

func (productRepo productrepo) GetOptions()([]model.Product, *httperors.HttpError){

	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	product := model.Product{}
	products := []model.Product{}
	GormDB.Model(&product).Find(&products)
	return products, nil
}