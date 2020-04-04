package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
)

//done
func (c *initAPI) GetCustomerHanler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	limit := r.FormValue("limit")
	id := r.FormValue("id")
	name := r.FormValue("name")

	payload := ListCustomer {
		Limit: limit,
		Search: &SearchValue{
			Id:   id,
			Name: name,
		},
	}

	resp, err := c.getCustomer(ctx, &payload)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "faild-convert-json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

//done
func (c *initAPI) InsertCustomerHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	customerName := r.FormValue("customerName")
	phoneNumber := r.FormValue("phoneNumber")
	emailAddress := r.FormValue("emailAddress")
	customerType := r.FormValue("customerType")

	resp, err := c.AddCustomer(ctx, &CustomerItem{
		CustomerName: customerName,
		PhoneNumber: phoneNumber,
		EmailAddress: emailAddress,
		CustomerType: customerType,
	})
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
		http.Error(w, "faild-convert-json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

//done
func (c *initAPI) InsertAddressHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	address := r.FormValue("address")
	postNumber := r.FormValue("postNumber")
	addressType := r.FormValue("addressType")
	id := r.FormValue("customerId")

	resp, err := c.InsertAddressInfo(ctx, &AddressItem{
		Id: id,
		Address: address,
		PostNumber: postNumber,
		AddressType: addressType,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "faild-convert-json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

//done
func (c *initAPI) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	buyer := r.FormValue("buyer")
	product := r.FormValue("product")
	store := r.FormValue("storeId")
	paymentType:= r.FormValue("paymentType")
	productLength := r.FormValue("productLength")

	resp, err := c.CreateOrder(ctx, &OrderItem{
		Buyer: buyer,
		Product: product,
		StoreId: store,
		PaymentType: paymentType,
		ProductLength: productLength,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "faild-convert-json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

//done
func (c *initAPI) ProcessOrderHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	id := r.FormValue("orderId")
	status := r.FormValue("orderStatus")

	resp, err := c.UpdateOrder(ctx, &OrderUpdate{
		OrderId: id,
		OrderStatus: status,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "faild-convert-json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

//done
func (c *initAPI) CreateStoreHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	name := r.FormValue("storeName")
	id := r.FormValue("customerId")

	resp, err := c.CreateStore(ctx, &StoreItem{
		Name: name,
		CustomerId: id,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "faild-convert-json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

//done
func (c *initAPI) GetproductHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	limit := r.FormValue("limit")
	id := r.FormValue("id")
	store := r.FormValue("storeId")
	productName := r.FormValue("productName")
	order := r.FormValue("order")
	orderBy := r.FormValue("orderBy")

	resp, err := c.GetProduct(ctx, &ProductList{
		Limit: limit,
		Search: &SearchProductValue{
			Id: id,
			StoreId: store,
			ProductName: productName,
		},
		Order: order,
		OrderBy: orderBy,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "faild-convert-json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

//done
func (c *initAPI) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	storeId := r.FormValue("storeId")
	productName := r.FormValue("productName")
	stock := r.FormValue("stock")
	amount := r.FormValue("amount")

	resp, err := c.CreateProduct(ctx, &ProductItem{
		StoreId: storeId,
		ProductName: productName,
		Stock: stock,
		Amount: amount,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "faild-convert-json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (c *initAPI) initDB() {

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	port, err := strconv.Atoi(dbPort)
	if err != nil {
		log.Println(err.Error())
		return
	}

	dbConfig := &pgx.ConnConfig{
		Port:     uint16(port),
		Host:     dbHost,
		User:     dbUser,
		Password: dbPass,
		Database: dbName,
	}

	connection := pgx.ConnPoolConfig{
		ConnConfig:     *dbConfig,
		MaxConnections: 5,
	}

	c.Db, err = pgx.NewConnPool(connection)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func StartHTTP() http.Handler {
	api, err := createAPI()
	if err != nil {
		log.Println(err)
		return nil
	}

	api.initDB()

	r := mux.NewRouter()
	//get customer list
	r.HandleFunc("/api/customer/list", api.GetCustomerHanler).Methods("GET")
	//create customer handler
	r.HandleFunc("/api/customer/create", api.InsertCustomerHandler).Methods("POST")
	//create customer address
	r.HandleFunc("/api/address/create", api.InsertAddressHandler).Methods("POST")
	//create order
	r.HandleFunc("/api/order", api.CreateOrderHandler).Methods("POST")
	//do payment orders
	r.HandleFunc("/api/order/update", api.ProcessOrderHandler).Methods("POST")
	//create store
	r.HandleFunc("/api/store/create", api.CreateStoreHandler).Methods("POST")
	//product list
	r.HandleFunc("/api/product/list", api.GetproductHandler).Methods("GET")
	//create product
	r.HandleFunc("/api/product/create", api.CreateProductHandler).Methods("POST")
	return r
}
