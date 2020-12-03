package model

import (
	"github.com/myrachanto/accounting/httperors"
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
func (search Search) Validate() *httperors.HttpError{ 
	if search.Column == "" {
		return httperors.NewNotFoundError("Invalid Column")
	}
	if search.Direction == "" {
		return httperors.NewNotFoundError("Invalid Direction")
	}
	if search.Search_column == "" {
		return httperors.NewNotFoundError("Invalid Search column")
	}
	if search.Search_operator == "" {
		return httperors.NewNotFoundError("Invalid Search operator")
	}	
	if search.Per_page <= 0 {
		return httperors.NewNotFoundError("Invalid Per page")
	}
	return nil
}