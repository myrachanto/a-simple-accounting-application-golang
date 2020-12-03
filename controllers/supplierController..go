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
//SupplierController ...
var (
	SupplierController supplierController = supplierController{}
)
type supplierController struct{ } 
/////////controllers/////////////////
func (controller supplierController) Create(c echo.Context) error {
	supplier := &model.Supplier{}
	supplier.Name = c.FormValue("name")
	supplier.Company = c.FormValue("company")
	supplier.Phone = c.FormValue("phone")
	supplier.Address = c.FormValue("address")
	supplier.Email = c.FormValue("email")
	supplier.Password = c.FormValue("password")

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
		filePath := "./public/imgs/suppliers/" + pic.Filename
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
		
	supplier.Picture = pic.Filename
	s, err1 := service.Supplierservice.Create(supplier)
	if err1 != nil {
		return c.JSON(err1.Code, err1)
	}
	return c.JSON(http.StatusCreated, s)
}
func (controller supplierController) Login(c echo.Context) error {
	supplier := &model.Loginsupplier{}
	if err := c.Bind(supplier); err != nil {
		httperror := httperors.NewBadRequestError("Invalid json body")
		return c.JSON(httperror.Code, httperror)
	}	
	
	auth, problem := service.Supplierservice.Login(supplier)
	if problem != nil {
		fmt.Println(problem)
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, auth)	
}
func (controller supplierController) Forgot(c echo.Context) error {
	email := c.FormValue("email")
	s, problem := service.Supplierservice.Forgot(email)
	if problem != nil {
		fmt.Println(problem)
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, s)	
}
func (controller supplierController) Logout(c echo.Context) error {
	token := string(c.Param("token"))
	problem := service.Supplierservice.Logout(token)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, "succeessifully logged out")	
}
func (controller supplierController) GetAll(c echo.Context) error {
	//  suppliers := []model.Supplier{} 
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
	suppliers, err3 := service.Supplierservice.GetAll(search)
	if err3 != nil {
		return c.JSON(err3.Code, err3)
	}
	return c.JSON(http.StatusOK, suppliers)
} 
func (controller supplierController) GetOne(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	fmt.Println(id)
	supplier, problem := service.Supplierservice.GetOne(id)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, supplier)	
}

func (controller supplierController) Update(c echo.Context) error {
		
	supplier := &model.Supplier{}
	supplier.Name = c.FormValue("name")
	supplier.Company = c.FormValue("company")
	supplier.Phone = c.FormValue("phone")
	supplier.Address = c.FormValue("address")
	supplier.Email = c.FormValue("email")

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
		filePath := "./public/imgs/suppliers/" + pic.Filename
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
		
	supplier.Picture = pic.Filename
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	updatedsupplier, problem := service.Supplierservice.Update(id, supplier)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, updatedsupplier)
}

func (controller supplierController) ViewReport(c echo.Context) error {
	options, problem := service.Supplierservice.ViewReport()
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, options)	
}
func (controller supplierController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	success, failure := service.Supplierservice.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code, failure)
	}
	return c.JSON(success.Code, success)
		
}
