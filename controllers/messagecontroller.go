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
//MessageController ...
var (
	MessageController messageController = messageController{}
)
type messageController struct{ }
/////////controllers/////////////////
func (controller messageController) Create(c echo.Context) error {
	message := &model.Message{}
	if err := c.Bind(message); err != nil {
		httperror := httperors.NewBadRequestError("Invalid json body")
		return c.JSON(httperror.Code, httperror)
	}
	createdmessage, err1 := service.MessageService.Create(message)
	if err1 != nil {
		return c.JSON(err1.Code, err1)
	}
	return c.JSON(http.StatusCreated, createdmessage)
}
func (controller messageController) GetAll(c echo.Context) error {
	
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
	
	results, err3 := service.MessageService.GetAll(search, page,pagesize)
	if err3 != nil {
		return c.JSON(err3.Code, err3)
	}
	return c.JSON(http.StatusOK, results)
} 
func (controller messageController) GetOne(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	fmt.Println(id)
	message, problem := service.MessageService.GetOne(id)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, message)	
}

func (controller messageController) Update(c echo.Context) error {
		
	message :=  &model.Message{}
	if err := c.Bind(message); err != nil {
		httperror := httperors.NewBadRequestError("Invalid json body")
		return c.JSON(httperror.Code, httperror)
	}	
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	updatedmessage, problem := service.MessageService.Update(id, message)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, updatedmessage)
}

func (controller messageController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	success, failure := service.MessageService.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code, failure)
	}
	return c.JSON(success.Code, success)
		
}
