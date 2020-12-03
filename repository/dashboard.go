package repository

import (
	"log"
	"os"
	"github.com/joho/godotenv"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	// "github.com/myrachanto/accounting/support"
)
//Dashboardrepo repo
var (
	Dashboardrepo dashboardrepo = dashboardrepo{}
)
///curtesy to gorm
type dashboardrepo struct{}
//////////////
////////////TODO user id///////////
/////////////////////////////////////////
func (dashboardRepo dashboardrepo) View()(*model.Dashboard, *httperors.HttpError) {
	error1 := godotenv.Load()
	if error1 != nil {
		log.Fatal("Error loading .env file in routes")
	}
	//headers
	Invoicename := os.Getenv("Invoicename")
	Customername := os.Getenv("Customersname")
	Username := os.Getenv("Usersname")
	Categoryname := os.Getenv("Categoryname")
	Majorcatname := os.Getenv("Majorcatname")
	Productsname := os.Getenv("Productsname")
	//staements
	Invoicestatement := os.Getenv("Invoice")
	Customerstatement := os.Getenv("Customers")
	Userstatement := os.Getenv("Users")
	Categorystatement := os.Getenv("Category")
	Majorcatstatement := os.Getenv("Majorcat")
	Productsstatement := os.Getenv("Products")
	//icons
	Invoicestatementicon := os.Getenv("Invoiceicon")
	Customerstatementicon := os.Getenv("Customersicon")
	Userstatementicon := os.Getenv("Usersicon")
	Categorystatementicon := os.Getenv("Categoryicon")
	Majorcatstatementicon := os.Getenv("Majorcaticon")
	Productsstatementicon := os.Getenv("Productsicon")
	dashboard := model.Dashboard{}
	product,err1 := Productrepo.All()
	if err1 != nil {
		return nil, err1
	}
	category,err2 := Categoryrepo.All() 
	if err2 != nil {
		return nil, err2
	}
	majorcat,err3 := Majorcategoryrepo.All()
	if err3 != nil {
		return nil, err3
	}
	customers,err4 := Customerrepo.All()
	if err4 != nil {
		return nil, err4
	}
	invoices,err5 := Invoicerepo.All()
	if err5 != nil {
		return nil, err5
	}

	users,err6 := Userrepo.All()
	if err6 != nil {
		return nil, err6
	}
	dashboard.Categorys.Name = Categoryname
	dashboard.Categorys.Total = len(category)
	dashboard.Categorys.Description = Categorystatement
	dashboard.Categorys.Icon = Categorystatementicon
	/////////////products///////////
	dashboard.Products.Name = Productsname
	dashboard.Products.Total = len(product)
	dashboard.Products.Description =Productsstatement
	dashboard.Products.Icon = Productsstatementicon
	/////////////products///////////
	dashboard.Majorcategorys.Name = Majorcatname
	dashboard.Majorcategorys.Total = len(majorcat)
	dashboard.Majorcategorys.Description =Majorcatstatement
	dashboard.Majorcategorys.Icon = Majorcatstatementicon
	/////////////products///////////
	dashboard.Invoices.Name = Invoicename
	dashboard.Invoices.Total = len(invoices)
	dashboard.Invoices.Description =Invoicestatement
	dashboard.Invoices.Icon = Invoicestatementicon
	/////////////products///////////
	dashboard.Users.Name = Username
	dashboard.Users.Total = len(users)
	dashboard.Users.Description =Userstatement
	dashboard.Users.Icon = Userstatementicon
	/////////////products///////////
	dashboard.Customers.Name = Customername
	dashboard.Customers.Total = len(customers)
	dashboard.Customers.Description =Customerstatement
	dashboard.Customers.Icon = Customerstatementicon
	return &dashboard, nil
}
func (dashboardRepo dashboardrepo) Email() (*model.Email, *httperors.HttpError) {
	
	email := model.Email{}
	email.Email = "Business@gmail.com"
	email.To = "example@gmail.com"
	email.Subject = "RE:"
	email.Message = "this is the email message body"
	customers,err4 := Customerrepo.All()
	if err4 != nil {
		return nil, err4
	}
	email.Customers = customers
	
	return &email, nil
}

