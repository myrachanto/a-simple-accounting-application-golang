package controllers

import(
	"net/http"
	"github.com/labstack/echo"
	"github.com/myrachanto/accounting/service"
)
 //SalesController ...
var (
	SalesController salesController = salesController{}
)
type salesController struct{ }
/////////controllers/////////////////
func (controller salesController) View(c echo.Context) error {
		
	dated := c.QueryParam("dated")
	searchq2 := c.QueryParam("searchq2")
	searchq3 := c.QueryParam("searchq3")
	createdsales, err1 := service.Salesservice.View(dated,searchq2,searchq3)
	if err1 != nil {
		return c.JSON(err1.Code, err1)
	}
	return c.JSON(http.StatusOK, createdsales)
}
func (controller salesController) Purchases(c echo.Context) error {
		
	dated := c.QueryParam("dated")
	searchq2 := c.QueryParam("searchq2")
	searchq3 := c.QueryParam("searchq3")
	createdsales, err1 := service.Salesservice.Purchases(dated,searchq2,searchq3)
	if err1 != nil {
		return c.JSON(err1.Code, err1)
	}
	return c.JSON(http.StatusOK, createdsales)
}
func (controller salesController) Pl(c echo.Context) error {
		
	dated := c.QueryParam("dated")
	searchq2 := c.QueryParam("searchq2")
	searchq3 := c.QueryParam("searchq3")
	createdsales, err1 := service.Salesservice.Pl(dated,searchq2,searchq3)
	if err1 != nil {
		return c.JSON(err1.Code, err1)
	}
	return c.JSON(http.StatusOK, createdsales)
}
func (controller salesController) Supplierstement(c echo.Context) error {
	code := c.Param("code")
	dated := c.QueryParam("dated")
	searchq2 := c.QueryParam("searchq2")
	searchq3 := c.QueryParam("searchq3")
	createdsales, err1 := service.Salesservice.Supplierstement(code,dated,searchq2,searchq3)
	if err1 != nil {
		return c.JSON(err1.Code, err1)
	}
	return c.JSON(http.StatusOK, createdsales)
}
func (controller salesController) Customerstement(c echo.Context) error {
	code := c.Param("code")
	dated := c.QueryParam("dated")
	searchq2 := c.QueryParam("searchq2")
	searchq3 := c.QueryParam("searchq3")
	createdsales, err1 := service.Salesservice.Customerstement(code,dated,searchq2,searchq3)
	if err1 != nil {
		return c.JSON(err1.Code, err1)
	}
	return c.JSON(http.StatusOK, createdsales)
}
// func (controller salesController) Customer(c echo.Context) error {
		
// 	dated := c.QueryParam("dated")
// 	searchq2 := c.QueryParam("searchq2")
// 	searchq3 := c.QueryParam("searchq3")
// 	createdsales, err1 := service.Salesservice.Customer(dated,searchq2,searchq3)
// 	if err1 != nil {
// 		return c.JSON(err1.Code, err1)
// 	}
// 	return c.JSON(http.StatusOK, createdsales)
// }
// func (controller salesController) Email(c echo.Context) error {
// 	saless, err3 := service.salesservice.Email()
// 	if err3 != nil {
// 		return c.JSON(err3.Code, err3)
// 	}
// 	return c.JSON(http.StatusOK, saless)
// } 
// func (controller salesController) Send(c echo.Context) error {
// 	sales, problem := service.salesservice.Send()
// 	if problem != nil {
// 		return c.JSON(problem.Code, problem)
// 	}
// 	return c.JSON(http.StatusOK, sales)	
// }
