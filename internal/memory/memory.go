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

Апдейт: оно теперь массив, так что не очень релевантно, но тем не менее.
*/

var once sync.Once

type MemStore struct {
	mu     sync.Mutex
	orders []models.OrderStorage
}

func NewStorage(m models.Storage) models.Storage {
	once.Do(func() {
		m = &MemStore{}

		m.RestoreData()
	})

	return m
}
