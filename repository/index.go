package repository

import (
	// "log"
	// "os"
	// "github.com/joho/godotenv"
  "gorm.io/gorm"
  "gorm.io/driver/sqlite"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
)
//IndexRepo
var (
	IndexRepo indexRepo = indexRepo{}
	Operator = map[string]string{"all":"all","equal_to":"=","not_equal_to":"<>","less_than":"<",
"greater_than":">","less_than_or_equal_to":"<=","greater_than_ro_equal_to":">=",
"like":"like","between":"between","in":"in","not_in":"not_in"}
)

const (
	Layout = "2006-01-02"
	layoutUS  = "January 2, 2006"
)

///curtesy to gorm
type indexRepo struct{}
func init(){
	GormDB, err1 := gorm.Open(sqlite.Open("accounting.db"), &gorm.Config{})
    if err1 != nil {
        panic("failed to connect database")
    }
	GormDB.AutoMigrate(&model.Cart{})
	GormDB.AutoMigrate(&model.Category{})
	GormDB.AutoMigrate(&model.Customer{})
	GormDB.AutoMigrate(&model.CustomnerAuth{})
	GormDB.AutoMigrate(&model.Supplier{})
	GormDB.AutoMigrate(&model.SupplierAuth{})
	
	GormDB.AutoMigrate(&model.User{})
	GormDB.AutoMigrate(&model.Auth{})
	GormDB.AutoMigrate(&model.Message{})
	GormDB.AutoMigrate(&model.Nortification{})	
	/////product/////////////////////////
	GormDB.AutoMigrate(&model.Product{})
	GormDB.AutoMigrate(&model.Majorcategory{})
	GormDB.AutoMigrate(&model.Category{})
	GormDB.AutoMigrate(&model.Subcategory{})
	GormDB.AutoMigrate(&model.Price{})
	GormDB.AutoMigrate(&model.Tax{})
	GormDB.AutoMigrate(&model.Discount{})
	/////customer/////////////////////////
	GormDB.AutoMigrate(&model.Cart{})
	GormDB.AutoMigrate(&model.Customer{})
	GormDB.AutoMigrate(&model.Invoice{})
	GormDB.AutoMigrate(&model.Paymentform{})
	GormDB.AutoMigrate(&model.Transaction{})
	GormDB.AutoMigrate(&model.DebtTransaction{})
	/////supplier/////////////////////////
	GormDB.AutoMigrate(&model.SInvoice{})
	GormDB.AutoMigrate(&model.STransaction{})
	GormDB.AutoMigrate(&model.Scart{})
	/////accounts/////////////////////////
	GormDB.AutoMigrate(&model.Payment{})
	GormDB.AutoMigrate(&model.Receipt{})
	GormDB.AutoMigrate(&model.Payrectrasan{})
	GormDB.AutoMigrate(&model.Expence{})
	GormDB.AutoMigrate(&model.Expencetrasan{})
	GormDB.AutoMigrate(&model.Asset{})
	GormDB.AutoMigrate(&model.Asstrans{}) 
	GormDB.AutoMigrate(&model.Liability{})
	GormDB.AutoMigrate(&model.Liatran{})
	return 
}
func (indexRepo indexRepo) Getconnected() (GormDB *gorm.DB, err *httperors.HttpError) {
	// err1 := godotenv.Load()
	// if err1 != nil {
	// 	log.Fatal("Error loading .env file in routes")
	// }
	// dbuser := os.Getenv("DbUsername")
	// DbName := os.Getenv("DbName")
	// dbURI := dbuser+"@/"+DbName+"?charset=utf8&parseTime=True&loc=Local"
	// GormDB, err2 := gorm.Open("mysql", dbURI)
	// if err2 != nil {
	// 	return nil, httperors.NewNotFoundError("No Mysql db connection")
	// }
	// GormDB, err1 := gorm.Open("sqlite3", "accounting.db")

	// GormDB, err1 := gorm.Open(sqlite.Open("accounting.db"), &gorm.Config{})
	GormDB, err1 := gorm.Open(sqlite.Open("accounting.db"), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt: true,
	})
    if err1 != nil {
        panic("failed to connect database")
    }
	return GormDB, nil
}
func (indexRepo indexRepo) DbClose(GormDB *gorm.DB) {
	// defer GormDB.Close()
}