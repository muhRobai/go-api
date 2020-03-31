package api

import (
	"encoding/json"
	"net/http"
	"context"
	"strconv"
	"log"
	"io"
	"os"
	
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
)

func parseCustomer(body io.ReadCloser) (*CustomerItem, error) {
	decoder := json.NewDecoder(body)

	var payload *CustomerItem
	err := decoder.Decode(payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (c *initAPI) GetCustomerHanler(w http.ResponseWriter, r *http.Request){
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	limit := r.FormValue("limit")
	id := r.FormValue("id")
	name := r.FormValue("name")

	payload := ListCustomer{
		Limit: limit,
		Search: &SearchValue{
			Id: id,
			Name: name,
		},
	}

	resp, err := c.getCustomer(ctx, &payload)

	if err != nil {
		log.Println(err)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return 
	}
	w.Header().Set("Content-Type","application/json")
	w.Write(data)
}

func (c *initAPI) InsertCustomerHandler(w http.ResponseWriter, r *http.Request){
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	payload, err := parseCustomer(r.Body)
	resp, err := c.AddCustomer(ctx, payload)
	if err != nil {
		log.Println(err)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
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
		Port: uint16(port),
		Host: dbHost,
		User: dbUser,
		Password: dbPass,
		Database: dbName,
	}

	connection := pgx.ConnPoolConfig{
		ConnConfig: *dbConfig,
		MaxConnections: 5,
	}

	c.Db, err = pgx.NewConnPool(connection)
	if err != nil {
		log.Println(err.Error())
		return 
	}
}

func StartHTTP() http.Handler{
	api, err := createAPI()
	if err != nil {
		log.Println(err)
		return nil
	}

	api.initDB()

	r := mux.NewRouter()
	r.HandleFunc("/api/list", api.GetCustomerHanler).Methods("POST")

	return r
}