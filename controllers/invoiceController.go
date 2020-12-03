package controllers

import(
	"fmt"
	"strconv"	
	"net/http"
	"github.com/labstack/echo"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/service"
	"github.com/myrachanto/accounting/support"
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
	invoice.Terms = c.FormValue("terms")
	invoice.Instructions = c.FormValue("instructions")
	createdinvoice, err1 := service.Invoiceservice.Create(invoice)
	if err1 != nil {
		return c.JSON(err1.Code, err1)
	}
	return c.JSON(http.StatusCreated, createdinvoice)
}
func (controller invoiceController) GetAll(c echo.Context) error {
	invoices := []model.Invoice{}
	column := string(c.QueryParam("column"))
	direction := string(c.QueryParam("direction"))
	search_column := string(c.QueryParam("search_column"))
	search_operator := string(c.QueryParam("search_operator"))
	search_query_1 := string(c.QueryParam("search_query_1"))
	search_query_2 := string(c.QueryParam("search_query_2"))
	per_page, err := strconv.Atoi(c.QueryParam("per_page"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid per number")
		return c.JSON(httperror.Code, httperror)
	}
	fmt.Println("------------------------")
	search := &support.Search{Column:column, Direction:direction,Search_column:search_column,Search_operator:search_operator,Search_query_1:search_query_1,Search_query_2:search_query_2,Per_page:per_page}
	
	invoices, err3 := service.Invoiceservice.GetAll(invoices,search)
	if err3 != nil {
		return c.JSON(err3.Code, err3)
	}
	return c.JSON(http.StatusOK, invoices)
}
func (controller invoiceController) GetCredit(c echo.Context) error {
	invoices := []model.Invoice{}
	column := string(c.QueryParam("column"))
	direction := string(c.QueryParam("direction"))
	search_column := string(c.QueryParam("search_column"))
	search_operator := string(c.QueryParam("search_operator"))
	search_query_1 := string(c.QueryParam("search_query_1"))
	search_query_2 := string(c.QueryParam("search_query_2"))
	per_page, err := strconv.Atoi(c.QueryParam("per_page"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid per number")
		return c.JSON(httperror.Code, httperror)
	}
	fmt.Println("------------------------")
	search := &support.Search{Column:column, Direction:direction,Search_column:search_column,Search_operator:search_operator,Search_query_1:search_query_1,Search_query_2:search_query_2,Per_page:per_page}
	
	invoices, err3 := service.Invoiceservice.GetCredit(invoices,search)
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
