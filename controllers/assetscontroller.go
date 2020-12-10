package controllers

import(
	"fmt"
	"os"
	"io"
	"strconv"	
	"net/http"
	"github.com/labstack/echo"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/service"
)
 //AssetController ...
var (
	AssetController assetController = assetController{}
)
type assetController struct{ }
/////////controllers/////////////////
func (controller assetController) Create(c echo.Context) error {
	asset := &model.Asset{}
	
	
	asset.Name = c.FormValue("name")
	asset.Description = c.FormValue("description")
	asset.Ownership = c.FormValue("ownership")
	asset.Depreciationtype = c.FormValue("depreciationtype")
	asset.Liscence = c.FormValue("liscence")
	asset.Pincode = c.FormValue("pin")
	asset.Usercode = c.FormValue("usercode")
	asset.Assetcode = c.FormValue("code")
	d, err := strconv.ParseFloat(c.FormValue("depreciationrate"), 64)
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid selling price")
		return c.JSON(httperror.Code, httperror)
	}
	asset.Depreciationrate = d
	s, err := strconv.ParseFloat(c.FormValue("price"), 64)
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid selling price")
		return c.JSON(httperror.Code, httperror)
	}
	asset.Price = s
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
		filePath := "./public/imgs/assets/" + pic.Filename
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
		
	asset.Picture = pic.Filename
	createdasset, err1 := service.AssetService.Create(asset)
	if err1 != nil {
		return c.JSON(err1.Code, err1)
	}
	return c.JSON(http.StatusCreated, createdasset)
}

func (controller assetController) View(c echo.Context) error {
	code, problem := service.AssetService.View()
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, code)	
}
func (controller assetController) GetAll(c echo.Context) error {
	
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
		
		results, err3 := service.AssetService.GetAll(search, page,pagesize)
		if err3 != nil {
			return c.JSON(err3.Code, err3)
		}
		return c.JSON(http.StatusOK, results)
} 
func (controller assetController) GetOne(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	fmt.Println(id)
	asset, problem := service.AssetService.GetOne(id)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, asset)	
}

func (controller assetController) Update(c echo.Context) error {
		
	asset :=  &model.Asset{}
	if err := c.Bind(asset); err != nil {
		httperror := httperors.NewBadRequestError("Invalid json body")
		return c.JSON(httperror.Code, httperror)
	}	
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	updatedasset, problem := service.AssetService.Update(id, asset)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, updatedasset)
}

func (controller assetController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	success, failure := service.AssetService.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code, failure)
	}
	return c.JSON(success.Code, success)
		
}
