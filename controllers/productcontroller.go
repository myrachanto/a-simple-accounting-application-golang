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
	ProductController productController = productController{}
)
type productController struct{ }
/////////controllers/////////////////
func (controller productController) Create(c echo.Context) error {
	product := &model.Product{}
	
	product.Name = c.FormValue("name")
	product.Title = c.FormValue("title")
	product.Description = c.FormValue("description")
	product.Category = c.FormValue("category")
	product.Usercode = c.FormValue("usercode")

	s, err := strconv.ParseFloat(c.FormValue("sprice"), 64)
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid selling price")
		return c.JSON(httperror.Code, httperror)
	}
	product.SPrice = s
	// // Multipart form
		// Source

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
		filePath := "./public/imgs/products/" + pic.Filename
		// Destination
		dst, err4 := os.Create(filePath)
		if err4 != nil {
			httperror := httperors.NewBadRequestError("the Directory mess")
			return c.JSON(httperror.Code, err4)
		}
		defer dst.Close()
		//  copy
		if _, err = io.Copy(dst, src); err != nil {
			if err2 != nil {
				httperror := httperors.NewBadRequestError("error filling")
				return c.JSON(httperror.Code, httperror)
			}
		}
		
	product.Picture = pic.Filename
	createdproduct, err1 := service.Productservice.Create(product)
	if err1 != nil {
		return c.JSON(err1.Code, err1)
	} 
	return c.JSON(http.StatusCreated, createdproduct)
}

func (controller productController) View(c echo.Context) error {
	options, problem := service.Productservice.View()
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, options)	
}
func (controller productController) SearchProduct(c echo.Context) error {

	search := string(c.QueryParam("search"))
	// page, err := strconv.Atoi(c.QueryParam("page"))
	// if err != nil {
	// 	httperror := httperors.NewBadRequestError("Invalid page number")
	// 	return c.JSON(httperror.Code, httperror)
	// }
	// pagesize, err := strconv.Atoi(c.QueryParam("pagesize"))
	// if err != nil {
	// 	httperror := httperors.NewBadRequestError("Invalid pagesize")
	// 	return c.JSON(httperror.Code, httperror)
	// }
	options, problem := service.Productservice.ProductSearch(search)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, options)	
}
func (controller productController) GetProducts(c echo.Context) error {
	products := []model.Product{}	
	column := string(c.QueryParam("column"))
	name := string(c.QueryParam("name"))
	search_column := string(c.QueryParam("search_column"))
	search_operator := string(c.QueryParam("search_operator"))
	search_query_1 := string(c.QueryParam("search_query_1"))
	per_page, err := strconv.Atoi(c.QueryParam("per_page"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid per number")
		return c.JSON(httperror.Code, httperror)
	}
	fmt.Println("------------------------")
	search := &support.Productsearch{Column:column,Name:name,Search_column:search_column,Search_operator:search_operator,Search_query_1:search_query_1,Per_page:per_page}
	
	products, err3 := service.Productservice.GetProducts(products, search)
	if err3 != nil {
		return c.JSON(err3.Code, err3)
	}
	return c.JSON(http.StatusOK, products)
} 
func (controller productController) GetAll(c echo.Context) error {
	
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
	
	results, err3 := service.Productservice.GetAll(search, page,pagesize)
	if err3 != nil {
		return c.JSON(err3.Code, err3)
	}
	return c.JSON(http.StatusOK, results)
} 
func (controller productController) GetOne(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	fmt.Println(id)
	product, problem := service.Productservice.GetOne(id)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, product)	
}

func (controller productController) Update(c echo.Context) error {
		
	product :=  &model.Product{}
	
	product.Name = c.FormValue("name")
	product.Title = c.FormValue("title")
	product.Description = c.FormValue("description")
	product.Category = c.FormValue("category")
	s := c.FormValue("sprice")
	b, err := strconv.ParseFloat(s, 64)
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid buying price")
		return c.JSON(httperror.Code, httperror)
	}

	fmt.Println(product)
	fmt.Println(s)
	// s, err := strconv.ParseFloat(c.FormValue("sprice"), 64)
	// if err != nil {
	// 	httperror := httperors.NewBadRequestError("Invalid selling price")
	// 	return c.JSON(httperror.Code, httperror)
	// }
	product.BPrice = b
	// product.SPrice = s
	// // Multipart form
		// Source

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
		filePath := "./public/imgs/products/" + pic.Filename
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
		
	product.Picture = pic.Filename
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	updatedproduct, problem := service.Productservice.Update(id, product)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	
	return c.JSON(http.StatusOK, updatedproduct)
}

func (controller productController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	success, failure := service.Productservice.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code, failure)
	}
	return c.JSON(success.Code, success)
		
}
