package memory

import (
	"sync"

	"github.com/chnmk/order-info-l0/internal/models"
)

/*
TODO: написать объяснительную.

Зачем нужны новые ключи, почему нельзя получать данные по orders_uid:
	- явный порядок появления данных
	- в веб-интерфейсе тоже получать по orders_uid? открывать логи и смотреть копировать оттуда заказа?
	или проще всё-таки написать id=1 и на этом всё
*/

var (
	DATA Storage
	once sync.Once
)

type Storage interface {
	init()
	HandleMessage([]byte)
	AddOrder(models.Order)
	Read(int) models.Order
	RestoreData()
}

type MemStore struct {
	mu         sync.Mutex
	currentkey int
	orders     map[int]models.Order
}

func (d *MemStore) init() {
	d.orders = make(map[int]models.Order)
}

func NewStorage() Storage {
	once.Do(func() {
		DATA = &MemStore{}
		DATA.init()
	})

	return DATA
}
