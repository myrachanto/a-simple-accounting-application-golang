package repository

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/support"
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
	sInvoice.Title = "sales"
	sInvoice.Cn = false
	sInvoice.Status = "invoice"
	sInvoice.Description = "Sale of goods and services"

	transactions, e := Scartrepo.ScarttoTransaction(code)
	if e != nil {
		return "", e
	}
	if sInvoice.Suppliername == "undefined" && sInvoice.Suppliername == "" {
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
				tx4.Model(&product).Where("name = ?", product.Name).Update("quantity", remaining)
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
func (sInvoiceRepo sInvoicerepo) PaidsInvoices() (t []model.SInvoice, r *httperors.HttpError) {

	sInvoice := model.SInvoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&sInvoice).Where("status = ?", "paid").Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

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
func (sInvoiceRepo sInvoicerepo) SupplierCreditsbycode(code string) (t []model.SInvoice, r *httperors.HttpError) {

	sInvoice := model.SInvoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&sInvoice).Where("suppliercode = ? AND status = ?", code, "credit").Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

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

func (sInvoiceRepo sInvoicerepo) SuppliersInvoicebycode(code string) (t []model.SInvoice, r *httperors.HttpError) {

	sInvoice := model.SInvoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&sInvoice).Where("suppliercode = ? AND status = ?", code, "invoice").Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (sInvoiceRepo sInvoicerepo) GetAll(sInvoices []model.SInvoice, search *support.Search) ([]model.SInvoice, *httperors.HttpError) {

	results, err1 := sInvoiceRepo.Search(search, sInvoices)
	if err1 != nil {
		return nil, err1
	}
	return results, nil
}
func (sInvoiceRepo sInvoicerepo) GetCredit(sInvoices []model.SInvoice, search *support.Search) ([]model.SInvoice, *httperors.HttpError) {

	results, err1 := sInvoiceRepo.Search(search, sInvoices)
	if err1 != nil {
		return nil, err1
	}
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
func (sInvoiceRepo sInvoicerepo) Search(Ser *support.Search, sInvoices []model.SInvoice) ([]model.SInvoice, *httperors.HttpError) {
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	sInvoice := model.SInvoice{}
	switch Ser.Search_operator {
	case "all":
		GormDB.Model(&sInvoice).Order(Ser.Column + " " + Ser.Direction).Find(&sInvoices)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////

		break
	case "equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column + " " + Ser.Direction).Find(&sInvoices)

		break
	case "not_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column + " " + Ser.Direction).Find(&sInvoices)

		break
	case "less_than":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column + " " + Ser.Direction).Find(&sInvoices)

		break
	case "greater_than":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column + " " + Ser.Direction).Find(&sInvoices)

		break
	case "less_than_or_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column + " " + Ser.Direction).Find(&sInvoices)

		break
	case "greater_than_ro_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column + " " + Ser.Direction).Find(&sInvoices)

		break
	case "in":
		// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1, ",")
		fmt.Println(s)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column + " " + Ser.Direction).Find(&sInvoices)

		break
	case "not_in":
		//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1, ",")
		GormDB.Not(Ser.Search_column, s).Order(Ser.Column + " " + Ser.Direction).Find(&sInvoices)

		// break;
	case "like":
		// fmt.Println(Ser.Search_query_1)
		if Ser.Search_query_1 == "all" {
			//db.Order("name DESC")
			GormDB.Order(Ser.Column + " " + Ser.Direction).Find(&sInvoices)

		} else {

			GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column + " " + Ser.Direction).Find(&sInvoices)
		}
		break
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column + " " + Ser.Direction).Find(&sInvoices)

		break
	default:
		return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)

	return sInvoices, nil
}
