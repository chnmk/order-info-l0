package memory

import (
	"log/slog"
	"slices"
	"sync"

	"github.com/chnmk/order-info-l0/internal/database"
	"github.com/chnmk/order-info-l0/internal/models"
	"github.com/jackc/pgx/v5"
)

var DATA MemStore

type MemStore struct {
	mu         sync.Mutex
	currentkey int
	orders     map[int]models.Order
}

func (d *MemStore) Init() {
	DATA.orders = make(map[int]models.Order)
}

// Добавляет заказ value в память.
func (d *MemStore) AddOrder(value models.Order) {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, ok := d.orders[d.currentkey]
	if !ok {
		d.orders[d.currentkey] = value
		slog.Info("added order to memory storage")
	} else {
		slog.Error("failed to add order to memory storage: key already exists")
	}

	database.InsertOrder(database.DB, value, d.currentkey)
	d.currentkey++
}

// Забирает все данные из БД, устанавливает значение currentkey на максимальное id заказа из БД.
func (d *MemStore) RestoreData(db *pgx.Conn) {
	d.mu.Lock()
	defer d.mu.Unlock()

	slog.Info("restoring data from DB...")

	ids := database.GetOrdersIDs(db)
	if len(ids) == 0 {
		slog.Info("no data found in DB, restoring canceled")
		return
	}

	// TODO: у нас будут интерфейсы, так что скорее всего че-т поумнее придумаем.
	for _, id := range ids {
		key, order := database.SelectOrderById(db, id)
		d.orders[key] = order
	}

	d.currentkey = slices.Max(ids)

	slog.Info("data successfully restored")
}
