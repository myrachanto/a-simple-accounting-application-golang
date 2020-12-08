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
 //PaymentformController ...
var (
	PaymentformController paymentformController = paymentformController{}
)
type paymentformController struct{ }
/////////controllers/////////////////
func (controller paymentformController) Create(c echo.Context) error {
	paymentform := &model.Paymentform{}
	
	if err := c.Bind(paymentform); err != nil {
		httperror := httperors.NewBadRequestError("Invalid json body")
		return c.JSON(httperror.Code, httperror)
	}	
	createdpaymentform, err1 := service.Paymentformservice.Create(paymentform)
	if err1 != nil {
		return c.JSON(err1.Code, err1)
	}
	return c.JSON(http.StatusCreated, createdpaymentform)
}
func (controller paymentformController) GetAll(c echo.Context) error {
	search := string(c.QueryParam("q"))
	paymentforms, err3 := service.Paymentformservice.GetAll(search)
	if err3 != nil {
		return c.JSON(err3.Code, err3)
	}
	return c.JSON(http.StatusOK, paymentforms)
} 
func (controller paymentformController) GetOne(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	fmt.Println(id)
	paymentform, problem := service.Paymentformservice.GetOne(id)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, paymentform)	
}

func (controller paymentformController) Update(c echo.Context) error {
		
	paymentform :=  &model.Paymentform{}
	if err := c.Bind(paymentform); err != nil {
		httperror := httperors.NewBadRequestError("Invalid json body")
		return c.JSON(httperror.Code, httperror)
	}	
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	updatedpaymentform, problem := service.Paymentformservice.Update(id, paymentform)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, updatedpaymentform)
}

func (controller paymentformController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	success, failure := service.Paymentformservice.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code, failure)
	}
	return c.JSON(success.Code, success)
		
}
