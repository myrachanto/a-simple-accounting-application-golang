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
 //SubcategoryController ...
var (
	SubcategoryController subcategoryController = subcategoryController{}
)
type subcategoryController struct{ }
/////////controllers/////////////////
func (controller subcategoryController) Create(c echo.Context) error {
	subcategory := &model.Subcategory{}
	
	if err := c.Bind(subcategory); err != nil {
		httperror := httperors.NewBadRequestError("Invalid json body")
		return c.JSON(httperror.Code, httperror)
	}	
	createdsubcategory, err1 := service.Subcategoryservice.Create(subcategory)
	if err1 != nil {
		return c.JSON(err1.Code, err1)
	}
	return c.JSON(http.StatusCreated, createdsubcategory)
}
func (controller subcategoryController) GetAll(c echo.Context) error {

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
	
	results, err3 := service.Subcategoryservice.GetAll(search, page,pagesize)
	if err3 != nil {
		return c.JSON(err3.Code, err3)
	}
	return c.JSON(http.StatusOK, results)
} 
func (controller subcategoryController) GetOne(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	fmt.Println(id)
	subcategory, problem := service.Subcategoryservice.GetOne(id)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, subcategory)	
}

func (controller subcategoryController) Update(c echo.Context) error {
		
	subcategory :=  &model.Subcategory{}
	if err := c.Bind(subcategory); err != nil {
		httperror := httperors.NewBadRequestError("Invalid json body")
		return c.JSON(httperror.Code, httperror)
	}	
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	updatedsubcategory, problem := service.Subcategoryservice.Update(id, subcategory)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, updatedsubcategory)
}

func (controller subcategoryController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	success, failure := service.Subcategoryservice.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code, failure)
	}
	return c.JSON(success.Code, success)
		
}
