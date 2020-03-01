package api

import (
	"strconv"
	"log"
	"errors"
	"time"
	"fmt"
	"context"
)

func(c *initAPI) getCustomer(ctx context.Context, req *ListCustomer) (*ListCustomer, error) {
	var err error
	limit := 5
	if req.Limit != "" {
		limit, err = strconv.Atoi(req.Limit)
		if err != nil {
			log.Println("invalid-parse-limit")
			return nil, errors.New("invalid-parse-limit")
		}		
	}

	rows, err := c.Db.Query(`
		SELECT * FROM 
		customer LIMIT $1
	`, limit)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	
	var item []CustomerItem

	for rows.Next() {
		var data CustomerItem
		var date time.Time
		err := rows.Scan(&data.Id, &data.CustomerName, &date)

		if err != nil {
			log.Println(err)
			return nil, errors.New("unable-get-data")
		}

		data.BirthDate = date.Unix()
		item = append(item, data)
	}

	result := &ListCustomer{
		List: item,
		Total: int32(len(item)),
	}

	return result, nil
}

func(c *initAPI) AddCustomer(ctx context.Context, req *CustomerItem) (*CustomerItem, error){
	if req.CustomerName == "" {
		return nil, errors.New("missing-name")
	}

	if req.BirthDate == 0 {
		return nil, errors.New("missing-birthdate")
	}

	birthDate := time.Unix(req.BirthDate, 0).Format("2006-01-02")

	var id string
	err := c.Db.QueryRow(`
		INSERT INTO customers 
			(name, birthdate)
		VALUES ($1, $2) RETURNING id
	`, req.CustomerName, birthDate).Scan(&id)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &CustomerItem{
		Id: id,
	}, nil
}

func(c *initAPI) UpdateCustomers(ctx context.Context, req *CustomerItem) error {
	if req == nil {
		return errors.New("missing-request")
	}

	queryStr := "UPDATE customer"
	var params []interface{}
	value := 0

	if req.CustomerName != "" {
		value++
		queryStr = fmt.Sprintf("%s SET name=$%d", queryStr, value)
		params = append(params, req.CustomerName)
	}

	if req.BirthDate != 0 {
		value++
		birthDate := time.Unix(req.BirthDate, 0).Format("2006-01-02")
		queryStr = fmt.Sprintf("%s SET birthdate=$%d", queryStr, value)
		params = append(params, birthDate)
	}

	if len(params) < 1 {
		return errors.New("missing-request")
	}

	_, err := c.Db.Exec(queryStr, params...)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func(c *initAPI) DeleteCustomer(ctx context.Context, req *CustomerItem) error {
	if req.Id == "" {
		return errors.New("missing-request")
	}

	_, err := c.Db.Exec(`
		DELETE FROM customer WHERE id=$1
	`, req.Id)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
} 