package model

import ()
//Sales representation
type Sales struct {
	Sales SalesModule `json:"sales"`
	GrossProfit SalesModule `json:"grossprofit"`
	Transactions []Transaction `json:"transactions"`
	DebtTransactions []DebtTransaction `json:"debtors"`
	Debts SalesModule  `json:"debts"`
	PaidInvoices SalesModule `json:"paid"`
	CreditNotes SalesModule `json:"creditnotes"`
}
//SalesModule structure of dashboard items
type SalesModule struct{
	Name string 
	Total float64 
	Description string
	Icon string 
}
// //Email structure
// type Email struct{
// 	Email string
// 	To string
// 	Subject string
// 	Message string
// 	Customers []Customer
// }