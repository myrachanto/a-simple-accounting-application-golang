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
//ExpencetrasanController ...
var (
	ExpencetrasanController expencetrasanController = expencetrasanController{}
)
type expencetrasanController struct{ }
/////////controllers/////////////////
func (controller expencetrasanController) Create(c echo.Context) error {
	expencetrasan := &model.Expencetrasan{}
	
	expencetrasan.Name = c.FormValue("name")
	expencetrasan.Title = c.FormValue("title")
	expencetrasan.Description = c.FormValue("description")
	expencetrasan.Code = c.FormValue("code")
	expencetrasan.Usercode = c.FormValue("usercode")

	s, err := strconv.ParseFloat(c.FormValue("amount"), 64)
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid selling price")
		return c.JSON(httperror.Code, httperror)
	}
	expencetrasan.Amount = s
	createdexpencetrasan, err1 := service.ExpencetrasanService.Create(expencetrasan)
	if err1 != nil {
		return c.JSON(err1.Code, err1)
	}
	return c.JSON(http.StatusCreated, createdexpencetrasan)
}
func (controller expencetrasanController) CreateExp(c echo.Context) error {
	expencetrasan := &model.Expencetrasan{}
	
	expencetrasan.Name = c.FormValue("name")
	expencetrasan.Title = c.FormValue("title")
	expencetrasan.Description = c.FormValue("description")
	expencetrasan.Code = c.FormValue("code")
	expencetrasan.Type = c.FormValue("status")
	expencetrasan.Usercode = c.FormValue("usercode")

	s, err := strconv.ParseFloat(c.FormValue("amount"), 64)
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid selling price")
		return c.JSON(httperror.Code, httperror)
	}
	expencetrasan.Amount = s
	createdexpencetrasan, err1 := service.ExpencetrasanService.CreateExp(expencetrasan)
	if err1 != nil {
		return c.JSON(err1.Code, err1)
	}
	return c.JSON(http.StatusCreated, createdexpencetrasan)
}
func (controller expencetrasanController) View(c echo.Context) error {
	code := c.Param("code")
	options, problem := service.ExpencetrasanService.View(code)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, options)	
}
func (controller expencetrasanController) ViewExp(c echo.Context) error {
	options, problem := service.ExpencetrasanService.ViewExp()
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, options)	
}

func (controller expencetrasanController) UpdateTrans(c echo.Context) error {
		
	name := c.FormValue("name")
	code := c.FormValue("code")
	updatedcart, problem := service.ExpencetrasanService.UpdateTrans(name, code)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, updatedcart)
}
func (controller expencetrasanController) ViewReport(c echo.Context) error {
		dated := c.QueryParam("dated")
		searchq2 := c.QueryParam("searchq2")
		searchq3 := c.QueryParam("searchq3")
	options, problem := service.ExpencetrasanService.ViewReport(dated,searchq2,searchq3)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, options)	
}
func (controller expencetrasanController) GetAll(c echo.Context) error {
	
	search := string(c.QueryParam("q"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid page number")
		return c.JSON(httperror.Code, httperror)
	}
	pagesize, err := strconv.Atoi(c.QueryParam("pagesize"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid pagesize")
		return c.JSON(httperror.Code, httperror)
	}
	
	results, err3 := service.ExpencetrasanService.GetAll(search, page,pagesize)
	if err3 != nil {
		return c.JSON(err3.Code, err3)
	}
	return c.JSON(http.StatusOK, results)
} 
func (controller expencetrasanController) GetOne(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	fmt.Println(id)
	expencetrasan, problem := service.ExpencetrasanService.GetOne(id)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, expencetrasan)	
}

func (controller expencetrasanController) Update(c echo.Context) error {
		
	expencetrasan :=  &model.Expencetrasan{}
	if err := c.Bind(expencetrasan); err != nil {
		httperror := httperors.NewBadRequestError("Invalid json body")
		return c.JSON(httperror.Code, httperror)
	}	
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	updatedexpencetrasan, problem := service.ExpencetrasanService.Update(id, expencetrasan)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, updatedexpencetrasan)
}

func (controller expencetrasanController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	success, failure := service.ExpencetrasanService.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code, failure)
	}
	return c.JSON(success.Code, success)
		
}
