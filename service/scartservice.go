package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
)
//Scartservice ..
var (
	Scartservice scartservice = scartservice{}

) 
type scartservice struct {
	
}

func (service scartservice) Create(scart *model.Scart) (string, *httperors.HttpError) {
	car, err1 := r.Scartrepo.Create(scart)
	if err1 != nil {
		return "", err1
	}
	 return car, nil

}
func (service scartservice) View(code string) ([]model.Scart, *httperors.HttpError) {
	options, err1 := r.Scartrepo.View(code)
	if err1 != nil {
		return nil, err1
	}
	return options, nil
}

func (service scartservice) Getcredits(code string) ([]model.Transaction, *httperors.HttpError) {
	options, err1 := r.Scartrepo.Getcredits(code)
	if err1 != nil {
		return nil, err1
	}
	return options, nil
}


func (service scartservice) GetcreditsList(code string) ([]model.Transaction, *httperors.HttpError) {
	options, err1 := r.Scartrepo.GetcreditsList(code)
	if err1 != nil {
		return nil, err1
	}
	return options, nil 
}
func (service scartservice) GetOne(id int) (*model.Scart, *httperors.HttpError) {
	scart, err1 := r.Scartrepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return scart, nil
}

func (service scartservice) GetAll(scarts []model.Scart) ([]model.Scart, *httperors.HttpError) {
	scarts, err := r.Scartrepo.GetAll(scarts)
	if err != nil {
		return nil, err
	}
	return scarts, nil
}

func (service scartservice) Update(qty float64, name, code string) (string, *httperors.HttpError) {
	scart, err1 := r.Scartrepo.Update(qty, name,code)
	if err1 != nil {
		return "", err1
	}
	 
	return scart, nil
}
func (service scartservice) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Scartrepo.Delete(id)
		return success, failure
}
func (service scartservice) DeleteALL(code string) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Scartrepo.DeleteAll(code)
		return success, failure
}
///////deleting a batch////////////////////

//db.Where("age = ?", 20).Delete(&User{})