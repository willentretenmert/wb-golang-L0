package main

import (
	"fmt"
	"github.com/jackc/pgx"
	"os"
	"sync"
)

func insertOrder(conn *pgx.Conn, order Order) error {

	tx, err := conn.Begin()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to begin transaction: %v\n", err)
		os.Exit(1)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	var exists bool
	err2 := conn.QueryRow("SELECT EXISTS(SELECT 1 FROM orders WHERE order_uid=$1)", order.OrderUID).Scan(&exists)

	if exists {
		return err2
	}
	//adding order
	_, err = tx.Exec(`
        INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        `,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.Shardkey,
		order.SmID,
		order.DateCreated,
		order.OofShard,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to insert order: %v\n", err)
		os.Exit(1)
	}

	_, err = tx.Exec(`
        INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        `,
		order.OrderUID,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to insert delivery: %v\n", err)
		os.Exit(1)
	}

	_, err = tx.Exec(`
        INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        `,
		order.OrderUID,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDt,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to insert payment: %v\n", err)
		os.Exit(1)
	}

	for i := 0; i < len(order.Items); i++ {
		_, err = tx.Exec(`
        INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
        `,
			order.OrderUID,
			order.Items[i].ChrtID,
			order.Items[i].TrackNumber,
			order.Items[i].Price,
			order.Items[i].RID,
			order.Items[i].Name,
			order.Items[i].Sale,
			order.Items[i].Size,
			order.Items[i].TotalPrice,
			order.Items[i].NmID,
			order.Items[i].Brand,
			order.Items[i].Status,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to insert item (%d): %v\n", i, err)
			os.Exit(1)
		}
	}

	err = tx.Commit()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to commit transaction: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Order and associated data added successfully.")

	return err
}

func loadOrder(conn *pgx.Conn, orderUID string) (Order, error) {
	var order Order
	order.Items = []Item{}

	query := `SELECT
    o.order_uid, o.track_number, o.entry, o.locale,
    o.internal_signature, o.customer_id, o.delivery_service,
    o.shardkey, o.sm_id, o.date_created, o.oof_shard,
    d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
    p.transaction, p.request_id, p.currency, p.provider,
    p.amount, p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee,
    i.chrt_id, i.track_number, i.price, i.rid, i.name,
    i.sale, i.size, i.total_price, i.nm_id, i.brand, i.status
FROM
    orders o
        JOIN
    delivery d ON o.order_uid = d.order_uid
        JOIN
    payment p ON o.order_uid = p.order_uid
        JOIN
    items i ON o.order_uid = i.order_uid
WHERE
        o.order_uid = $1;`
	rows, err := conn.Query(query, orderUID)
	if err != nil {
		return order, err
	}
	defer rows.Close()

	for rows.Next() {
		var item Item
		var delivery Delivery
		var payment Payment
		// заполните прочие поля, которые вам нужны

		err = rows.Scan(
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.Shardkey,
			&order.SmID,
			&order.DateCreated,
			&order.OofShard,

			&delivery.Name,
			&delivery.Phone,
			&delivery.Zip,
			&delivery.City,
			&delivery.Address,
			&delivery.Region,
			&delivery.Email,

			&payment.Transaction,
			&payment.RequestID,
			&payment.Currency,
			&payment.Provider,
			&payment.Amount,
			&payment.PaymentDt,
			&payment.Bank,
			&payment.DeliveryCost,
			&payment.GoodsTotal,
			&payment.CustomFee,

			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.RID,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			return order, err
		}

		if len(order.Items) == 0 {
			order.Delivery = delivery
			order.Payment = payment
		}

		order.Items = append(order.Items, item)
	}

	if err = rows.Err(); err != nil {
		return order, err
	}

	return order, nil
}

func loadCacheFromDB(conn *pgx.Conn, cache *sync.Map) error {
	// Запрашиваем список уникальных orderUID из таблицы orders.
	rows, err := conn.Query("SELECT DISTINCT order_uid FROM orders")
	if err != nil {
		return err
	}
	defer rows.Close()

	var orderUIDs []string
	for rows.Next() {
		var orderUID string
		if err := rows.Scan(&orderUID); err != nil {
			return err
		}
		orderUIDs = append(orderUIDs, orderUID)
	}

	for _, uid := range orderUIDs {
		order, err := loadOrder(conn, uid)
		if err != nil {
			return err
		}
		cache.Store(uid, order)
		fmt.Printf("Order %s loaded into cache\n", uid)
	}

	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}
