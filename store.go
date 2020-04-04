package api

import (
	"context"
	"errors"
)

func (c *initAPI) CreateStore(ctx context.Context, req *StoreItem) (*StoreId, error) {
	if req.CustomerId == "" {
		return nil, errors.New("missing-customer-id")
	}

	if req.Name == "" {
		return nil, errors.New("missing-store-name")
	}

	var id string
	err := c.Db.QueryRow(`
		INSERT INTO store (customer_id, name) VALUES ($1, $2) RETURNING id
	`, req.CustomerId, req.Name).Scan(&id)

	if err != nil {
		return nil, err
	}

	return &StoreId{
		Id: id,
	}, nil
}
