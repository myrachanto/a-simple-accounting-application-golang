package repository

import (
	"fmt"
	"log"
	"os"
	"time"
	"strconv"
	"github.com/joho/godotenv"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
)
//Supplierrepo ...
var (
	Supplierrepo supplierrepo = supplierrepo{}
)

///curtesy to gorm
type supplierrepo struct{}

func (supplierRepo supplierrepo) Create(supplier *model.Supplier) (string, *httperors.HttpError) {
	if err := supplier.Validate(); err != nil {
		return "", err
	}
	ok, err1 := supplier.ValidatePassword(supplier.Password)
	if !ok {
		return "", err1
	}
	ok = supplier.ValidateEmail(supplier.Email)
	if !ok {
		return "", httperors.NewNotFoundError("Your email format is wrong!")
	}
	ok = supplierRepo.supplierExist(supplier.Email)
	if ok {
		return "", httperors.NewNotFoundError("Your email already exists!")
	}
	hashpassword, err2 := supplier.HashPassword(supplier.Password)
	if err2 != nil {
		return "", err2
	}
	supplier.Password = hashpassword
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	
	code, x := supplierRepo.GeneCode()
	if x != nil {
		return "", x
	}
	supplier.Suppliercode = code
	GormDB.Create(&supplier)
	IndexRepo.DbClose(GormDB)
	return "supplier created successifully", nil
}
func (supplierRepo supplierrepo) Login(asupplier *model.Loginsupplier) (*model.SupplierAuth, *httperors.HttpError) {
	if err := asupplier.Validate(); err != nil {
		return nil, err
	}
	ok := supplierRepo.supplierExist(asupplier.Email)
	if !ok {
		return nil, httperors.NewNotFoundError("Your email does not exists!")
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	supplier := model.Supplier{}
	GormDB.Model(&supplier).Where("email = ?", asupplier.Email).First(&supplier)
	ok = supplier.Compare(asupplier.Password, supplier.Password)
	if !ok {
		return nil, httperors.NewNotFoundError("wrong email password combo!")
	}
	tk := &model.SupplierToken{
		SupplierID: supplier.ID,
		Name: supplier.Name, 
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: model.ExpiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading key")
	}
	encyKey := os.Getenv("EncryptionKey")
	tokenString, error := token.SignedString([]byte(encyKey))
	if error != nil {
		fmt.Println(error)
	}
	// messages ,e := supplierRepo.UnreadMessages(supplier.ID)
	// if e != nil {
	// 	return nil, e
	// }
	// norti ,e := supplierRepo.UnreadNortifications(supplier.ID)
	// if e != nil {
	// 	return nil, e
	// }
	auth := &model.SupplierAuth{SupplierID:supplier.ID, Name:supplier.Name, Token:tokenString}
	GormDB.Create(&auth)
	IndexRepo.DbClose(GormDB)
	
	return auth, nil
}
func (supplierRepo supplierrepo) Logout(token string) (*httperors.HttpError) {
	auth := model.CustomnerAuth{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return err1
	}
	res := GormDB.First(&auth, "token =?", token)
	if res.Error != nil {
		return httperors.NewNotFoundError("Something went wrong logging out!")
	 }
	
	GormDB.Model(&auth).Where("token =?", token).First(&auth)
	
	GormDB.Delete(auth)
	IndexRepo.DbClose(GormDB)
	
	return  nil
}
func (supplierRepo supplierrepo) Forgot(email string) (string, *httperors.HttpError) {
	ok := supplierRepo.supplierExist(email)
	if !ok {
		return "", httperors.NewNotFoundError("That Email does not exists with our records!")
	}
	
	return "Email sent!", nil
}
func (supplierRepo supplierrepo) GetOne(id int,dated,searchq2,searchq3 string) (*model.Supplierdetails, *httperors.HttpError) {
	ok := supplierRepo.SupplierExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("supplier with that id does not exists!")
	}   

	supplier := Supplierrepo.Getsupplierbyid(id)
	// invoices, e := SInvoicerepo.SuppliersInvoice(supplier.Name)
	invoices, e := SInvoicerepo.SuppliersInvoicebycode(supplier.Suppliercode,dated,searchq2,searchq3)
	if e != nil { 
		return nil, e
	}
	// credits, er := SInvoicerepo.SupplierCredits(supplier.Name)
	credits, er := SInvoicerepo.SupplierCreditsbycode(supplier.Suppliercode,dated,searchq2,searchq3)
	if er != nil {
		return nil, er
	}
	payments, er3 := Paymentrepo.AllSearchSupplier(supplier.Suppliercode,dated,searchq2,searchq3)
	if er3 != nil {
		return nil, er3
	}
	return &model.Supplierdetails{
		Supplier: supplier,
		SInvoices: invoices,
		Grns: credits,
		Payment:payments,
	}, nil
}
func (supplierRepo supplierrepo)GeneCode() (string, *httperors.HttpError) {
	supplier := model.Supplier{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	err := GormDB.Last(&supplier)
	if err.Error != nil {
		var c1 uint = 1
		code := "SupplierCode"+strconv.FormatUint(uint64(c1), 10)
		return code, nil
	 }
	c1 := supplier.ID + 1
	code := "SupplierCode"+strconv.FormatUint(uint64(c1), 10)
	IndexRepo.DbClose(GormDB)
	return code, nil
	
}
func (supplierRepo supplierrepo) GetOptions()([]model.Supplier, *httperors.HttpError){

	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	supplier := model.Supplier{}
	suppliers := []model.Supplier{}
	GormDB.Model(&supplier).Find(&suppliers)
	return suppliers, nil
}
func (supplierRepo supplierrepo) GetAll(search string, page,pagesize int) ([]model.Supplier, *httperors.HttpError) {
	results := []model.Supplier{}
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
func (supplierRepo supplierrepo) All() (t []model.Supplier, r *httperors.HttpError) {

	supplier := model.Supplier{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&supplier).Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (supplierRepo supplierrepo) AllDebts(dated,searchq2,searchq3 string) (results []model.CreditTransaction, r *httperors.HttpError) {

	now := time.Now()
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}

	if dated != "custom"{
		if dated == "In the last 24hrs"{
			d := now.AddDate(0, 0, -1)
			GormDB.Where("updated_at > ? AND status = ?", d,"notpaid").Find(&results)
		}
		if dated == "In the last 7days"{
			d := now.AddDate(0, 0, -7)
			GormDB.Where("updated_at > ? AND status = ?", d,"notpaid").Find(&results)
		}
		if dated == "In the last 15day"{
			d := now.AddDate(0, 0, -15)
			GormDB.Where("updated_at > ? AND status = ?", d,"notpaid").Find(&results)
		}
		if dated == "In the last 30days"{
			d := now.AddDate(0, 0, -30)
			GormDB.Where("updated_at > ? AND status = ?", d,"notpaid").Find(&results)
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
		GormDB.Where("status = ? AND updated_at BETWEEN ? AND ?","notpaid", start, end).Find(&results)
	}
	IndexRepo.DbClose(GormDB)
	return results, nil

}
func (supplierRepo supplierrepo)SupplierExistbycode(code string) bool {
	supplier := model.Supplier{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	GormDB.Where("suppliercode = ? ", code).First(&supplier)
	if supplier.ID == 0 {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}

func (customerRepo customerrepo)GetsupplierwithCode(code string) *model.Supplier {
	supplier := model.Supplier{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil
	}
	GormDB.Where("suppliercode = ? ", code).First(&supplier)
	if supplier.Name == "" {
	   return nil
	}
	IndexRepo.DbClose(GormDB)
	return &supplier
	
}
// func (supplierRepo supplierrepo) GetAll(search *support.Search) ([]interface{}, *httperors.HttpError) {
// 	supplier := model.Supplier{}
// 	// suppliers := []model.Supplier{}
// 	// results, err1 := supplierRepo.Search(search, supplier)
// 	 results, err1 := support.SearchQuery(search, supplier)
// 	if err1 != nil {
// 			return nil, err1
// 		}
// 	return results, nil 
// }

func (supplierRepo supplierrepo) Update(id int, supplier *model.Supplier) (*model.Supplier, *httperors.HttpError) {
	ok := supplierRepo.SupplierExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("supplier with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}

	GormDB.Model(&supplier).Where("id = ?", id).Save(&supplier)
	
	IndexRepo.DbClose(GormDB)

	return supplier, nil
}
func (supplierRepo supplierrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := supplierRepo.SupplierExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("supplier with that id does not exists!")
	}
	supplier := model.Supplier{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	} 
	GormDB.Model(&supplier).Where("id = ?", id).First(&supplier)
	GormDB.Delete(supplier)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (supplierRepo supplierrepo)SupplierExistByname(name string) bool {
	supplier := model.Supplier{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	GormDB.Where("name = ? ", name).First(&supplier)
	if supplier.Name == "" {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (supplierRepo supplierrepo)Getsupplier(name string) *model.Supplier {
	supplier := model.Supplier{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil
	}
	GormDB.Where("name = ? ", name).First(&supplier) 
	if supplier.Name == "" {
	   return nil
	}
	IndexRepo.DbClose(GormDB)
	return &supplier
	
}
func (supplierRepo supplierrepo)Getsupplierwithcode(code string) *model.Supplier {
	supplier := model.Supplier{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil
	}
	GormDB.Where("suppliercode = ? ", code).First(&supplier)
	if supplier.Name == "" {
	   return nil
	}
	IndexRepo.DbClose(GormDB)
	return &supplier
	
}
func (supplierRepo supplierrepo)Getsupplierbyid(id int) *model.Supplier {
	supplier := model.Supplier{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil
	}
	GormDB.Where("id = ? ", id).First(&supplier)
	if supplier.Name == "" {
	   return nil
	}
	IndexRepo.DbClose(GormDB)
	return &supplier
	
}
func (supplierRepo supplierrepo) AllSearch(dated,searchq2,searchq3 string) (results []model.Supplier, r *httperors.HttpError) {

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
func (supplierRepo supplierrepo) ViewReport(dated,searchq2,searchq3 string) (*model.SupplierView, *httperors.HttpError) {
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	suppliers, er := Supplierrepo.AllSearch(dated,searchq2,searchq3)
	if er != nil {
		return nil, er
	}
	supplier := model.Supplier{}
	lastsuppliers := []model.Supplier{}
	now := time.Now()
	lastWeek := now.AddDate(0, 0, -7)
	today := now.AddDate(0, 0, -1)
	GormDB.Model(&supplier).Where("updated_at > ?", lastWeek).Find(&lastsuppliers)
	
	todaysuppliers := []model.Supplier{}
	GormDB.Model(&supplier).Where("updated_at > ?", today).Find(&todaysuppliers)
	
	z := model.SupplierView{}
	z.Suppliers = suppliers
	z.AllSuppliers.Name = "All suppliers"
	z.AllSuppliers.Total = float64(len(suppliers))
	z.AllSuppliers.Description = "Total suppliers registered"
	//////////////////////////////////////////////////////////////
	z.Lastweek.Name = "Last 7 days suppliers"
	z.Lastweek.Total = float64(len(lastsuppliers))
	z.Lastweek.Description = "Total suppliers registered for the last seven days"
	///////////////////////////////////////////////////////////////
	z.Todays.Name = "Todays suppliers"
	z.Todays.Total = float64(len(todaysuppliers))
	z.Todays.Description = "Total suppliers registered Today"
	
	IndexRepo.DbClose(GormDB)
	return &z, nil
}
func (supplierRepo supplierrepo)supplierExist(email string) bool {
	supplier := model.Supplier{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	GormDB.Where("name = ? ", email).First(&supplier)
	if supplier.Name == "" {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (supplierRepo supplierrepo)SupplierExistByid(id int) bool {
	supplier := model.Supplier{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&supplier, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}