package model

import (
	"time"
)
//CustomerState customer statement structure
type CustomerState struct {
	Description string `json:"description"`
	Tax float64 `json:"tax"`
	Discount float64 `json:"discount"`
	Total float64 `json:"total"`
	Balance float64 `json:"balance"`
	Dated time.Time `json:"dated"`
	Amount string `json:"amount"`
	Status string `json:"status"`
	CreatedAt string `json:"created_at"`
}
//SupplierState customer statement structure
type SupplierState struct {
	Description string `json:"description"`
	Tax float64 `json:"tax"`
	Discount float64 `json:"discount"`
	Total float64 `json:"total"`
	Balance float64 `json:"balance"`
	Dated time.Time `json:"dated"`
	Amount string `json:"amount"`
	Status string `json:"status"`
	CreatedAt string `json:"created_at"`
}
//Customerstates view
type Customerstates struct{
 Statements []CustomerState `json:"statements"`
 Customer *Customer `json:"customer"`
 Datos string `json:"datos"`
 Query1 string `json:"query1"`
 Query2 string `json:"query2"`
}
//Supplierstates view
type Supplierstates struct{
	Statements []CustomerState `json:"statements"`
	Supplier *Supplier `json:"supplier"`
	Datos string `json:"datos"`
	Query1 string `json:"query1"`
	Query2 string `json:"query2"`
 }