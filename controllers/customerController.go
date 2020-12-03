package controllers

import(
	"fmt"
	"strconv"	
	
	"io"
	"os"
	"net/http"
	"github.com/labstack/echo"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/service"
	"github.com/myrachanto/accounting/support"
) 
  
var (
	CustomerController customerController = customerController{}
)
type customerController struct{ } 
/////////controllers/////////////////
func (controller customerController) Create(c echo.Context) error {
	customer := &model.Customer{}
	customer.Name = c.FormValue("name")
	customer.Company = c.FormValue("company")
	customer.Phone = c.FormValue("phone")
	customer.Address = c.FormValue("address")
	customer.Email = c.FormValue("email")
	customer.Password = c.FormValue("password")

	   pic, err2 := c.FormFile("picture")
	//    fmt.Println(pic.Filename)
	   if err2 != nil {
				httperror := httperors.NewBadRequestError("Invalid picture")
				return c.JSON(httperror.Code, err2)
			}	
		src, err := pic.Open()
		if err != nil {
			httperror := httperors.NewBadRequestError("the picture is corrupted")
			return c.JSON(httperror.Code, err)
		}	
		defer src.Close()
		filePath := "./public/imgs/customers/" + pic.Filename
		// Destination
		dst, err4 := os.Create(filePath)
		if err4 != nil {
			httperror := httperors.NewBadRequestError("the Directory mess")
			return c.JSON(httperror.Code, err4)
		}
		defer dst.Close()
		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			if err2 != nil {
				httperror := httperors.NewBadRequestError("error filling")
				return c.JSON(httperror.Code, httperror)
			}
		}
		
	customer.Picture = pic.Filename
	s, err1 := service.Customerservice.Create(customer)
	if err1 != nil {
		return c.JSON(err1.Code, err1)
	}
	return c.JSON(http.StatusCreated, s)
}
func (controller customerController) Login(c echo.Context) error {
	customer := &model.Logincustomer{}
	// auth := &model.CustomnerAuth{}
	if err := c.Bind(customer); err != nil {
		httperror := httperors.NewBadRequestError("Invalid json body")
		return c.JSON(httperror.Code, httperror)
	}	
	
	auth, problem := service.Customerservice.Login(customer)
	if problem != nil {
		fmt.Println(problem)
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, auth)	
}
func (controller customerController) Forgot(c echo.Context) error {
	email := c.FormValue("email")
	s, problem := service.Customerservice.Forgot(email)
	if problem != nil {
		fmt.Println(problem)
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, s)	
}
func (controller customerController) Logout(c echo.Context) error {
	token := string(c.Param("token"))
	problem := service.Customerservice.Logout(token)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, "succeessifully logged out")	
}
func (controller customerController) GetAll(c echo.Context) error {
	//  customers := []model.Customer{} 
	column := string(c.QueryParam("column"))
	direction := string(c.QueryParam("direction"))
	search_column := string(c.QueryParam("search_column"))
	search_operator := string(c.QueryParam("search_operator"))
	search_query_1 := string(c.QueryParam("search_query_1"))
	search_query_2 := string(c.QueryParam("search_query_2"))
	per_page, err := strconv.Atoi(c.QueryParam("per_page"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid per_page number")
		return c.JSON(httperror.Code, httperror)
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid page number")
		return c.JSON(httperror.Code, httperror)
	}
	fmt.Println("------------------------")
	search := &support.Search{Column:column, Direction:direction,Search_column:search_column,Search_operator:search_operator,Search_query_1:search_query_1,Search_query_2:search_query_2,Per_page:per_page,Page:page}
	Customers, err3 := service.Customerservice.GetAll(search)
	if err3 != nil {
		return c.JSON(err3.Code, err3)
	}
	return c.JSON(http.StatusOK, Customers)
} 
func (controller customerController) GetOne(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	fmt.Println(id)
	customer, problem := service.Customerservice.GetOne(id)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, customer)	
}

func (controller customerController) Update(c echo.Context) error {
		
	customer := &model.Customer{}
	customer.Name = c.FormValue("name")
	customer.Company = c.FormValue("company")
	customer.Phone = c.FormValue("phone")
	customer.Address = c.FormValue("address")
	customer.Email = c.FormValue("email")

	   pic, err2 := c.FormFile("picture")
	//    fmt.Println(pic.Filename)
	   if err2 != nil {
				httperror := httperors.NewBadRequestError("Invalid picture")
				return c.JSON(httperror.Code, err2)
			}	
		src, err := pic.Open()
		if err != nil {
			httperror := httperors.NewBadRequestError("the picture is corrupted")
			return c.JSON(httperror.Code, err)
		}	
		defer src.Close()
		filePath := "./public/imgs/customers/" + pic.Filename
		// Destination
		dst, err4 := os.Create(filePath)
		if err4 != nil {
			httperror := httperors.NewBadRequestError("the Directory mess")
			return c.JSON(httperror.Code, err4)
		}
		defer dst.Close()
		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			if err2 != nil {
				httperror := httperors.NewBadRequestError("error filling")
				return c.JSON(httperror.Code, httperror)
			}
		}
		
	customer.Picture = pic.Filename
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	updatedcustomer, problem := service.Customerservice.Update(id, customer)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, updatedcustomer)
}

func (controller customerController) ViewReport(c echo.Context) error {
	options, problem := service.Customerservice.ViewReport()
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, options)	
}
func (controller customerController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	success, failure := service.Customerservice.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code, failure)
	}
	return c.JSON(success.Code, success)
		
}
