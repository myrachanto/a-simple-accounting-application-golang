package controllers

import(
	"net/http"
	"github.com/labstack/echo"
	"github.com/myrachanto/accounting/service"
)
 
var (
	SalesController salesController = salesController{}
)
type salesController struct{ }
/////////controllers/////////////////
func (controller salesController) View(c echo.Context) error {
		
	createdsales, err1 := service.Salesservice.View()
	if err1 != nil {
		return c.JSON(err1.Code, err1)
	}
	return c.JSON(http.StatusOK, createdsales)
}
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
