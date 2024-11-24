package db_model

import (
	"context"
	"log/slog"
)

/*
Строка для создания таблицы с информацией о заказах.

Таблица связывается с одной записью в таблице "delivery" через поле delivery_id и одной записью в таблице "payment" через поле payment_id.
Связывается с одной или несколькими записями в таблице items через таблицу itemsbind полем order_uid.
*/
var createOrders = `CREATE TABLE IF NOT EXISTS orders (
	id INTEGER PRIMARY KEY,
	order_uid VARCHAR(255) NOT NULL UNIQUE,
	track_number VARCHAR(255),
	entry VARCHAR(255),
	delivery_id INTEGER,
	payment_id INTEGER,
	locale VARCHAR(255),
	internal_signature VARCHAR(255),
	customer_id VARCHAR(255),
	delivery_service VARCHAR(255),
	shardkey VARCHAR(255),
	sm_id INTEGER,
	date_created VARCHAR(255),
    oof_shard VARCHAR(255)
	)`

/*
Строка для создания таблицы с информацией о доставке заказов.

Каждая запись в таблице orders связывается с одной записью в этой таблице.
*/
var createDelivery = `CREATE TABLE IF NOT EXISTS delivery (
	id serial PRIMARY KEY, 
	name VARCHAR(255),
	phone VARCHAR(255),
	zip VARCHAR(255),
	city VARCHAR(255),
	address VARCHAR(255),
	region VARCHAR(255),
	email VARCHAR(255)
	)`

/*
Строка для создания таблицы с информацией о платежах.

Каждая запись в таблице orders связывается с одной записью в этой таблице.
*/
var createPayments = `CREATE TABLE IF NOT EXISTS payments (
	id serial PRIMARY KEY,
	transaction VARCHAR(255),
	request_id VARCHAR(255),
	currency VARCHAR(255),
	provider VARCHAR(255),
	amount INTEGER,
	payment_dt INTEGER,
	bank VARCHAR(255),
	delivery_cost INTEGER,
	goods_total INTEGER,
	custom_fee INTEGER
	)`

/*
Строка для создания таблицы с информацией о товаре.

Через таблицу itemsbind полем id несколько предметов из этой таблицы могут быть связаны с одним заказом из таблицы orders.
*/
var createItems = `CREATE TABLE IF NOT EXISTS items (
	id serial PRIMARY KEY,
	chrt_id INTEGER,
	track_number VARCHAR(255),
	price INTEGER,
	rid VARCHAR(255),
	name VARCHAR(255),
	sale INTEGER,
	size VARCHAR(255),
	total_price INTEGER,
	nm_id INTEGER,
	brand VARCHAR(255),
	status INTEGER
	)`

/*
Строка для создания таблицы Itemsbind.

Эта таблица связывает один заказ из таблицы orders с несколькими предметами из таблицы items.
*/
var createItemsbind = `CREATE TABLE IF NOT EXISTS itemsbind (
	id serial PRIMARY KEY,
	order_id INTEGER,
	item_id INTEGER
	)`

/*
Создаёт отсутствующие таблицы в базе данных.

Не использует индексы из-за потенциально большого количества операций записи.
*/
func (db *PostgresDB) CreateTables() {
	queries := []string{createOrders, createDelivery, createPayments, createItems, createItemsbind}

	for _, q := range queries {
		_, err := db.DB.Exec(context.Background(), q)
		if err != nil {
			slog.Error("Failed to create table " + err.Error())
		}
	}
}
