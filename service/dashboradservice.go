package service

import (
	// "fmt"
	// "github.com/myrachanto/accounting/support"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
)

var (
	Dashboardservice dashboardservice = dashboardservice{}

) 
type dashboardservice struct {
	
}

func (service dashboardservice) View() (*model.Dashboard, *httperors.HttpError) {
	dashboard, err1 := r.Dashboardrepo.View()
	if err1 != nil {
		return nil, err1
	}
	 return dashboard, nil

}
func (service dashboardservice) Email() (*model.Email, *httperors.HttpError) {
	dashboard, err1 := r.Dashboardrepo.Email()
	if err1 != nil {
		return nil, err1
	}
	return dashboard, nil
}

// func (service dashboardservice) Send() (*model.Email, *httperors.HttpError) {
// 	dashboards, err := r.Dashboardrepo.Send()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return dashboards, nil
// }
//db.Where("age = ?", 20).Delete(&User{})