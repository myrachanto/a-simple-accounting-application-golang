package repository

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"github.com/joho/godotenv"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
)
//Userrepo ...
var (
	Userrepo userrepo = userrepo{}
)

///curtesy to gorm
type userrepo struct{}

func (userRepo userrepo) Create(user *model.User) (string, *httperors.HttpError) {
	if err := user.Validate(); err != nil {
		return "", err
	}
	ok, err1 := user.ValidatePassword(user.Password)
	if !ok {
		return "", err1
	}
	ok = user.ValidateEmail(user.Email)
	if !ok {
		return "", httperors.NewNotFoundError("Your email format is wrong!")
	}
	ok = userRepo.UserExist(user.Email)
	if ok {
		return "", httperors.NewNotFoundError("Your email already exists!")
	}
	
	hashpassword, err2 := user.HashPassword(user.Password)
	if err2 != nil {
		return "", err2
	}
	user.Password = hashpassword
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	code, x := userRepo.GeneCode()
	if x != nil {
		return "", x
	}
	user.Usercode = code
	// fmt.Println(user)
	GormDB.Create(&user)
	IndexRepo.DbClose(GormDB)
	return "user created successifully", nil
}
func (userRepo userrepo) Login(auser *model.LoginUser) (*model.Auth, *httperors.HttpError) {
	if err := auser.Validate(); err != nil {
		return nil, err
	}
	ok := userRepo.UserExist(auser.Email)
	if !ok {
		return nil, httperors.NewNotFoundError("Your email does not exists!")
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	user := model.User{}

	GormDB.Model(&user).Where("email = ?", auser.Email).First(&user)
	ok = user.Compare(auser.Password, user.Password)
	if !ok {
		return nil, httperors.NewNotFoundError("wrong email password combo!")
	}
	tk := &model.Token{
		UserID: user.ID,
		UName: user.UName,
		Admin:user.Admin,
		Usercode:user.Usercode,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: model.ExpiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading key")
	}
	encyKey := os.Getenv("EncryptionKey")
	tokenString, error := token.SignedString([]byte(encyKey))
	if error != nil {
		fmt.Println(error)
	}
	// messages ,e := userRepo.UnreadMessages(user.ID)
	// if e != nil {
	// 	return nil, e
	// }
	// norti ,e := userRepo.UnreadNortifications(user.ID)
	// if e != nil {
	// 	return nil, e
	// }
	auth := &model.Auth{UserID:user.ID, UName:user.UName, Usercode:user.Usercode, Admin:user.Admin, Picture:user.Picture, Token:tokenString}
	GormDB.Create(&auth)
	IndexRepo.DbClose(GormDB)
	
	return auth, nil
}

func (userRepo userrepo) All() (t []model.User, r *httperors.HttpError) {

	user := model.User{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&user).Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (userRepo userrepo) Logout(token string) (*httperors.HttpError) {
	auth := model.Auth{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return err1
	}
	res := GormDB.First(&auth, "token =?", token)
	if res.Error != nil {
		return httperors.NewNotFoundError("Something went wrong logging out!")
	 }
	
	GormDB.Model(&auth).Where("token =?", token).First(&auth)
	
	GormDB.Delete(auth)
	IndexRepo.DbClose(GormDB)
	
	return  nil
}
func (userRepo userrepo)GeneCode() (string, *httperors.HttpError) {
	user := model.User{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	err := GormDB.Last(&user)
	if err.Error != nil {
		var c1 uint = 1
		code := "UserCode"+strconv.FormatUint(uint64(c1), 10)
		return code, nil
	 }
	c1 := user.ID + 1
	code := "UserCode"+strconv.FormatUint(uint64(c1), 10)
	IndexRepo.DbClose(GormDB)
	return code, nil
	
}
func (userRepo userrepo) GetOne(code,dated,searchq2,searchq3 string) (*model.UserProfile, *httperors.HttpError) { 
	ok := userRepo.UserExistbycode(code)
	if !ok {
		return nil, httperors.NewNotFoundError("User with that code does not exists!")
	}
	user := model.User{} 
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil { 
		return nil, err1
	}
	
	GormDB.Model(&user).Where("usercode = ?", code).First(&user)
	IndexRepo.DbClose(GormDB)
	sent, err := Messagerepo.Sent(code,dated,searchq2,searchq3)
	if err != nil {
		return nil,err 
	}
	inbox, err2 := Messagerepo.Inbox(code,dated,searchq2,searchq3)
	if err2 != nil {
		return nil,err2
	}
	users, err3 := Userrepo.All()
	if err3 != nil {
		return nil,err3
	}
	return &model.UserProfile{
		User:user,
		Sent:sent,
		Inbox:inbox,
		Users:users,
	}, nil
}
func (userRepo userrepo)UserExistbycode(code string) bool {
	u := model.User{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	GormDB.Where("usercode = ?", code).First(&u)
	if u.ID == 0 {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (userRepo userrepo)Userbycode(code string) *model.User {
	u := model.User{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil
	}
	GormDB.Where("usercode = ?", code).First(&u)
	if u.ID == 0 {
	   return nil
	}
	IndexRepo.DbClose(GormDB)
	return &u
	
}
func (userRepo userrepo) GetAll(search string, page,pagesize int) ([]model.User, *httperors.HttpError) {
	results := []model.User{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	if search == ""{
		GormDB.Find(&results)
	}
	// db.Scopes(Paginate(r)).Find(&users)
	GormDB.Scopes(Paginate(page,pagesize)).Where("name LIKE ?", "%"+search+"%").Or("email LIKE ?", "%"+search+"%").Or("company LIKE ?", "%"+search+"%").Find(&results)

	IndexRepo.DbClose(GormDB)
	return results, nil
}

func (userRepo userrepo) UpdateRole(id int,role, usercode string) (string, *httperors.HttpError) {
	user := model.User{}
	ok := Userrepo.UserExistByid(id)
	if !ok {
		return "", httperors.NewNotFoundError("customer with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}

	// "employee",
	// "supervisor",
	// "admin",
	if usercode == "employee"{
		GormDB.Model(&user).Where("id = ?", id).Update("employee",true)
	}
	if usercode == "supervisor"{
		GormDB.Model(&user).Where("id = ?", id).Update("supervisor",true)
	}
	GormDB.Model(&user).Where("id = ?", id).Update("admin",true)
	IndexRepo.DbClose(GormDB)

	return "user updated succesifully", nil
}
func (userRepo userrepo) Update(id int, user *model.User) (*model.User, *httperors.HttpError) {
	ok := userRepo.UserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("User with that id does not exists!")
	}
	
	hashpassword, err2 := user.HashPassword(user.Password)
	if err2 != nil {
		return nil, err2
	}
	user.Password = hashpassword
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	User := model.User{}
	uuser := model.User{}
	
	GormDB.Model(&User).Where("id = ?", id).First(&uuser)
	if user.FName  == "" {
		user.FName = uuser.FName
	}
	if user.LName  == "" {
		user.LName = uuser.LName
	}
	if user.UName  == "" {
		user.UName = uuser.UName
	}
	if user.Phone  == "" {
		user.Phone = uuser.Phone
	}
	if user.Address  == "" {
		user.Address = uuser.Address
	}
	if user.Picture  == "" {
		user.Picture = uuser.Picture
	}
	if user.Email  == "" {
		user.Email = uuser.Email
	}
	if user.Admin  == false {
		user.Admin = true
	}
	GormDB.Save(&user)
	
	IndexRepo.DbClose(GormDB)

	return user, nil
}
func (userRepo userrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := userRepo.UserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("User with that id does not exists!")
	}
	user := model.User{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&user).Where("id = ?", id).First(&user)
	GormDB.Delete(user)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (userRepo userrepo)UserExist(email string) bool {
	user := model.User{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&user, "email =?", email)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (userRepo userrepo)UserExistByid(id int) bool {
	user := model.User{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&user, "id =?", id)
	if res.Error != nil{
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
// func (userRepo userrepo)UnreadMessages(id uint)  (int, *httperors.HttpError)  {
// 	messages := []model.Message{}
// 	GormDB, err1 := IndexRepo.Getconnected()
// 	if err1 != nil {
// 		return 0, err1
// 	}
// 	GormDB.Where("id = ? AND read = ? ", id, false).Find(&messages)	
// 	 c := 0
// 	 for i, _:= range messages{
// 		 c += i
// 	 }
// 	IndexRepo.DbClose(GormDB)
// 	return c, nil
	
// }
// func (userRepo userrepo)UnreadNortifications(id uint)  (int, *httperors.HttpError)  {
// 	ns := []model.Nortification{}
// 	GormDB, err1 := IndexRepo.Getconnected()
// 	if err1 != nil {
// 		return 0, err1
// 	}
// 	GormDB.Where("id = ? AND read = ? ", id, false).Find(&ns)	
// 	 c := 0
// 	 for i, _:= range ns{
// 		 c += i
// 	 }
// 	IndexRepo.DbClose(GormDB)
// 	return c, nil
	
// }
