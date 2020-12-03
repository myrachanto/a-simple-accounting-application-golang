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
func (controller expencetrasanController) View(c echo.Context) error {
	code := c.Param("code")
	options, problem := service.ExpencetrasanService.View(code)
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
	options, problem := service.ExpencetrasanService.ViewReport()
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, options)	
}
func (controller expencetrasanController) GetAll(c echo.Context) error {
	expencetrasans := []model.Expencetrasan{}
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
	expencetrasans, err3 := service.ExpencetrasanService.GetAll(expencetrasans,search)
	if err3 != nil {
		return c.JSON(err3.Code, err3)
	}
	return c.JSON(http.StatusOK, expencetrasans)
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
