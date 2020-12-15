package repository

import (
	"time"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
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
	user := Userrepo.Userbycode(message.Tousercode)
	message.Name = user.UName
	message.Read = false
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


func (messageRepo messagerepo) Sent(code,dated,searchq2,searchq3 string) (results []model.Message, r *httperors.HttpError) {

	now := time.Now()
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}

	if dated != "custom"{
		if dated == "In the last 24hrs"{
			d := now.AddDate(0, 0, -1)
			GormDB.Where("updated_at > ? AND fromusercode = ?", d,code).Find(&results)
		}
		if dated == "In the last 7days"{
			d := now.AddDate(0, 0, -7)
			GormDB.Where("updated_at > ? AND fromusercode = ?", d,code).Find(&results)
		}
		if dated == "In the last 15day"{
			d := now.AddDate(0, 0, -15)
			GormDB.Where("updated_at > ? AND fromusercode = ?", d,code).Find(&results)
		}
		if dated == "In the last 30days"{
			d := now.AddDate(0, 0, -30)
			GormDB.Where("updated_at > ? AND fromusercode = ?", d,code).Find(&results)
		}
	}
	if dated == "custom"{
		start,err := time.Parse(Layout,searchq2)
		if err != nil {
			return nil, httperors.NewNotFoundError("Something went wrong parsing date1!")
		}
		end,err1 := time.Parse(Layout,searchq3)
		if err1 != nil {
			return nil, httperors.NewNotFoundError("Something went wrong parsing date1!")
		}
		GormDB.Where("fromusercode = ? AND updated_at BETWEEN ? AND ?",code, start, end).Find(&results)
	}
	IndexRepo.DbClose(GormDB)
	return results, nil

}


func (messageRepo messagerepo) Inbox(code,dated,searchq2,searchq3 string) (results []model.Message, r *httperors.HttpError) {

	now := time.Now()
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}

	if dated != "custom"{
		if dated == "In the last 24hrs"{
			d := now.AddDate(0, 0, -1)
			GormDB.Where("updated_at > ? AND tousercode = ?", d,code).Find(&results)
		}
		if dated == "In the last 7days"{
			d := now.AddDate(0, 0, -7)
			GormDB.Where("updated_at > ? AND tousercode = ?", d,code).Find(&results)
		}
		if dated == "In the last 15day"{
			d := now.AddDate(0, 0, -15)
			GormDB.Where("updated_at > ? AND tousercode = ?", d,code).Find(&results)
		}
		if dated == "In the last 30days"{
			d := now.AddDate(0, 0, -30)
			GormDB.Where("updated_at > ? AND tousercode = ?", d,code).Find(&results)
		}
	}
	if dated == "custom"{
		start,err := time.Parse(Layout,searchq2)
		if err != nil {
			return nil, httperors.NewNotFoundError("Something went wrong parsing date1!")
		}
		end,err1 := time.Parse(Layout,searchq3)
		if err1 != nil {
			return nil, httperors.NewNotFoundError("Something went wrong parsing date1!")
		}
		GormDB.Where("tousercode = ? AND updated_at BETWEEN ? AND ?",code, start, end).Find(&results)
	}
	IndexRepo.DbClose(GormDB)
	return results, nil

}
func (messageRepo messagerepo) GetAllUnread() (*model.MessageUnread, *httperors.HttpError) {
	results := []model.Message{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	} 
	// db.Scopes(Paginate(r)).Find(&users)
	GormDB.Where("read = ?", false).Find(&results)

	IndexRepo.DbClose(GormDB)
	return &model.MessageUnread{
		Num: len(results),
		Messages: results,
	}, nil
}
func (messageRepo messagerepo) GetAll(dated,searchq2,searchq3 string) (results []model.Message, e *httperors.HttpError) {
	
	now := time.Now()
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	if dated != "custom"{ 
		if dated == "In the last 24hrs"{
			d := now.AddDate(0, 0, -1)
			GormDB.Where("updated_at > ?", d).Find(&results)
		}
		if dated == "In the last 7days"{
			d := now.AddDate(0, 0, -7)
			GormDB.Where("updated_at > ?",d).Find(&results)
		}
		if dated == "In the last 15day"{
			d := now.AddDate(0, 0, -15)
			GormDB.Where("updated_at > ?",d).Find(&results)
		}
		if dated == "In the last 30days"{
			d := now.AddDate(0, 0, -30)
			GormDB.Where("updated_at > ?",d).Find(&results)
		}
	}
	if dated == "custom"{
		start,err := time.Parse(Layout,searchq2)
		if err != nil {
			return nil, httperors.NewNotFoundError("Something went wrong parsing date1!")
		}
		end,err1 := time.Parse(Layout,searchq3)
		if err1 != nil {
			return nil, httperors.NewNotFoundError("Something went wrong parsing date1!")
		}
		GormDB.Where("updated_at BETWEEN ? AND ?",start, end).Find(&results)
	}
	IndexRepo.DbClose(GormDB)
	return results, nil
}
func (messageRepo messagerepo) AllSearch(dated,searchq2,searchq3 string) (results []model.Message, r *httperors.HttpError) {

	now := time.Now()
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	if dated != "custom"{ 
		if dated == "In the last 24hrs"{
			d := now.AddDate(0, 0, -1)
			GormDB.Where("updated_at > ?", d).Find(&results)
		}
		if dated == "In the last 7days"{
			d := now.AddDate(0, 0, -7)
			GormDB.Where("updated_at > ?",d).Find(&results)
		}
		if dated == "In the last 15day"{
			d := now.AddDate(0, 0, -15)
			GormDB.Where("updated_at > ?",d).Find(&results)
		}
		if dated == "In the last 30days"{
			d := now.AddDate(0, 0, -30)
			GormDB.Where("updated_at > ?",d).Find(&results)
		}
	}
	if dated == "custom"{
		start,err := time.Parse(Layout,searchq2)
		if err != nil {
			return nil, httperors.NewNotFoundError("Something went wrong parsing date1!")
		}
		end,err1 := time.Parse(Layout,searchq3)
		if err1 != nil {
			return nil, httperors.NewNotFoundError("Something went wrong parsing date1!")
		}
		GormDB.Where("updated_at BETWEEN ? AND ?",start, end).Find(&results)
	}
	IndexRepo.DbClose(GormDB)
	return results, nil

}
func (messageRepo messagerepo) Update(id int) (*model.Message, *httperors.HttpError) {
	message :=  &model.Message{}
	ok := messageRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("message with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&message).Where("id = ?", id).Update("read",true)
	
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
