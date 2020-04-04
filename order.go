package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"log"

	"github.com/jackc/pgx"
)

func (c *initAPI) UpdateOrder(ctx context.Context, req *OrderUpdate) (*OrderId, error) {
	if req.OrderId == "" {
		return nil, errors.New("missing-order-id")
	}

	status := parseStatusOrder(req.OrderStatus)

	if status == 0 {
		return nil, errors.New("unknow-status")
	}

	var id string
	err := c.Db.QueryRow(`
		UPDATE orders SET status=$1 WHERE id = $2 RETURNING id
	`, status, req.OrderId).Scan(&id)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if id == "" {
		return nil, errors.New("unable-update-order")
	}

	return &OrderId{
		Id: id,
	}, nil
}

func (c *initAPI) CreateOrder(ctx context.Context, req *OrderItem) (*OrderResponse, error) {
	if req.Buyer == "" {
		return nil, errors.New("missing-buyer")
	}

	if req.Product == "" {
		return nil, errors.New("missing-product")
	}

	if req.StoreId == "" {
		return nil, errors.New("missing-store-id")
	}

	payment := parsePaymentType(req.PaymentType)
	if payment == 0 {
		return nil, errors.New("missing-payment-type")
	}

	resp, err := c.getCustomer(ctx, &ListCustomer{
		Limit: "1",
		Search: &SearchValue{
			Id: req.Buyer,
		},
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(resp.List) < 1 {
		return nil, errors.New("customer-not-found")
	}

	customer := resp.List[0]

	address, err := c.GetAddressInfo(ctx, &CustomerId{
		Id: customer.Id,
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(address.List) < 1 {
		return nil, errors.New("address-missing")
	}

	tx, err := c.Db.Begin()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer tx.Rollback()

	status := parseStatusOrder("NEW")

	now := time.Now()
	effectiveTime := now
	expiryTime := now.Add(72 * time.Hour)
	numberVA := fmt.Sprintf("%s%s", generateVANumber(), customer.PhoneNumber)

	paymentCode, err := json.Marshal(struct {
		ProcessorCode string `json:"processorCode"`
		ProductLength string `json:"productLength"`
		StoreId       string `josn:"storeId"`
		PaymentItem   string `json:"paymentItem"`
	}{numberVA, req.ProductLength, req.StoreId, req.Product})

	var id string
	err = tx.QueryRow(`
		INSERT INTO orders (buyer, created_time, status, payment_type, payment_code, effective_time, expiry_time)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
	`, req.Buyer, time.Now(), status, payment, paymentCode, effectiveTime, expiryTime).Scan(&id)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if id == "" {
		return nil, errors.New("unable-create-order")
	}

	product, err := strconv.Atoi(req.ProductLength)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = updateProduct(tx, req.Product, req.StoreId, product)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &OrderResponse{
		NumberVA:      numberVA,
		ExpiryTime:    expiryTime.Unix(),
		EffectiveTime: effectiveTime.Unix(),
		OrderId:       id,
	}, nil
}

func parseStatusOrder(payload string) int32 {
	return orderStatus_value[payload]
}

func parsePaymentType(payload string) int32 {
	return paymentType_value[payload]
}

func generateVANumber() string {
	var VA string
	for i := 0; i < 5; i++ {
		VA = fmt.Sprintf("%s%d", VA, rand.Intn(10))
	}
	return VA
}

func updateProduct(tx *pgx.Tx, productId, storeId string, productLength int) error {
	var stock int
	err := tx.QueryRow(`
		SELECT stock FROM product WHERE id=$1 AND store_id= $2
	`, productId, storeId).Scan(&stock)

	if err != nil {
		return err
	}

	stockProduct := stock - productLength
	if stockProduct < 1 {
		return errors.New("stock-product-not-found")
	}

	tag, err := tx.Exec(`
		UPDATE product SET stock= $1 WHERE id= $2
	`, stockProduct, productId)

	if err != nil {
		return err
	}

	if tag.RowsAffected() != 1 {
		return errors.New("unable-to-update-product")
	}

	return nil

}
