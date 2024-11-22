package database

import (
	"context"
	"log"
	"log/slog"

	"github.com/chnmk/order-info-l0/internal/models"
	"github.com/jackc/pgx/v5"
)

/*
Вставляет в таблицу данные о заказе и связывает их с самыми новыми записями в delivery и payment.

Запросы q_insert_delivery и q_insert_item должны быть выполнены раньше.
Будем пока считать что никакой конкурентности нет. При необходимости воспользуемся транзакциями.
*/
var q_insert_order = `
	INSERT INTO orders(order_uid, track_number, entry, locale, internal_signature, customer_id, 
		delivery_service, shardkey, sm_id, date_created, oof_shard, delivery_id, payment_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	RETURNING order_uid
`
var q_insert_delivery = `
	INSERT INTO delivery(name, phone, zip, city, address, region, email)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id
`
var q_insert_payment = `
	INSERT INTO payments(transaction, request_id, currency, provider, amount, payment_dt, bank, 
	delivery_cost, goods_total, custom_fee)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id
`
var q_insert_item = `
	INSERT INTO items(chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	RETURNING id
`

// Должен быть выполнен после запроса q_insert_item (как минимум пока транзакции не используются).
var q_insert_itemsbind = `
	INSERT INTO itemsbind(order_uid, item_id)
	VALUES ($1, $2)
	RETURNING id
`

func InsertOrder(db *pgx.Conn, order models.Order) {
	slog.Info("adding order to database...")

	row := db.QueryRow(context.Background(), q_insert_delivery,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	)

	var delivery_id int
	err := row.Scan(&delivery_id)
	if err != nil {
		log.Fatalf("Failed to insert data: %v\n", err)
	}

	row = db.QueryRow(context.Background(), q_insert_payment,
		order.Payment.Transaction,
		order.Payment.Request_id,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.Payment_dt,
		order.Payment.Bank,
		order.Payment.Delivery_cost,
		order.Payment.Goods_total,
		order.Payment.Custom_fee,
	)

	var payment_id int
	err = row.Scan(&payment_id)
	if err != nil {
		log.Fatalf("Failed to insert data: %v\n", err)
	}

	row = db.QueryRow(context.Background(), q_insert_order,
		order.Order_uid,
		order.Track_number,
		order.Entry,
		order.Locale,
		order.Internal_signature,
		order.Customer_id,
		order.Delivery_service,
		order.Shardkey,
		order.Sm_id,
		order.Date_created,
		order.Oof_shard,
		delivery_id,
		payment_id,
	)

	// order_row используется только для проверки ответа
	var order_id string
	err = row.Scan(&order_id)
	if err != nil {
		log.Fatalf("Failed to insert data: %v\n", err)
	}

	for _, i := range order.Items {
		row = db.QueryRow(context.Background(), q_insert_item,
			i.Chrt_id,
			i.Track_number,
			i.Price,
			i.Rid,
			i.Name,
			i.Sale,
			i.Size,
			i.Total_price,
			i.Nm_id,
			i.Brand,
			i.Status,
		)

		var item_id int
		err = row.Scan(&item_id)
		if err != nil {
			log.Fatalf("Failed to insert data: %v\n", err)
		}

		row = db.QueryRow(context.Background(), q_insert_itemsbind, order.Order_uid, item_id)
		// bind_id используется только для проверки ответа
		var bind_id int
		err = row.Scan(&bind_id)
		if err != nil {
			log.Fatalf("Failed to insert data: %v\n", err)
		}
	}

	slog.Info("order successfully added to database")
}
