package repository

import (
	"fmt"
	"time"
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
func (salesRepo salesrepo) View(dated,searchq2,searchq3 string)(*model.Sales, *httperors.HttpError) {
sales := model.Sales{}
	paidinvoices,err4 := Invoicerepo.PaidInvoices(dated,searchq2,searchq3)
	if err4 != nil {
		return nil, err4
	}
	debts,err3 := Customerrepo.AllDebts(dated,searchq2,searchq3)
	if err3 != nil {
		return nil, err3
	}
	transactions,err2 := Transactionrepo.Allsearch(dated,searchq2,searchq3)
	if err2 != nil {
		return nil, err2
	}
	sal,err5 := Transactionrepo.Sales(dated,searchq2,searchq3)
	if err5 != nil {
		return nil, err5
	}
	var sale float64 = 0 
	var cost float64 = 0
	for _,i := range sal {
		sale += i.Total
		cost += i.Cost
	}
	
	gp := sale - cost

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
	sales.Sales.Total = sale
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
func (salesRepo salesrepo) Purchases(dated,searchq2,searchq3 string)(*model.Purchases, *httperors.HttpError) {
	purchases := model.Purchases{}
	invoices,err5 := SInvoicerepo.AllSearch(dated,searchq2,searchq3)
	if err5 != nil {
		return nil, err5
	}
	paidinvoices,err4 := SInvoicerepo.PaidsInvoices(dated,searchq2,searchq3)
	if err4 != nil {
		return nil, err4
	}
	debts,err3 := Supplierrepo.AllDebts(dated,searchq2,searchq3)
	if err3 != nil {
		return nil, err3
	}
	transactions,err2 := STransactionrepo.Allsearch(dated,searchq2,searchq3)
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
func (salesRepo salesrepo) PL(dated,searchq2,searchq3 string)(*model.Pl, *httperors.HttpError) {
	sales,err5 := Transactionrepo.Sales(dated,searchq2,searchq3)
	if err5 != nil {
		return nil, err5
	}
	var sale float64 = 0 
	var cost float64 = 0
	for _,i := range sales {
		sale += i.Total
		cost += i.Cost
	}
	directex,er := Expencetrasanrepo.Alldirect(dated,searchq2,searchq3)
	if er != nil {
		return nil, er
	}
	var dex float64 = 0
	for _,de := range directex {
		dex += de.Amount
	}

	indirectex,er1 := Expencetrasanrepo.Allindirect(dated,searchq2,searchq3)
	if er1 != nil {
		return nil, er1
	}
	var idex float64 = 0
	for _,ide := range indirectex {
		idex += ide.Amount
	}
	otherex,er3 := Expencetrasanrepo.Allother(dated,searchq2,searchq3)
	if er3 != nil {
		return nil, er3
	}
	var oex float64 = 0
	for _,oe := range otherex {
		oex += oe.Amount
	}
	if dated == "custom" {
		dato := searchq2 +" to "+ searchq3
		return &model.Pl{
			Sales:sale,
			Costofsale:cost,
			DirectExpence:dex,
			InDirectExpence:idex,
			OtherExpence:oex,
			Dated: dato,
		}, nil
	}
	return &model.Pl{
		Sales:sale,
		Costofsale:cost,
		DirectExpence:dex,
		InDirectExpence:idex,
		OtherExpence:oex, 
		Dated: dated,
	}, nil
}

func (salesRepo salesrepo) Supplierstement(code,dated,searchq2,searchq3 string)(*model.Supplierstates, *httperors.HttpError) {
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	results,err := Salesrepo.AllsearchSupplier(code,dated,searchq2,searchq3)
	if err != nil {
		return nil, err
	}

	supplier := Supplierrepo.Getsupplierwithcode(code)

	IndexRepo.DbClose(GormDB)
	return &model.Supplierstates{
		Statements:results,
		Supplier:supplier,
		Datos:dated,
		Query1:searchq2,
		Query2:searchq3,
	}, nil 
}
func (salesRepo salesrepo) Customerstement(code,dated,searchq2,searchq3 string)(*model.Customerstates, *httperors.HttpError) {
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	results,err := Salesrepo.Allsearch(code,dated,searchq2,searchq3)
	if err != nil {
		return nil, err
	}

	customer := Customerrepo.GetcustomerwithCode(code)

	IndexRepo.DbClose(GormDB)
	return &model.Customerstates{
		Statements:results,
		Customer:customer,
		Datos:dated,
		Query1:searchq2,
		Query2:searchq3,
	}, nil 
}

func (salesRepo salesrepo) AllsearchSupplier(code,dated,searchq2,searchq3 string) (results []model.CustomerState, r *httperors.HttpError) {

	now := time.Now()
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}

	if dated != "custom"{
		if dated == "In the last 24hrs"{
			d := now.AddDate(0, 0, -1)
			GormDB.Table("s_invoices").Select("s_invoices.description, s_invoices.tax, s_invoices.discount, s_invoices.total, s_invoices.balance,s_invoices.dated, payments.amount, payments.status,payments.created_at").Joins("left join payments on payments.itemcode = s_invoices.suppliercode").Where("s_invoices.updated_at > ? AND s_invoices.suppliercode = ?", d, code).Find(&results)
		}
		if dated == "In the last 7days"{
			d := now.AddDate(0, 0, -7)
			GormDB.Table("s_invoices").Select("s_invoices.description, s_invoices.tax, s_invoices.discount, s_invoices.total,s_invoices.balance,s_invoices.dated, payments.amount, payments.status,payments.created_at").Joins("left join payments on payments.itemcode = s_invoices.suppliercode").Where("s_invoices.updated_at > ? AND s_invoices.suppliercode = ?", d, code).Find(&results)
		}
		if dated == "In the last 15day"{
			d := now.AddDate(0, 0, -15)
			GormDB.Table("s_invoices").Select("s_invoices.description, s_invoices.tax, s_invoices.discount, s_invoices.total,s_invoices.balance,s_invoices.dated, payments.amount, payments.status,payments.created_at").Joins("left join payments on payments.itemcode = s_invoices.suppliercode").Where("s_invoices.updated_at > ? AND s_invoices.suppliercode = ?", d, code).Find(&results)
		}
		if dated == "In the last 30days"{
			d := now.AddDate(0, 0, -30)
			GormDB.Table("s_invoices").Select("s_invoices.description, s_invoices.tax, s_invoices.discount, s_invoices.total,s_invoices.balance,s_invoices.dated, payments.amount, payments.status,payments.created_at").Joins("left join payments on payments.itemcode = s_invoices.suppliercode").Where("s_invoices.updated_at > ? AND s_invoices.suppliercode = ?", d, code).Find(&results)
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
		GormDB.Table("s_invoices").Select("s_invoices.description, s_invoices.tax, s_invoices.discount, s_invoices.total,s_invoices.balance,s_invoices.dated, payments.amount, payments.status,payments.created_at").Joins("left join payments on payments.itemcode = s_invoices.suppliercode").Where("s_invoices.suppliercode = ? AND s_invoices.updated_at BETWEEN ? AND ?",code, start, end).Find(&results)
	}
	IndexRepo.DbClose(GormDB)
	return results, nil

}
func (salesRepo salesrepo) Allsearch(code,dated,searchq2,searchq3 string) (results []model.CustomerState, r *httperors.HttpError) {

	now := time.Now()
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}

	if dated != "custom"{
		if dated == "In the last 24hrs"{
			d := now.AddDate(0, 0, -1)
			GormDB.Table("invoices").Select("invoices.description, invoices.tax, invoices.discount, invoices.total, invoices.balance,invoices.dated, receipts.amount, receipts.status,receipts.created_at").Joins("left join receipts on receipts.customercode = invoices.customercode").Where("invoices.updated_at > ? AND invoices.customercode = ?", d, code).Find(&results)
		}
		if dated == "In the last 7days"{
			d := now.AddDate(0, 0, -7)
			GormDB.Table("invoices").Select("invoices.description, invoices.tax, invoices.discount, invoices.total,invoices.balance,invoices.dated, receipts.amount, receipts.status,receipts.created_at").Joins("left join receipts on receipts.customercode = invoices.customercode").Where("invoices.updated_at > ? AND invoices.customercode = ?", d, code).Find(&results)
		}
		if dated == "In the last 15day"{
			d := now.AddDate(0, 0, -15)
			GormDB.Table("invoices").Select("invoices.description, invoices.tax, invoices.discount, invoices.total,invoices.balance,invoices.dated, receipts.amount, receipts.status,receipts.created_at").Joins("left join receipts on receipts.customercode = invoices.customercode").Where("invoices.updated_at > ? AND invoices.customercode = ?", d, code).Find(&results)
		}
		if dated == "In the last 30days"{
			d := now.AddDate(0, 0, -30)
			GormDB.Table("invoices").Select("invoices.description, invoices.tax, invoices.discount, invoices.total,invoices.balance,invoices.dated, receipts.amount, receipts.status,receipts.created_at").Joins("left join receipts on receipts.customercode = invoices.customercode").Where("invoices.updated_at > ? AND invoices.customercode = ?", d, code).Find(&results)
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
		GormDB.Table("invoices").Select("invoices.description, invoices.tax, invoices.discount, invoices.total,invoices.balance,invoices.dated, receipts.amount, receipts.status,receipts.created_at").Joins("left join receipts on receipts.customercode = invoices.customercode").Where("invoices.customercode = ? AND invoices.updated_at BETWEEN ? AND ?",code, start, end).Find(&results)
	}
	IndexRepo.DbClose(GormDB)
	return results, nil

}
// func (salesRepo salesrepo) Purchases(dated,searchq2,searchq3 string)(*model.Purchases, *httperors.HttpError) {
 
// 	return &purchases, nil
// }
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

