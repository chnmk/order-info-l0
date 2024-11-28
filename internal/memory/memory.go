package memory

import (
	"sync"

	cfg "github.com/chnmk/order-info-l0/internal/config"
	"github.com/chnmk/order-info-l0/internal/models"
)

// Ожидается, что хранилище будет создано только один раз.
// Подстраховка от неожиданного поведения.
var once sync.Once

// Имплементация интерфейса models.Storage
type MemStore struct {
	mu     sync.Mutex
	orders []models.OrderStorage
	maxId  int
}

// Возвращает новое хранилище. При необходимости восстанавливает данные из БД.
func NewStorage(m models.Storage) models.Storage {
	once.Do(func() {
		m = &MemStore{}

		if cfg.RestoreData {
			m.RestoreData()
		}
	})

	return m
}
