package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
)
//MessageService ...
var (
	MessageService messageService = messageService{}

) 
type messageService struct {
	
}

func (service messageService) Create(message *model.Message) (*model.Message, *httperors.HttpError) {
	if err := message.Validate(); err != nil {
		return nil, err
	}	
	message, err1 := r.Messagerepo.Create(message)
	if err1 != nil {
		return nil, err1
	}
	 return message, nil

}
func (service messageService) GetOne(id int) (*model.Message, *httperors.HttpError) {
	message, err1 := r.Messagerepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return message, nil
}

func (service messageService) GetAllUnread() (*model.MessageUnread, *httperors.HttpError) {
	results, err := r.Messagerepo.GetAllUnread()
	return results, err
}
func (service messageService) GetAll(dated,searchq2,searchq3 string) ([]model.Message, *httperors.HttpError) {
	results, err := r.Messagerepo.GetAll(dated,searchq2,searchq3)
	return results, err
}
func (service messageService) Update(id int) (*model.Message, *httperors.HttpError) {
	message, err1 := r.Messagerepo.Update(id)
	if err1 != nil {
		return nil, err1
	}
	
	return message, nil
}
func (service messageService) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Messagerepo.Delete(id)
		return success, failure
}
