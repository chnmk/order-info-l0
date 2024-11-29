package memory

import (
	"context"
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
func NewStorage(ctx context.Context, m models.Storage) models.Storage {
	once.Do(func() {
		m = &MemStore{}

		// Восстанавливает данные.
		if cfg.RestoreData {
			m.RestoreData(ctx)
		}

		// Создаёт пул обработчиков сообщений.
		for i := 0; i < cfg.MemoryHandlerGoroutines; i++ {
			cfg.ExitWg.Add(1)
			go cfg.Data.HandleMessage(ctx)
		}

		// Создаёт обработчик устаревших сообщений.
		cfg.ExitWg.Add(1)
		go cfg.Data.ClearData(ctx)
	})

	return m
}
