package repository

import (
	"fmt"
	"strings"
	"strconv"
	"time"
  "gorm.io/gorm"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/support"
)
//Invoicerepo..
var (
	Invoicerepo invoicerepo = invoicerepo{}
)

///curtesy to gorm
type invoicerepo struct{}

func (invoiceRepo invoicerepo) Create(invoice *model.Invoice) (string, *httperors.HttpError) {
	code := invoice.Code
	t,r := Cartrepo.SumTotal(code);if r != nil {
		return "", r
	}
	
	exps,er := Expencetrasanrepo.GetExpencesByCode(code);if er != nil {
		return "", er
	}
	var ep float64 = 0
	for _, exp := range exps{
		ep += exp.Amount
	}
	var tex float64 = t.Total + ep
	fmt.Println(tex, t.Total, ep)
	invoice.Expences = ep
	invoice.Total = tex
	// cart := Cartrepo.Getcustomerwithcode(code)
	// invoice.Customername = cart.Customername
	customername := invoice.Customername 
	customer := Customerrepo.Getcustomer(customername)
	now := time.Now()
	invoice.CustomerID = customer.ID
	invoice.Dated = now
	invoice.Duedate = now.AddDate(0, 1, 0)
	invoice.Discount = t.Discount
	invoice.Tax = t.Tax
	invoice.Subtotal = t.Subtotal
	invoice.Title = "sales" 
	invoice.Cn = false
	invoice.Status = "invoice"
	invoice.Description = "Sale of goods and services"
	
	transactions,e := Cartrepo.CarttoTransaction(code);if e != nil {
		return "", e
	}
	
	if invoice.Customername == "undefined" && invoice.Customername == "" {
		return "", httperors.NewNotFoundError("Please choose a Customer name!")
		
	}
	if invoice.Customername == "undefined"{
		cart := Cartrepo.Getcustomerwithcode(code)
		invoice.Customername = cart.Customername
	}
	// model.Transactions = tr 
	debtTransaction := model.DebtTransaction{Code: code, Description: "Goods sold", Customername:customername, Amount: t.Total}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}

	carts, err7 := Cartrepo.Updateproductqty(code)
	if err7 != nil {
		return "", err7
	}
	////////////begin transaction/////////////////////
	GormDB.Transaction(func(tx *gorm.DB) error {
		
		fmt.Println("level 1")
		tx.Create(&invoice)

		for _, trans := range transactions {
			tx.Transaction(func(tx2 *gorm.DB) error {
				trans.Credit = false
				trans.Status = "invoice"
				fmt.Println("level 2")
				tx2.Create(&trans)
				return nil
			})
		}
	
		tx.Transaction(func(tx3 *gorm.DB) error { 
		
			fmt.Println("level 3")
			tx3.Create(&debtTransaction)
			return nil
		})
		for _, c := range carts {
			product := Productrepo.Productqty(c.Name)
			remaining := product.Quantity - c.Quantity
			tx.Transaction(func(tx4 *gorm.DB) error {
				fmt.Println("level 4")
				tx4.Model(&product).Where("name = ?", product.Name).Update("quantity",remaining)
				return nil
			})
		}
		
			return nil
		})
	Cartrepo.DeleteAll(code)
	IndexRepo.DbClose(GormDB)
	return "invoice created succesifully", nil
}

func (invoiceRepo invoicerepo) updatedinvoice(code, status string)  *httperors.HttpError{
	invoice := model.Invoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return err1
	}
	GormDB.Where("code = ?", code).Find(&invoice)
	invoice.Status = status
	GormDB.Save(&invoice)
	IndexRepo.DbClose(GormDB)
	return nil
}
func (invoiceRepo invoicerepo) View() (*model.Cinvoiceoptions, *httperors.HttpError) {
	CIOptions := &model.Cinvoiceoptions{}

	customers,err1 := Customerrepo.All()
	if err1 != nil {
		return nil, httperors.NewNotFoundError("You got an error fetching customers")
	}
	taxs,err2 := Taxrepo.All()
	if err2 != nil {
		return nil, httperors.NewNotFoundError("You got an error fetching customers")
	}
 
	products,err3 := Productrepo.All()
	if err3 != nil {
		return nil, httperors.NewNotFoundError("You got an error fetching customers")
	}
	prices,err5 := Pricerepo.All()
	if err5 != nil {
		return nil, httperors.NewNotFoundError("You got an error fetching customers")
	}	
	discounts,err6 := Discountrepo.All()
	if err6 != nil {
		return nil, httperors.NewNotFoundError("You got an error fetching customers")
	}
	paymentforms,err7 := Paymentformrepo.All()
	if err7 != nil {
		return nil, httperors.NewNotFoundError("You got an error fetching customers")
	}
	expences,err8 := Expencerepo.All()
	if err8 != nil {
		return nil, httperors.NewNotFoundError("You got an error fetching customers")
	}
	code,err4 := invoiceRepo.GeneCode()
	if err4 != nil {
		return nil, httperors.NewNotFoundError("You got an error fetching customers")
	}
	CIOptions.Code = code
	CIOptions.Customers = customers
	CIOptions.Taxs = taxs
	CIOptions.Products = products
	CIOptions.Prices = prices
	CIOptions.Discounts = discounts
	CIOptions.Paymentform = paymentforms
	CIOptions.Expences = expences
	return CIOptions, nil
} 
func (invoiceRepo invoicerepo) GetOne(code string) (*model.InvoiceView, *httperors.HttpError) {
	ok := invoiceRepo.InvoiceExistByCode(code)
	if !ok {
		return nil, httperors.NewNotFoundError("invoice with that code does not exists!")
	}
	invoice := model.Invoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil { 
		return nil, err1
	} 

	transactions, e := Transactionrepo.GetTransactionsinvoice(code)
	if e != nil {
		return nil, e
	}	
	expences, er := Expencetrasanrepo.Getexpencestransactions(code)
	if er != nil {
		return nil, er
	}
	credits, er2 := Transactionrepo.GetTransactionscredit(code)
	if er2 != nil {
		return nil, er2
	}
	
	GormDB.Model(&invoice).Where("code = ?", code).First(&invoice)
	IndexRepo.DbClose(GormDB)
	customer := Customerrepo.Getcustomer(invoice.Customername)
	return &model.InvoiceView{
		Invoice: invoice,
		Customer: customer,
		Transactions: transactions,
		ExpencesTrans:expences,
		Credits: credits,
	}, nil
}
func (invoiceRepo invoicerepo) All() (t []model.Invoice, r *httperors.HttpError) {

	invoice := model.Invoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&invoice).Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (invoiceRepo invoicerepo) InvoiceByCustomer(name string) (t []model.Invoice, r *httperors.HttpError) {

	invoice := model.Invoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&invoice).Where("customername = ?", name).Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (invoiceRepo invoicerepo) PaidInvoices() (t []model.Invoice, r *httperors.HttpError) {

	invoice := model.Invoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&invoice).Where("status = ?", "paid").Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (invoiceRepo invoicerepo) CustomerCredits(name string) (t []model.Invoice, r *httperors.HttpError) {

	invoice := model.Invoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&invoice).Where("customername = ? AND status = ?", name, "credit").Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (invoiceRepo invoicerepo) Customerinvoice(name string) (t []model.Invoice, r *httperors.HttpError) {

	invoice := model.Invoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&invoice).Where("customername = ? AND status = ?", name, "invoice").Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (invoiceRepo invoicerepo) GetAll(invoices []model.Invoice,search *support.Search) ([]model.Invoice, *httperors.HttpError) {
	
	results, err1 := invoiceRepo.Search(search, invoices)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}
func (invoiceRepo invoicerepo) GetCredit(invoices []model.Invoice,search *support.Search) ([]model.Invoice, *httperors.HttpError) {
	
	results, err1 := invoiceRepo.Search(search, invoices)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}
func (invoiceRepo invoicerepo) GetCreditNotes(search string) ([]model.Invoice, *httperors.HttpError) {
	credits := []model.Invoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Where("name LIKE ?", "%"+ search +"%").Or("title LIKE ?", "%"+ search +"%").Or("description LIKE ?", "%"+ search +"%").Find(&credits)
	if err1 != nil { 
			return nil, err1
		}
	return credits, nil
}
func (invoiceRepo invoicerepo) Update(code string) (string, *httperors.HttpError) {
	
	invoice := model.Invoice{}
	transactions := []model.Transaction{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}

	ainvoice := model.Invoice{}
	
	GormDB.Where("code = ?", code).First(&ainvoice)
		
	customername := ainvoice.Customername
	customer := Customerrepo.Getcustomer(customername)
	now := time.Now()
	invoice.CustomerID = customer.ID
	invoice.Dated = now
	invoice.Code = code
	invoice.Duedate = now.AddDate(0, 0, 1)
	invoice.Title = "Credit"
	invoice.Status = "credit"
	invoice.Description = "Credit of goods and services"
	GormDB.Where("code = ? AND credit = ? AND status = ?", code, false, "pending").Find(&transactions)
	var tax float64 = 0
	var discount float64 = 0
	var total float64 = 0
	for _, i := range transactions {
		tax += i.Tax
		discount += i.Discount
		total += i.Total
	}
	
fmt.Println(tax, discount, total)
	debtTransaction := model.DebtTransaction{Code: code, Description: "Goods Credited", Customername:customername, Amount: -total}

	invoice.Discount = discount
	invoice.Customername = customername
	invoice.Tax = tax
	invoice.Subtotal = (total-tax+discount)
	invoice.Total = total
	invoice.Cn = true
	trans := model.Transaction{}
		////////////begin transaction/////////////////////
		GormDB.Transaction(func(tx *gorm.DB) error {
		
			fmt.Println("level 1")
			tx.Create(&invoice)
		
			tx.Transaction(func(tx2 *gorm.DB) error { 
			
				fmt.Println("level 2")
				tx2.Create(&debtTransaction)
				return nil
			})
			for _, c := range transactions {
				product := Productrepo.Productqty(c.Productname)
				remaining := product.Quantity + c.Quantity
				tx.Transaction(func(tx3 *gorm.DB) error {
					fmt.Println("level 3")
					tx3.Model(&product).Where("name = ?", product.Name).Update("quantity",remaining)
					return nil
				})
			}
			for _, t := range transactions {
				tx.Transaction(func(tx4 *gorm.DB) error {
					fmt.Println("level 4")
					tx4.Model(&trans).Where("code = ? AND productname = ? AND total > ?", code, t.Productname, 0).Select("credit", "status").UpdateColumns(model.Transaction{Credit: true, Status: "credit"})
					return nil
				})
			}
			//////////////end of transaction///////////////

			return nil
		})
	IndexRepo.DbClose(GormDB)

	return "item credited successifully", nil
}
func (invoiceRepo invoicerepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := invoiceRepo.invoiceUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("invoice with that id does not exists!")
	}
	invoice := model.Invoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&invoice).Where("id = ?", id).First(&invoice)
	GormDB.Delete(invoice)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (invoiceRepo invoicerepo)invoiceUserExistByid(id int) bool {
	invoice := model.Invoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&invoice, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (invoiceRepo invoicerepo)InvoiceExistByCode(code string) bool {
	invoice := model.Invoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}

	res := GormDB.First(&invoice, "code =?", code)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB) 
	return true
	
}
func (invoiceRepo invoicerepo)GeneCode() (string, *httperors.HttpError) {
	invoice := model.Invoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	err := GormDB.Last(&invoice)
	if err.Error != nil {
		var c1 uint = 1
		code := "CustomerInvNo"+strconv.FormatUint(uint64(c1), 10)
		return code, nil
	 }
	c1 := invoice.ID + 1
	code := "CustomerInvNo"+strconv.FormatUint(uint64(c1), 10)
	IndexRepo.DbClose(GormDB)
	return code, nil
	
}
func (invoiceRepo invoicerepo) Search(Ser *support.Search, invoices []model.Invoice)([]model.Invoice, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	invoice := model.Invoice{}
	switch(Ser.Search_operator){
	case "all":
	GormDB.Model(&invoice).Order(Ser.Column+" "+Ser.Direction).Find(&invoices)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
	
	break;
	case "equal_to":
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&invoices);
	
	break;
	case "not_equal_to":
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&invoices);	
	
	break;
	case "less_than" :
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&invoices);	
	
	break;
	case "greater_than":
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&invoices);	
	
	break;
	case "less_than_or_equal_to":
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&invoices);	
	
	break;
	case "greater_than_ro_equal_to":
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&invoices);	
	
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&invoices);
	
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
	GormDB.Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&invoices);
	
	// break;
case "like":
	// fmt.Println(Ser.Search_query_1)
	if Ser.Search_query_1 == "all" {
			//db.Order("name DESC")
	GormDB.Order(Ser.Column+" "+Ser.Direction).Find(&invoices)

	}else {

		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&invoices);
	}
break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&invoices);
	
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return invoices, nil
}
