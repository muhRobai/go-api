package api

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"
)

func (c *initAPI) GetProduct(ctx context.Context, req *ProductList) (*ProductList, error) {

	query := `SELECT id, product_name, stock, store_id, amount FROM product`

	if req.Search != nil {
		query = fmt.Sprintf("%s WHERE TRUE", query)
	}

	var params []interface{}
	count := 0

	if req.Search.Id != "" {
		count++
		query = fmt.Sprintf("%s AND id=$%d", query, count)
		params = append(params, req.Search.Id)
	}

	if req.Search.StoreId != "" {
		count++
		query = fmt.Sprintf("%s AND store_id=$%d", query, count)
		params = append(params, req.Search.StoreId)
	}

	if req.Search.ProductName != "" {
		count++
		query = fmt.Sprintf("%s AND product_name LIKE $%d", query, count)
		params = append(params, req.Search.ProductName)
	}

	orderBy := "created_time"
	order := "DESC"

	if req.OrderBy != "" {
		orderBy = req.OrderBy
	}

	if req.Order != "" {
		order = req.Order
	}

	query = fmt.Sprintf("%s ORDER BY %s %s ", query, orderBy, order)
	
	count++
	limit := 50
	var err error
	if req.Limit != "" {
		limit, err = strconv.Atoi(req.Limit)
		if err != nil {
			return nil, err
		}
	}

	query = fmt.Sprintf("%s LIMIT $%d", query, count)
	params = append(params, limit)

	rows, err := c.Db.Query(query, params...)

	if err != nil {
		return nil, err
	}

	var items []*ProductItem

	for rows.Next() {
		var data ProductItem
		var stock int
		err = rows.Scan(&data.Id, &data.ProductName, &stock, &data.StoreId, &data.Amount)
		if err != nil {
			return nil, err
		}

		data.Stock = strconv.Itoa(stock)
		items = append(items, &data)
	}

	return &ProductList{
		List:  items,
		Limit: req.Limit,
	}, nil
}

func (c *initAPI) CreateProduct(ctx context.Context, req *ProductItem) (*ProductItem, error) {
	if req.StoreId == "" {
		return nil, errors.New("missing-store-id")
	}

	if req.Stock == "" {
		return nil, errors.New("missing-stock")
	}

	if req.Amount == "" {
		return nil, errors.New("missing-amount")
	}

	if req.ProductName == "" {
		return nil, errors.New("missing-product-name")
	}

	stock, err := strconv.Atoi(req.Stock)
	if err != nil {
		return nil, err
	}

	var id string
	err = c.Db.QueryRow(`
		INSERT INTO product (store_id, product_name, stock, amount, created_time) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`, req.StoreId, req.ProductName, stock, req.Amount, time.Now()).Scan(&id)

	if err != nil {
		return nil, err
	}

	return &ProductItem{
		Id: id,
	}, nil
}
