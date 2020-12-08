package repository

import (
	"fmt"
	// "log"
	// "os"
	// "github.com/joho/godotenv"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	// "github.com/myrachanto/accounting/support"
)
//Salesrepo repo
var (
	Salesrepo salesrepo = salesrepo{}
)
///curtesy to gorm
type salesrepo struct{}
//////////////
////////////TODO user id///////////
/////////////////////////////////////////
func (salesRepo salesrepo) View()(*model.Sales, *httperors.HttpError) {
	sales := model.Sales{}
	invoices,err5 := Invoicerepo.All()
	if err5 != nil {
		return nil, err5
	}
	paidinvoices,err4 := Invoicerepo.PaidInvoices()
	if err4 != nil {
		return nil, err4
	}
	debts,err3 := Customerrepo.AllDebts()
	if err3 != nil {
		return nil, err3
	}
	transactions,err2 := Transactionrepo.All()
	if err2 != nil {
		return nil, err2
	}

	var to float64 = 0
	var tax float64 = 0
	var discount float64 = 0
	for _, s := range invoices {
		to += s.Total
		tax += s.Tax
		discount += s.Discount
	}
	gp := to - tax-discount

	var dt float64 = 0
	for _, d := range debts {
		dt += d.Amount
	}
	var pi float64 = 0
	for _, d := range paidinvoices {
		pi += d.Total
	}
	////sales/////////////
	sales.Sales.Name = "Sales"
	sales.Sales.Total = to
	sales.Sales.Description = "Total sales"
	sales.Sales.Icon = ""
	////grossprofit/////////////
	sales.GrossProfit.Name = "Gross Profit"
	sales.GrossProfit.Total = gp
	sales.GrossProfit.Description = "Gross profit recorded"
	sales.GrossProfit.Icon = ""
	////Customers/////////////
	sales.Debts.Name = "Debtors"
	sales.Debts.Total = dt
	sales.Debts.Description = "Total Debtors registered"
	sales.Debts.Icon = ""
	////Invoices/////////////
	sales.PaidInvoices.Name = "Paid Invoices"
	sales.PaidInvoices.Total = pi
	sales.PaidInvoices.Description = "Total Amount paid"
	sales.PaidInvoices.Icon = ""

	////Transaction/////////////
	fmt.Println(transactions)
	sales.Transactions = transactions
	////debtTransaction/////////////
	sales.DebtTransactions = debts
	return &sales, nil
}
func (salesRepo salesrepo) Purchases()(*model.Purchases, *httperors.HttpError) {
	purchases := model.Purchases{}
	invoices,err5 := SInvoicerepo.All()
	if err5 != nil {
		return nil, err5
	}
	paidinvoices,err4 := SInvoicerepo.PaidsInvoices()
	if err4 != nil {
		return nil, err4
	}
	debts,err3 := Supplierrepo.AllDebts()
	if err3 != nil {
		return nil, err3
	}
	transactions,err2 := STransactionrepo.All()
	if err2 != nil {
		return nil, err2
	}

	var to float64 = 0
	var tax float64 = 0
	var discount float64 = 0
	for _, s := range invoices {
		to += s.Total
		tax += s.Tax
		discount += s.Discount
	}

	var dt float64 = 0
	for _, d := range debts {
		dt += d.Amount
	}
	var pi float64 = 0
	for _, d := range paidinvoices {
		pi += d.Total
	}
	////sales/////////////
	purchases.Purchases.Name = "Purchases"
	purchases.Purchases.Total = to
	purchases.Purchases.Description = "Total Purchases"
	purchases.Purchases.Icon = ""
	////grossprofit/////////////
	////Customers/////////////
	purchases.Creditors.Name = "Creditors"
	purchases.Creditors.Total = dt
	purchases.Creditors.Description = "Total Creditors registered"
	purchases.Creditors.Icon = ""
	////Invoices/////////////
	purchases.PaidInvoices.Name = "Paid Invoices"
	purchases.PaidInvoices.Total = pi
	purchases.PaidInvoices.Description = "Total Amount paid"
	purchases.PaidInvoices.Icon = ""

	////Transaction/////////////
	fmt.Println(transactions)
	purchases.STransactions = transactions
	////debtTransaction/////////////
	purchases.CreditTransaction = debts
	return &purchases, nil
}
// func (salesRepo salesrepo) Email() (*model.Email, *httperors.HttpError) {
	
// 	email := model.Email{}
// 	email.Email = "Business@gmail.com"
// 	email.To = "example@gmail.com"
// 	email.Subject = "RE:"
// 	email.Message = "this is the email message body"
// 	customers,err4 := Customerrepo.All()
// 	if err4 != nil {
// 		return nil, err4
// 	}
// 	email.Customers = customers
	
// 	return &email, nil
// }

