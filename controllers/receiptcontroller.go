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
	"github.com/myrachanto/accounting/support"
)
 
var ( 
	ReceiptController receiptController =  receiptController{}
)
type receiptController struct{ }
/////////controllers/////////////////
func (controller receiptController) Create(c echo.Context) error {
	receipt := &model.Receipt{}
	receipt.CustomerName = c.FormValue("customername")
	receipt.Status = c.FormValue("status")
	receipt.Description = c.FormValue("description")
	receipt.Type = c.FormValue("type")
	d := c.FormValue("clearancedate")

	s, err := strconv.ParseFloat(c.FormValue("amount"), 64)
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid selling price")
		return c.JSON(httperror.Code, httperror)
	}
	receipt.Amount = s
layout := "01/02/2006"
	t, er := time.Parse(layout, d)
	if er != nil {
		httperror := httperors.NewBadRequestError("Invalid Date")
		return c.JSON(httperror.Code, httperror)
	}
	receipt.ClearanceDate = t
		createdreceipt, err1 := service.Receiptservice.Create(receipt)
		if err1 != nil {
			return c.JSON(err1.Code, err1)
		}
	return c.JSON(http.StatusCreated, createdreceipt)
}
func (controller receiptController) GetAll(c echo.Context) error {
	receipts := []model.Receipt{}
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
	
	receipts, err3 := service.Receiptservice.GetAll(receipts,search)
	if err3 != nil {
		return c.JSON(err3.Code, err3)
	}
	return c.JSON(http.StatusOK, receipts)
} 
func (controller receiptController) GetOne(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	fmt.Println(id)
	receipt, problem := service.Receiptservice.GetOne(id)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, receipt)	
}

func (controller receiptController) Update(c echo.Context) error {
		
	receipt :=  &model.Receipt{}
	if err := c.Bind(receipt); err != nil {
		httperror := httperors.NewBadRequestError("Invalid json body")
		return c.JSON(httperror.Code, httperror)
	}	
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	updatedreceipt, problem := service.Receiptservice.Update(id, receipt)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, updatedreceipt)
}

func (controller receiptController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	success, failure := service.Receiptservice.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code, failure)
	}
	return c.JSON(success.Code, success)
		
}
