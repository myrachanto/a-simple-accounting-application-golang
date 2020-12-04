package repository

import (
	"fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
)
//Cartrepo...
var (
	Cartrepo cartrepo = cartrepo{}
)
//Totals ...
type Totals struct {
	Tax float64
	Discount float64
	Subtotal float64
	Total float64
}
///curtesy to gorm
type cartrepo struct{}
//////////////
////////////TODO user id///////////
/////////////////////////////////////////
func (cartRepo cartrepo) Create(cart *model.Cart) (string, *httperors.HttpError) {
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	// code, err := Invoicerepo.GeneCode()
	// if err != nil {
	// 	return nil, err
	// }
	// cart.Code = code
	
	if cart.Quantity <= 0 {
		return "", httperors.NewNotFoundError("please Add a number bigger than zero!")
	}
	tx := cart.Tax
	dis := cart.Discount
	fmt.Println(tx,dis)

	grossamount := cart.Quantity * cart.SPrice
	taxamount := tx/100 * grossamount
	discountamount := dis/100 * grossamount
	// fmt.Println(grossamount,taxamount,discountamount)
	// fmt.Println(cart)
	cart.Total = grossamount - discountamount + taxamount
	code := cart.Code
	cart.Tax = taxamount
	name := cart.Name
	cart.Discountpercent = dis
	cart.Taxpercent = tx
	cart.Discount = discountamount

	ok := Productrepo.Productexist(name)
	if ok == false {
		return "", httperors.NewNotFoundError("That product does not exist!")
	}
	ok = Invoicerepo.InvoiceExistByCode(code)
	if ok == true {
		return "", httperors.NewNotFoundError("That invoice is already saved!")
	}
	ok = cartRepo.cartproductexist(name,code)
	if ok == true {
		return "", httperors.NewNotFoundError("That product is already saved!")
	}
	if cart.Customername == "" {
		return "", httperors.NewNotFoundError("Please choose a customer name!")
	}
	ok = Customerrepo.customerExistByname(cart.Customername)
	if ok != true {
		return "", httperors.NewNotFoundError("Please choose a customer name!")
	}
	ok = cartRepo.bindcarttocustomer(cart.Customername,code)
	if ok != true {
		return "", httperors.NewNotFoundError("Please select the same customer in this invoice!")
	}
	product := Productrepo.Productqty(name)
	if cart.Quantity > product.Quantity {
		return "", httperors.NewNotFoundError("please that more than we have in stock")
	}
	GormDB.Create(&cart)
	IndexRepo.DbClose(GormDB)
	return "Item added successifully to the cart", nil
}
func (cartRepo cartrepo) View(code string) ([]model.Cart, *httperors.HttpError) {
	mc, e := cartRepo.Getcarts(code)
	if e != nil{
		return nil, e
	}
	return mc, nil
}
func (cartRepo cartrepo) Getcarts(code string) (t []model.Cart, e *httperors.HttpError) {

	cart := model.Cart{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&cart).Where("code = ?", code).Find(&t)
	IndexRepo.DbClose(GormDB)
	
	return t, nil
}
func (cartRepo cartrepo) Getcredits(code string) (t []model.Transaction, err *httperors.HttpError) {
	
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	fmt.Println(code)
	GormDB.Where("code = ? AND credit = ? ", code, false).Find(&t)
	IndexRepo.DbClose(GormDB)
	
	return t, nil
}
func (cartRepo cartrepo) GetcreditsList(code string) (t []model.Transaction, err *httperors.HttpError) {
	
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	fmt.Println(code)
	GormDB.Where("code = ? AND credit = ? AND status = ? ", code, true, "credit").Find(&t)
	IndexRepo.DbClose(GormDB)
	
	return t, nil
}
func (cartRepo cartrepo) GetOne(id int) (*model.Cart, *httperors.HttpError) {
	ok := cartRepo.cartUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("cart with that id does not exists!")
	}
	cart := model.Cart{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&cart).Where("id = ?", id).First(&cart)
	IndexRepo.DbClose(GormDB)
	
	return &cart, nil
}
func (cartRepo cartrepo) All() (t []model.Cart, r *httperors.HttpError) {

	cart := model.Cart{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&cart).Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (cartRepo cartrepo) GetAll(carts []model.Cart) ([]model.Cart, *httperors.HttpError) {
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	cart := model.Cart{}
	GormDB.Model(&cart).Find(&carts)
	
	IndexRepo.DbClose(GormDB)
	if len(carts) == 0 {
		return nil, httperors.NewNotFoundError("No results found!")
	}
	return carts, nil
}

func (cartRepo cartrepo) Update(qty float64, name,code string) (string, *httperors.HttpError) {
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
	tr := model.Transaction{}
	trps := []model.Transaction{}
	trpsc := []model.Transaction{}
	

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

	transact := model.Transaction{}
	GormDB.Where("Productname = ? AND code = ? AND credit = ? AND status = ?" , name, code, false, "invoice").First(&transact)
	fmt.Println(transact)
	// disco = transact.Discount

	// grossamount := (qty * transact.Price)
	var grossamount float64 = -(qty * transact.Price)
	var discountamount float64 = qty * (transact.Discount/transact.Quantity)
	var taxamount float64 = -(qty * (transact.Tax/transact.Quantity))
	// cart.Total = grossamount - discountamount + taxamount
	 var tot float64 = grossamount + discountamount + taxamount
	fmt.Println(grossamount, discountamount,taxamount, qty ,">>>>>>>>>>>>>>>>>>>>>>>>>>>")
	transaction := model.Transaction{}
	transaction.Total = tot
	transaction.Subtotal = grossamount
	transaction.Tax = taxamount
	transaction.Discount = discountamount
	transaction.Code = code
	transaction.Productname = name
	transaction.Title = "credit invoice quantity"
	transaction.Quantity = qty
	transaction.Price = transact.Price
	transaction.Credit = false
	transaction.Status = "pending"
	fmt.Println(transaction)
	GormDB.Create(&transaction)
	IndexRepo.DbClose(GormDB) 
	return "Item added successifully to the cart", nil
}
func (cartRepo cartrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := cartRepo.cartUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("cart with that id does not exists!") 
	}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}

	cart := model.Cart{}
	GormDB.Where("id = ?", id).Delete(&cart)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (cartRepo cartrepo) DeleteAll(code string) (*httperors.HttpSuccess, *httperors.HttpError) {
	transaction := model.Transaction{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1 
	}
	GormDB.Where("code = ? AND credit = ? AND status = ?", code, false, "pending").Delete(&transaction)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (cartRepo cartrepo)cartUserExistByid(id int) bool {
	cart := model.Cart{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&cart, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (cartRepo cartrepo)customerexist(name string) bool {
	cart := model.Cart{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	GormDB.Where("customername = ? ", name).First(&cart)
	if cart.Name == "" {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (cartRepo cartrepo)cartproductexist(name, code string) bool {
	cart := model.Cart{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	GormDB.Where("name = ? AND code = ?", name, code).First(&cart)
	if cart.Name == "" {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}

func (cartRepo cartrepo)Getcustomerwithcode(code string) *model.Cart {
	cart := model.Cart{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil
	}
	GormDB.Where("code = ?", code).First(&cart)
	if cart.Name == "" {
	   return nil
	}
	IndexRepo.DbClose(GormDB)
	return &cart
	
}
func (cartRepo cartrepo)bindcarttocustomer(name, code string) bool {
	// cart := model.Cart{}
	carts := []model.Cart{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	GormDB.Where("code = ?", code).Find(&carts)
	for _, g := range carts {
		if g.Customername != name {
			return false
		}
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (cartRepo cartrepo)SumTotal(code string) (Total *Totals, err *httperors.HttpError) {
	carts := []model.Cart{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil,httperors.NewNotFoundError("db connection failed!")
	}
	GormDB.Where("code = ?", code).Find(&carts)
	IndexRepo.DbClose(GormDB)

	var	tax, disc, subt,tot float64 = 0,0,0,0
	for _, t := range carts {
		tax += t.Tax
		disc += t.Discount
		subt += t.Subtotal
		tot += t.Total
	}
	tots := &Totals{Tax:tax,Discount:disc,Subtotal:subt,Total:tot}
	return tots,nil
}
func (cartRepo cartrepo)CarttoTransaction(code string) (tr []model.Transaction, err *httperors.HttpError) {
	carts := []model.Cart{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil,httperors.NewNotFoundError("db connection failed!")
	}
	GormDB.Where("code = ?", code).Find(&carts)
	IndexRepo.DbClose(GormDB)
	for _, c := range carts {
		trans := model.Transaction{Productname :c.Name, Title:"Product sale", Quantity: c.Quantity, Price: c.SPrice,Tax:c.Tax, Code:code, Subtotal:c.Subtotal, Discount:c.Discount,Total:c.Total}
		tr = append(tr, trans)
	}
	return tr,nil
}
func (cartRepo cartrepo)Updateproductqty(code string) ([]model.Cart, *httperors.HttpError) {
	carts := []model.Cart{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil,httperors.NewNotFoundError("db connection failed!")
	}
	GormDB.Where("code = ?", code).Find(&carts)

	IndexRepo.DbClose(GormDB)
	return carts, nil
}
