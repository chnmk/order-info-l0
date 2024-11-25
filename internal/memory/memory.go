package memory

import (
	"sync"

	cfg "github.com/chnmk/order-info-l0/internal/config"
	"github.com/chnmk/order-info-l0/internal/models"
)

/*
TODO: написать объяснительную.

Зачем нужны новые ключи, почему нельзя получать данные по orders_uid:
	- явный порядок появления данных
	- в веб-интерфейсе тоже получать по orders_uid? открывать логи и смотреть копировать оттуда заказа?
	или проще всё-таки написать id=1 и на этом всё
*/

var once sync.Once

type MemStore struct {
	mu         sync.Mutex
	currentkey int
	orders     map[int]models.Order
}

func (m *MemStore) Init() {
	m.orders = make(map[int]models.Order)
}

func NewStorage() {
	once.Do(func() {
		cfg.Data = &MemStore{}
		cfg.Data.Init()
	})
}
