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
 //CartController ..
var (
	CartController cartController = cartController{}
)

type cartController struct{ }
/////////controllers/////////////////
func (controller cartController) Create(c echo.Context) error {
	cart := &model.Cart{}
	cart.Customername = c.FormValue("customername")
	cart.Name = c.FormValue("name")
	cart.Code = c.FormValue("code")
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
	cart.Quantity = q
	cart.SPrice = s
	cart.Tax =t
	cart.Discount = d
	createdcart, err1 := service.Cartservice.Create(cart)
	if err1 != nil {
		return c.JSON(err1.Code, err1)
	}
	return c.JSON(http.StatusCreated, createdcart)
}

func (controller cartController) View(c echo.Context) error {
	code := c.Param("code")
	options, problem := service.Cartservice.View(code)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, options)	
}
func (controller cartController) Getcredits(c echo.Context) error {
	code := c.Param("code")
	options, problem := service.Cartservice.Getcredits(code)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, options)	 
}
func (controller cartController) GetcreditsList(c echo.Context) error {
	code := c.Param("code")
	options, problem := service.Cartservice.GetcreditsList(code)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, options)	
}

func (controller cartController) GetOne(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	fmt.Println(id)
	cart, problem := service.Cartservice.GetOne(id)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, cart)	
}

func (controller cartController) Updatetrans(c echo.Context) error {
		
	name := c.FormValue("name")
	code := c.FormValue("code")
	qty, err := strconv.ParseFloat(c.FormValue("quantity"), 64)
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid quantity")
		return c.JSON(httperror.Code, httperror)
	}
	fmt.Println(qty)
	updatedcart, problem := service.Cartservice.Update( qty, name, code)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, updatedcart)
}

func (controller cartController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	
	success, failure := service.Cartservice.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code, failure)
	}
	return c.JSON(success.Code, success)
		
}
func (controller cartController) DeleteAll(c echo.Context) error {
	code := c.Param("code")
	success, failure := service.Cartservice.DeleteALL(code)
	if failure != nil {
		return c.JSON(failure.Code, failure)
	}
	return c.JSON(success.Code, success)
		
}

