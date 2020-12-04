package controllers

import(
	"fmt"
	"strconv"	
	"net/http"
	"time"
	"github.com/labstack/echo"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/service"
)
 //PaymentController ...
var ( 
	PaymentController paymentController =  paymentController{}
)
type paymentController struct{ }
/////////controllers/////////////////
func (controller paymentController) Create(c echo.Context) error {
	payment := &model.Payment{}
	payment.SupplierName = c.FormValue("suppliername")
	payment.Status = c.FormValue("status")
	payment.Description = c.FormValue("description")
	payment.Type = c.FormValue("type")
	payment.Code = c.FormValue("code")
	d := c.FormValue("clearancedate")
	fmt.Println(d)

	s, err := strconv.ParseFloat(c.FormValue("amount"), 64)
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid selling price")
		return c.JSON(httperror.Code, httperror)
	}
	payment.Amount = s
  t, er := time.Parse(layoutISO, d)
	if er != nil {
		httperror := httperors.NewBadRequestError("Invalid Date")
		return c.JSON(httperror.Code, httperror)
	}
	payment.ClearanceDate = t
		createdpayment, err1 := service.Paymentservice.Create(payment)
		if err1 != nil {
			return c.JSON(err1.Code, err1)
		}
	return c.JSON(http.StatusCreated, createdpayment)
}

func (controller paymentController) ViewReport(c echo.Context) error {
	options, problem := service.Paymentservice.ViewReport()
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, options)	
}
func (controller paymentController) Updatepayments(c echo.Context) error {
		
	code := c.FormValue("code")
	status := c.FormValue("status")
	updatedcart, problem := service.Paymentservice.Updatepayments(code,status)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, updatedcart)
}
func (controller paymentController) View(c echo.Context) error {
	code, problem := service.Paymentservice.View()
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, code)	
}
func (controller paymentController) GetAll(c echo.Context) error {
	
	payments, err3 := service.Paymentservice.GetAll()
	if err3 != nil {
		return c.JSON(err3.Code, err3)
	}
	return c.JSON(http.StatusOK, payments)
} 
func (controller paymentController) GetOne(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	fmt.Println(id)
	payment, problem := service.Paymentservice.GetOne(id)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, payment)	
}

func (controller paymentController) Update(c echo.Context) error {
		
	payment :=  &model.Payment{}
	if err := c.Bind(payment); err != nil {
		httperror := httperors.NewBadRequestError("Invalid json body")
		return c.JSON(httperror.Code, httperror)
	}	
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	updatedpayment, problem := service.Paymentservice.Update(id, payment)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, updatedpayment)
}

func (controller paymentController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	success, failure := service.Paymentservice.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code, failure)
	}
	return c.JSON(success.Code, success)
		
}
