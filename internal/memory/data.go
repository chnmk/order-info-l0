package memory

import (
	"sync"

	"github.com/chnmk/order-info-l0/internal/models"
)

var DATA MemStore

type MemStore struct {
	mu     sync.Mutex
	orders map[int]models.Order
}

func (d *MemStore) Init() {
	DATA.orders = make(map[int]models.Order)
}

func (d *MemStore) AddOrder(key int, value models.Order) {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, ok := d.orders[key]
	if !ok {
		d.orders[key] = value
	}
}
