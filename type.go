package api

import (
	"github.com/jackc/pgx"
)

type initAPI struct {
	Db *pgx.ConnPool
}

type ListCustomer struct {
	Limit  string `json:"limit"`
	Page   string `json:"page"`
	List   []*CustomerItem
	Total  int32        `json:"total"`
	Search *SearchValue `json:"search"`
}

type SearchValue struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type CustomerItem struct {
	Id           string `json:"id"`
	CustomerName string `json:"customerName"`
	PhoneNumber  string `json:"phonenumber"`
	EmailAddress string `json:"emailaddres"`
	CustomerType string  `json:"customertype"`
}

type OrderItem struct {
	Buyer         string `json:"buyer"`
	Product       string `json:"product"`
	StoreId       string `json:"storeId"`
	Status        int32  `json:"status"`
	PaymentType   string `json:"paymentType"`
	PaymentStatus string `json:"paymentStatus"`
	ProductLength string `json:"productLength"`
}

type OrderResponse struct {
	NumberVA      string `json:"numberVA"`
	ExpiryTime    int64  `json:"expiryTime"`
	EffectiveTime int64  `json:"effectiveTime"`
	OrderId       string `json:"orderId"`
}

type Product struct {
	Id          string `json:"idProduct"`
	ProductName string `json:"productName"`
	Stock       string `json:"stock"`
}

type AddressItem struct {
	Id          string `json:"customerId"`
	Address     string `json:"address"`
	PostNumber  string `json:"postNumber"`
	AddressType string `json:"addressType"`
}

type AddressList struct {
	List []*AddressItem `json:"list"`
}

type AddressResponse struct {
	Id string `json:"customerId"`
}

type CustomerId struct {
	Id string `json:"customerId"`
}

type OrderUpdate struct {
	OrderId     string `json:"orderId"`
	OrderStatus string `json:"orderStatus"`
}

type OrderId struct {
	Id string `json:"orderId"`
}

type ProductItem struct {
	Id          string `json:"id"`
	StoreId     string `json:"storeId"`
	ProductName string `json:"ProductName"`
	Stock       string `json:"stock"`
	Amount      string `json:"amount"`
}

type SearchProductValue struct {
	Id          string `json:"id"`
	StoreId     string `json:"storeId"`
	ProductName string `json:"productName"`
}

type ProductList struct {
	Limit   string              `json:"limit"`
	Search  *SearchProductValue `json:"search"`
	OrderBy string              `json:"orderBy"`
	Order   string              `json:"order"`
	List    []*ProductItem      `json:"list"`
}

type StoreItem struct {
	Name string `json:"name"`
	CustomerId string `json:"customerId"`
}

type StoreId struct {
	Id string `json:"id"`
}

//enum for status payment
var orderStatus_value = map[string]int32{
	"UNKNOW_ORDER_STATUS": 0,
	"NEW":                 1,
	"COMPLATE":            3,
	"ERROR":               4,
	"PENDING":             5,
}

//enum payment type
var paymentType_value = map[string]int32{
	"UNKNOW_PAYMENT_TYPE": 0,
	"SUBSCRIBE":           1,
	"ONCE":                2,
}

//enum address type
var addressType_value = map[string]int32{
	"UNKNOW_ADDRESS_TYPE": 0,
	"HOME":                1,
	"OFFICE":              2,
	"STORE":               3,
}

//enum customer type
var customerType_value = map[string]int32{
	"UNKNOW_CUSTOMER_TYPE": 0,
	"BUYER":                1,
	"SELLER":               2,
}
