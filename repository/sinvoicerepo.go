package repository

import (
	"fmt"
	"strconv"
	"time"

	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"gorm.io/gorm"
)

//sInvoicerepo..supplier invoice supplier
var (
	SInvoicerepo sInvoicerepo = sInvoicerepo{}
)

///curtesy to gorm
type sInvoicerepo struct{}

func (sInvoiceRepo sInvoicerepo) Create(sInvoice *model.SInvoice) (string, *httperors.HttpError) {
	code := sInvoice.Code
	t, r := Scartrepo.SumTotal(code) 
	if r != nil {
		return "", r
	}

	exps, er := Expencetrasanrepo.GetExpencesByCode(code)
	if er != nil {
		return "", er
	}
	var ep float64 = 0
	for _, exp := range exps {
		ep += exp.Amount
	}
	var tex float64 = t.Total + ep
	fmt.Println(tex, t.Total, ep)
	sInvoice.Expences = ep
	sInvoice.Total = tex
	// cart := Scartrepo.Getsupplierwithcode(code)
	// sInvoice.Suppliername = cart.suppliername
	suppliername := sInvoice.Suppliername
	supplier := Supplierrepo.Getsupplier(suppliername)
	now := time.Now()
	sInvoice.SupplierID = supplier.ID
	sInvoice.Dated = now
	sInvoice.Duedate = now.AddDate(0, 1, 0)
	sInvoice.Discount = t.Discount
	sInvoice.Tax = t.Tax
	sInvoice.Subtotal = t.Subtotal
	sInvoice.Title = "Purchases"
	sInvoice.Balance = tex
	sInvoice.Cn = false
	sInvoice.Status = "invoice"
	sInvoice.Paidstatus = "notpaid"
	sInvoice.AllPaidstatus = "notpaid"
	sInvoice.Description = "Sale of goods and services"

	transactions, e := Scartrepo.ScarttoTransaction(code)
	if e != nil {
		return "", e
	}
	if sInvoice.Suppliername == "" {
		return "", httperors.NewNotFoundError("Please choose a Supplier name!")
		
	}
	if sInvoice.Suppliername == "undefined" {
		cart := Scartrepo.Getsupplierwithcode(code)
		sInvoice.Suppliername = cart.Suppliername
	}
	// model.Transactions = tr
	creditTransaction := model.CreditTransaction{Code: code, Description: "Goods Bought", Suppliername: suppliername, Amount: -t.Total}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}

	carts, err7 := Scartrepo.Updateproductqty(code)
	if err7 != nil {
		return "", err7
	}
	////////////begin transaction/////////////////////
	GormDB.Transaction(func(tx *gorm.DB) error {

		fmt.Println("level 1")
		tx.Create(&sInvoice)

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
			tx3.Create(&creditTransaction)
			return nil
		})
		for _, c := range carts {
			product := Productrepo.Productqty(c.Name)
			remaining := product.Quantity + c.Quantity
			tx.Transaction(func(tx4 *gorm.DB) error {
				fmt.Println("level 4")
				tx4.Model(&product).Where("name = ?", product.Name).Updates(model.Product{Quantity: remaining, Bprice: c.BPrice})
				return nil
			})
		}

		return nil
	})
	Scartrepo.DeleteAll(code)
	IndexRepo.DbClose(GormDB)
	return "Invoice created successifully", nil
}

func (sInvoiceRepo sInvoicerepo) updatedsInvoice(code, status string) *httperors.HttpError {
	sInvoice := model.SInvoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return err1
	}
	GormDB.Where("code = ?", code).Find(&sInvoice)
	sInvoice.Status = status
	GormDB.Save(&sInvoice)
	IndexRepo.DbClose(GormDB)
	return nil
}
func (sInvoiceRepo sInvoicerepo) View() (*model.Sinvoiceoptions, *httperors.HttpError) {
	SIOptions := &model.Sinvoiceoptions{}

	suppliers, err1 := Supplierrepo.All()
	if err1 != nil {
		return nil, httperors.NewNotFoundError("You got an error fetching suppliers")
	}
	taxs, err2 := Taxrepo.All()
	if err2 != nil {
		return nil, httperors.NewNotFoundError("You got an error fetching suppliers")
	}

	products, err3 := Productrepo.All()
	if err3 != nil {
		return nil, httperors.NewNotFoundError("You got an error fetching suppliers")
	}
	prices, err5 := Pricerepo.All()
	if err5 != nil {
		return nil, httperors.NewNotFoundError("You got an error fetching suppliers")
	}
	discounts, err6 := Discountrepo.All()
	if err6 != nil {
		return nil, httperors.NewNotFoundError("You got an error fetching suppliers")
	}
	paymentforms, err7 := Paymentformrepo.All()
	if err7 != nil {
		return nil, httperors.NewNotFoundError("You got an error fetching suppliers")
	}
	expences, err8 := Expencerepo.All()
	if err8 != nil {
		return nil, httperors.NewNotFoundError("You got an error fetching suppliers")
	}
	code, err4 := sInvoiceRepo.GeneCode()
	if err4 != nil {
		return nil, httperors.NewNotFoundError("something went wrong with the code")
	}
	SIOptions.Code = code
	SIOptions.Suppliers = suppliers
	SIOptions.Taxs = taxs
	SIOptions.Products = products
	SIOptions.Prices = prices
	SIOptions.Discounts = discounts
	SIOptions.Paymentform = paymentforms
	SIOptions.Expences = expences
	return SIOptions, nil
}
func (sInvoiceRepo sInvoicerepo) InvoiceByCustomercodenotpaid(code string) (t []model.SInvoice, r *httperors.HttpError) {

	invoice := model.SInvoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&invoice).Where("suppliercode = ? AND paidstatus = ?", code, "notpaid").Or("suppliercode = ? AND paidstatus = ?", code, "partialpaid").Find(&t)
	// fmt.Println(t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (sInvoiceRepo sInvoicerepo) GetOne(code string) (*model.SInvoiceView, *httperors.HttpError) {
	ok := sInvoiceRepo.SInvoiceExistByCode(code)
	if !ok {
		return nil, httperors.NewNotFoundError("sInvoice with that code does not exists!")
	}
	sInvoice := model.SInvoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}

	transactions, e := STransactionrepo.GetsTransactionsinvoice(code)
	if e != nil {
		return nil, e
	}
	expences, er := Expencetrasanrepo.Getexpencestransactions(code)
	if er != nil {
		return nil, er
	}
	grns, er2 := STransactionrepo.GetsTransactionscredit(code)
	if er2 != nil {
		return nil, er2
	}

	GormDB.Model(&sInvoice).Where("code = ?", code).First(&sInvoice)
	IndexRepo.DbClose(GormDB)
	// supplier := Supplierrepo.Getsupplier(sInvoice.Suppliername)
	supplier := Supplierrepo.Getsupplierwithcode(sInvoice.Suppliercode)
	return &model.SInvoiceView{
		SInvoice:      sInvoice,
		Supplier:      supplier,
		STransactions: transactions,
		ExpencesTrans: expences,
		Grn:           grns,
	}, nil
}
func (sInvoiceRepo sInvoicerepo) All() (t []model.SInvoice, r *httperors.HttpError) {

	sInvoice := model.SInvoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&sInvoice).Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (sInvoiceRepo sInvoicerepo) AllSearch(dated,searchq2,searchq3 string) (results []model.SInvoice, r *httperors.HttpError) {

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
func (sInvoiceRepo sInvoicerepo) SInvoiceBysupplier(name string) (t []model.SInvoice, r *httperors.HttpError) {

	sInvoice := model.SInvoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&sInvoice).Where("suppliername = ?", name).Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (sInvoiceRepo sInvoicerepo) SInvoiceBysuppliercode(code string) (t []model.SInvoice, r *httperors.HttpError) {

	sInvoice := model.SInvoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&sInvoice).Where("suppliercode = ?", code).Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (sInvoiceRepo sInvoicerepo) PaidsInvoices(dated,searchq2,searchq3 string) (results []model.SInvoice, r *httperors.HttpError) {
	now := time.Now()
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}

	if dated != "custom"{
		if dated == "In the last 24hrs"{
			d := now.AddDate(0, 0, -1)
			GormDB.Where("updated_at > ? AND status = ?", d,"paid").Find(&results)
		}
		if dated == "In the last 7days"{
			d := now.AddDate(0, 0, -7)
			GormDB.Where("updated_at > ? AND status = ?", d,"paid").Find(&results)
		}
		if dated == "In the last 15day"{
			d := now.AddDate(0, 0, -15)
			GormDB.Where("updated_at > ? AND status = ?", d,"paid").Find(&results)
		}
		if dated == "In the last 30days"{
			d := now.AddDate(0, 0, -30)
			GormDB.Where("updated_at > ? AND status = ?", d,"paid").Find(&results)
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
		GormDB.Where("status = ? AND updated_at BETWEEN ? AND ?","paid", start, end).Find(&results)
	}
	IndexRepo.DbClose(GormDB)
	return results, nil

}
func (sInvoiceRepo sInvoicerepo) SupplierCredits(name string) (t []model.SInvoice, r *httperors.HttpError) {

	sInvoice := model.SInvoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&sInvoice).Where("suppliername = ? AND status = ?", name, "credit").Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (sInvoiceRepo sInvoicerepo) SupplierCreditsbycode(code,dated,searchq2,searchq3 string) (results []model.SInvoice, r *httperors.HttpError) {

	
	now := time.Now()
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}

	if dated != "custom"{
		if dated == "In the last 24hrs"{
			d := now.AddDate(0, 0, -1)
			GormDB.Where("suppliercode = ? AND updated_at > ? AND status = ?", code, d,"credit").Find(&results)
		}
		if dated == "In the last 7days"{
			d := now.AddDate(0, 0, -7)
			GormDB.Where("suppliercode = ? AND updated_at > ? AND status = ?", code, d,"credit").Find(&results)
		}
		if dated == "In the last 15day"{
			d := now.AddDate(0, 0, -15)
			GormDB.Where("suppliercode = ? AND updated_at > ? AND status = ?", code, d,"credit").Find(&results)
		}
		if dated == "In the last 30days"{
			d := now.AddDate(0, 0, -30)
			GormDB.Where("suppliercode = ? AND updated_at > ? AND status = ?", code, d,"credit").Find(&results)
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
		GormDB.Where("suppliercode = ? AND status = ? AND updated_at BETWEEN ? AND ?",code,"credit", start, end).Find(&results)
	}
	IndexRepo.DbClose(GormDB)
	return results, nil

}
func (sInvoiceRepo sInvoicerepo) SuppliersInvoice(name string) (t []model.SInvoice, r *httperors.HttpError) {

	sInvoice := model.SInvoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&sInvoice).Where("suppliername = ? AND status = ?", name, "invoice").Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}

func (sInvoiceRepo sInvoicerepo) SuppliersInvoicebycode(code,dated,searchq2,searchq3 string) (results []model.SInvoice, r *httperors.HttpError) {

	now := time.Now()
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	if dated != "custom"{
		if dated == "In the last 24hrs"{
			d := now.AddDate(0, 0, -1)
			GormDB.Where("suppliercode = ? AND updated_at > ? AND status = ?", code, d,"invoice").Find(&results)
		}
		if dated == "In the last 7days"{
			d := now.AddDate(0, 0, -7)
			GormDB.Where("suppliercode = ? AND updated_at > ? AND status = ?", code, d,"invoice").Find(&results)
		}
		if dated == "In the last 15day"{
			d := now.AddDate(0, 0, -15)
			GormDB.Where("suppliercode = ? AND updated_at > ? AND status = ?", code, d,"invoice").Find(&results)
		}
		if dated == "In the last 30days"{
			d := now.AddDate(0, 0, -30)
			GormDB.Where("suppliercode = ? AND updated_at > ? AND status = ?", code, d,"invoice").Find(&results)
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
		GormDB.Where("suppliercode = ? AND status = ? AND updated_at BETWEEN ? AND ?",code,"invoice", start, end).Find(&results)
	}
	IndexRepo.DbClose(GormDB)
	return results, nil

}
func (sInvoiceRepo sInvoicerepo) GetAll(search,dated,searchq2,searchq3 string) ([]model.SInvoice, *httperors.HttpError) {
	results := []model.SInvoice{}
	now := time.Now()
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	fmt.Println(search,dated,searchq2,searchq3)
	if search == "" && dated == ""{
		GormDB.Where("status = ?","invoice").Find(&results)
	}
	if search != "" && dated == ""{
		GormDB.Where("suppliername LIKE ? AND status = ?", "%"+search+"%", "invoice").Or("code LIKE ? AND status = ?", "%"+search+"%", "invoice").Or("title LIKE ? AND status = ?", "%"+search+"%", "invoice").Or("description LIKE ? AND status = ?", "%"+search+"%", "invoice").Find(&results)
	}
	if search != "" && dated != "custom"{
		if dated == "In the last 24hrs"{
			d := now.AddDate(0, 0, -1)
			GormDB.Where("suppliername LIKE ? AND dated > ? AND status = ?", "%"+search+"%",d,"invoice").Or("code LIKE ? AND dated > ? AND status = ?", "%"+search+"%",d,"invoice").Or("title LIKE ? AND dated > ? AND status = ?", "%"+search+"%",d,"invoice").Or("description LIKE ? AND dated > ? AND status = ?", "%"+search+"%",d,"invoice").Find(&results)
		}
		if dated == "In the last 7days"{
			d := now.AddDate(0, 0, -7)
			GormDB.Where("suppliername LIKE ? AND dated > ? AND status = ?", "%"+search+"%",d,"invoice").Or("code LIKE ? AND dated > ? AND status = ?", "%"+search+"%",d,"invoice").Or("title LIKE ? AND dated > ? AND status = ?", "%"+search+"%",d,"invoice").Or("description LIKE ? AND dated > ? AND status = ?", "%"+search+"%",d,"invoice").Find(&results)
		}
		if dated == "In the last 15day"{
			d := now.AddDate(0, 0, -15)
			GormDB.Where("suppliername LIKE ? AND dated > ? AND status = ?", "%"+search+"%",d,"invoice").Or("code LIKE ? AND dated > ? AND status = ?", "%"+search+"%",d,"invoice").Or("title LIKE ? AND dated > ? AND status = ?", "%"+search+"%",d,"invoice").Or("description LIKE ? AND dated > ? AND status = ?", "%"+search+"%",d,"invoice").Find(&results)
		}
		if dated == "In the last 30days"{
			d := now.AddDate(0, 0, -30)
			GormDB.Where("suppliername LIKE ? AND dated > ? AND status = ?", "%"+search+"%",d,"invoice").Or("code LIKE ? AND dated > ? AND status = ?", "%"+search+"%",d,"invoice").Or("title LIKE ? AND dated > ? AND status = ?", "%"+search+"%",d,"invoice").Or("description LIKE ? AND dated > ? AND status = ?", "%"+search+"%",d,"invoice").Find(&results)
		}
	}
	if search == "" && dated != "custom"{
		if dated == "In the last 24hrs"{
			d := now.AddDate(0, 0, -1)
			GormDB.Where("dated > ? AND status = ?", d,"invoice").Or("dated > ? AND status = ?",d,"invoice").Or("dated > ? AND status = ?", d,"invoice").Or("dated > ? AND status = ?",d,"invoice").Find(&results)
		}
		if dated == "In the last 7days"{
			d := now.AddDate(0, 0, -7)
			GormDB.Where("dated > ? AND status = ?",d,"invoice").Or("dated > ? AND status = ?",d,"invoice").Or("dated > ? AND status = ?",d,"invoice").Or("dated > ? AND status = ?", d,"invoice").Find(&results)
		}
		if dated == "In the last 15day"{
			d := now.AddDate(0, 0, -15)
			GormDB.Where("dated > ? AND status = ?", d,"invoice").Or("dated > ? AND status = ?",d,"invoice").Or("dated > ? AND status = ?", d,"invoice").Or("dated > ? AND status = ?", d,"invoice").Find(&results)
		}
		if dated == "In the last 30days"{
			d := now.AddDate(0, 0, -30)
			GormDB.Where("dated > ? AND status = ?",d,"invoice").Or("dated > ? AND status = ?",d,"invoice").Or("dated > ? AND status = ?",d,"invoice").Or("dated > ? AND status = ?",d,"invoice").Find(&results)
		}
	}
	if search != "" && dated == "custom"{
		start,err := time.Parse(Layout,searchq2)
		if err != nil {
			return nil, httperors.NewNotFoundError("Something went wrong parsing date1!")
		}
		end,err1 := time.Parse(Layout,searchq3)
		if err1 != nil {
			return nil, httperors.NewNotFoundError("Something went wrong parsing date1!")
		}
		GormDB.Where("suppliername LIKE ? AND status = ? AND dated BETWEEN ? AND ?", "%"+search+"%","invoice",start, end).Or("code LIKE ? AND status = ? AND dated BETWEEN ? AND ?", "%"+search+"%","invoice",start, end).Or("title LIKE ? AND status = ? AND dated BETWEEN ? AND ?", "%"+search+"%","invoice",start, end).Or("description LIKE ? AND status = ? AND dated BETWEEN ? AND ?", "%"+search+"%","invoice",start, end).Find(&results)
	}
	if search == "" && dated == "custom"{
		start,err := time.Parse(Layout,searchq2)
		if err != nil {
			return nil, httperors.NewNotFoundError("Something went wrong parsing date1!")
		}
		end,err1 := time.Parse(Layout,searchq3)
		if err1 != nil {
			return nil, httperors.NewNotFoundError("Something went wrong parsing date1!")
		}
		GormDB.Where("status = ? AND dated BETWEEN ? AND ?","invoice",start, end).Or("status = ? AND dated BETWEEN ? AND ?","invoice",start, end).Or("status = ? AND dated BETWEEN ? AND ?","invoice",start, end).Or("status = ? AND dated BETWEEN ? AND ?","invoice",start, end).Find(&results)
	}

	IndexRepo.DbClose(GormDB)
	return results, nil
}
func (sInvoiceRepo sInvoicerepo) GetCredit(search,dated,searchq2,searchq3 string) ([]model.SInvoice, *httperors.HttpError) {
	results := []model.SInvoice{}
	now := time.Now()
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	if search == "" && dated == ""{
		GormDB.Where("status = ?","credit").Find(&results)
	}
	if search != "" && dated == ""{
		GormDB.Where("suppliername LIKE ? AND status = ?", "%"+search+"%", "credit").Or("code LIKE ? status = ?", "%"+search+"%", "credit").Or("title LIKE ? status = ?", "%"+search+"%", "credit").Or("description LIKE ? status = ?", "%"+search+"%", "credit").Find(&results)
	}
	if search != "" && dated != "custom"{
		if dated == "In the last 24hrs"{
			d := now.AddDate(0, 0, -1)
			GormDB.Where("suppliername LIKE ? AND dated > ? AND status = ?", "%"+search+"%",d,"credit").Or("code LIKE ? AND dated > ? AND status = ?", "%"+search+"%",d,"credit").Or("title LIKE ? AND dated > ? AND status = ?", "%"+search+"%",d,"credit").Or("description LIKE ? AND dated > ? AND status = ?", "%"+search+"%",d,"credit").Find(&results)
		}
		if dated == "In the last 7days"{
			d := now.AddDate(0, 0, -7)
			GormDB.Where("suppliername LIKE ? AND dated > ? AND status = ?", "%"+search+"%",d,"credit").Or("code LIKE ? AND dated > ? AND status = ?", "%"+search+"%",d,"credit").Or("title LIKE ? AND dated > ? AND status = ?", "%"+search+"%",d,"credit").Or("description LIKE ? AND dated > ? AND status = ?", "%"+search+"%",d,"credit").Find(&results)
		}
		if dated == "In the last 15day"{
			d := now.AddDate(0, 0, -15)
			GormDB.Where("suppliername LIKE ? AND dated > ? AND status = ?", "%"+search+"%",d,"credit").Or("code LIKE ? AND dated > ? AND status = ?", "%"+search+"%",d,"credit").Or("title LIKE ? AND dated > ? AND status = ?", "%"+search+"%",d,"credit").Or("description LIKE ? AND dated > ? AND status = ?", "%"+search+"%",d,"credit").Find(&results)
		}
		if dated == "In the last 30days"{
			d := now.AddDate(0, 0, -30)
			GormDB.Where("suppliername LIKE ? AND dated > ? AND status = ?", "%"+search+"%",d,"credit").Or("code LIKE ? AND dated > ? AND status = ?", "%"+search+"%",d,"credit").Or("title LIKE ? AND dated > ? AND status = ?", "%"+search+"%",d,"credit").Or("description LIKE ? AND dated > ? AND status = ?", "%"+search+"%",d,"credit").Find(&results)
		}
	}
	if search != "" && dated == ""{
		GormDB.Where("status = ?", "%"+search+"%", "credit").Or("code LIKE ? status = ?", "%"+search+"%", "credit").Or("title LIKE ? status = ?", "%"+search+"%", "credit").Or("description LIKE ? status = ?", "%"+search+"%", "credit").Find(&results)
	}
	if search == "" && dated != "custom"{
		if dated == "In the last 24hrs"{
			d := now.AddDate(0, 0, -1)
			GormDB.Where("dated > ? AND status = ?",d,"credit").Or("dated > ? AND status = ?",d,"credit").Or("dated > ? AND status = ?",d,"credit").Or("dated > ? AND status = ?",d,"credit").Find(&results)
		}
		if dated == "In the last 7days"{
			d := now.AddDate(0, 0, -7)
			GormDB.Where("dated > ? AND status = ?",d,"credit").Or("dated > ? AND status = ?",d,"credit").Or("dated > ? AND status = ?",d,"credit").Or("dated > ? AND status = ?",d,"credit").Find(&results)
		}
		if dated == "In the last 15day"{
			d := now.AddDate(0, 0, -15)
			GormDB.Where("dated > ? AND status = ?",d,"credit").Or("dated > ? AND status = ?",d,"credit").Or("dated > ? AND status = ?",d,"credit").Or("dated > ? AND status = ?",d,"credit").Find(&results)
		}
		if dated == "In the last 30days"{
			d := now.AddDate(0, 0, -30)
			GormDB.Where("dated > ? AND status = ?",d,"credit").Or("dated > ? AND status = ?",d,"credit").Or("dated > ? AND status = ?",d,"credit").Or("dated > ? AND status = ?",d,"credit").Find(&results)
		}
	}
	if search != "" && dated == "custom"{
		start,err := time.Parse(Layout,searchq2)
		if err != nil {
			return nil, httperors.NewNotFoundError("Something went wrong parsing date1!")
		}
		end,err1 := time.Parse(Layout,searchq3)
		if err1 != nil {
			return nil, httperors.NewNotFoundError("Something went wrong parsing date1!")
		}
		GormDB.Where("suppliername LIKE ? AND status = ? AND dated BETWEEN ? AND ?", "%"+search+"%","credit",start, end).Or("code LIKE ? AND status = ? AND dated BETWEEN ? AND ?", "%"+search+"%","credit",start, end).Or("title LIKE ? AND status = ? AND dated BETWEEN ? AND ?", "%"+search+"%","credit",start, end).Or("description LIKE ? AND status = ? AND dated BETWEEN ? AND ?", "%"+search+"%","credit",start, end).Find(&results)
	}

	if search == "" && dated == "custom"{
		start,err := time.Parse(Layout,searchq2)
		if err != nil {
			return nil, httperors.NewNotFoundError("Something went wrong parsing date1!")
		}
		end,err1 := time.Parse(Layout,searchq3)
		if err1 != nil {
			return nil, httperors.NewNotFoundError("Something went wrong parsing date1!")
		}
		GormDB.Where("status = ? AND dated BETWEEN ? AND ?","credit",start, end).Or("status = ? AND dated BETWEEN ? AND ?","credit",start, end).Or("status = ? AND dated BETWEEN ? AND ?","credit",start, end).Or("status = ? AND dated BETWEEN ? AND ?","credit",start, end).Find(&results)
	}
	IndexRepo.DbClose(GormDB)
	return results, nil
}
func (sInvoiceRepo sInvoicerepo) GetCreditNotes(search string) ([]model.SInvoice, *httperors.HttpError) {
	credits := []model.SInvoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Where("name LIKE ?", "%"+search+"%").Or("title LIKE ?", "%"+search+"%").Or("description LIKE ?", "%"+search+"%").Find(&credits)
	if err1 != nil {
		return nil, err1
	}
	return credits, nil
}
func (sInvoiceRepo sInvoicerepo) Update(code string) (string, *httperors.HttpError) {

	sInvoice := model.SInvoice{}
	transactions := []model.STransaction{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}

	asInvoice := model.SInvoice{}

	GormDB.Where("code = ?", code).First(&asInvoice)

	suppliername := asInvoice.Suppliername
	supplier := Supplierrepo.Getsupplier(suppliername)
	now := time.Now()
	sInvoice.SupplierID = supplier.ID
	sInvoice.Dated = now
	sInvoice.Code = code
	sInvoice.Duedate = now.AddDate(0, 0, 1)
	sInvoice.Title = "Credit"
	sInvoice.Status = "credit"
	sInvoice.Description = "Return of goods and services"
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
	debtTransaction := model.CreditTransaction{Code: code, Description: "Goods Returned to supplier", Suppliername: suppliername, Amount: total}

	sInvoice.Discount = discount
	sInvoice.Suppliername = suppliername
	sInvoice.Tax = tax
	sInvoice.Subtotal = (total - tax + discount)
	sInvoice.Total = total
	sInvoice.Cn = true
	trans := model.STransaction{}
	////////////begin transaction/////////////////////
	GormDB.Transaction(func(tx *gorm.DB) error {

		fmt.Println("level 1")
		tx.Create(&sInvoice)

		tx.Transaction(func(tx2 *gorm.DB) error {

			fmt.Println("level 2")
			tx2.Create(&debtTransaction)
			return nil
		})
		for _, c := range transactions {
			product := Productrepo.Productqty(c.Productname)
			remaining := product.Quantity - c.Quantity
			tx.Transaction(func(tx3 *gorm.DB) error {
				fmt.Println("level 3")
				tx3.Model(&product).Where("name = ?", product.Name).Update("quantity", remaining)
				return nil
			})
		}
		for _, t := range transactions {
			tx.Transaction(func(tx4 *gorm.DB) error {
				fmt.Println("level 4")
				tx4.Model(&trans).Where("code = ? AND productname = ? AND total < ?", code, t.Productname, 0).Select("credit", "status").UpdateColumns(model.Transaction{Credit: true, Status: "credit"})
				return nil
			})
		}
		//////////////end of transaction///////////////

		return nil
	})
	IndexRepo.DbClose(GormDB)

	return "item credited successifully", nil
}
func (sInvoiceRepo sInvoicerepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := sInvoiceRepo.SInvoiceUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("sInvoice with that id does not exists!")
	}
	sInvoice := model.SInvoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&sInvoice).Where("id = ?", id).First(&sInvoice)
	GormDB.Delete(sInvoice)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (sInvoiceRepo sInvoicerepo) SInvoiceUserExistByid(id int) bool {
	sInvoice := model.SInvoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&sInvoice, "id =?", id)
	if res.Error != nil {
		return false
	}
	IndexRepo.DbClose(GormDB)
	return true

}
func (sInvoiceRepo sInvoicerepo) SInvoiceExistByCode(code string) bool {
	sInvoice := model.SInvoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}

	res := GormDB.First(&sInvoice, "code =?", code)
	if res.Error != nil {
		return false
	}
	IndexRepo.DbClose(GormDB)
	return true

}
func (sInvoiceRepo sInvoicerepo) GeneCode() (string, *httperors.HttpError) {
	sInvoice := model.SInvoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	err := GormDB.Last(&sInvoice)
	if err.Error != nil {
		var c1 uint = 1
		code := "SupplierInvNo" + strconv.FormatUint(uint64(c1), 10)
		return code, nil
	}
	c1 := sInvoice.ID + 1
	code := "SupplierInvNo" + strconv.FormatUint(uint64(c1), 10)
	IndexRepo.DbClose(GormDB)
	return code, nil

}
func (sInvoiceRepo sInvoicerepo)GetInvoicebyCode(code string) *model.SInvoice {
	invoice := model.SInvoice{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil
	}
	GormDB.Where("code = ? ", code).First(&invoice)
	if invoice.ID == 0 {
	   return nil
	}
	IndexRepo.DbClose(GormDB)
	return &invoice
	
}