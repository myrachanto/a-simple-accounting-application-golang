package routes

import (
	"log"
	"os"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/myrachanto/accounting/controllers"
	jwt "github.com/dgrijalva/jwt-go"
)
//StoreAPI =>entry point to routes
func StoreAPI(){

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file in routes")
	}
	PORT := os.Getenv("PORT")
	key := os.Getenv("EncryptionKey")

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover()) 
	e.Use(middleware.CORS())

	e.Static("/", "public")
	JWTgroup := e.Group("/api/")
	JWTgroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS256",
		SigningKey: []byte(key),
	}))
// e.Set(id, 1)

	// // Get retrieves data from the context.
	// Get(key string) interface{}

	// // Set saves data in the context.
	// Set(key string, val interface{}) 

	// admin := e.Group("admin/")
	// admin.Use(isAdmin)
	/////////////////////////////////////////////////////////////////////////////////////
	////////////////////////needs more info ////////////////////////////////////////////
	///////////////////////////////////////////////////////////////////////////////////
	// var IsLoggedIn = middleware.JWTWithConfig(middleware.JWTConfig{
	// 	SigningMethod: "HS256",
	// 	SigningKey: []byte(key),
	// })
	//JwtG := e.Group("/users")
	// JwtG.Use(middleware.JWT([]byte(key)))
	// Routes
	//e.GET("/is-loggedin", h.private, IsLoggedIn)
	//e.POST("/register", IsLoggedIn,isAdmin,isEmployee,isSupervisor, controllers.UserController.Create)
	e.POST("/register", controllers.UserController.Register)
	e.POST("/login", controllers.UserController.Login)
	JWTgroup.GET("logout/:token", controllers.UserController.Logout)
	JWTgroup.GET("users", controllers.UserController.GetAll)
	JWTgroup.GET("users/:id", controllers.UserController.GetOne)
	JWTgroup.PUT("users/role/:id", controllers.UserController.UpdateRole)
	JWTgroup.PUT("users/:id", controllers.UserController.Update)
	JWTgroup.DELETE("users/:id", controllers.UserController.Delete)

	///////////dashboard/////////////////////////////	
	JWTgroup.GET("dashboard", controllers.DashboardController.View)
	JWTgroup.GET("email/create", controllers.DashboardController.Email)
	// JWTgroup.POST("email/create", controllers.DashboardController.Send)
	//e.DELETE("loggoutall/:id", controllers.UserController.DeleteALL) logout all accounts
	///////////message/////////////////////////////	
	JWTgroup.POST("messages", controllers.MessageController.Create)
	JWTgroup.GET("messages", controllers.MessageController.GetAll)
	JWTgroup.GET("messages/:id", controllers.MessageController.GetOne)
	JWTgroup.PUT("messages/:id", controllers.MessageController.Update)
	JWTgroup.DELETE("messages/:id", controllers.MessageController.Delete)
	///////////nortifications/////////////////////////////	
	JWTgroup.POST("nortifications", controllers.NortificationController.Create)
	JWTgroup.GET("nortifications", controllers.MessageController.GetAll) 
	JWTgroup.GET("nortifications/:id", controllers.MessageController.GetOne)
	JWTgroup.PUT("nortifications/:id", controllers.MessageController.Update)
	JWTgroup.DELETE("nortifications/:id", controllers.MessageController.Delete)
	///////////category/////////////////////////////	
	JWTgroup.GET("categorys/view", controllers.CategoryController.View)
	JWTgroup.POST("categorys", controllers.CategoryController.Create)
	JWTgroup.GET("categorys", controllers.CategoryController.GetAll)
	JWTgroup.GET("categorys/:id", controllers.CategoryController.GetOne)
	JWTgroup.PUT("categorys/:id", controllers.CategoryController.Update)
	JWTgroup.DELETE("categorys/:id", controllers.CategoryController.Delete)
	///////////majorcategory/////////////////////////////	
	JWTgroup.POST("majorcategory", controllers.MCategoryController.Create)
	JWTgroup.GET("majorcategory", controllers.MCategoryController.GetAll)
	JWTgroup.GET("majorcategory/:id", controllers.MCategoryController.GetOne)
	JWTgroup.PUT("majorcategory/:id", controllers.MCategoryController.Update)
	JWTgroup.DELETE("majorcategory/:id", controllers.MCategoryController.Delete)
	///////////paymentform/////////////////////////////	
	JWTgroup.POST("paymentform", controllers.PaymentformController.Create)
	JWTgroup.GET("paymentform", controllers.PaymentformController.GetAll)
	JWTgroup.GET("paymentform/:id", controllers.PaymentformController.GetOne)
	JWTgroup.PUT("paymentform/:id", controllers.PaymentformController.Update)
	JWTgroup.DELETE("paymentform/:id", controllers.PaymentformController.Delete)
	///////////subcategory/////////////////////////////	
	JWTgroup.POST("subcategory", controllers.SubcategoryController.Create)
	JWTgroup.GET("subcategory", controllers.SubcategoryController.GetAll)
	JWTgroup.GET("subcategory/:id", controllers.SubcategoryController.GetOne)
	JWTgroup.PUT("subcategory/:id", controllers.SubcategoryController.Update)
	JWTgroup.DELETE("subcategory/:id", controllers.SubcategoryController.Delete)
	///////////subcategory/////////////////////////////	
	JWTgroup.GET("products/view", controllers.ProductController.View)
	JWTgroup.POST("products", controllers.ProductController.Create)
	e.GET("productsearch", controllers.ProductController.GetProducts)
	JWTgroup.GET("products", controllers.ProductController.GetAll)
	JWTgroup.GET("searchproducts", controllers.ProductController.SearchProduct)
	e.GET("products/:id", controllers.ProductController.GetOne)
	JWTgroup.GET("products/:id", controllers.ProductController.GetOne)
	JWTgroup.PUT("products/quantity/:id", controllers.ProductController.UpdateQty)
	JWTgroup.PUT("products/:id", controllers.ProductController.Update)
	JWTgroup.DELETE("products/:id", controllers.ProductController.Delete)
	///////////cart/////////////////////////////	
	JWTgroup.POST("carts", controllers.CartController.Create)
	JWTgroup.GET("carts/view/:code", controllers.CartController.View)
	JWTgroup.GET("carts/:id", controllers.CartController.GetOne)
	// JWTgroup.PUT("carts/:id", controllers.CartController.Update)
	JWTgroup.POST("carts/credit", controllers.CartController.Updatetrans) 
	JWTgroup.DELETE("carts/cancel/:code", controllers.CartController.DeleteAll)
	JWTgroup.DELETE("carts/delete/:id", controllers.CartController.Delete)
	JWTgroup.GET("carts/credits/:code", controllers.CartController.Getcredits)
	JWTgroup.GET("carts/creditslist/:code", controllers.CartController.GetcreditsList)
	///////////////////////////////////carts
	///////////////////customer module///////////////////////////////////////
	///////////Invoice/////////////////////////////	////////////////////////
	JWTgroup.POST("customer", controllers.CustomerController.Create)
	JWTgroup.GET("customer", controllers.CustomerController.GetAll)
	JWTgroup.GET("customer/report", controllers.CustomerController.ViewReport)
	JWTgroup.GET("customer/:id", controllers.CustomerController.GetOne)
	JWTgroup.PUT("customer/:id", controllers.CustomerController.Update)
	JWTgroup.DELETE("customer/:id", controllers.CustomerController.Delete)
	///////////Invoice/////////////////////////////	
	JWTgroup.POST("invoicescart", controllers.InvoiceController.CreateCart)
	JWTgroup.GET("invoice/view", controllers.InvoiceController.View)
	JWTgroup.POST("invoice", controllers.InvoiceController.Create)
	JWTgroup.GET("invoice", controllers.InvoiceController.GetAll)
	JWTgroup.GET("invoice/:id", controllers.InvoiceController.GetOne)
	JWTgroup.POST("invoice/credit", controllers.InvoiceController.Credit) 
	JWTgroup.GET("invoice/credit", controllers.InvoiceController.GetCredit) 
	// JWTgroup.PUT("invoice/:id", controllers.InvoiceController.Update) 
	// JWTgroup.DELETE("invoice/:id", controllers.InvoiceController.Delete)
	///////////trasanctions/////////////////////////////	
	JWTgroup.POST("trasanctions", controllers.TransactionController.Create)
	JWTgroup.GET("trasanctions", controllers.TransactionController.GetAll)
	JWTgroup.GET("trasanctions/:id", controllers.TransactionController.GetOne)
	JWTgroup.PUT("trasanctions/:id", controllers.TransactionController.Update)
	JWTgroup.DELETE("trasanctions/:id", controllers.TransactionController.Delete)
	//////////////////////////////////////////////////////////////////////////
	///////////////////supplier module///////////////////////////////////////
	///////////Invoice/////////////////////////////	//////////////////////////

	JWTgroup.POST("supplier", controllers.SupplierController.Create)
	JWTgroup.GET("supplier", controllers.SupplierController.GetAll)
	JWTgroup.GET("supplier/report", controllers.SupplierController.ViewReport)
	JWTgroup.GET("supplier/:id", controllers.SupplierController.GetOne)
	JWTgroup.PUT("supplier/:id", controllers.SupplierController.Update)
	JWTgroup.DELETE("supplier/:id", controllers.SupplierController.Delete)
	///////////supplier Invoice/////////////////////////////	

	JWTgroup.POST("sinvoicescart", controllers.SInvoiceController.Createscart)
	JWTgroup.GET("sinvoice/view", controllers.SInvoiceController.View)
	JWTgroup.POST("sinvoice", controllers.SInvoiceController.Create)
	JWTgroup.GET("sinvoice", controllers.SInvoiceController.GetAll)
	JWTgroup.GET("sinvoice/:id", controllers.SInvoiceController.GetOne)
	JWTgroup.POST("sinvoice/credit", controllers.SInvoiceController.Credit) 
	JWTgroup.GET("sinvoice/credit", controllers.SInvoiceController.GetCredit) 
	// JWTgroup.PUT("sinvoice/:id", controllers.SinvoiceController.Update)
	// JWTgroup.DELETE("sinvoice/:id", controllers.SinvoiceController.Delete)
	///////////trasanctions/////////////////////////////	
	JWTgroup.POST("strasanctions", controllers.STransactionController.Create)
	JWTgroup.GET("strasanctions", controllers.STransactionController.GetAll)
	JWTgroup.GET("strasanctions/:id", controllers.STransactionController.GetOne)
	JWTgroup.PUT("strasanctions/:id", controllers.STransactionController.Update)
	JWTgroup.DELETE("strasanctions/:id", controllers.STransactionController.Delete)
	//////////////////////////////////////////////////////////////////////////
	///////////////////finance module///////////////////////////////////////
	///////////payments/////////////////////////////////////////////////////
	JWTgroup.POST("payments", controllers.PaymentController.Create)
	JWTgroup.GET("payments/view", controllers.PaymentController.View)
	JWTgroup.GET("payments", controllers.PaymentController.GetAll)
	JWTgroup.GET("payments/report", controllers.PaymentController.ViewReport)
	JWTgroup.POST("payments/transaction", controllers.PaymentController.Updatepayments)
	JWTgroup.GET("payments/cleared", controllers.PaymentController.ViewCleared)
	JWTgroup.GET("payments/cleared/:code", controllers.PaymentController.ViewInvoices)
	JWTgroup.POST("payments/cleared", controllers.PaymentController.AddPaymentsTrans)
	JWTgroup.GET("payments/:id", controllers.PaymentController.GetOne)
	// JWTgroup.PUT("payments/:id", controllers.PaymentController.Update)
	// JWTgroup.DELETE("payments/:id", controllers.PaymentController.Delete)
	///////////receipts/////////////////////////////////////////////////////
	JWTgroup.POST("receipts", controllers.ReceiptController.Create)
	JWTgroup.GET("receipts/view", controllers.ReceiptController.View)
	JWTgroup.GET("receipts", controllers.ReceiptController.GetAll) 
	JWTgroup.GET("receipts/report", controllers.ReceiptController.ViewReport)
	JWTgroup.GET("receipts/cleared", controllers.ReceiptController.ViewCleared)
	JWTgroup.GET("receipts/cleared/:customercode", controllers.ReceiptController.ViewInvoices)
	JWTgroup.POST("receipts/cleared", controllers.ReceiptController.AddReceiptTrans)
	JWTgroup.POST("receipts/transaction", controllers.ReceiptController.UpdateReceipts)
	JWTgroup.POST("receipts/allocate", controllers.ReceiptController.AddReceiptTrans)
	JWTgroup.GET("receipts/:id", controllers.ReceiptController.GetOne)
	// JWTgroup.PUT("receipts/:id", controllers.ReceiptController.Update) 
	// JWTgroup.DELETE("receipts/:id", controllers.ReceiptController.Delete)
	///////////payrecpt/////////////////////////////////////////////////////
	
	JWTgroup.GET("Viewspayrecpt", controllers.PayrectrasanController.View)
	JWTgroup.POST("payrecpt", controllers.PayrectrasanController.Create)
	// JWTgroup.GET("payrecpt", controllers.PayrectrasanController.GetAll)
	JWTgroup.GET("payrecpt/:id", controllers.PayrectrasanController.GetOne)
	///////////Assets///////////////////////////////////////////////////// 
	JWTgroup.GET("assets/view", controllers.AssetController.View)
	JWTgroup.POST("assets", controllers.AssetController.Create)
	JWTgroup.GET("assets", controllers.AssetController.GetAll)
	JWTgroup.GET("assets/:id", controllers.AssetController.GetOne)
	///////////Assets/////////////////////////////////////////////////////
	JWTgroup.POST("assetstransactions", controllers.AsstransController.Create)
	JWTgroup.GET("assetstransactions", controllers.AsstransController.GetAll)
	JWTgroup.GET("assetstransactions/:id", controllers.AsstransController.GetOne)
	///////////Assets///////////////////////////////////////////////////// 
	JWTgroup.GET("liability/view", controllers.LiabilityController.View)
	JWTgroup.POST("liability", controllers.LiabilityController.Create)
	JWTgroup.GET("liability", controllers.LiabilityController.GetAll)
	JWTgroup.GET("liability/:id", controllers.LiabilityController.GetOne)
	///////////Assets/////////////////////////////////////////////////////
	JWTgroup.POST("liatransanctions", controllers.LiatranController.Create)
	JWTgroup.GET("liatransanctions", controllers.LiatranController.GetAll)
	JWTgroup.GET("liatransanctions/:id", controllers.LiatranController.GetOne)
	///////////Expence/////////////////////////////////////////////////////
	JWTgroup.POST("expence", controllers.ExpenceController.Create)
	JWTgroup.GET("expence", controllers.ExpenceController.GetAll)
	JWTgroup.GET("expence/:id", controllers.ExpenceController.GetOne)
	JWTgroup.PUT("expence/:id", controllers.ExpenceController.Update)
	JWTgroup.PUT("expence/treans/:id", controllers.ExpenceController.Update)
	JWTgroup.DELETE("expence/:id", controllers.ExpenceController.Delete)
	///////////expencetans/////////////////////////////////////////////////////
	JWTgroup.POST("expencetransanctions", controllers.ExpencetrasanController.Create)
	JWTgroup.POST("expencetransanctions/create", controllers.ExpencetrasanController.CreateExp)
	JWTgroup.GET("expencetransanctions", controllers.ExpencetrasanController.GetAll)
	JWTgroup.GET("expencetransanctions/report", controllers.ExpencetrasanController.ViewReport)
	JWTgroup.GET("expencetransanctions/view", controllers.ExpencetrasanController.ViewExp)
	JWTgroup.GET("expencetransanctions/views/:code", controllers.ExpencetrasanController.View)
	JWTgroup.GET("expencetransanctions/views", controllers.ExpencetrasanController.ViewExp)
	JWTgroup.GET("expencetransanctions/:id", controllers.ExpencetrasanController.GetOne)
	JWTgroup.POST("expence/transaction", controllers.ExpencetrasanController.UpdateTrans)
	JWTgroup.DELETE("expencetransanctions/:id", controllers.ExpencetrasanController.Delete)
	//////////////////////////////////////////////////////////////////////////ViewReport
	///////////////////Miscellenous module///////////////////////////////////////
	///////////prices/////////////////////////////////////////////////////
	JWTgroup.GET("prices/view", controllers.PriceController.View)
	JWTgroup.POST("prices", controllers.PriceController.Create)
	JWTgroup.GET("prices", controllers.PriceController.GetAll)
	JWTgroup.GET("prices/:id", controllers.PriceController.GetOne)
	JWTgroup.PUT("prices/:id", controllers.PriceController.Update)
	JWTgroup.DELETE("prices/:id", controllers.PriceController.Delete)
	///////////tax/////////////////////////////////////////////////////
	JWTgroup.POST("tax", controllers.TaxController.Create)
	JWTgroup.GET("tax", controllers.TaxController.GetAll)
	JWTgroup.GET("tax/:id", controllers.TaxController.GetOne)
	JWTgroup.PUT("tax/:id", controllers.TaxController.Update)
	JWTgroup.DELETE("tax/:id", controllers.TaxController.Delete)
	///////////discounts//////////////////////////////////////////
	JWTgroup.POST("discounts", controllers.DiscountController.Create)
	JWTgroup.GET("discounts", controllers.DiscountController.GetAll)
	JWTgroup.GET("discounts/:id", controllers.DiscountController.GetOne)
	JWTgroup.PUT("discounts/:id", controllers.DiscountController.Update)
	JWTgroup.DELETE("discounts/:id", controllers.DiscountController.Delete)
	///////////scart/////////////////////////////	

	JWTgroup.POST("scarts", controllers.ScartController.Create)
	JWTgroup.GET("scarts/view/:code", controllers.ScartController.View)
	JWTgroup.GET("scarts/:id", controllers.ScartController.GetOne)
	// JWTgroup.PUT("scart/:id", controllers.CartController.Update)
	JWTgroup.POST("scarts/credit", controllers.ScartController.Updatetrans) 
	JWTgroup.DELETE("scarts/cancel/:code", controllers.ScartController.DeleteAll)
	JWTgroup.DELETE("scarts/delete/:id", controllers.ScartController.Delete)
	JWTgroup.GET("scarts/credits/:code", controllers.ScartController.Getcredits)
	JWTgroup.GET("scarts/creditslist/:code", controllers.ScartController.GetcreditsList)
	/////////////////////////////////////////////////////////////////////
	///////////////////////////////////////////////////////////
	////////////////////reports////////////////////////////////////
	JWTgroup.GET("sales/dashboard", controllers.SalesController.View)
	JWTgroup.GET("purchases/dashboard", controllers.SalesController.Purchases)
	// JWTgroup.GET("email/create", controllers.DashboardController.Email)

	// Start server
	e.Logger.Fatal(e.Start(PORT))
}
func isAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("uname").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		isAdmin := claims["Admin"].(bool)
		if isAdmin == false {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}
func isSupervisor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("uname").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		isSupervisor := claims["Supervisor"].(bool)
		if isSupervisor == false {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}
func isEmployee(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("uname").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		isEmployee := claims["Employee"].(bool)
		if isEmployee == false {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}