package model

import ()
//Dashboard representation
type Dashboard struct {
	Products Module `json:"products"`
	Categorys Module `json:"categorys"`
	Majorcategorys Module `json:"majorcats"`
	Customers Module  `json:"customers"`
	Users Module `json:"users"`
	Invoices Module `json:"invoices"`
}
//Module structure of dashboard items
type Module struct{
	Name string
	Total int
	Description string
	Icon string
}
//Email structure
type Email struct{
	Email string
	To string
	Subject string
	Message string
	Customers []Customer
}