package memory

import (
	"context"
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

// Обрабатывает сообщение в горутине.
func HandleMessage(m []byte) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		var order models.Order

		err := json.Unmarshal(m, &order)
		if err != nil {
			slog.Info("failed to unmarshal, skipping")
			return
		}

		if ok := ValidateMsg(order); !ok {
			slog.Info("failed to validate, skipping")
			return
		}

		err = database.DB.InsertOrder(DATA.currentkey, m, context.Background())
		if err != nil {
			slog.Error("failed to add order: order already exists")
			return
		}

		DATA.AddOrder(order)
	}()

	wg.Wait()
	slog.Info("order handling finished")
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
