package api

import (
	"github.com/jackc/pgx"
)

type initAPI struct{
	Db 		*pgx.ConnPool
}

type ListCustomer struct{
	Limit 		string	`json:"limit"`
	Page 		string	`json:"page"`
	List 		[]CustomerItem
	Total		int32	`json:"total"`
	Search		*SearchValue `json:"search"`
}

type SearchValue struct{
	Id 			string	`json:"id"`
	Name		string	`json:"name"` 
}

type CustomerItem struct{	
	Id				string 	`json:"id"`
	CustomerName	string	`json:"customerName"`
	BirthDate		int64	`json:"birthDate"`
}
