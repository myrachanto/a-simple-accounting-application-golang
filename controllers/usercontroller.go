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
)
//UserController ..
var (
	UserController userController = userController{}
)
type userController struct{ }
/////////controllers/////////////////
func (controller userController) Register(c echo.Context) error {
	fmt.Println("endpoint called!")
	user := &model.User{}
	
	user.FName = c.FormValue("fname")
	user.UName = c.FormValue("uname")
	user.LName = c.FormValue("lname")
	user.Phone = c.FormValue("phone")
	user.Address = c.FormValue("address")
	user.Email = c.FormValue("email")
	user.Password = c.FormValue("password")

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
		filePath := "./public/imgs/users/" + pic.Filename
		// Destination
		dst, err4 := os.Create(filePath)
		if err4 != nil {
			httperror := httperors.NewBadRequestError("the Directory mess")
			return c.JSON(httperror.Code, err4)
		}
		defer dst.Close()
		// Copy
		
		
	user.Picture = pic.Filename
	createdUser, err1 := service.UserService.Create(user)
	if err1 != nil {
		return c.JSON(err1.Code, err1)
	}
	return c.JSON(http.StatusCreated, createdUser)
}
func (controller userController) Create(c echo.Context) error {
	user := &model.User{}
	user.FName = c.FormValue("fname")
	user.LName = c.FormValue("lname")
	user.UName = c.FormValue("uname")
	user.Phone = c.FormValue("phone")
	user.Address = c.FormValue("address")
	user.Email = c.FormValue("email")
	user.Password = c.FormValue("password")

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
		filePath := "./public/imgs/users/" + pic.Filename
		// Destination
		dst, err4 := os.Create(filePath)
		if err4 != nil {
			httperror := httperors.NewBadRequestError("the Directory mess")
			return c.JSON(httperror.Code, err4)
		}
		defer dst.Close()
		// Copy
		
		
	user.Picture = pic.Filename
	s, err1 := service.UserService.Create(user)
	if err1 != nil {
		return c.JSON(err1.Code, err1)
	}
	if _, err = io.Copy(dst, src); err != nil {
		if err2 != nil {
			httperror := httperors.NewBadRequestError("error filling")
			return c.JSON(httperror.Code, httperror)
		}
	}
	return c.JSON(http.StatusCreated, s)
}
func (controller userController) Login(c echo.Context) error {
	user := &model.LoginUser{}
	auth := &model.Auth{}
	
	user.Email = c.FormValue("email")
	user.Password = c.FormValue("password")
	auth, problem := service.UserService.Login(user)
	if problem != nil {
		fmt.Println(problem)
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, auth)	
}
func (controller userController) Logout(c echo.Context) error {
	token := string(c.Param("token"))
	problem := service.UserService.Logout(token)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, "succeessifully logged out")	
}
func (controller userController) GetAll(c echo.Context) error {
	
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
	
	results, err3 := service.UserService.GetAll(search, page,pagesize)
	if err3 != nil {
		return c.JSON(err3.Code, err3)
	}
	return c.JSON(http.StatusOK, results)
} 
func (controller userController) GetOne(c echo.Context) error {
		
	code := c.Param("code")
	dated := c.QueryParam("dated")
	searchq2 := c.QueryParam("searchq2")
	searchq3 := c.QueryParam("searchq3")
	// fmt.Println(c.QueryParams)
	user, problem := service.UserService.GetOne(code,dated,searchq2,searchq3)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, user)	
}

func (controller userController) UpdateRole(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	role := c.FormValue("role")
	usercode := c.FormValue("usercode")
	updateduser, problem := service.UserService.UpdateRole(id,role, usercode)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, updateduser)
}
func (controller userController) Update(c echo.Context) error {
	user :=  &model.User{}
	user.FName = c.FormValue("fname")
	user.LName = c.FormValue("lname")
	user.UName = c.FormValue("uname")
	user.Phone = c.FormValue("phone")
	user.Address = c.FormValue("address")
	user.Email = c.FormValue("email")
	user.Password = "njenga456"
	// user.Password = c.FormValue("password")
	user.Admin = true

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
		filePath := "./public/imgs/users/" + pic.Filename
		// Destination
		dst, err4 := os.Create(filePath)
		if err4 != nil {
			httperror := httperors.NewBadRequestError("the Directory mess")
			return c.JSON(httperror.Code, err4)
		}
		defer dst.Close()
		// Copy
		
	user.Picture = pic.Filename
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	updateduser, problem := service.UserService.Update(id, user)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	if _, err = io.Copy(dst, src); err != nil {
		if err2 != nil {
			httperror := httperors.NewBadRequestError("error filling")
			return c.JSON(httperror.Code, httperror)
		}
	}
	
	return c.JSON(http.StatusOK, updateduser)
}

func (controller userController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	success, failure := service.UserService.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code, failure)
	}
	return c.JSON(success.Code, success)
		
}