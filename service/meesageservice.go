package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
	"github.com/myrachanto/accounting/support"
)

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

func (service messageService) GetAll(messages []model.Message,search *support.Search) ([]model.Message, *httperors.HttpError) {
	messages, err := r.Messagerepo.GetAll(messages,search)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (service messageService) Update(id int, message *model.Message) (*model.Message, *httperors.HttpError) {
	message, err1 := r.Messagerepo.Update(id, message)
	if err1 != nil {
		return nil, err1
	}
	
	return message, nil
}
func (service messageService) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Messagerepo.Delete(id)
		return success, failure
}
