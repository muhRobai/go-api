package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
)

func (c *initAPI) getCustomer(ctx context.Context, req *ListCustomer) (*ListCustomer, error) {
	var err error

	strQuery := `SELECT id, name, phone_number, email, customer_type from customer`

	var params []interface{}
	count := 0
	if req.Search != nil {
		strQuery = fmt.Sprintf("%s WHERE TRUE", strQuery)
	}

	if req.Search.Id != "" {
		count++
		strQuery = fmt.Sprintf("%s AND id=$%d", strQuery, count)
		params = append(params, req.Search.Id)
	}

	if req.Search.Name != "" {
		count++
		strQuery = fmt.Sprintf("%s AND name=$%d", strQuery, count)
		params = append(params, req.Search.Name)
	}

	limit := 5
	if req.Limit != "" {
		limit, err = strconv.Atoi(req.Limit)
		if err != nil {
			log.Println("invalid-parse-limit")
			return nil, errors.New("invalid-parse-limit")
		}
	}

	strQuery = fmt.Sprintf("%s LIMIT %d", strQuery, limit)
	log.Println(strQuery)
	rows, err := c.Db.Query(strQuery, params...)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var item []*CustomerItem

	for rows.Next() {
		var data CustomerItem
		var type_ int
		err := rows.Scan(&data.Id, &data.CustomerName, &data.PhoneNumber, &data.EmailAddress, &type_)

		if err != nil {
			log.Println(err)
			return nil, errors.New("unable-get-data")
		}
		data.CustomerType = strconv.Itoa(type_)
		item = append(item, &data)
	}

	result := &ListCustomer{
		List:  item,
		Total: int32(len(item)),
	}

	return result, nil
}

func (c *initAPI) AddCustomer(ctx context.Context, req *CustomerItem) (*CustomerItem, error) {
	if req == nil {
		return nil, errors.New("missing-request")
	}

	if req.CustomerName == "" {
		return nil, errors.New("missing-name")
	}

	if req.PhoneNumber == "" {
		return nil, errors.New("missing-phone-number")
	}

	if req.EmailAddress == "" {
		return nil, errors.New("missing-email-addres")
	}

	type_ := "BUYER"
	if req.CustomerType != "" {
		type_ = req.CustomerType
	}

	customerType := parseCustomerType(type_)

	var id string
	err := c.Db.QueryRow(`
		INSERT INTO customer
			(name, phone_number, email, customer_type)
		VALUES ($1, $2, $3, $4) RETURNING id
	`, req.CustomerName, req.PhoneNumber, req.EmailAddress, customerType).Scan(&id)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &CustomerItem{
		Id: id,
	}, nil
}

//soon
func (c *initAPI) UpdateCustomers(ctx context.Context, req *CustomerItem) error {
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

func (c *initAPI) DeleteCustomer(ctx context.Context, req *CustomerItem) error {
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

func (c *initAPI) InsertAddressInfo(ctx context.Context, req *AddressItem) (*AddressResponse, error) {
	if req.Id == "" {
		return nil, errors.New("missing-id-customer")
	}

	if req.Address == "" {
		return nil, errors.New("missing-address")
	}

	if req.PostNumber == "" {
		return nil, errors.New("missing-post-number")
	}

	address := parseAdddressType(req.AddressType)

	if address == 0 {
		return nil, errors.New("address-type-not-found")
	}

	var id string
	err := c.Db.QueryRow(`
		INSERT INTO address_info (customer_id, address, post_number, address_type) 
		VALUES ($1, $2, $3, $4) RETURNING customer_id
	`, req.Id, req.Address, req.PostNumber, address).Scan(&id)

	if err != nil {
		return nil, err
	}

	if id == "" {
		return nil, errors.New("unable-insert-address")
	}

	return &AddressResponse{
		Id: id,
	}, nil
}

func (c *initAPI) GetAddressInfo(ctx context.Context, req *CustomerId) (*AddressList, error) {
	if req.Id == "" {
		return nil, errors.New("missing-customer-id")
	}

	rows, err := c.Db.Query(`
		SELECT customer_id, 
			address, 
			post_number, 
			address_type 
		FROM address_info
		WHERE customer_id = $1
	`, req.Id)

	if err != nil {
		return nil, err
	}

	var item []*AddressItem
	for rows.Next() {
		var data AddressItem
		var type_ int
		err = rows.Scan(&data.Id,
			&data.Address,
			&data.PostNumber,
			&type_,
		)

		if err != nil {
			return nil, err
		}
		data.AddressType = strconv.Itoa(type_)
		item = append(item, &data)
	}

	return &AddressList{
		List: item,
	}, nil

}

func parseAdddressType(payload string) int32 {
	return addressType_value[payload]
}

func parseCustomerType(payload string) int32 {
	return customerType_value[payload]
}
