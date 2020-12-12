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
func (dashboardRepo dashboardrepo) View(dated,searchq2,searchq3 string)(*model.Dashboard, *httperors.HttpError) {
	error1 := godotenv.Load()
	if error1 != nil {
		log.Fatal("Error loading .env file in routes")
	}
	//headers
	Salesname := os.Getenv("Salesname")
	Purchasesname := os.Getenv("Purchasesname")
	Paymentsname := os.Getenv("Paymentsname")
	Receiptsname := os.Getenv("Receiptsname")
	Expencesname := os.Getenv("Expencesname")
	Customerssname := os.Getenv("Customerssname")
	Suppliersname := os.Getenv("Suppliersname")
	//staements
	Sales := os.Getenv("Sales")
	Purchases := os.Getenv("Purchases")
	Payments := os.Getenv("Payments")
	Receipts := os.Getenv("Receipts")
	Expences := os.Getenv("Expences")
	Customers := os.Getenv("Customers")
	Suppliers := os.Getenv("Suppliers")
	//icons
	Salesicon := os.Getenv("Salesicon")
	Customersicon := os.Getenv("Customersicon")
	Purchasesicon := os.Getenv("Purchasesicon")
	Receiptsicon := os.Getenv("Receiptsicon")
	Expencesicon := os.Getenv("Expencesicon")
	Suppliersicon := os.Getenv("Suppliersicon")
	Paymentsicon := os.Getenv("Paymentsicon")
	dashboard := model.Dashboard{}
	customerss,err1 := Customerrepo.AllSearch(dated,searchq2,searchq3)
	if err1 != nil {
		return nil, err1
	}
	supplierss,err2 := Supplierrepo.AllSearch(dated,searchq2,searchq3) 
	if err2 != nil {
		return nil, err2
	}
	saless,err3 := Transactionrepo.Allsearch(dated,searchq2,searchq3)
	if err3 != nil {
		return nil, err3
	}
	expencess,err4 := Expencetrasanrepo.AllSearch(dated,searchq2,searchq3)
	if err4 != nil {
		return nil, err4
	}
	purchasess,err5 := STransactionrepo.Allsearch(dated,searchq2,searchq3)
	if err5 != nil {
		return nil, err5
	}

	receiptss,err6 := Receiptrepo.AllSearch(dated,searchq2,searchq3)
	if err6 != nil {
		return nil, err6
	}

	paymentss,err6 := Paymentrepo.AllSearch(dated,searchq2,searchq3)
	if err6 != nil {
		return nil, err6
	}
	var sal float64 = 0
	for _,s := range saless {
		sal += s.Total
	}
	var pur float64 = 0
	for _,p := range purchasess {
		pur += p.Total
	}

	var pey float64 = 0
	for _,pe := range paymentss {
		pey += pe.Amount
	}

	var rec float64 = 0
	for _,re := range receiptss {
		rec += re.Amount
	}

	var ex float64 = 0
	for _,x := range expencess {
		ex += x.Amount
	}
	dashboard.Sales.Name = Salesname
	dashboard.Sales.Total = sal
	dashboard.Sales.Description = Sales
	dashboard.Sales.Icon = Salesicon
	/////////////products///////////
	dashboard.Purchases.Name = Purchasesname
	dashboard.Purchases.Total = pur
	dashboard.Purchases.Description =Purchases
	dashboard.Purchases.Icon = Purchasesicon
	/////////////products///////////
	dashboard.Payments.Name = Paymentsname
	dashboard.Payments.Total = pey
	dashboard.Payments.Description =Payments
	dashboard.Payments.Icon = Paymentsicon
	/////////////products///////////
	dashboard.Receipts.Name = Receiptsname
	dashboard.Receipts.Total = rec
	dashboard.Receipts.Description =Receipts
	dashboard.Receipts.Icon = Receiptsicon
	/////////////products///////////
	dashboard.Expences.Name = Expencesname
	dashboard.Expences.Total = ex
	dashboard.Expences.Description =Expences
	dashboard.Expences.Icon = Expencesicon
	/////////////products///////////
	dashboard.Customers.Name = Customerssname
	dashboard.Customers.Total = float64(len(customerss))
	dashboard.Customers.Description =Customers
	dashboard.Customers.Icon = Customersicon

	dashboard.Suppliers.Name = Suppliersname
	dashboard.Suppliers.Total = float64(len(supplierss))
	dashboard.Suppliers.Description =Suppliers
	dashboard.Suppliers.Icon = Suppliersicon
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

