package db_model

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/chnmk/order-info-l0/internal/models"
	"github.com/jackc/pgx/v5"
)

// Добавляет в таблицу данные о заказе и связывает их с самыми новыми записями в delivery и payment.
//
// Запросы q_insert_delivery и q_insert_item должны быть выполнены раньше.
var q_insert_order = `
	INSERT INTO orders(id, order_uid, track_number, entry, locale, internal_signature, customer_id, 
		delivery_service, shardkey, sm_id, date_created, oof_shard, delivery_id, payment_id)
	VALUES (@id, @order_uid, @track_number, @entry, @locale, @internal_signature, @customer_id, 
		@delivery_service, @shardkey, @sm_id, @date_created, @oof_shard, @delivery_id, @payment_id)
	RETURNING id
`

// Добавляет данные в таблицу delivery.
var q_insert_delivery = `
	INSERT INTO delivery(name, phone, zip, city, address, region, email)
	VALUES (@name, @phone, @zip, @city, @address, @region, @email)
	RETURNING id
`

// Добавляет данные в таблицу payments.
var q_insert_payment = `
	INSERT INTO payments(transaction, request_id, currency, provider, amount, payment_dt, bank, 
	delivery_cost, goods_total, custom_fee)
	VALUES (@transaction, @request_id, @currency, @provider, @amount, @payment_dt, @bank, 
	@delivery_cost, @goods_total, @custom_fee)
	RETURNING id
`

// Добавляет данные в таблицу items.
var q_insert_item = `
	INSERT INTO items(chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
	VALUES (@chrt_id, @track_number, @price, @rid, @name, @sale, @size, @total_price, @nm_id, @brand, @status)
	RETURNING id
`

// Добавляет данные в таблицу itemsbind.
//
// Запрос должен быть выполнен после запроса q_insert_item.
var q_insert_itemsbind = `
	INSERT INTO itemsbind(order_id, item_id)
	VALUES (@order_id, @item_id)
	RETURNING id
`

// Пробует добавить заказ в БД, возвращает ошибку только в случае если заказ с таким order_uid уже существует.
func (db *PostgresDB) InsertOrder(key int, m []byte, ctx context.Context) error {
	// Чтобы не передавать кучу данных во все вторичные реализации БД, приходится опять декодировать заказ здесь.
	var order models.Order
	err := json.Unmarshal(m, &order)
	if err != nil {
		slog.Error("Failed to marshal order during db insert operation: " + err.Error())
		return nil
	}

	slog.Info("starting insert transaction...")

	tx, err := db.DB.Begin(ctx)
	if err != nil {
		slog.Error("Failed to begin transaction: " + err.Error())
		return nil
	}
	defer tx.Rollback(ctx)

	args := pgx.NamedArgs{
		"name":    order.Delivery.Name,
		"phone":   order.Delivery.Phone,
		"zip":     order.Delivery.Zip,
		"city":    order.Delivery.City,
		"address": order.Delivery.Address,
		"region":  order.Delivery.Region,
		"email":   order.Delivery.Email,
	}
	row := tx.QueryRow(context.Background(), q_insert_delivery, args)

	var delivery_id int
	err = row.Scan(&delivery_id)
	if err != nil {
		slog.Error("Failed to insert data: " + err.Error())
		return nil
	}

	args = pgx.NamedArgs{
		"transaction":   order.Payment.Transaction,
		"request_id":    order.Payment.Request_id,
		"currency":      order.Payment.Currency,
		"provider":      order.Payment.Provider,
		"amount":        order.Payment.Amount,
		"payment_dt":    order.Payment.Payment_dt,
		"bank":          order.Payment.Bank,
		"delivery_cost": order.Payment.Delivery_cost,
		"goods_total":   order.Payment.Goods_total,
		"custom_fee":    order.Payment.Custom_fee,
	}
	row = tx.QueryRow(context.Background(), q_insert_payment, args)

	var payment_id int
	err = row.Scan(&payment_id)
	if err != nil {
		slog.Error("Failed to insert data: " + err.Error())
		return nil
	}

	args = pgx.NamedArgs{
		"id":                 key,
		"order_uid":          order.Order_uid,
		"track_number":       order.Track_number,
		"entry":              order.Entry,
		"locale":             order.Locale,
		"internal_signature": order.Internal_signature,
		"customer_id":        order.Customer_id,
		"delivery_service":   order.Delivery_service,
		"shardkey":           order.Shardkey,
		"sm_id":              order.Sm_id,
		"date_created":       order.Date_created,
		"oof_shard":          order.Oof_shard,
		"delivery_id":        delivery_id,
		"payment_id":         payment_id,
	}
	row = tx.QueryRow(context.Background(), q_insert_order, args)

	var order_id int
	err = row.Scan(&order_id)
	if err != nil {
		if strings.Contains(err.Error(), "(SQLSTATE 23505)") {
			return err
		} else {
			slog.Error("Failed to insert data: " + err.Error())
		}
	}

	for _, i := range order.Items {
		args = pgx.NamedArgs{
			"chrt_id":      i.Chrt_id,
			"track_number": i.Track_number,
			"price":        i.Price,
			"rid":          i.Rid,
			"name":         i.Name,
			"sale":         i.Sale,
			"size":         i.Size,
			"total_price":  i.Total_price,
			"nm_id":        i.Nm_id,
			"brand":        i.Brand,
			"status":       i.Status,
		}
		row = tx.QueryRow(context.Background(), q_insert_item, args)

		var item_id int
		err = row.Scan(&item_id)
		if err != nil {
			slog.Error("Failed to insert data: " + err.Error())
			return nil
		}

		args = pgx.NamedArgs{
			"order_id": order_id,
			"item_id":  item_id,
		}
		row = tx.QueryRow(context.Background(), q_insert_itemsbind, args)

		// bind_id используется только для проверки ответа
		var bind_id int
		err = row.Scan(&bind_id)
		if err != nil {
			slog.Error("Failed to insert data: " + err.Error())
			return nil
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		slog.Error("Failed to commit transaction: " + err.Error())
		return nil
	}

	slog.Info("insert transaction committed")
	return nil
}
