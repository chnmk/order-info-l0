package database

import (
	"context"
	"log"

	"github.com/chnmk/order-info-l0/internal/models"
	"github.com/jackc/pgx/v5"
)

/*
Строка для получения данных из таблицы orders.

Не получает данные из таблиц delivery, payments и items.
*/
var q_orders = `
	SELECT order_uid, track_number, entry,
	locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
	FROM orders 
	WHERE order_uid = $1
`

/*
Строка для получения данных из таблицы delivery.
*/
var q_delivery = `
	SELECT name, phone, zip, city, address, region, email
	FROM delivery d
	LEFT JOIN orders o
		ON d.id = o.delivery_id
	WHERE o.order_uid = $1
`

/*
Строка для получения данных из таблицы payments.
*/
var q_payments = `
	SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
	FROM payments p
	LEFT JOIN orders o
		ON p.id = o.payment_id
	WHERE o.order_uid = $1
`

/*
Строка для получения данных из items.
*/
var q_items = `
	SELECT chrt_id, i.track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
	FROM items i
	LEFT JOIN itemsbind ib
		on i.id = ib.item_id
	LEFT JOIN orders o
		on ib.order_uid = o.order_uid
	WHERE o.order_uid = $1
`

/*
// Строка для получения всех данных. Не может использоваться в таком виде, т.к. из items необходимо получить несколько записей.
var q_full = `
	SELECT o.order_uid, o.track_number, o.entry,
	d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
	p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee,
	i.chrt_id, i.track_number, i.price, i.rid, i.name, i.sale, i.size, i.total_price, i.nm_id, i.brand, i.status,
	o.locale, o.internal_signature, o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard
	FROM orders o
	LEFT JOIN delivery d
		ON o.delivery_id = d.id
	LEFT JOIN payment p
		ON o.payment_id = p.id
	LEFT JOIN itemsbind ib
		on o.order_uid = ib.order_uid
	LEFT JOIN items i
		on ib.item_id = i.id
	WHERE o.order_uid = $1
`

// Строка для получения всех данных кроме items.
var q_no_items = `
	SELECT o.order_uid, o.track_number, o.entry,
	d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
	p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee,
	o.locale, o.internal_signature, o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard
	FROM orders o
	LEFT JOIN delivery d
		ON o.delivery_id = d.id
	LEFT JOIN payment p
		ON o.payment_id = p.id
	WHERE o.order_uid = $1
`
*/

// Возвращает один заказ из базы данных по его order_uid.
func SelectOrderById(db *pgx.Conn, order_uid string) models.Order {
	var order models.Order

	// Получение данных из orders, delivery, payments. Порядок в Scan должен соответствовать порядку полей в запросе.
	err := db.QueryRow(context.Background(), q_orders, order_uid).Scan(
		&order.Order_uid,
		&order.Track_number,
		&order.Entry,
		&order.Locale,
		&order.Internal_signature,
		&order.Customer_id,
		&order.Delivery_service,
		&order.Shardkey,
		&order.Sm_id,
		&order.Date_created,
		&order.Oof_shard,
	)

	if err != nil {
		log.Fatalf("QueryRow failed: %v\n", err)
	}

	err = db.QueryRow(context.Background(), q_delivery, order_uid).Scan(
		&order.Delivery.Name,
		&order.Delivery.Phone,
		&order.Delivery.Zip,
		&order.Delivery.City,
		&order.Delivery.Address,
		&order.Delivery.Region,
		&order.Delivery.Email,
	)

	if err != nil {
		log.Fatalf("QueryRow failed: %v\n", err)
	}

	err = db.QueryRow(context.Background(), q_payments, order_uid).Scan(
		&order.Payment.Transaction,
		&order.Payment.Request_id,
		&order.Payment.Currency,
		&order.Payment.Provider,
		&order.Payment.Amount,
		&order.Payment.Payment_dt,
		&order.Payment.Bank,
		&order.Payment.Delivery_cost,
		&order.Payment.Goods_total,
		&order.Payment.Custom_fee,
	)

	if err != nil {
		log.Fatalf("QueryRow failed: %v\n", err)
	}

	// Получение данных из items. Порядок в Scan должен соответствовать порядку полей в запросе.
	rows, err := db.Query(context.Background(), q_items, order_uid)
	if err != nil {
		log.Fatalf("QueryRow failed: %v\n", err)
	}

	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var i models.Item
		err = rows.Scan(
			&i.Chrt_id,
			&i.Track_number,
			&i.Price,
			&i.Rid,
			&i.Name,
			&i.Sale,
			&i.Size,
			&i.Total_price,
			&i.Nm_id,
			&i.Brand,
			&i.Status,
		)
		if err != nil {
			log.Fatalf("QueryRow failed: %v\n", err)
		}
		items = append(items, i)
	}

	if rows.Err() != nil {
		log.Fatalf("QueryRow failed: %v\n", err)
	}

	order.Items = items

	log.Println("SelectOrderById: success")
	return order
}
