package model

import ()

type Pl struct {
	Sales float64 `json:"sales"`
	Costofsale float64 `json:"cost"`
	DirectExpence float64 `json:"directExpences"`
	InDirectExpence float64 `json:"indirectExpences"`
	OtherExpence float64 `json:"otherExpences"`
	Dated string `json:"dated"`
}