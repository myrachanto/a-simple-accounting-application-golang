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
 //SInvoiceController controller
var (
	SInvoiceController sInvoiceController = sInvoiceController{}
)
type sInvoiceController struct{ }
/////////controllers/////////////////
func (controller sInvoiceController) Createscart(c echo.Context) error {
	scart := &model.Scart{}
	scart.Name = c.FormValue("name")
	scart.Code = c.FormValue("code")
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
	scart.Quantity = b
	scart.BPrice = s
	fmt.Println(scart)
	createdscart, err1 := service.Scartservice.Create(scart)
	if err1 != nil {
		return c.JSON(err1.Code, err1)
	}
	return c.JSON(http.StatusCreated, createdscart)
}
func (controller sInvoiceController) Create(c echo.Context) error {
	sInvoice := &model.SInvoice{}
	
	sInvoice.Suppliername = c.FormValue("suppliername")
	sInvoice.Code = c.FormValue("code")
	sInvoice.Terms = c.FormValue("terms")
	sInvoice.Instructions = c.FormValue("instructions")
	sInvoice.Usercode = c.FormValue("usercode")
	sInvoice.Suppliercode = c.FormValue("suppliercode")
	createdsInvoice, err1 := service.SInvoiceservice.Create(sInvoice)
	if err1 != nil {
		return c.JSON(err1.Code, err1)
	}
	return c.JSON(http.StatusCreated, createdsInvoice)
}
func (controller sInvoiceController) GetAll(c echo.Context) error {
	search := string(c.QueryParam("search"))
	dated := string(c.QueryParam("dated"))
	searchq2 := string(c.QueryParam("searchq2"))
	searchq3 := string(c.QueryParam("searchq3"))

	sInvoices, err3 := service.SInvoiceservice.GetAll(search,dated,searchq2,searchq3)
	if err3 != nil {
		return c.JSON(err3.Code, err3)
	}
	return c.JSON(http.StatusOK, sInvoices)
}
func (controller sInvoiceController) GetCredit(c echo.Context) error {
	search := string(c.QueryParam("search"))
	dated := string(c.QueryParam("dated"))
	searchq2 := string(c.QueryParam("searchq2"))
	searchq3 := string(c.QueryParam("searchq3"))
	
	sInvoices, err3 := service.SInvoiceservice.GetCredit(search,dated,searchq2,searchq3)
	if err3 != nil {
		return c.JSON(err3.Code, err3)
	}
	return c.JSON(http.StatusOK, sInvoices)
}
 func (controller sInvoiceController) View(c echo.Context) error {
	code, problem := service.SInvoiceservice.View()
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, code)	
}
func (controller sInvoiceController) GetOne(c echo.Context) error {
	code := c.Param("id")
	
	sInvoice, problem := service.SInvoiceservice.GetOne(code)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, sInvoice)	
}

func (controller sInvoiceController) Credit(c echo.Context) error {
		
	
	code := c.FormValue("code")
	updatedsInvoice, problem := service.SInvoiceservice.Update(code)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, updatedsInvoice)
}

func (controller sInvoiceController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	success, failure := service.SInvoiceservice.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code, failure)
	}
	return c.JSON(success.Code, success)
		
}
