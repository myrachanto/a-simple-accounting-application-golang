package repository

import (
	"fmt"
	"strings"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/support"
)
//Messagerepo ... 
var (
	Messagerepo messagerepo = messagerepo{}
)

///curtesy to gorm
type messagerepo struct{}

func (messageRepo messagerepo) Create(message *model.Message) (*model.Message, *httperors.HttpError) {
	if err := message.Validate(); err != nil {
		return nil, err
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&message)
	IndexRepo.DbClose(GormDB)
	return message, nil
}
func (messageRepo messagerepo) GetOne(id int) (*model.Message, *httperors.HttpError) {
	ok := messageRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("message with that id does not exists!")
	}
	message := model.Message{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&message).Where("id = ?", id).First(&message)
	IndexRepo.DbClose(GormDB)
	
	return &message, nil
}

func (messageRepo messagerepo) GetAll(messages []model.Message,search *support.Search) ([]model.Message, *httperors.HttpError) {
	results, err1 := messageRepo.Search(search, messages)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}

func (messageRepo messagerepo) Update(id int, message *model.Message) (*model.Message, *httperors.HttpError) {
	ok := messageRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("message with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	amessage := model.Message{}
	
	GormDB.Model(&message).Where("id = ?", id).First(&amessage)
	if message.Name  == "" {
		message.Name = amessage.Name
	}
	if message.Title  == "" {
		message.Title = amessage.Title
	}
	if message.Description  == "" {
		message.Description = amessage.Description
	}
	GormDB.Save(&message)
	
	IndexRepo.DbClose(GormDB)

	return message, nil
}
func (messageRepo messagerepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := messageRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("Product with that id does not exists!")
	}
	message := model.Message{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&message).Where("id = ?", id).First(&message)
	GormDB.Delete(message)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (messageRepo messagerepo)ProductUserExistByid(id int) bool {
	message := model.Message{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&message, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}

func (messageRepo messagerepo) Search(Ser *support.Search, messages []model.Message)([]model.Message, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	message := model.Message{}
	switch(Ser.Search_operator){
	case "all":
		GormDB.Model(&message).Order(Ser.Column+" "+Ser.Direction).Find(&messages)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
	
	break;
	case "equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&messages);
	
	break;
	case "not_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&messages);	
	
	break;
	case "less_than" :
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&messages);	
	
	break;
	case "greater_than":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&messages);	
	
	break;
	case "less_than_or_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&messages);	
	
	break;
	case "greater_than_ro_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&messages);	
	
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&messages);
	
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		GormDB.Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&messages);
	
	// break;
	case "like":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&messages);
	
	break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&messages);
	
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return messages, nil
}