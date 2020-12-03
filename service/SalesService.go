package service

import (
	// "fmt"
	// "github.com/myrachanto/accounting/support"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
)
//Salesservice report service
var (
	Salesservice salesservice = salesservice{}

) 
type salesservice struct {
	
}

func (service salesservice) View() (*model.Sales, *httperors.HttpError) {
	sales, err1 := r.Salesrepo.View()
	if err1 != nil {
		return nil, err1
	}
	 return sales, nil

}
// func (service salesservice) Email() (*model.Email, *httperors.HttpError) {
// 	sales, err1 := r.salesrepo.Email()
// 	if err1 != nil {
// 		return nil, err1
// 	}
// 	return sales, nil
// }

// func (service salesservice) Send() (*model.Email, *httperors.HttpError) {
// 	saless, err := r.salesrepo.Send()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return saless, nil
// }
//db.Where("age = ?", 20).Delete(&User{})