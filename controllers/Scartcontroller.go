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
 //scartController ..
var (
	ScartController scartController = scartController{}
)

type scartController struct{ }
/////////controllers/////////////////
func (controller scartController) Create(c echo.Context) error {
	scart := &model.Scart{}
	scart.Suppliername = c.FormValue("suppliername")
	scart.Name = c.FormValue("name")
	scart.Code = c.FormValue("code")
	scart.Suppliercode = c.FormValue("suppliercode")
	scart.Usercode = c.FormValue("usercode")
	q, err := strconv.ParseFloat(c.FormValue("quantity"), 64)
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid quantity")
		return c.JSON(httperror.Code, httperror)
	}
	s, err := strconv.ParseFloat(c.FormValue("sprice"), 64)
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid selling price")
		return c.JSON(httperror.Code, httperror)
	}
	t, err := strconv.ParseFloat(c.FormValue("tax"), 64)
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid tax ")
		return c.JSON(httperror.Code, httperror)
	}
	d, err := strconv.ParseFloat(c.FormValue("discount"), 64)
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid discount")
		return c.JSON(httperror.Code, httperror)
	}
	scart.Quantity = q
	scart.BPrice = s
	scart.Tax =t
	scart.Discount = d
	createdscart, err1 := service.Scartservice.Create(scart)
	if err1 != nil {
		return c.JSON(err1.Code, err1)
	}
	return c.JSON(http.StatusCreated, createdscart)
}

func (controller scartController) View(c echo.Context) error {
	code := c.Param("code")
	options, problem := service.Scartservice.View(code)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, options)	
}
func (controller scartController) Getcredits(c echo.Context) error {
	code := c.Param("code")
	options, problem := service.Scartservice.Getcredits(code)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, options)	 
}
func (controller scartController) GetcreditsList(c echo.Context) error {
	code := c.Param("code")
	options, problem := service.Scartservice.GetcreditsList(code)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, options)	
}

func (controller scartController) GetOne(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	fmt.Println(id)
	scart, problem := service.Scartservice.GetOne(id)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, scart)	
}

func (controller scartController) Updatetrans(c echo.Context) error {
		
	name := c.FormValue("name")
	code := c.FormValue("code")
	qty, err := strconv.ParseFloat(c.FormValue("quantity"), 64)
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid quantity")
		return c.JSON(httperror.Code, httperror)
	}
	fmt.Println(qty)
	updatedscart, problem := service.Scartservice.Update( qty, name, code)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, updatedscart)
}

func (controller scartController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	
	success, failure := service.Scartservice.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code, failure)
	}
	return c.JSON(success.Code, success)
		
}
func (controller scartController) DeleteAll(c echo.Context) error {
	code := c.Param("code")
	success, failure := service.Scartservice.DeleteALL(code)
	if failure != nil {
		return c.JSON(failure.Code, failure)
	}
	return c.JSON(success.Code, success)
		
}

