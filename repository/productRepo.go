package repository

import (
	"fmt"
	"strings"
	"strconv"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/support"
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

	code, x := Paymentrepo.GeneCode()
	if x != nil {
		return nil, x
	}
	product.Productcode = code
	GormDB.Create(&product)
	IndexRepo.DbClose(GormDB)
	return product, nil
}
func (productRepo productrepo) GetOne(id int) (*model.Product, *httperors.HttpError) {
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
	
	return &product, nil
}

func (productRepo productrepo) View() ([]model.Category, *httperors.HttpError) {
	mc, e := Categoryrepo.All()
	if e != nil{
		return nil, e
	}
	return mc, nil
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
func (productRepo productrepo) GetProducts(products []model.Product,search *support.Productsearch) ([]model.Product, *httperors.HttpError) {
	results, err1 := productRepo.SearchFront(search, products)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}
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
func (productRepo productrepo) Search(Ser *support.Search, products []model.Product)([]model.Product, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	product := model.Product{}
	switch(Ser.Search_operator){
	case "all":
		GormDB.Model(&product).Order(Ser.Column+" "+Ser.Direction).Find(&products)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		
		
	break;
	case "equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);
		
	break;
	case "not_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);	
		
	break;
	case "less_than" :
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);	
		
	break;
	case "greater_than":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);	
		
	break;
	case "less_than_or_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);	
		
	break;
	case "greater_than_ro_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);	
		
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&products);
		
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		GormDB.Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&products);
		
	// break;
	case "like":
		// fmt.Println(Ser.Search_query_1)
		if Ser.Search_query_1 == "all" {
				//db.Order("name DESC")
		GormDB.Order(Ser.Column+" "+Ser.Direction).Find(&products)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
	

		}else {

			GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&products);
			
		}
	break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&products);
		
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return products, nil
}

func (productRepo productrepo) SearchFront(Ser *support.Productsearch, products []model.Product)([]model.Product, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	product := model.Product{}
	switch(Ser.Search_operator){
	case "all":
		GormDB.Model(&product).Order(Ser.Column+" "+Ser.Direction).Find(&products)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		
	break;
	case "equal_to":
		///product = %as% AND price (>/=/<=/>=/between)
		//db.Where("name = ? AND age >= ? ", "myrachanto", "28").Find(&users)
		//db.Where("name LIKE ?", "%a%").Find(&users)
		GormDB.Where(Ser.Column+" LIKE ? AND " +Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Name+"%", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);
		
	break;
	case "not_equal_to":
		GormDB.Where(Ser.Column+" LIKE ? AND " +Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Name+"%", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);	
		
	break;
	case "less_than" :
		fmt.Println(Ser)
		GormDB.Where(Ser.Column+" LIKE ? AND " +Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Name+"%", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);	
		
	break;
	case "greater_than":
		GormDB.Where(Ser.Column+" LIKE ? AND " +Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Name+"%", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);	
		
	break;
	case "less_than_or_equal_to":
		GormDB.Where(Ser.Column+" LIKE ? AND " +Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Name+"%", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);	
		
	break;
	case "greater_than_ro_equal_to":
		GormDB.Where(Ser.Column+" LIKE ? AND " +Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Name+"%", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);	
		
	break;
	case "like":
		// fmt.Println(Ser.Search_query_1)
		if Ser.Search_query_1 == "all" {
				//db.Order("name DESC")
		GormDB.Order(Ser.Column+" "+Ser.Direction).Find(&products)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		

		}else {

			GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&products);
		
		}
	break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return products, nil
}