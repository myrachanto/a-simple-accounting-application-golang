package repository

import (
	"fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
)
//Scartrepo...supplier cart repository
var (
	Scartrepo scartrepo = scartrepo{}
)
//STotals ... structure
type STotals struct {
	Tax float64
	Discount float64
	Subtotal float64
	Total float64
}
///curtesy to gorm
type scartrepo struct{}
//////////////
////////////TODO user id///////////
/////////////////////////////////////////
func (scartRepo scartrepo) Create(scart *model.Scart) (string, *httperors.HttpError) {
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	// code, err := Invoicerepo.GeneCode()
	// if err != nil {
	// 	return nil, err
	// }
	// scart.Code = code
	
	if scart.Quantity <= 0 {
		return "", httperors.NewNotFoundError("please Add a number bigger than zero!")
	}
	tx := scart.Tax
	dis := scart.Discount
	fmt.Println(tx,dis)

	grossamount := scart.Quantity * scart.BPrice
	taxamount := tx/100 * grossamount
	discountamount := dis/100 * grossamount
	// fmt.Println(grossamount,taxamount,discountamount)
	// fmt.Println(scart)
	scart.Total = grossamount - discountamount + taxamount
	code := scart.Code
	scart.Tax = taxamount
	name := scart.Name
	scart.Discountpercent = dis
	scart.Taxpercent = tx
	scart.Discount = discountamount

	ok := Productrepo.Productexist(name)
	if ok == false {
		return "", httperors.NewNotFoundError("That product does not exist!")
	}
	ok = Invoicerepo.InvoiceExistByCode(code)
	if ok == true {
		return "", httperors.NewNotFoundError("That invoice is already saved!")
	}
	ok = scartRepo.scartproductexist(name,code)
	if ok == true {
		return "", httperors.NewNotFoundError("That product is already saved!")
	} 
	if scart.Suppliername == "" {
		return "", httperors.NewNotFoundError("Please choose a Supplier name!")
	}
	ok = Supplierrepo.SupplierExistByname(scart.Suppliername)
	if ok != true {
		return "", httperors.NewNotFoundError("Please choose a Supplier name!")
	}
	ok = Scartrepo.bindscarttosupplier(scart.Suppliername,code)
	if ok != true {
		return "", httperors.NewNotFoundError("Please select the same Supplier in this invoice!")
	}
	// product := Productrepo.Productqty(name)
	// if scart.Quantity > product.Quantity {
	// 	return "", httperors.NewNotFoundError("please that more than we have in stock")
	// }
	fmt.Println(scart)
	GormDB.Create(&scart)
	IndexRepo.DbClose(GormDB)
	return "Item added successifully to the scart", nil
}
func (scartRepo scartrepo) View(code string) ([]model.Scart, *httperors.HttpError) {
	mc, e := scartRepo.Getscarts(code)
	if e != nil{
		return nil, e
	}
	return mc, nil
}
func (scartRepo scartrepo) Getscarts(code string) (t []model.Scart, e *httperors.HttpError) {

	scart := model.Scart{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&scart).Where("code = ?", code).Find(&t)
	IndexRepo.DbClose(GormDB)
	
	return t, nil
}
func (scartRepo scartrepo) Getcredits(code string) (t []model.Transaction, err *httperors.HttpError) {
	
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	fmt.Println(code)
	GormDB.Where("code = ? AND credit = ? ", code, false).Find(&t)
	IndexRepo.DbClose(GormDB)
	
	return t, nil
}
func (scartRepo scartrepo) GetcreditsList(code string) (t []model.Transaction, err *httperors.HttpError) {
	
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	fmt.Println(code)
	GormDB.Where("code = ? AND credit = ? AND status = ? ", code, true, "credit").Find(&t)
	IndexRepo.DbClose(GormDB)
	
	return t, nil
}
func (scartRepo scartrepo) GetOne(id int) (*model.Scart, *httperors.HttpError) {
	ok := scartRepo.scartUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("scart with that id does not exists!")
	}
	scart := model.Scart{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&scart).Where("id = ?", id).First(&scart)
	IndexRepo.DbClose(GormDB)
	
	return &scart, nil
}
func (scartRepo scartrepo) All() (t []model.Scart, r *httperors.HttpError) {

	scart := model.Scart{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&scart).Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (scartRepo scartrepo) GetAll(scarts []model.Scart) ([]model.Scart, *httperors.HttpError) {
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	scart := model.Scart{}
	GormDB.Model(&scart).Find(&scarts)
	
	IndexRepo.DbClose(GormDB)
	if len(scarts) == 0 {
		return nil, httperors.NewNotFoundError("No results found!")
	}
	return scarts, nil
}

func (scartRepo scartrepo) Update(qty float64, name,code string) (string, *httperors.HttpError) {
	// ok := Invoicerepo.invoiceUserExistByid(id)
	// if !ok {
	// 	return "", httperors.NewNotFoundError("invoice with that id does not exists!")
	// }
	
	ok := Invoicerepo.InvoiceExistByCode(code)
	if ok == false {
		return "", httperors.NewNotFoundError("Something went wrong with the invoice crediting!")
	}
	if qty <= 0 {
		return "", httperors.NewNotFoundError("please Add a number bigger than zero!")
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	tr := model.STransaction{}
	trps := []model.STransaction{}
	trpsc := []model.STransaction{}
	

	GormDB.Where("Productname = ? AND code = ? AND credit = ? AND status = ?", name, code, false, "invoice").First(&tr)
	if qty > tr.Quantity {  
		return "", httperors.NewNotFoundError("please the quantity is bigger than what the invoice has!")
	}
	GormDB.Where("Productname = ? AND code = ? AND credit = ? AND status = ?", name, code, true, "credit").Find(&trpsc)
	var qtysc float64 = 0
	for _, tsc := range trps {
		qtysc += tsc.Quantity
	}
	if qty > (tr.Quantity - qtysc) {
		return "", httperors.NewNotFoundError("please the quantity is bigger than what the invoice  has and pending credit has !")
	}
	//evaluate the effect of an array
	GormDB.Where("Productname = ? AND code = ? AND credit = ? AND status = ?", name, code, false, "pending").Find(&trps)
	var qtys float64 = 0
	for _, ts := range trps {
		qtys += ts.Quantity
	}
	if qty > (tr.Quantity - qtysc - qtys) {
		return "", httperors.NewNotFoundError("please the quantity is bigger than what the invoice  has, credit and pending credit has !")
	}

	transact := model.STransaction{}
	GormDB.Where("Productname = ? AND code = ? AND credit = ? AND status = ?" , name, code, false, "invoice").First(&transact)
	fmt.Println(transact)
	// disco = transact.Discount

	// grossamount := (qty * transact.Price)
	var grossamount float64 = -(qty * transact.Price)
	var discountamount float64 = qty * (transact.Discount/transact.Quantity)
	var taxamount float64 = -(qty * (transact.Tax/transact.Quantity))
	// scart.Total = grossamount - discountamount + taxamount
	 var tot float64 = grossamount + discountamount + taxamount
	fmt.Println(grossamount, discountamount,taxamount, qty ,">>>>>>>>>>>>>>>>>>>>>>>>>>>")
	transaction := model.STransaction{}
	transaction.Total = tot
	transaction.Subtotal = grossamount
	transaction.Tax = taxamount
	transaction.Discount = discountamount
	transaction.Code = code
	transaction.Productname = name
	transaction.Title = "Goods returned "
	transaction.Quantity = qty
	transaction.Price = transact.Price
	transaction.Credit = false
	transaction.Status = "pending"
	fmt.Println(transaction)
	GormDB.Create(&transaction)
	IndexRepo.DbClose(GormDB) 
	return "Item added successifully to the scart", nil
}
func (scartRepo scartrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := scartRepo.scartUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("scart with that id does not exists!") 
	}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}

	scart := model.Scart{}
	GormDB.Where("id = ?", id).Delete(&scart)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (scartRepo scartrepo) DeleteAll(code string) (*httperors.HttpSuccess, *httperors.HttpError) {
	transaction := model.STransaction{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1 
	}
	GormDB.Where("code = ? AND credit = ? AND status = ?", code, false, "pending").Delete(&transaction)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (scartRepo scartrepo)scartUserExistByid(id int) bool {
	scart := model.Scart{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&scart, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (scartRepo scartrepo)customerexist(name string) bool {
	scart := model.Scart{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	GormDB.Where("Suppliername = ? ", name).First(&scart)
	if scart.Name == "" {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (scartRepo scartrepo)scartproductexist(name, code string) bool {
	scart := model.Scart{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	GormDB.Where("name = ? AND code = ?", name, code).First(&scart)
	if scart.Name == "" {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}

func (scartRepo scartrepo)Getsupplierwithcode(code string) *model.Scart {
	scart := model.Scart{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil
	}
	GormDB.Where("code = ?", code).First(&scart)
	if scart.Name == "" {
	   return nil
	}
	IndexRepo.DbClose(GormDB)
	return &scart
	
}
func (scartRepo scartrepo)bindscarttosupplier(name, code string) bool {
	// scart := model.Scart{}
	scarts := []model.Scart{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	GormDB.Where("code = ?", code).Find(&scarts)
	for _, g := range scarts {
		if g.Suppliername != name {
			return false
		}
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (scartRepo scartrepo)SumTotal(code string) (Total *Totals, err *httperors.HttpError) {
	scarts := []model.Scart{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil,httperors.NewNotFoundError("db connection failed!")
	}
	GormDB.Where("code = ?", code).Find(&scarts)
	IndexRepo.DbClose(GormDB)

	var	tax, disc, subt,tot float64 = 0,0,0,0
	for _, t := range scarts {
		tax += t.Tax
		disc += t.Discount
		subt += t.Subtotal
		tot += t.Total
	}
	tots := &Totals{Tax:tax,Discount:disc,Subtotal:subt,Total:tot}
	return tots,nil
}
func (scartRepo scartrepo)ScarttoTransaction(code string) (tr []model.STransaction, err *httperors.HttpError) {
	scarts := []model.Scart{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil,httperors.NewNotFoundError("db connection failed!")
	}
	GormDB.Where("code = ?", code).Find(&scarts)
	IndexRepo.DbClose(GormDB)
	for _, c := range scarts {
		trans := model.STransaction{Productname :c.Name, Title:"Product sale", Quantity: c.Quantity, Price: c.BPrice,Tax:c.Tax, Code:code, Subtotal:c.Subtotal, Discount:c.Discount,Total:c.Total}
		tr = append(tr, trans)
	}
	return tr,nil
}
func (scartRepo scartrepo)Updateproductqty(code string) ([]model.Scart, *httperors.HttpError) {
	scarts := []model.Scart{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil,httperors.NewNotFoundError("db connection failed!")
	}
	GormDB.Where("code = ?", code).Find(&scarts)

	IndexRepo.DbClose(GormDB)
	return scarts, nil
}
