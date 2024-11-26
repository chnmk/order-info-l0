package db_model

import (
	"context"
	"log/slog"

	"github.com/chnmk/order-info-l0/internal/models"
	"github.com/jackc/pgx/v5"
)

// Строка для получения данных из таблицы orders.
//
// Не получает данные из таблиц delivery, payments и items.
var q_orders = `
	SELECT id, order_uid, track_number, entry,
	locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
	FROM orders 
	WHERE id = @id
`

// Строка для получения данных из таблицы delivery.
var q_delivery = `
	SELECT name, phone, zip, city, address, region, email
	FROM delivery d
	LEFT JOIN orders o
		ON d.id = o.delivery_id
	WHERE o.id = @id
`

// Строка для получения данных из таблицы payments.
var q_payments = `
	SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
	FROM payments p
	LEFT JOIN orders o
		ON p.id = o.payment_id
	WHERE o.id = @id
`

// Строка для получения данных из items.
var q_items = `
	SELECT chrt_id, i.track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
	FROM items i
	LEFT JOIN itemsbind ib
		on i.id = ib.item_id
	LEFT JOIN orders o
		on ib.order_id = o.id
	WHERE o.id = @id
`

// Возвращает один заказ из базы данных по его order_uid.
func (db *PostgresDB) SelectOrderById(id int) (int, models.Order) {
	var key int
	var order models.Order
	args := pgx.NamedArgs{"id": id}

	// Получение данных из orders, delivery, payments. Порядок в Scan должен соответствовать порядку полей в запросе.
	err := db.Conn.QueryRow(context.Background(), q_orders, args).Scan(
		&key,
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
		slog.Error("QueryRow failed: " + err.Error())
	}

	err = db.Conn.QueryRow(context.Background(), q_delivery, args).Scan(
		&order.Delivery.Name,
		&order.Delivery.Phone,
		&order.Delivery.Zip,
		&order.Delivery.City,
		&order.Delivery.Address,
		&order.Delivery.Region,
		&order.Delivery.Email,
	)

	if err != nil {
		slog.Error("QueryRow failed: " + err.Error())
	}

	err = db.Conn.QueryRow(context.Background(), q_payments, args).Scan(
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
		slog.Error("QueryRow failed: " + err.Error())
	}

	// Получение данных из items. Порядок в Scan должен соответствовать порядку полей в запросе.
	rows, err := db.Conn.Query(context.Background(), q_items, args)
	if err != nil {
		slog.Error("QueryRow failed: " + err.Error())
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
			slog.Error("QueryRow failed: " + err.Error())
		}
		items = append(items, i)
	}

	if rows.Err() != nil {
		slog.Error("QueryRow failed: " + err.Error())
	}

	order.Items = items

	slog.Info("SelectOrderById: success")
	return key, order
}
