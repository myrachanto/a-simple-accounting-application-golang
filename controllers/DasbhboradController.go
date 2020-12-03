package controllers

import(
	"net/http"
	"github.com/labstack/echo"
	"github.com/myrachanto/accounting/service"
)
 
var (
	DashboardController dashboardController = dashboardController{}
)
type dashboardController struct{ }
/////////controllers/////////////////
func (controller dashboardController) View(c echo.Context) error {
		
	createddashboard, err1 := service.Dashboardservice.View()
	if err1 != nil {
		return c.JSON(err1.Code, err1)
	}
	return c.JSON(http.StatusOK, createddashboard)
}
func (controller dashboardController) Email(c echo.Context) error {
	dashboards, err3 := service.Dashboardservice.Email()
	if err3 != nil {
		return c.JSON(err3.Code, err3)
	}
	return c.JSON(http.StatusOK, dashboards)
} 
// func (controller dashboardController) Send(c echo.Context) error {
// 	dashboard, problem := service.Dashboardservice.Send()
// 	if problem != nil {
// 		return c.JSON(problem.Code, problem)
// 	}
// 	return c.JSON(http.StatusOK, dashboard)	
// }
