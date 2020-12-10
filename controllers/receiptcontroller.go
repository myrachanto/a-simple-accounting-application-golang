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
 //ReceiptController ...
var ( 
	ReceiptController receiptController =  receiptController{}
)
const (
	layoutISO = "2006-01-02"
	layoutUS  = "January 2, 2006"
)
type receiptController struct{ }
/////////controllers/////////////////
func (controller receiptController) Create(c echo.Context) error {
	receipt := &model.Receipt{}
	receipt.CustomerName = c.FormValue("customername")
	receipt.Status = c.FormValue("status")
	receipt.Description = c.FormValue("description")
	receipt.Type = c.FormValue("type")
	receipt.Code = c.FormValue("code")
	receipt.Customercode = c.FormValue("customercode")
	receipt.Usercode = c.FormValue("usercode")
	receipt.ChequeNo = c.FormValue("chequeno")
	ex := c.FormValue("expirydate")
	d := c.FormValue("clearancedate")
	fmt.Println(d)

	s, err := strconv.ParseFloat(c.FormValue("amount"), 64)
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid selling price")
		return c.JSON(httperror.Code, httperror)
	}
	receipt.Amount = s
  t, er := time.Parse(layoutISO, d)
	if er != nil { 
		httperror := httperors.NewBadRequestError("Invalid Date")
		return c.JSON(httperror.Code, httperror)
	}
	receipt.ClearanceDate = t
	tx, ers := time.Parse(layoutISO, ex)
	if ers != nil { 
		httperror := httperors.NewBadRequestError("Invalid Date")
		return c.JSON(httperror.Code, httperror)
	}
	receipt.Expirydate = tx
		createdreceipt, err1 := service.Receiptservice.Create(receipt)
		if err1 != nil {
			return c.JSON(err1.Code, err1)
		}
	return c.JSON(http.StatusCreated, createdreceipt)
}
func (controller receiptController) UpdateReceipts(c echo.Context) error {
		
	code := c.FormValue("code")
	status := c.FormValue("status")
	updatedcart, problem := service.Receiptservice.UpdateReceipts(code,status)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, updatedcart)
}
func (controller receiptController) AddReceiptTrans(c echo.Context) error {
		
	clientcode := c.FormValue("customercode")
	invoicecode := c.FormValue("invoicecode")
	usercode := c.FormValue("usercode")
	receiptcode := c.FormValue("receiptcode")
	fmt.Println(receiptcode)
	amount, err := strconv.ParseFloat(c.FormValue("amount"), 64)
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid amount")
		return c.JSON(httperror.Code, httperror)
	}
	updatedcart, problem := service.Receiptservice.AddReceiptTrans(clientcode,invoicecode,usercode,receiptcode,amount)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, updatedcart)
}

func (controller receiptController) ViewReport(c echo.Context) error {
	options, problem := service.Receiptservice.ViewReport()
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, options)	
}

func (controller receiptController) ViewCleared(c echo.Context) error {
	options, problem := service.Receiptservice.ViewCleared()
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, options)	
}
func (controller receiptController) ViewInvoices(c echo.Context) error {
	customercode := c.Param("customercode")
	fmt.Println(customercode)
	invoices, problem := service.Receiptservice.ViewInvoices(customercode)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, invoices)	
}
func (controller receiptController) View(c echo.Context) error {
	code, problem := service.Receiptservice.View()
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, code)	
}
func (controller receiptController) GetAll(c echo.Context) error {
	
	receipts, err3 := service.Receiptservice.GetAll()
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
