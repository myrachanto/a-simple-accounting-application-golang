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
	
	message.Title = c.FormValue("title")
	message.Description = c.FormValue("description")
	message.Tousercode = c.FormValue("tousercode")
	message.Fromusercode = c.FormValue("fromusercode")
	createdmessage, err1 := service.MessageService.Create(message)
	if err1 != nil {
		return c.JSON(err1.Code, err1)
	}
	return c.JSON(http.StatusCreated, createdmessage)
}
func (controller messageController) GetAllUnread(c echo.Context) error {
	 
	results, err3 := service.MessageService.GetAllUnread()
	if err3 != nil {
		return c.JSON(err3.Code, err3)
	}
	return c.JSON(http.StatusOK, results)
}
func (controller messageController) GetAll(c echo.Context) error {
	 
	
	dated := c.QueryParam("dated")
	searchq2 := c.QueryParam("searchq2")
	searchq3 := c.QueryParam("searchq3")
	results, err3 := service.MessageService.GetAll(dated,searchq2,searchq3)
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
		
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	updatedmessage, problem := service.MessageService.Update(id)
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
