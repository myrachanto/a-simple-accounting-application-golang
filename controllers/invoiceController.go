package controllers

import(
	"fmt"
	"strconv"	
	"net/http"
	"github.com/labstack/echo"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/service"
)
 //InvoiceController controller
var (
	InvoiceController invoiceController = invoiceController{}
)
type invoiceController struct{ }
/////////controllers/////////////////
func (controller invoiceController) CreateCart(c echo.Context) error {
	cart := &model.Cart{}
	cart.Name = c.FormValue("name")
	cart.Code = c.FormValue("code")
	cart.Usercode = c.FormValue("usercode")
	cart.Customercode = c.FormValue("customercode")
	sprice := c.FormValue("sprice")
	quantity := c.FormValue("quantity")
	fmt.Println(sprice, quantity, ">>>>>>>>>>>")
	b, err := strconv.ParseFloat(c.FormValue("quantity"), 64)
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid buying price")
		return c.JSON(httperror.Code, httperror)
	}
	s, err := strconv.ParseFloat(c.FormValue("sprice"), 64)
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid selling price")
		return c.JSON(httperror.Code, httperror)
	}
	cart.Quantity = b
	cart.SPrice = s
	fmt.Println(cart)
	createdcart, err1 := service.Cartservice.Create(cart)
	if err1 != nil {
		return c.JSON(err1.Code, err1)
	}
	return c.JSON(http.StatusCreated, createdcart)
}
func (controller invoiceController) Create(c echo.Context) error {
	invoice := &model.Invoice{}
	
	invoice.Customername = c.FormValue("customername")
	invoice.Code = c.FormValue("code")
	invoice.Usercode = c.FormValue("usercode")
	invoice.Customercode = c.FormValue("customercode")
	invoice.Terms = c.FormValue("terms")
	invoice.Instructions = c.FormValue("instructions")
	createdinvoice, err1 := service.Invoiceservice.Create(invoice)
	if err1 != nil {
		return c.JSON(err1.Code, err1)
	}
	return c.JSON(http.StatusCreated, createdinvoice)
}
func (controller invoiceController) GetAll(c echo.Context) error {
	search := string(c.QueryParam("search"))
	dated := string(c.QueryParam("dated"))
	searchq2 := string(c.QueryParam("searchq2"))
	searchq3 := string(c.QueryParam("searchq3"))
	
	invoices, err3 := service.Invoiceservice.GetAll(search,dated,searchq2,searchq3)
	if err3 != nil {
		return c.JSON(err3.Code, err3)
	}
	return c.JSON(http.StatusOK, invoices)
}
func (controller invoiceController) GetCredit(c echo.Context) error {
	search := string(c.QueryParam("search"))
	dated := string(c.QueryParam("dated"))
	searchq2 := string(c.QueryParam("searchq2"))
	searchq3 := string(c.QueryParam("searchq3"))
	invoices, err3 := service.Invoiceservice.GetCredit(search,dated,searchq2,searchq3)
	if err3 != nil {
		return c.JSON(err3.Code, err3)
	}
	return c.JSON(http.StatusOK, invoices)
}
 func (controller invoiceController) View(c echo.Context) error {
	code, problem := service.Invoiceservice.View()
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, code)	
}
func (controller invoiceController) GetOne(c echo.Context) error {
	code := c.Param("id")
	
	invoice, problem := service.Invoiceservice.GetOne(code)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, invoice)	
}

func (controller invoiceController) Credit(c echo.Context) error {
		
	
	code := c.FormValue("code")
	updatedinvoice, problem := service.Invoiceservice.Update(code)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, updatedinvoice)
}

func (controller invoiceController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	success, failure := service.Invoiceservice.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code, failure)
	}
	return c.JSON(success.Code, success)
		
}
