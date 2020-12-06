package model

import ()
//Dashboard representation
type Dashboard struct {
	Sales Module `json:"sales"`
	Purchases Module `json:"purchases"`
	Payments Module `json:"payments"`
	Receipts Module  `json:"receipts"`
	Expences Module `json:"expences"`
	Customers Module `json:"customers"`
	Suppliers Module `json:"suppliers"`
}
//Module structure of dashboard items
type Module struct{
	Name string
	Total float64
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