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
//Customerrepo ...
var (
	Customerrepo customerrepo = customerrepo{}
)

///curtesy to gorm
type customerrepo struct{}

func (customerRepo customerrepo) Create(customer *model.Customer) (string, *httperors.HttpError) {
	if err := customer.Validate(); err != nil {
		return "", err
	}
	ok, err1 := customer.ValidatePassword(customer.Password)
	if !ok {
		return "", err1
	}
	ok = customer.ValidateEmail(customer.Email)
	if !ok {
		return "", httperors.NewNotFoundError("Your email format is wrong!")
	}
	ok = customerRepo.customerExist(customer.Email)
	if ok {
		return "", httperors.NewNotFoundError("Your email already exists!")
	}
	hashpassword, err2 := customer.HashPassword(customer.Password)
	if err2 != nil {
		return "", err2
	}
	customer.Password = hashpassword
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	code, x := customerRepo.GeneCode()
	if x != nil {
		return "", x
	}
	customer.Customercode = code
	GormDB.Create(&customer)
	IndexRepo.DbClose(GormDB)
	return "customer created successifully", nil
}
func (customerRepo customerrepo) Login(acustomer *model.Logincustomer) (*model.CustomnerAuth, *httperors.HttpError) {
	if err := acustomer.Validate(); err != nil {
		return nil, err
	}
	ok := customerRepo.customerExist(acustomer.Email)
	if !ok {
		return nil, httperors.NewNotFoundError("Your email does not exists!")
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	customer := model.Customer{}
	GormDB.Model(&customer).Where("email = ?", acustomer.Email).First(&customer)
	ok = customer.Compare(acustomer.Password, customer.Password)
	if !ok {
		return nil, httperors.NewNotFoundError("wrong email password combo!")
	}
	tk := &model.CustomerToken{
		CustomerID: customer.ID,
		Name: customer.Name,
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
	// messages ,e := customerRepo.UnreadMessages(customer.ID)
	// if e != nil {
	// 	return nil, e
	// }
	// norti ,e := customerRepo.UnreadNortifications(customer.ID)
	// if e != nil {
	// 	return nil, e
	// }
	auth := &model.CustomnerAuth{CustomerID:customer.ID, Name:customer.Name, Token:tokenString}
	GormDB.Create(&auth)
	IndexRepo.DbClose(GormDB)
	
	return auth, nil
}
func (customerRepo customerrepo) Logout(token string) (*httperors.HttpError) {
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
func (customerRepo customerrepo) Forgot(email string) (string, *httperors.HttpError) {
	ok := customerRepo.customerExist(email)
	if !ok {
		return "", httperors.NewNotFoundError("That Email does not exists with our records!")
	}
	
	return "Email sent!", nil
}
func (customerRepo customerrepo) GetOne(id int,dated,searchq2,searchq3 string) (*model.Customerdetails, *httperors.HttpError) {
	ok := customerRepo.customerExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("customer with that id does not exists!")
	} 
	fmt.Println(dated,searchq2,searchq3)
	customer := model.Customer{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1 
	}
	GormDB.Model(&customer).Where("id = ?", id).First(&customer)
	IndexRepo.DbClose(GormDB)
	// invoices, e := Invoicerepo.Customerinvoice(customer.Name)
	invoices, e := Invoicerepo.Customerinvoicebycode(customer.Customercode,dated,searchq2,searchq3)
	if e != nil {
		return nil, e
	}
	// credits, er := Invoicerepo.CustomerCredits(customer.Name)
	credits, er := Invoicerepo.CustomerCreditsbycode(customer.Customercode,dated,searchq2,searchq3)
	if er != nil {
		return nil, er
	}
	return &model.Customerdetails{
		Customer: customer,
		Invoices: invoices,
		Credits: credits,
	}, nil
}
func (customerRepo customerrepo) GetOptions()([]model.Customer, *httperors.HttpError){

	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	customer := model.Customer{}
	customers := []model.Customer{}
	GormDB.Model(&customer).Find(&customers)
	return customers, nil
}
func (customerRepo customerrepo) GetAll(search string, page,pagesize int) ([]model.Customer, *httperors.HttpError) {
	results := []model.Customer{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	if search == ""{
		GormDB.Find(&results)
	}
	// db.Scopes(Paginate(r)).Find(&users)
	GormDB.Scopes(Paginate(page,pagesize)).Where("name LIKE ?", "%"+search+"%").Or("email LIKE ?", "%"+search+"%").Or("company LIKE ?", "%"+search+"%").Find(&results)

	IndexRepo.DbClose(GormDB)
	return results, nil
}
func (customerRepo customerrepo) All() (t []model.Customer, r *httperors.HttpError) {

	customer := model.Customer{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&customer).Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (customerRepo customerrepo) AllSearch(dated,searchq2,searchq3 string) (results []model.Customer, r *httperors.HttpError) {

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
func (customerRepo customerrepo)GeneCode() (string, *httperors.HttpError) {
	customer := model.Customer{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	err := GormDB.Last(&customer)
	if err.Error != nil {
		var c1 uint = 1
		code := "CustomerCode"+strconv.FormatUint(uint64(c1), 10)
		return code, nil
	 }
	c1 := customer.ID + 1
	code := "CustomerCode"+strconv.FormatUint(uint64(c1), 10)
	IndexRepo.DbClose(GormDB)
	return code, nil
	
}
func (customerRepo customerrepo) AllDebts(dated,searchq2,searchq3 string) (results []model.DebtTransaction, r *httperors.HttpError) {

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
// func (customerRepo customerrepo) GetAll(search *support.Search) ([]interface{}, *httperors.HttpError) {
// 	customer := model.Customer{}
// 	// customers := []model.Customer{}
// 	// results, err1 := customerRepo.Search(search, customer)
// 	 results, err1 := support.SearchQuery(search, customer)
// 	if err1 != nil {
// 			return nil, err1
// 		}
// 	return results, nil 
// }

func (customerRepo customerrepo) Update(id int, customer *model.Customer) (*model.Customer, *httperors.HttpError) {
	ok := customerRepo.customerExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("customer with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&customer).Where("id = ?", id).Save(&customer)
	
	IndexRepo.DbClose(GormDB)

	return customer, nil
}
func (customerRepo customerrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := customerRepo.customerExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("customer with that id does not exists!")
	}
	customer := model.Customer{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	} 
	GormDB.Model(&customer).Where("id = ?", id).First(&customer)
	GormDB.Delete(customer)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (customerRepo customerrepo)customerExistByname(name string) bool {
	customer := model.Customer{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	GormDB.Where("name = ? ", name).First(&customer)
	if customer.Name == "" {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (customerRepo customerrepo)Getcustomer(name string) *model.Customer {
	customer := model.Customer{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil
	}
	GormDB.Where("name = ? ", name).First(&customer)
	if customer.Name == "" {
	   return nil
	}
	IndexRepo.DbClose(GormDB)
	return &customer
	
}

func (customerRepo customerrepo)GetcustomerwithCode(code string) *model.Customer {
	customer := model.Customer{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil
	}
	GormDB.Where("customercode = ? ", code).First(&customer)
	if customer.Name == "" {
	   return nil
	}
	IndexRepo.DbClose(GormDB)
	return &customer
	
}
func (customerRepo customerrepo) ViewReport(dated,searchq2,searchq3 string) (*model.CustomerView, *httperors.HttpError) {
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	custs, er := Customerrepo.AllSearch(dated,searchq2,searchq3)
	if er != nil {
		return nil, er
	}
	customer := model.Customer{}
	lastcustomers := []model.Customer{}
	now := time.Now()
	lastWeek := now.AddDate(0, 0, -7)
	today := now.AddDate(0, 0, -1)
	GormDB.Model(&customer).Where("updated_at > ?", lastWeek).Find(&lastcustomers)
	
	todaycustomers := []model.Customer{}
	GormDB.Model(&customer).Where("updated_at > ?", today).Find(&todaycustomers)
	
	z := model.CustomerView{}
	z.Customers = custs
	z.AllCustomers.Name = "All customers"
	z.AllCustomers.Total = float64(len(custs))
	z.AllCustomers.Description = "Total Customers registered"
	//////////////////////////////////////////////////////////////
	z.Lastweek.Name = "Last 7 days Customers"
	z.Lastweek.Total = float64(len(lastcustomers))
	z.Lastweek.Description = "Total Customers registered for the last seven days"
	///////////////////////////////////////////////////////////////
	z.Todays.Name = "Todays Customers"
	z.Todays.Total = float64(len(todaycustomers))
	z.Todays.Description = "Total Customers registered Today"
	
	IndexRepo.DbClose(GormDB)
	return &z, nil
}
func (customerRepo customerrepo)customerExist(email string) bool {
	customer := model.Customer{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	GormDB.Where("name = ? ", email).First(&customer)
	if customer.Name == "" {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (customerRepo customerrepo)CustomerExistbycode(customercode string) bool {
	customer := model.Customer{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	GormDB.Where("customercode = ? ", customercode).First(&customer)
	if customer.ID == 0 {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (customerRepo customerrepo)customerExistByid(id int) bool {
	customer := model.Customer{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&customer, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
