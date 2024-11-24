package memory

import (
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/chnmk/order-info-l0/internal/database"
	"github.com/chnmk/order-info-l0/internal/models"
)

/*
TODO: написать объяснительную.

Зачем нужны новые ключи, почему нельзя получать данные по orders_uid:
	- явный порядок появления данных
	- в веб-интерфейсе тоже получать по orders_uid? открывать логи и смотреть копировать оттуда заказа?
	или проще всё-таки написать id=1 и на этом всё
*/

var DATA MemStore

type MemStore struct {
	mu         sync.Mutex
	currentkey int
	orders     map[int]models.Order
}

func (d *MemStore) Init() {
	DATA.orders = make(map[int]models.Order)
}

func UnmarshalBytes(m []byte) {
	var order models.Order
	err := json.Unmarshal(m, &order)
	if err != nil {
		slog.Info("failed to unmarshal, skipping")
	} else {
		if ok := ValidateMsg(order); !ok {
			slog.Info("failed to validate, skipping")
		} else {
			// slog.Info(order)
			DATA.AddOrder(order)
		}
	}
}

// Проверяет что нужные поля не пустые и соответствуют нашим требованиям.
//
// Пока что нам точно нужны те данные, которые выводятся в веб-интерфейсе.
func ValidateMsg(order models.Order) bool {
	if order.Order_uid == "" ||
		order.Delivery.Name == "" ||
		order.Delivery.City == "" ||
		order.Delivery.Address == "" ||
		order.Delivery.Phone == "" ||
		len(order.Items) < 1 {

		return false
	}

	for _, i := range order.Items {
		if i.Chrt_id == 0 ||
			i.Name == "" ||
			i.Total_price == 0 {
			return false
		}
	}

	return true
}

// Добавляет заказ value в память.
func (d *MemStore) AddOrder(value models.Order) {
	d.mu.Lock()
	defer d.mu.Unlock()

	slog.Info("adding order to memory storage...")

	_, ok := d.orders[d.currentkey]
	if ok {
		slog.Error("failed to add order: id already exists")
		return
	}

	err := database.DB.InsertOrder(value, d.currentkey)
	if err != nil {
		slog.Error("failed to add order: order already exists")
		return
	}

	d.orders[d.currentkey] = value
	slog.Info("added order to memory storage")
	d.currentkey++
}
