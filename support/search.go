package support

import (
	"fmt"
	"strings"
	"log"
	"os"
	"reflect"
	"github.com/joho/godotenv"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/myrachanto/accounting/httperors"
	// "github.com/myrachanto/astore-v2.0/model"
	"github.com/vcraescu/go-paginator" 
	"github.com/vcraescu/go-paginator/adapter"
	 
)
type Search struct{
	Column string
	Direction string
	Search_column string
	Search_operator string
	Search_query_1 string
	Search_query_2 string
	Per_page int
	Page int

}
type Productsearch struct{
	Column string
	Name string
	Direction string
	Search_column string
	Search_operator string
	Search_query_1 string
	Per_page int

}

var (
	IndexRepo indexRepo = indexRepo{}
	Operator = map[string]string{"all":"all","equal_to":"=","not_equal_to":"<>","less_than":"<",
"greater_than":">","less_than_or_equal_to":"<=","greater_than_ro_equal_to":">=",
"like":"like","between":"between","in":"in","not_in":"not_in"}
)

///curtesy to gorm
type indexRepo struct{}
func Getconnected() (GormDB *gorm.DB, err *httperors.HttpError) {
	err1 := godotenv.Load()
	if err1 != nil {
		log.Fatal("Error loading .env file in routes")
	}
	dbuser := os.Getenv("DbUsername")
	DbName := os.Getenv("DbName")
	dbURI := dbuser+"@/"+DbName+"?charset=utf8&parseTime=True&loc=Local"
	GormDB, err2 := gorm.Open("mysql", dbURI)
	if err2 != nil {
		return nil, httperors.NewNotFoundError("No Mysql db connection")
	}
	return GormDB, nil
}
func DbClose(GormDB *gorm.DB) {
	defer GormDB.Close()
}
func SearchQuery(Ser *Search,  val interface{})(values []interface{}, err *httperors.HttpError){
	value := reflect.TypeOf(val)
	// switch(value){
	// case "all":
	// }
	fmt.Println(value)
	GormDB, err1 := Getconnected()
	if err1 != nil {
		return nil, err1
	}
	switch(Ser.Search_operator){
	case "all":
		//db.Order("name DESC")
		q := GormDB.Model(&value).Order(Ser.Column+" "+Ser.Direction).Find(&values)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		if err3 := p.Results(&values); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "equal_to":
		q := GormDB.Model(&value).Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&values);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		if err3 := p.Results(&values); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "not_equal_to":
		q := GormDB.Model(&value).Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&values);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		if err3 := p.Results(&values); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "less_than" :
		/////////need research on how to query this data///////////////////////////////////////////
		// db.Find(&users, User{Age: 20})
		// q := GormDB.Find(&values, value{Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1}).Order(Ser.Column+" "+Ser.Direction);
		q := GormDB.Model(&value).Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&values);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		if err3 := p.Results(&values); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "greater_than":
		q := GormDB.Model(&value).Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&values);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		if err3 := p.Results(&values); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "less_than_or_equal_to":
		q := GormDB.Model(&value).Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&values);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		if err3 := p.Results(&values); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "greater_than_ro_equal_to":
		q := GormDB.Model(&value).Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&values);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		if err3 := p.Results(&values); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		q := GormDB.Model(&value).Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&values);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		if err3 := p.Results(&values); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		q := GormDB.Model(&value).Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&values);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		if err3 := p.Results(&values); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	// break;
	case "like":
		q := GormDB.Model(&value).Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&values);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		if err3 := p.Results(&values); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		q := GormDB.Model(&value).Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&values);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		if err3 := p.Results(&values); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	DbClose(GormDB)
	
	return values, nil
}
// func (search Search) Validate() *httperors.HttpError{ 
// 	if search.Column == "" {
// 		return httperors.NewNotFoundError("Invalid Column")
// 	}
// 	if search.Direction == "" {
// 		return httperors.NewNotFoundError("Invalid Direction")
// 	}
// 	if search.Search_column == "" {
// 		return httperors.NewNotFoundError("Invalid Search column")
// 	}
// 	if search.Search_operator == "" {
// 		return httperors.NewNotFoundError("Invalid Search operator")
// 	}	
// 	if search.Per_page <= 0 {
// 		return httperors.NewNotFoundError("Invalid Per page")
// 	}
// 	return nil
// }