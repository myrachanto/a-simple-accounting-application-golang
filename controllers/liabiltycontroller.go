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
 //LiabilityController controlls the liability
var (
	LiabilityController liabilityController = liabilityController{}
)
type liabilityController struct{ }
/////////controllers/////////////////
func (controller liabilityController) Create(c echo.Context) error {
	liability := &model.Liability{}
	
	if err := c.Bind(liability); err != nil {
		httperror := httperors.NewBadRequestError("Invalid json body")
		return c.JSON(httperror.Code, httperror)
	}	
	createdliability, err1 := service.Liabilityservice.Create(liability)
	if err1 != nil {
		return c.JSON(err1.Code, err1)
	}
	return c.JSON(http.StatusCreated, createdliability)
}
func (controller liabilityController) GetAll(c echo.Context) error {
	
	search := string(c.QueryParam("q"))
	
	liabilitys, err3 := service.Liabilityservice.GetAll(search)
	if err3 != nil {
		return c.JSON(err3.Code, err3)
	}
	return c.JSON(http.StatusOK, liabilitys)
} 
func (controller liabilityController) GetOne(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	fmt.Println(id)
	liability, problem := service.Liabilityservice.GetOne(id)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, liability)	
}

func (controller liabilityController) Update(c echo.Context) error {
		
	liability :=  &model.Liability{}
	if err := c.Bind(liability); err != nil {
		httperror := httperors.NewBadRequestError("Invalid json body")
		return c.JSON(httperror.Code, httperror)
	}	
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	updatedliability, problem := service.Liabilityservice.Update(id, liability)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, updatedliability)
}

func (controller liabilityController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	success, failure := service.Liabilityservice.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code, failure)
	}
	return c.JSON(success.Code, success)
		
}
