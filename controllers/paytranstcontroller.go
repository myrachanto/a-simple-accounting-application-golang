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
 
var (
	PayrectrasanController payrectrasanController = payrectrasanController{}
)
type payrectrasanController struct{ }
/////////controllers/////////////////
func (controller payrectrasanController) Create(c echo.Context) error {
	payrectrasan := &model.Payrectrasan{}
	
	if err := c.Bind(payrectrasan); err != nil {
		httperror := httperors.NewBadRequestError("Invalid json body")
		return c.JSON(httperror.Code, httperror)
	}	
	createdpayrectrasan, err1 := service.PayrectrasanService.Create(payrectrasan)
	if err1 != nil {
		return c.JSON(err1.Code, err1)
	}
	return c.JSON(http.StatusCreated, createdpayrectrasan)
}

func (controller payrectrasanController) Updatepayments(c echo.Context) error {
		
	code := c.FormValue("code")
	status := c.FormValue("status") 
	amount, err := strconv.ParseFloat(c.FormValue("amount"), 64)
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid selling price")
		return c.JSON(httperror.Code, httperror)
	}
	updatedcart, problem := service.PayrectrasanService.Updatepayments(amount,code,status)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, updatedcart)
}
func (controller payrectrasanController) View(c echo.Context) error {
	code, problem := service.PayrectrasanService.View()
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, code)	
}
func (controller payrectrasanController) GetOne(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	fmt.Println(id)
	payrectrasan, problem := service.PayrectrasanService.GetOne(id)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, payrectrasan)	
}

func (controller payrectrasanController) Update(c echo.Context) error {
		
	payrectrasan :=  &model.Payrectrasan{}
	if err := c.Bind(payrectrasan); err != nil {
		httperror := httperors.NewBadRequestError("Invalid json body")
		return c.JSON(httperror.Code, httperror)
	}	
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	updatedpayrectrasan, problem := service.PayrectrasanService.Update(id, payrectrasan)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, updatedpayrectrasan)
}

func (controller payrectrasanController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	success, failure := service.PayrectrasanService.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code, failure)
	}
	return c.JSON(success.Code, success)
		
}
